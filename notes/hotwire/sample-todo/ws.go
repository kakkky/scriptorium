package main

import (
	"context"
	"net/http"
	"time"

	"github.com/coder/websocket"
)

// handleWS は Turbo Stream を WebSocket で配信する。
// クライアントは HTML に <turbo-stream-source src="ws://.../events/ws"> を置くことで接続する。
func handleWS(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"localhost:8080", "127.0.0.1:8080"},
	})
	if err != nil {
		return
	}
	defer c.CloseNow()

	ctx := r.Context()
	ch := hub.Subscribe()
	defer hub.Unsubscribe(ch)

	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return
			}
			writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			err := c.Write(writeCtx, websocket.MessageText, []byte(msg))
			cancel()
			if err != nil {
				return
			}
		case <-ctx.Done():
			c.Close(websocket.StatusNormalClosure, "")
			return
		}
	}
}
