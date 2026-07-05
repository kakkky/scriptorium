# Hotwire vs HTMX と hotwire-go の意義

## TL;DR

- Go には既に HTMX の実装補助 package (`donseba/go-htmx`, 586★) が存在する
- それでも Hotwire を Go でやる意義はある。理由は **哲学の差** と **サーバー側規約の厚さの差**
- 「小さい部分更新なら HTMX で十分」だが、「大規模かつ長期運用に耐える規約を欲しい層」には Hotwire が刺さる。この層向けの Go package は現状ゼロ

---

## 前提: HTMX 側の Go エコシステム

`donseba/go-htmx` (<https://github.com/donseba/go-htmx>) が既に存在する。

- 提供機能
  - リクエスト判定: `IsHxRequest`, `IsHxBoosted`, `IsHxHistoryRestoreRequest`
  - レスポンス helper: `HX-Swap`, `HX-Trigger`, `HX-Push-Url`, `HX-Re-Target` などのヘッダ書き込み
  - `NewSwap()` / `NewTrigger()` で挙動を型付きで組み立て
  - `TriggerSuccess` / `Error` / `Info` / `Warning` などの通知イベント発火 helper
  - SSE マネージャ (HTMX SSE extension 対応)
  - `io.Writer` 実装で既存の Go 処理に埋め込みやすい
- 586★、v1.13.1 (2026/04)、活発なメンテナンス
- **サーバー側規約の吸収は「HTMX が薄い分、それに合わせて薄い」**。ヘッダ操作の抽象化が主で、テンプレートエンジンやレンダリング責務は持たない

つまり **Go × HTMX には既に成熟した選択肢がある**。ここが hotwire-go を「作る意義がない」と言われかねない出発点。

---

## Hotwire の技術的優位性 (Hotwire vs HTMX 議論の整理)

Hotwire 公式フォーラム "Hotwire vs HTMX Comparison" (<https://discuss.hotwired.dev/t/hotwire-vs-htmx-comparison/1614>) と、これまでの調査 (`what-is.md`, `investigation.md`) を踏まえて整理する。

### 1. プログレッシブエンハンスメントを設計思想の中核に据えている

- Hotwire 開発チーム (sam) 本人の発言: 「HTML と CSS でできることは活用し、必要なとき JS 層を追加する」がコア哲学
- **JS が無効でも一次的な動作を保証**する設計。SEO / アクセシビリティ / 初回描画に有利
- HTMX は `hx-boost` で似たことはできるが、思想の中心には据えていない
- 意味合い: 「ノーJS で動くことを前提にした後、Turbo が上乗せする」設計 = **サーバー側で完全な HTML を返せば済む** ことが保証される。Go の SSR 中心の書き味と極めて相性が良い

### 2. 3 層の名前付き抽象 (Drive / Frames / Streams)

- HTMX: `hx-get` / `hx-target` / `hx-swap` という**汎用プリミティブ**の集合
- Hotwire: 「ページ遷移」「フレーム」「差分ストリーム」という**名前の付いた抽象**を最初から分けている
- 意味合い: **抽象の粒度が粗い分、書き味が揃う**。特に「サーバー主導の DOM 差分を broadcast する」ユースケースで、Turbo Streams の 9 action は明示的な語彙になっており、HTMX の `hx-swap-oob` の自作構成に比べて設計コストが低い

### 3. JS を書き足す時の連続性 (Stimulus)

- フォーラム投稿 #5 (chromonav, 2021/02/03) の観察と一致
  - 「Hotwire はトランジションとキャッシュ管理が堅牢で期待通り動く」
  - 「HTMX はキャッシング未装備」「Web Components 互換に難あり (パースエラーが多い)」
  - chromonav 本人は最終的に Hotwire + Stimulus を選択
- 差の実像:
  - HTMX で属性表現の限界を超えると、hyperscript か素の JS (Alpine.js を足す構成が多い) に逃げる
  - Hotwire は Stimulus が「JS を書く場所と作法」を用意している = Controller / Action / Target / Value / Outlet で拡張パスが決まっている
- **「シンプルなうちは HTMX、複雑化すると Hotwire の方が書き味の天井が高い」** が実務観察

### 4. ネイティブアプリへの延伸 (Turbo Native / Strada)

- HTMX には無い軸
- iOS / Android の WebView をベースにネイティブと橋渡し
- **サーバー側は同じ Turbo Streams / Frames を返せばよい** = hotwire-go を作ればモバイル対応も自動的に付いてくる
- Go サーバー × Web + Native アプリを同一の HTML over the wire 契約で扱える構図が組める

### 5. 「JS API の柔軟性と成熟度」 (chromonav の実案件選択理由)

- **JS API の柔軟性**
  - Turbo: `Turbo.visit()`, `Turbo.cache.clear()`, `Turbo.session.drive` などのプログラマブル API
  - ライフサイクルイベントが細かい: `turbo:before-visit` / `turbo:before-fetch-request` / `turbo:before-fetch-response` / `turbo:before-render` / `turbo:frame-load` / `turbo:submit-start` など、フェーズごとにフック可能
  - Stimulus: Controller / Value change callback / Target connect-disconnect / Outlet で **UI 挙動が増えても書き味が揃う**
- **成熟度**
  - Turbolinks (2012) からの 10 年以上の運用実績
  - Basecamp / HEY / Shopify などプロダクション使用実績
  - 大規模運用でのエッジケース (History API / キャッシュ / フォーム再送 / CSRF / Frame と Stream の interplay) が枯れている

---

## Hotwire の弱点 (フェアに)

議論を偏らせないために chromonav 投稿 #5 で挙がった Hotwire 側の課題も残す:

- JS API が全て public に公開されているわけではなく、カスタムイベントトリガーには Stimulus 経由が事実上必須
- ファイルサイズは HTMX より大きい
- 公式ドキュメントは「ガイド」中心で、複雑な実装ではソースコード読解が必要になる場面がある

これらは **Go サーバー側の package を作る話とは直交** している (JS 側の課題)。hotwire-go の存在意義を毀損しない。

---

## hotwire-go を作る意義

上記優位性のうち、**サーバー側 package として価値を生む部分**を抽出する。

### 1. サーバー側規約が厚い = package で吸収する価値が大きい

- HTMX: 素の HTML 断片を返せば良い + ヘッダ操作。**規約が薄いので package の面積も薄い** (`go-htmx` が実際そうなっている)
- Hotwire: 以下を守らないと動かない
  - `text/vnd.turbo-stream.html` の Content-Type
  - `<turbo-stream action="...">` の 9 action 語彙と `<template>` 内包規則
  - `<turbo-frame id="...">` の ID 一致抽出規則
  - `Turbo-Frame` ヘッダ / `Accept` ヘッダによるレスポンス切替
  - フォーム後の 303 See Other 規約
  - SSE / WebSocket 経由の `<turbo-stream>` broadcast プロトコル
- **これらを毎回自作するのは辛い** (`investigation.md:124-131` で確認済み)
- Go 用の統合 package は現時点で存在しない (`investigation.md:1-11`)

### 2. Go の SSR 志向と Hotwire 哲学が噛み合う

- Hotwire のプログレッシブエンハンスメント哲学 = 「サーバーは完全な HTML を返す」を前提にする
- Go は `html/template` + `net/http` で完全な SSR を素直に書ける言語
- 相性は思想レベルで良い。`motivation.md:1-11` の「Go の SSR で完結させたい」動機と直結

### 3. モバイル対応が「サーバー変更ゼロ」で付いてくる

- Turbo Native / Strada は `turbo-ios` / `turbo-android` が既にあり、**サーバー側は Web と同じ Turbo Streams / Frames を返すだけ**
- hotwire-go を作れば Go サーバー × iOS / Android アプリ の構成が組める。HTMX ではこの軸を取れない

### 4. `donseba/go-htmx` との棲み分けが明確

- go-htmx: HTMX 前提。ヘッダ helper + SSE + トリガー通知
- hotwire-go (仮): Turbo 前提。**Frames / Streams テンプレート helper + Content-Type ルーティング + SSE broker + 9 action 型組み立て**
- 対象が全く違うので競合しない。「Go で SSR + HTML over the wire を書きたい人」に対する **選択肢を増やす** ポジション

### 5. scope.md の内容が既に MVP に十分近い

- `scope.md` で優先順位まで整理済み (`scope.md:140-147`)
  1. リクエスト判定 helper (`FrameID`, `WantsStream`) と分岐 dispatcher
  2. Turbo Frames の template helper と単体レンダ helper
  3. Turbo Streams の 9 action helper と `text/vnd.turbo-stream.html` writer
  4. Turbo Drive の 303 redirect / cache-control 系 helper
  5. SSE hub による broadcast
  6. Stimulus 用の薄い helper (需要が出たら)
- 実装可能性は高い。net/http + html/template だけで 1〜4 は書け、5 も外部依存不要

---

## 結論

- HTMX / go-htmx は「シンプルな部分更新を素早く書きたい層」に対する成熟解として既に存在する
- Hotwire は **哲学 (プログレッシブエンハンスメント) と抽象の粒度 (Drive/Frames/Streams) とスケール可能性 (Stimulus) と Native 対応** で異なる価値提案を持つ
- そのうち **サーバー側規約の厚さ** は Go package として吸収する価値が最も大きい部分であり、現状 Go にその選択肢は無い
- hotwire-go は「HTMX の代替」ではなく、**「Rails の turbo-rails 相当を Go で最小依存に提供する」** ポジションを取れば、明確なニーズと空白に対して意義が成立する

---

## 参考

- Hotwire 公式フォーラム "Hotwire vs HTMX Comparison": <https://discuss.hotwired.dev/t/hotwire-vs-htmx-comparison/1614>
- 同スレッド chromonav 投稿 #5: <https://discuss.hotwired.dev/t/hotwire-vs-htmx-comparison/1614/5>
- HTMX 公式: <https://htmx.org/>
- donseba/go-htmx: <https://github.com/donseba/go-htmx>
- Hotwire 公式: <https://hotwired.dev/>
- Turbo Handbook: <https://turbo.hotwired.dev/handbook/introduction>
- Stimulus Handbook: <https://stimulus.hotwired.dev/handbook/introduction>
- 関連ノート: `motivation.md` / `what-is.md` / `investigation.md` / `scope.md`
