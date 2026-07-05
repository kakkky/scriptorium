golang製のhotwire packageは今の所なさそう。

いい根拠として、goでhotwire使ってみた！というリポジトリが散見されるが、色々頑張って実装してそうな感じがある。

https://github.com/cldmstr/gogohotwire

https://github.com/wolfeidau/hotwire-golang-website



以下、それぞれのリポジトリが何をしているか / hotwireをどう実装して何を実現しているか整理
------

## 1. cldmstr/gogohotwire

### 目的
- Go バックエンドで Hotwire (Turbo + Stimulus) を統合するテンプレートを示すデモアプリ
- 題材はドライバー / レースを扱う簡易的なレース管理システムだが、主眼は「Go で Hotwire を動かすときにハマりがちなポイントの実装例を示すこと」

### 使用スタック
- Web フレームワーク: **Echo**
- テンプレート: 標準 `html/template` に自作の `yield` 相当ヘルパを足してレイアウト合成を可能にしている
- フロント: Stimulus コントローラ (ドライバーのサイドバートグル、レース更新用)
- タスクランナー: `Taskfile.yaml`

### Hotwire の使い方
- **Turbo Frames**: eager-load / lazy-load の両パターン。フレーム外へのナビゲーションのターゲット指定も実演
- **Turbo Streams**: `prepend` / `update` アクションを使った差分更新
- **Turbo Stream over SSE**: レース実行中のリアルタイム更新を SSE でストリーム配信
- **Stimulus**: サイドバーのトグル、レース更新のトリガ等の UI 挙動

### Turbo Stream の配信経路
- フォーム送信のレスポンスとして HTML フラグメントを返す (通常経路)
- SSE でストリーミング (レースの実況更新)
- テンプレート応答時に MIME を切り替え

### 実装上のキモ
- Echo の HTML ハンドラに **カスタムレンダラ** を差し込み、`Content-Type: application/vnd.turbo-stream.html` を適切に返せるようにしている (README で明言されているポイント)
- `.stream.html` という拡張子で通常テンプレと Stream 用テンプレを区別
- `internal/template/template.go` に描画ロジック、`internal/template/tools.go` にヘルパ (yield 等)
- レースドメイン (`races/service_model.go`) の中で SSE サーバー側の実装が入っている

### ディレクトリ概観
```
├── assets/               # CSS, JS, Stimulus controllers
├── cmd/gogohotwire/      # エントリポイント
├── drivers/              # ドライバー ドメイン
├── internal/
│   ├── app/views/        # レイアウトテンプレ
│   └── template/         # カスタムレンダラ
├── races/                # レース ドメイン (SSE 実装含む)
└── Taskfile.yaml
```

### 学べる点 / 注意点
- 「Go 側で Hotwire を成立させるには、テンプレートのレイアウト合成 (yield) と Turbo Stream の Content-Type を自前で仕込む必要がある」ことが具体例として見える
- 逆に言うと **Hotwire 専用の Go パッケージは使っていない**。Echo + html/template + 手書きのカスタムレンダラという構成

---

## 2. wolfeidau/hotwire-golang-website

### 目的
- Go + Hotwire Turbo で「HTML Over The Wire」を成立させる実装パターンを示すデモサイト
- gogohotwire に比べるとやや包括的で、Turbo Drive / Frames / Streams / Stimulus と Hotwire の主要要素をひととおり触れている

### 使用スタック
- Web フレームワーク: **Echo**
- テンプレート: `html/template` + **`echo-go-templates`** (別ライブラリ) でレイアウト管理・埋め込みアセット対応
- アセットバンドル: **esbuild** + **`echo-esbuild-middleware`** による起動時自動配信 (JS 側は TypeScript)
- 依存管理: Node.js (フロントビルド用) + Go modules
- HTTPS 前提: `mkcert` で証明書を作り `localhost:9443` で立ち上げるスタイル

### Hotwire の使い方
- **Turbo Drive**: クライアントサイドナビゲーション
- **Turbo Frames**: スコープ付き DOM 置換
- **Turbo Streams**: **SSE と WebSocket の両方**を用意
- **Stimulus**: 挙動拡張

### Turbo Stream の配信経路
- **Server-Sent Events**: 片方向のサーバー → クライアント配信
- **WebSocket**: 双方向のリアルタイム通信
- 両方の実装が同居しているのが gogohotwire との明確な差分

### 実装上のキモ
- コアの実装は `internal/server/hotwire.go` に集約
- レイアウトは `views/layouts/base.html`。共通 CSS/JS は CDN 経由で引くスタイル
- テンプレート用に `echo-go-templates` を採用しており、gogohotwire のような yield の自作は不要。代わりにライブラリ側にレイアウトヘルパと Turbo Stream 応答向け Content-Type 対応が入っている
- esbuild 経由で JS を都度バンドルするので、Stimulus コントローラの開発体験が良い (hot-reload 対応)

### ディレクトリ概観
```
├── assets/                # 静的アセット
├── cmd/hotwire-server/    # エントリポイント
├── internal/              # コアロジック (hotwire.go を含む)
├── views/                 # HTML テンプレ (layouts/base.html 等)
├── Makefile
└── tools.go
```

### 学べる点 / 注意点
- gogohotwire が「素の html/template + 全部自作」なのに対し、こちらは **既存の Echo 向け補助ライブラリ (`echo-go-templates`, `echo-esbuild-middleware`) を組み合わせて Hotwire を回している**
- SSE と WebSocket の両方が実装されているので、Turbo Stream の配信経路を比較検討したいときの参考になる
- ただしこちらも **Hotwire 専用の Go パッケージを提供しているわけではない**。あくまで「既存の汎用ライブラリ + Hotwire の作法を Go 側で満たす」構成

---

## 3. 比較サマリ

- 両方とも **Echo + html/template** を土台にしている点は共通
- Turbo Stream の Content-Type (`application/vnd.turbo-stream.html`) を返す仕掛けは、どちらも「フレームワーク組み込みではなく、レンダラや補助ライブラリ側で手当てする」形になっている
- 配信経路:
  - gogohotwire は **SSE のみ**
  - wolfeidau は **SSE + WebSocket 両方**
- テンプレート合成:
  - gogohotwire は **yield 相当ヘルパを自作**
  - wolfeidau は **`echo-go-templates` に任せる**
- フロント資産:
  - gogohotwire は素の JS / Stimulus
  - wolfeidau は **esbuild で TypeScript バンドル**、mkcert で HTTPS 開発環境
- 結論として、両者とも「Go 側に Hotwire 専用パッケージがない」ことを裏付けており、Echo + html/template + 手書きの Content-Type ハンドリング + (必要なら) SSE / WebSocket 実装を組み合わせて成立させる、というのが現時点の Go における Hotwire の一般解と読める

## 4. 元の主張への含意 (「Go 製の Hotwire package はなさそう」への根拠づけ)

- 両リポジトリともに **Rails の `turbo-rails` 相当をパッケージ化していない**
- 具体的には次を自前で用意する必要がある:
  - Turbo Stream 用の Content-Type (`application/vnd.turbo-stream.html`) を返すレスポンスヘルパ
  - Turbo Stream アクション (`append`, `prepend`, `replace`, `update`, `remove` …) を発行するテンプレヘルパ
  - Turbo Stream over SSE / WebSocket のブローカ
  - レイアウト (`yield`) 相当のテンプレート合成
- これらを **一体で提供する Go パッケージが (少なくとも一般に知られている範囲では) 存在しない** → 元の「なさそう」判断と整合

---

## 5. 標準 `html/template` (+ `net/http`) で Hotwire を一通り実現するのに必要なもの

「フレームワーク (Echo) やライブラリ (`echo-go-templates` 等) を使わずに、標準 pkg だけでどこまでいけるか」の観点で、Hotwire の各要素ごとに整理する。

### 5-1. Turbo Drive (画面遷移の SPA 化)
- **サーバー側で必要なもの: 実質なし**
- クライアント側で `turbo.js` を読み込むだけで動く (CDN 経由 or 手動配置)
- サーバーは通常どおり完全な HTML を返せば OK。Turbo が `<a>` / `<form>` を横取りして fetch に置き換える
- 唯一の注意点: リダイレクトは **303 See Other** で返すのが Turbo 推奨 (POST 後のフォーム再送防止)

### 5-2. Turbo Frames (部分置換)
- **サーバー側で必要なもの: 実質なし**
- `<turbo-frame id="xxx">…</turbo-frame>` で囲った HTML フラグメントを返すだけ
- 標準 `html/template` の `{{template}}` / `{{block}}` で十分表現できる
- Lazy load 用に `<turbo-frame id="xxx" src="/foo">` を返して、`/foo` 側では対応する `<turbo-frame id="xxx">…</turbo-frame>` だけを返すハンドラを組む

### 5-3. Turbo Streams (フォームレスポンスとしての差分適用)
- **サーバー側で必要なもの (3 点)**:
  1. **Content-Type 判定**: リクエストヘッダの `Accept` に `text/vnd.turbo-stream.html` が含まれるかを見て、Turbo Stream 応答が可能かを判断
  2. **レスポンス Content-Type**: `text/vnd.turbo-stream.html; charset=utf-8` を返す
     ```go
     w.Header().Set("Content-Type", "text/vnd.turbo-stream.html; charset=utf-8")
     ```
  3. **`<turbo-stream>` タグを組み立てるテンプレ**: `action`, `target`(または `targets` で CSS セレクタ) を属性で持ち、中に `<template>…</template>` で差し込む HTML を入れる形
     ```html
     {{define "stream"}}
     <turbo-stream action="{{.Action}}" target="{{.Target}}">
       <template>{{.HTML}}</template>
     </turbo-stream>
     {{end}}
     ```
- 1 レスポンスに複数の `<turbo-stream>` を並べれば、複数箇所を同時に更新できる

### 5-4. Turbo Stream over SSE (サーバープッシュ)
- **サーバー側で必要なもの (4 点)**:
  1. SSE 用ハンドラでヘッダを設定:
     ```go
     w.Header().Set("Content-Type", "text/event-stream")
     w.Header().Set("Cache-Control", "no-cache")
     w.Header().Set("Connection", "keep-alive")
     ```
  2. `http.Flusher` を使って書き込みごとにフラッシュ
  3. データフォーマットは `data: <turbo-stream …>…</turbo-stream>\n\n` (改行 2 個で 1 メッセージ)。複数行になる場合は各行に `data:` を付ける
  4. **Pub/Sub / ファンアウト機構**: どのクライアントに何を送るかを管理する仕組み。標準 pkg だけで書くなら
     - `map[chan []byte]struct{}` + `sync.RWMutex` でサブスクライバ管理
     - `context.Context` の `Done()` で切断検知
     - ドメインロジック側からブロードキャストする関数を用意
- クライアント側: `<turbo-stream-source src="/events">` を HTML に置くだけで Turbo が SSE を張ってくれる

### 5-5. Turbo Stream over WebSocket (双方向)
- **サーバー側で必要なもの**:
  - **標準 pkg だけでは WebSocket サーバーを組めない** (`net/http` は WebSocket ネゴシエーションを提供しない)
  - 実用上は `gorilla/websocket` か `nhooyr.io/websocket` (最近は `coder/websocket`) を入れる必要がある
  - **ここが「標準 pkg 縛り」の唯一の破綻ポイント**
- 上記を入れれば、SSE と同様のブロードキャスタを WebSocket 版で作れば OK
- クライアント側: Turbo 単体では WebSocket 対応がなく、`@hotwired/turbo-rails` の `<turbo-cable-stream-source>` 相当のアダプタが必要 (自前で 30 行くらい書ける)

### 5-6. Stimulus (フロント挙動)
- **サーバー側で必要なもの: なし**
- クライアント側で `stimulus.js` と自作のコントローラ JS を配信するだけ
- `net/http` の `http.FileServer` / `http.ServeFile` で静的配信できる

### 5-7. レイアウト合成 (Rails の `yield` 相当)
- **標準 `html/template` の `{{block}}` で十分表現できる**
  ```html
  {{/* layouts/base.html */}}
  <html><body>{{block "content" .}}{{end}}</body></html>
  ```
  ```html
  {{/* pages/home.html */}}
  {{define "content"}}<h1>Home</h1>{{end}}
  ```
  ```go
  t := template.Must(template.ParseFiles("layouts/base.html", "pages/home.html"))
  t.ExecuteTemplate(w, "layouts/base.html", data)
  ```
- gogohotwire の「yield を自作した」というのは、複数の子テンプレを動的に差し込みたいケース (partial の入れ子) を扱いやすくするためのもので、単純な 1 段レイアウトなら `{{block}}` で足りる
- パーシャル系は `ParseGlob("partials/*.html")` で名前空間を張って `{{template "partials/foo" .}}` で呼び出す形が定番

### 5-8. まとめ: 標準 pkg 縛りで足りない / 自作が必要なもの

- **足りない (外部依存が要る)**
  - WebSocket サーバー実装 (`gorilla/websocket` 等)
- **標準 pkg で書けるが「書く必要がある」もの**
  - Turbo Stream 用の Content-Type ミドルウェア / レスポンスラッパ
  - `<turbo-stream>` を組むテンプレヘルパ (Action ごとに関数化しておくと使いやすい)
  - SSE ハンドラ + ブロードキャスタ (`map[chan][]byte` + `sync.RWMutex` + `http.Flusher`)
  - フォーム後のリダイレクトを 303 で返すヘルパ
  - CSRF (Turbo は `<meta name="csrf-token">` を読んで自動で `X-CSRF-Token` を付けてくれる。サーバー側は検証をどこかに用意)
- **標準 pkg だけでいける (自作不要)**
  - Turbo Drive
  - Turbo Frames
  - Stimulus 配信
  - レイアウト合成 (`{{block}}`)

### 5-9. 最小構成のイメージ

- `net/http` の `ServeMux` でルーティング
- `html/template` で `layouts/base.html` + ページテンプレ + `_stream.html` (Turbo Stream 用の小テンプレ) を管理
- `internal/hotwire` パッケージに:
  - `func StreamHeader(w http.ResponseWriter)` — Content-Type セット
  - `type Stream struct{ Action, Target, TemplateName string; Data any }`
  - `func RenderStreams(w http.ResponseWriter, t *template.Template, streams []Stream)` — 複数 stream をまとめて書き出し
  - `type Broker struct{...}` — SSE 用のファンアウト
- `assets/` に `turbo.js` / `stimulus.js` / 自作 JS を置いて `http.FileServer` で配信
- WebSocket が必要になったら `gorilla/websocket` を足す (それまでは SSE で十分なことが多い)
