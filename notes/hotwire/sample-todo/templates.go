package main

import (
	"bytes"
	"embed"
	"html/template"
	"net/http"
	"strings"
)

//go:embed views/*.html
var viewsFS embed.FS

func mustParse(files ...string) *template.Template {
	return template.Must(template.ParseFS(viewsFS, files...))
}

var (
	indexTmpl        = mustParse("views/layout.html", "views/index.html", "views/_item.html")
	aboutTmpl        = mustParse("views/layout.html", "views/about.html")
	itemFrameTmpl    = mustParse("views/_item.html")
	editFrameTmpl    = mustParse("views/_edit.html")
	streamCreateTmpl = mustParse("views/stream_create.html", "views/_item.html")
	streamToggleTmpl = mustParse("views/stream_toggle.html", "views/_item.html")
	streamDeleteTmpl = mustParse("views/stream_delete.html")
)

// wantsStream は Turbo からの form 送信で Turbo Stream 応答が受け入れ可能かを判定する。
func wantsStream(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept"), "text/vnd.turbo-stream.html")
}

// writeStreamHeader は Turbo Stream 応答の Content-Type をセットする。
func writeStreamHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/vnd.turbo-stream.html; charset=utf-8")
}

func render(w http.ResponseWriter, t *template.Template, name string, data any) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	}
	if err := t.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// renderStreamToString は Turbo Stream 断片を文字列として組み立てる (broadcast 用)。
func renderStreamToString(t *template.Template, name string, data any) string {
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, name, data); err != nil {
		return ""
	}
	return buf.String()
}
