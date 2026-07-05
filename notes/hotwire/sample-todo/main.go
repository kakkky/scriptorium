package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed assets/*.js
var assetsFS embed.FS

func main() {
	store.Create("Turbo Drive を試す")
	store.Create("Turbo Frame を試す")
	store.Create("Turbo Stream を試す")

	subAssets, err := fs.Sub(assetsFS, "assets")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", handleIndex)
	mux.HandleFunc("GET /about", handleAbout)
	mux.HandleFunc("POST /todos", handleCreate)
	mux.HandleFunc("GET /todos/{id}", handleShow)
	mux.HandleFunc("GET /todos/{id}/edit", handleEdit)
	mux.HandleFunc("POST /todos/{id}", handleUpdate)
	mux.HandleFunc("POST /todos/{id}/toggle", handleToggle)
	mux.HandleFunc("POST /todos/{id}/delete", handleDelete)
	mux.HandleFunc("GET /events/sse", handleSSE)
	mux.HandleFunc("GET /events/ws", handleWS)
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(subAssets))))

	addr := ":8080"
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
