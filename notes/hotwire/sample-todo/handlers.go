package main

import (
	"net/http"
	"strconv"
	"strings"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	render(w, indexTmpl, "layout", map[string]any{
		"Todos": store.List(),
	})
}

func handleAbout(w http.ResponseWriter, r *http.Request) {
	render(w, aboutTmpl, "layout", nil)
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
	title := strings.TrimSpace(r.FormValue("title"))
	if title == "" {
		http.Error(w, "title required", http.StatusBadRequest)
		return
	}
	todo := store.Create(title)

	// broadcast: 全 subscriber (sender 自身も含む) に配信
	hub.Broadcast(renderStreamToString(streamCreateTmpl, "stream_create", todo))

	if wantsStream(r) {
		// Turbo クライアント: 空 stream レスポンスで完了させ、実際の更新は SSE/WS 経由
		writeStreamHeader(w)
		return
	}
	// non-Turbo クライアント (JS 無効等): 通常のページ遷移
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleShow(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	todo, ok := store.Get(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	render(w, itemFrameTmpl, "todo_item", todo)
}

func handleEdit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	todo, ok := store.Get(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	render(w, editFrameTmpl, "todo_edit", todo)
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	title := strings.TrimSpace(r.FormValue("title"))
	done := r.FormValue("done") == "1"
	todo, ok := store.Update(id, title, done)
	if !ok {
		http.NotFound(w, r)
		return
	}

	// broadcast: 一覧の該当行を全員の画面で更新する
	hub.Broadcast(renderStreamToString(streamToggleTmpl, "stream_toggle", todo))

	// sender の Frame は表示モードに戻す (Frame 応答)
	render(w, itemFrameTmpl, "todo_item", todo)
}

func handleToggle(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	todo, ok := store.Toggle(id)
	if !ok {
		http.NotFound(w, r)
		return
	}

	hub.Broadcast(renderStreamToString(streamToggleTmpl, "stream_toggle", todo))

	if wantsStream(r) {
		writeStreamHeader(w)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	if !store.Delete(id) {
		http.NotFound(w, r)
		return
	}

	hub.Broadcast(renderStreamToString(streamDeleteTmpl, "stream_delete", struct{ ID int }{ID: id}))

	if wantsStream(r) {
		writeStreamHeader(w)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
