package main

import (
	"fmt"
	"net/http"
	"strings"
)

// handleSSE は Turbo Stream を Server-Sent Events で配信する。
// クライアントは HTML に <turbo-stream-source src="/events/sse"> を置くことで接続する。
func handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	ch := hub.Subscribe()
	defer hub.Unsubscribe(ch)

	// initial comment to flush headers immediately
	fmt.Fprint(w, ": connected\n\n")
	flusher.Flush()

	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return
			}
			// SSE 規約: 複数行を送るときは各行に "data: " 接頭辞を付け、末尾を \n\n で区切る
			for _, line := range strings.Split(msg, "\n") {
				fmt.Fprintf(w, "data: %s\n", line)
			}
			fmt.Fprint(w, "\n")
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}
