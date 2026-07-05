# hotwire-go sample: in-memory Todo

Go の `net/http` + `html/template` で **Turbo Drive / Turbo Frame / Turbo Stream (form 応答 + SSE broadcast + WebSocket broadcast) / Stimulus** の 5 要素を触るサンプル。データは in-memory (プロセス再起動で消える)。

WebSocket 部分だけ `github.com/coder/websocket` に依存。それ以外は Go 標準のみ。

## 起動

```
cd sample-todo
go run .
```

`http://localhost:8080/` を開く。

## 4 要素がどこで効いているか

### Turbo Drive

- **どこ**: ナビゲーションの `<a href="/">` / `<a href="/about">` のリンク遷移
- **確認方法**: DevTools の Network タブを開いてリンクをクリック。フルリロードではなく `fetch` でページが取れて、body が差し替わる
- **サーバー側の仕込み**: なし。turbo.js を読み込むだけで自動的に `<a>` / `<form>` を intercept する

### Turbo Frame

- **どこ**: 各 Todo アイテムを `<turbo-frame id="todo_N">` で囲んでいる (`views/_item.html`)。「編集」リンクを押すとその行だけが編集フォームに差し替わる
- **仕組み**:
  - 「編集」の `<a href="/todos/1/edit">` はフレーム内リンクなので、Turbo はフレーム内でナビゲーションする
  - サーバーが `/todos/1/edit` に `<turbo-frame id="todo_1">…edit form…</turbo-frame>` を返す (`handleEdit`)
  - 同じ id のフレームだけがクライアント側で差し替えられる。他の Todo 行や nav は動かない
- **サーバー側の仕込み**: `<turbo-frame id="…">` を含む HTML を返すだけ。判定不要

### Turbo Stream

- **どこ**: 追加 / チェックトグル / 削除
- **仕組み**:
  - form 送信時、Turbo が `Accept: text/vnd.turbo-stream.html` を積む (`wantsStream` で判定)
  - サーバーは `Content-Type: text/vnd.turbo-stream.html` で `<turbo-stream action="…" target="…">` を返す
  - 使用している action:
    - 追加 → `append` (`views/stream_create.html`) → `#todos` の末尾に item を挿入
    - トグル → `replace` (`views/stream_toggle.html`) → `#todo_N` を新しい item で置換
    - 削除 → `remove` (`views/stream_delete.html`) → `#todo_N` を DOM から削除
- **サーバー側の仕込み**:
  - `Accept` ヘッダ判定 (`wantsStream`)
  - `Content-Type` セット (`writeStreamHeader`)
  - `<turbo-stream>` テンプレート
  - non-Turbo フォールバック (JS 無効時) として `303 See Other` でリダイレクト

### Turbo Stream (broadcast: SSE / WebSocket)

- **どこ**: layout に `<turbo-stream-source src="/events/sse">` を置いてある。ページを開いたブラウザは自動で SSE 接続を張り、以降のサーバー側変更を全員リアルタイムで受け取る
- **確認方法**: **同じ URL を 2 つのブラウザ (or シークレットウィンドウ) で開き、片方で追加/トグル/削除する → もう片方も自動更新される**
- **仕組み**:
  - `hub.go` の `Hub` (pub/sub) に SSE / WS ハンドラが Subscribe
  - `handleCreate` / `handleToggle` / `handleDelete` / `handleUpdate` が `hub.Broadcast(...)` で全 subscriber に `<turbo-stream>` HTML 断片を配信
  - **sender 自身も自分の SSE/WS 経由で更新を受け取る** → form 応答は空の Stream レスポンスに留めている (重複防止)
- **経路の切り替え**: layout の `<turbo-stream-source>` の `src` を書き換える
  - SSE: `src="/events/sse"` (デフォルト)
  - WebSocket: `src="ws://localhost:8080/events/ws"` (`ws://` で始まると Turbo が自動で WS 接続に切り替える)

### Stimulus

- **どこ**: 追加フォームの「クリア」ボタン
- **仕組み**:
  - form に `data-controller="clear-input"`
  - input に `data-clear-input-target="input"`
  - button に `data-action="click->clear-input#clear"`
  - `ClearInputController#clear()` が呼ばれて input を空にする
- **サーバー側の仕込み**: `data-*` 属性を HTML に書くだけ。**JS 側でしかやることがない** (Stimulus のカバー面が薄い理由)

## サーバー側で自分で書いている「Hotwire 用の定型」

package 化する意義がある部分。今後 hotwire-go で吸収したい。

- リクエスト判定: `wantsStream(r)` (`Accept` ヘッダで Stream 受入可否を判定)
- レスポンスヘッダ: `writeStreamHeader(w)` (`Content-Type: text/vnd.turbo-stream.html`)
- `<turbo-stream action=… target=…>` の HTML 組み立て (現状はテンプレート内でベタ書き)
- form 送信後の 303 リダイレクト (non-Turbo フォールバック)
- **SSE ハンドラ**: ヘッダ 3 種 (`text/event-stream` / `no-cache` / `keep-alive`) + `http.Flusher` + `data:` 行組み立て + context 切断検知
- **WebSocket ハンドラ**: Origin 検証 + subscribe/unsubscribe + 書き込みタイムアウト
- **Pub/Sub hub**: subscribers 管理 + 遅い subscriber を drop する非ブロッキング配信 + goroutine 安全性

## 標準機能で書いている部分 (package で覆わない領域)

- レイアウト合成: `views/layout.html` の `{{block "content" .}}` + 各ページの `{{define "content"}}` (Go 標準)
- partial の呼び出し: `{{template "todo_item" .}}` (Go 標準)
- Turbo Frame の HTML: `<turbo-frame id="…">` を書くだけ (規約なし)

## ファイル構成

```
sample-todo/
├── go.mod
├── main.go                  // エントリポイント (mux 設定 + assets 配信)
├── store.go                 // Todo / Store (in-memory)
├── templates.go             // template パース + Stream helper
├── handlers.go              // HTTP ハンドラ
├── hub.go                   // Pub/Sub hub (SSE と WS が共有)
├── sse.go                   // SSE ハンドラ
├── ws.go                    // WebSocket ハンドラ
├── assets/
│   └── app.js               // Turbo import + Stimulus controller + form リセット
├── views/
│   ├── layout.html          // レイアウト (block "content") + <turbo-stream-source>
│   ├── index.html           // Todo 一覧 (define "content")
│   ├── about.html           // About ページ (define "content")
│   ├── _item.html           // Todo 単体 (turbo-frame 込み)
│   ├── _edit.html           // 編集フォーム (turbo-frame 込み)
│   ├── stream_create.html   // append stream
│   ├── stream_toggle.html   // replace stream
│   └── stream_delete.html   // remove stream
└── README.md
```

### 各 Go ファイルの責務

- **main.go**: エントリポイント。ServeMux 設定、`/assets/*.js` の静的配信、`/events/sse` と `/events/ws` の登録、シードデータ投入
- **store.go**: `Todo` 型と `Store` (`map + RWMutex`)。パッケージ変数 `store` を公開
- **templates.go**: `//go:embed` で views を読み、テンプレート変数を初期化。`wantsStream(r)` / `writeStreamHeader(w)` / `render(w, t, name, data)` / `renderStreamToString(t, name, data)` の共通 helper
- **handlers.go**: 8 個の HTTP ハンドラ (`handleIndex` / `handleAbout` / `handleCreate` / `handleShow` / `handleEdit` / `handleUpdate` / `handleToggle` / `handleDelete`)。状態変更系は `hub.Broadcast(...)` で全 subscriber に配信
- **hub.go**: `Hub` (Pub/Sub)。buffered channel と `sync.RWMutex` で構築。遅い subscriber を drop する非ブロッキング配信
- **sse.go**: `handleSSE`。`text/event-stream` ヘッダ + `http.Flusher` + `data:` 行組み立て + `r.Context().Done()` で切断検知
- **ws.go**: `handleWS`。`github.com/coder/websocket` で `Accept` → subscribe → 5 秒タイムアウトで書き込み

## 使ってみるとわかること

- Todos ページで:
  - Turbo Drive: /about リンクをクリック。ページ全体が置き換わるが scroll 位置以外の JS 状態は維持される
  - Turbo Frame: 「編集」→ その行だけ編集フォームに差し替わる。「キャンセル」で戻る
  - Turbo Stream: 追加ボタンで一覧末尾に item が append される。チェックボックスで打ち消し線がつく (replace)。削除ボタンで DOM から消える (remove)
  - Stimulus: 追加フォームの「クリア」ボタンで入力欄が空になる

## 依存

- Go 1.22 以上 (`http.ServeMux` のパターン matching, `PathValue`, `GET /{$}` 記法)
- Turbo / Stimulus は esm.sh の CDN から動的 import (Node.js / npm 不要)

## 制限

- in-memory なので再起動でデータが消える
- WebSocket の Origin 検証は `localhost:8080` / `127.0.0.1:8080` のみ許可 (本番前提なら `OriginPatterns` を絞る)
- broadcast は「全員に配信」しかしない。個別チャンネル (Todo リスト A / B で分ける等) は実装していない
