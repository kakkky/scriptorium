package main

import (
	"sync"
)

// Hub は Turbo Stream 断片を複数の subscriber に broadcast する pub/sub。
// SSE と WebSocket の両方で共有される。
type Hub struct {
	mu   sync.RWMutex
	subs map[chan string]struct{}
}

func NewHub() *Hub {
	return &Hub{subs: make(map[chan string]struct{})}
}

// Subscribe は buffered チャネルを返す。呼び出し側が Unsubscribe まで責任を持つ。
func (h *Hub) Subscribe() chan string {
	ch := make(chan string, 16)
	h.mu.Lock()
	h.subs[ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

func (h *Hub) Unsubscribe(ch chan string) {
	h.mu.Lock()
	if _, ok := h.subs[ch]; ok {
		delete(h.subs, ch)
		close(ch)
	}
	h.mu.Unlock()
}

// Broadcast は全 subscriber に非同期送信する。
// 遅い subscriber の buffer が詰まっていたらメッセージを drop して他をブロックしない。
func (h *Hub) Broadcast(msg string) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for ch := range h.subs {
		select {
		case ch <- msg:
		default:
			// buffer full: drop to prevent slow subscribers from blocking others
		}
	}
}

var hub = NewHub()
