# hotwire-go package が提供すべきスコープ

## 前提

- 対象: Go + `html/template` でサーバーサイド Hotwire アプリを書くための package
- 基本方針: **Hotwire の JS 側 (turbo.js / stimulus.js) には何も足さない**。ブラウザには本家 turbo/stimulus をそのまま読み込む
- package が担うのは **サーバー側で守るべき規約と、そこから生まれる定型を Go の書き味で吸収するレイヤ**

## リクエスト種別と分岐 (package 設計の背骨)

以下 3 種類を 1 箇所で分岐できるようにするのが package の中核。

```
リクエスト
  ├─ Accept: text/vnd.turbo-stream.html   → Turbo Streams レスポンス
  ├─ Turbo-Frame: <frame_id>              → Frame 単体レスポンス
  └─ どちらもなし                          → 通常のフル HTML (Drive 経由の遷移含む)
```

- 判定 helper: `hotwire.WantsStream(r) bool` / `hotwire.FrameID(r) (string, bool)`
- 分岐した後のレスポンス書き込み helper (後述) をセットで提供する

---

## Turbo Drive

サーバー側でやることがほぼ無い領域。package でカバーする面は最小限。

**カバーする**

- **POST 後のリダイレクト規約 helper**: `hotwire.RedirectAfterSubmit(w, r, url)` で 303 See Other を返す
  - Turbo Drive はフォーム送信の 200 応答を「フォームエラー再表示」と解釈する。成功時は 303 (or 302) が期待値
- **`Turbo-Location` レスポンスヘッダの helper**: XHR 内部リダイレクトで History API の URL を差し替えるためのヘッダ書き込み
- **preview キャッシュ制御メタタグの helper**: `{{ turboCacheControl "no-preview" }}` 等
- **turbo.js を読み込む `<script>` タグの helper** (お好みで): `{{ turboScripts }}`

**カバーしない**

- リンク/form の intercept、`<body>` 差し替え、`<head>` マージ、History API 操作 → 全部 turbo.js の仕事
- `data-turbo="false"` などの opt-out 属性 → HTML に書くだけなので helper 不要

---

## Turbo Frames

サーバー側の関与が増える領域。package の価値が最も出る場所の一つ。

**カバーする**

- **`<turbo-frame>` を吐くテンプレート helper**
  - `id` / `src` / `loading="lazy"` / `target="_top"` などの属性を安全に組み立てる
  - `{{ turboFrame "message_1" }}...{{ endTurboFrame }}` 相当の named template 群を用意
- **`Turbo-Frame` リクエストヘッダの判定 helper**
  - `hotwire.FrameID(r) (id string, ok bool)` でハンドラ内で分岐可能に
- **フレーム部分のみのレンダリング helper** (帯域削減 & レンダコスト削減)
  - `hotwire.RenderFrame(w, tmpl, "message_1", data)` で該当フレームだけを返す
  - 実装方針: 各 frame を独立した named template として持たせ、ヘッダを見て該当 template だけ実行する (フル HTML を DOM パースして抜き出すより素直)
- **lazy-load 先の frame 専用ハンドラを書くための骨組み**
  - `src` 先で「フル HTML の一部としても、frame 単体としても」レンダリングできる pattern を提供

**カバーしない**

- レスポンス HTML からの frame 抽出・DOM 差し込み → turbo.js の仕事
- `data-turbo-frame` 属性 → HTML に書くだけ

---

## Turbo Streams

サーバー主導の DOM 操作をサーバー言語 (Go) 側の型で書けるようにする領域。package の中で最もコード量が多くなる想定。

**カバーする**

- **`<turbo-stream>` を吐くテンプレート helper (9 action 分)**
  - `append` / `prepend` / `before` / `after` / `replace` / `update` / `remove` / `morph` / `refresh`
  - `target` (単一 DOM ID) と `targets` (CSS セレクタ) の両方に対応
  - `<template>` タグの内包を必須にする書き味 (書き忘れ事故を防ぐ)
  - 例: `hotwire.Append("messages", tmpl, data)` が `<turbo-stream action="append" target="messages"><template>...</template></turbo-stream>` を組み立てる
- **Content-Type の自動セット**
  - Stream レスポンス書き込み helper が `Content-Type: text/vnd.turbo-stream.html` を必ず付ける
- **複数 stream を 1 レスポンスに詰めるための writer**
  - `hotwire.StreamResponse{...}.Write(w)` で `<turbo-stream>` 断片を任意個並べて返せる
- **broadcast (リアルタイム更新) を SSE で流す helper**
  - `net/http` だけで書ける SSE 実装を優先
  - `hotwire.SSEHub` 的な pub/sub 骨組みと、`<turbo-stream-source src="/events">` に応答する handler
  - WebSocket 版は必要になったら別 sub-package で追加 (最初は SSE のみで良い)
- **`<turbo-stream-source>` を吐くテンプレート helper**

**カバーしない**

- ブラウザ側での DOM 挿入・置換・削除の実行 → turbo.js の仕事
- SSE/WebSocket 接続の維持 → ブラウザの EventSource / WebSocket API がやる

---

## Stimulus

サーバー側でやることは実質「HTML に `data-*` を書くだけ」なので、package で覆う面は非常に薄い。**最初はカバーしなくて良い**。

**カバーするかもしれない (v2 以降で検討)**

- `data-controller` / `data-action` / `data-*-target` / `data-*-value-*` を書くための template helper
  - HTML に直接書いても済むので、helper があってもなくても大差はない
  - 複数 controller の連結 (`data-controller="a b c"`) など、素で書くと面倒な部分だけ helper 化する余地はある
- Stimulus JS を読み込む `<script>` タグの helper

**カバーしない**

- Controller / Action / Target / Value の JS 側実装 → stimulus.js の仕事
- 状態管理 → DOM の data 属性に任せる

---

## 横断 / 共通

**カバーする**

- **レスポンス分岐ミドルウェア / dispatcher**
  - Rails の `respond_to do |format| ... end` に相当するもの
  - `hotwire.Respond(w, r, hotwire.Formats{ HTML: ..., Frame: ..., Stream: ... })` のような API 案
- **CSRF トークンの meta タグ helper** (使う人向け)
  - Turbo は `<meta name="csrf-token">` を見て自動でヘッダに載せる
- **turbo.js / stimulus.js の import map (or `<script>`) helper**

**カバーしない**

- 認証・セッション・DB 等の Web フレームワーク一般の関心事
- ルーティング → `net/http` / `chi` / `echo` 等の既存 mux に任せる

---

## API 設計方針 (メモ)

- ハンドラは `http.HandlerFunc` のままで書ける粒度に留める。フレームワークにしない
- `html/template` の上に薄く乗る。独自テンプレートエンジンは作らない
- テンプレート helper は `template.FuncMap` として提供 + named template でも書けるように両立
- Stream の action helper は「文字列連結ではなく、型で `<turbo-stream>` を組み立てる」方向にする (XSS 事故を減らす)
- SSE hub は net/http だけで完結させ、外部依存を持ち込まない

## 優先順位 (実装順の目安)

1. リクエスト判定 helper (`FrameID`, `WantsStream`) と分岐 dispatcher
2. Turbo Frames の template helper と単体レンダ helper
3. Turbo Streams の 9 action helper と `text/vnd.turbo-stream.html` writer
4. Turbo Drive の 303 redirect / cache-control 系 helper
5. SSE hub による broadcast
6. Stimulus 用の薄い helper (需要が出たら)
