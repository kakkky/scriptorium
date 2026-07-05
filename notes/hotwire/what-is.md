# Hotwire とは

## 一言で

- **HTML Over The Wire** = ページ更新を「JSON + クライアント側レンダリング」ではなく、**サーバーが返した HTML 片をブラウザに流し込む**ことで行うアプローチ
- 37signals (Basecamp) 発。Rails に標準組み込みだが、フロント側は素の JS/HTML 仕様の上に乗っているだけなので、サーバー言語は本来問わない

## 何を実現しているか

- SPA 的な UX (フルリロードなし・部分更新・リアルタイム反映) を、**サーバーサイドレンダリング中心**で実現する
- クライアント側にビジネスロジックと状態を極力持たせない (状態は DOM に持たせる)
- 書く JavaScript の量を大幅に減らす。書くとしても "sprinkles" 程度の小さな振る舞い

## 構成要素

- **Turbo**: ページ遷移・部分更新・リアルタイム更新を担う JS ライブラリ。中身は 3 つ (Drive / Frames / Streams)
- **Stimulus**: DOM に対して最小限の JS 振る舞いを宣言的に付けるためのフレームワーク
- **Hotwire Native** (旧 Strada): iOS/Android の WebView をベースにネイティブと橋渡しする仕組み。本ノートでは深追いしない

---

## 技術的にどう実現しているか

以下は「サーバーが何を返し、ブラウザで Turbo/Stimulus が何をやっているか」に絞って書く。

### Turbo Drive

- **実現するもの**: フルリロードなしのページ遷移 (SPA 的な体感)
- **仕組み**:
  - `<a>` クリックと `<form>` submit を JS で intercept する
  - `fetch` で **完全な HTML ページ** をサーバーから取得する
  - レスポンス HTML の `<body>` を丸ごと差し替え、`<head>` はマージする (script/link は再評価しない)
  - History API で URL を書き換える
  - `window` / `document` オブジェクトは維持されるので、JS のグローバル状態は保たれる
  - `data-turbo="false"` で個別に opt-out できる
- **サーバー側の責務**: 従来どおり **完全な HTML ページ** を返すだけ。専用の API 契約は不要

### Turbo Frames

- **実現するもの**: ページの一部だけを独立したナビゲーションコンテキストとして切り出し、そこだけ差し替える (SPA の部分ルーティング相当)
- **仕組み**:
  - HTML を `<turbo-frame id="foo">...</turbo-frame>` というカスタム要素で囲む
  - そのフレーム内のリンク/フォームからのリクエストは自動的にキャプチャされる
  - サーバーが返した HTML から **同じ id を持つ `<turbo-frame>` の中身だけを抽出**して、その場に差し込む
  - フレーム外の要素は無視されるので、「表示用ページ」と「編集用ページ」を同一 URL で返し分けても両方に流用できる
  - `src` 属性を書けば、フレームが DOM に載った瞬間にその URL から自動で中身を読みに行く (eager-load)
  - `loading="lazy"` を付ければ、フレームが画面に見えるまで読み込みを遅延する
  - フレーム内のリンクに `data-turbo-frame="_top"` を付ければ、フレーム内ではなくページ全体を遷移させられる
- **サーバー側の責務**: 通常のページ HTML に `<turbo-frame id="foo">...</turbo-frame>` が含まれていれば OK。抽出はクライアント側

```html
<!-- 一覧ページ -->
<turbo-frame id="message_1">
  <h1>My message title</h1>
  <a href="/messages/1/edit">Edit</a>
</turbo-frame>

<!-- /messages/1/edit のレスポンスにも同じ id のフレームが入っていれば、そこだけ差し替わる -->
<turbo-frame id="message_1">
  <form action="/messages/1" method="post">...</form>
</turbo-frame>
```

### Turbo Streams

- **実現するもの**: **サーバー主導の DOM 更新**。フォーム送信のレスポンスや、他クライアントへのリアルタイム反映など、「サーバー側の状態変化をブラウザの DOM 変化として押しつける」用途
- **仕組み**:
  - サーバーが `<turbo-stream>` 要素を含む HTML を返す
  - action ごとに、指定した DOM に対する DOM 操作を Turbo が実行する
  - 対応する **action は 9 種類** (公式ハンドブック)
    - `append` / `prepend`: 子要素の末尾/先頭に追加
    - `before` / `after`: 対象要素の前/後に挿入
    - `replace`: 要素そのものを置き換え
    - `update`: 要素の中身だけを置き換え
    - `remove`: 要素を削除
    - `morph`: DOM を morphing で置き換え (差分適用)
    - `refresh`: ページ全体をリフレッシュ
  - 対象は `target="dom-id"` (単一要素) または `targets="css-selector"` (複数要素) で指定
  - 挿入する HTML は必ず `<template>` の中に書く
  - フォーム送信のレスポンスとして返すときの Content-Type は **`text/vnd.turbo-stream.html`**
  - Turbo は form 送信時に `Accept` ヘッダにこの MIME を自動で入れるので、サーバーはそれを見て「Stream で返すか、通常の HTML で返すか」を切り替えられる
  - リアルタイム配信は `<turbo-stream-source src="...">` を `<body>` の中に置く
    - `ws://` / `wss://` なら WebSocket 接続
    - それ以外は EventSource (SSE) で接続
    - 開いたときにストリーム購読を開始し、閉じたら終了する
- **サーバー側の責務**:
  - フォーム応答: Content-Type を `text/vnd.turbo-stream.html` にして `<turbo-stream>` を返す
  - broadcast: WebSocket または SSE で同じ `<turbo-stream>` 断片を push する
  - つまり **「DOM に対する差分操作の粒度」をサーバー側で HTML として表現する**のが Turbo Streams の勘所

```html
<!-- append の例 -->
<turbo-stream action="append" target="messages">
  <template>
    <div id="message_5">Hello</div>
  </template>
</turbo-stream>
```

### Stimulus

- **実現するもの**: サーバー側で表現しづらいクライアント固有の振る舞い (トグル、コピー、フォームの補助バリデーション、モーダル制御など)
- **哲学**:
  - サーバーが返した HTML を「強化 (augment)」するだけの控えめなフレームワーク
  - `class` 属性が HTML と CSS を繋ぐブリッジであるのと同様に、`data-controller` 属性が HTML と JS を繋ぐブリッジ
  - **状態は原則として DOM (data 属性) に持たせる**。JS メモリ上に持たせない → サーバー側の再レンダリングと自然に整合する
- **仕組み** (4 つの概念):
  - **Controller**: `data-controller="hello"` を付けた要素に対応する JS クラス。DOM に現れた瞬間に自動で connect され、消えたら disconnect される (MutationObserver ベース)
  - **Action**: `data-action="click->hello#greet"` の形で、DOM イベントを Controller のメソッドに結線する
  - **Target**: `data-hello-target="name"` の形で、Controller から重要な子要素を参照できるようにする
  - **Value**: `data-hello-value-url="..."` の形で、Controller が読み書き・変更監視できる値を DOM に持たせる
- **サーバー側の責務**: HTML に必要な `data-*` 属性を付けて返すだけ。JS 側は Controller クラスを書くだけで、DOM との結線は完全に宣言的

```html
<div data-controller="clipboard" data-clipboard-value-url="/token">
  <input data-clipboard-target="source" type="text" value="abc123" readonly>
  <button data-action="click->clipboard#copy">Copy</button>
</div>
```

---

## サーバー ⇔ ブラウザで流れるもの (まとめ)

- **Turbo Drive**: 完全な HTML ページ (通常のレスポンス)
- **Turbo Frames**: 完全な HTML ページ。その中に該当 id の `<turbo-frame>` が入っていれば OK
- **Turbo Streams (form 応答)**: `text/vnd.turbo-stream.html` の `<turbo-stream>` 断片
- **Turbo Streams (broadcast)**: WebSocket / SSE で流れる `<turbo-stream>` 断片
- **Stimulus**: HTML に混ぜる `data-*` 属性

→ **サーバーが返すものは常に HTML**。JSON にシリアライズしてクライアント側で再度 HTML に組み直す往復をやらない。これが "HTML Over The Wire" の意味

## Go で作る観点でのポイント (短く)

- `html/template` で `<turbo-frame>` `<turbo-stream>` を出力できれば大枠は動く
- form 応答で Content-Type を `text/vnd.turbo-stream.html` に切り替える helper が要る
- broadcast は WebSocket でも SSE でも良いが、SSE のほうが `net/http` だけで書けて Go では素直
- リクエストヘッダ `Accept` に `text/vnd.turbo-stream.html` が含まれるか、`Turbo-Frame` ヘッダが付いているかで、「フル HTML / Stream / Frame」の返し分けを判定する

## 一次情報

- <https://hotwired.dev/>
- <https://turbo.hotwired.dev/handbook/introduction>
- <https://turbo.hotwired.dev/handbook/frames>
- <https://turbo.hotwired.dev/handbook/streams>
- <https://stimulus.hotwired.dev/handbook/introduction>
