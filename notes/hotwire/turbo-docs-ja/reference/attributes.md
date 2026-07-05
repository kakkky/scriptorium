> 原文: https://turbo.hotwired.dev/reference/attributes
> タイトル: Attributes and Meta Tags
> 最終取得: 2026-07-05

# Attributes and Meta Tags

## Data Attributes

要素の挙動を Turbo 用にカスタマイズするために、以下の data 属性を適用できる。

- `data-turbo="false"` は、リンクとフォーム、およびそれらの子孫に対して Turbo Drive を無効化する。祖先で無効化されているものを再度有効化するには `data-turbo="true"` を使う。ただし注意が必要である。Turbo Drive が無効化されている場合、ブラウザはリンククリックを通常どおりに扱うが、[native adapters](../handbook/06_native.md) はアプリを終了する可能性がある。なお、Turbo が[パスを無視している](../handbook/02_drive.md#ignored-paths)場合、`data-turbo="true"` を指定してもそれを強制的に有効化することはできない。

- `data-turbo-track="reload"` は、要素の HTML を追跡し、変更された際にフルページリロードを行う。典型的には [`script` および CSS `link` 要素を最新に保つ](../handbook/02_drive.md#reloading-when-assets-change)ために使われる。

- `data-turbo-track="dynamic"` は、要素の HTML を追跡し、HTML レスポンスに含まれない場合にその要素を削除する。典型的にはナビゲーション中に [`style` および `link` 要素を削除する](../handbook/02_drive.md#removing-assets-when-they-change)ために使われる。

- `data-turbo-frame` は、ナビゲーション対象となる Turbo Frame を指定する。詳細は [Frames のドキュメント](frames.md)を参照のこと。

- `data-turbo-preload` は、[Drive](../handbook/02_drive.md#preload-links-into-the-cache) に対して次ページのコンテンツを事前取得するよう指示する。

- `data-turbo-prefetch="false"` は、要素がホバーされた際の[リンクの prefetch](../handbook/02_drive.md#prefetching-links-on-hover) を無効化する。

- `data-turbo-action` は、[Visit](../handbook/02_drive.md#page-navigation-basics) のアクションをカスタマイズする。有効な値は `replace` または `advance` である。Turbo Frames と組み合わせて、[frame のナビゲーションをページ visit に昇格させる](../handbook/04_frames.md#promoting-a-frame-navigation-to-a-page-visit)ためにも使える。

- `data-turbo-permanent` は、[要素をページ読み込みをまたいで保持する](../handbook/07_building.md#persisting-elements-across-page-loads)。要素には一意な `id` 属性が必要である。また、[morphing を伴うページリフレッシュ](../handbook/03_page_refreshes.md)を使う際に、morph 対象から要素を除外する用途にも用いられる。

- `data-turbo-temporary` は、ドキュメントがキャッシュされる前に要素を削除し、リストア時に再表示されるのを防ぐ。

- `data-turbo-eval="false"` は、インラインの `script` 要素が Visit のたびに再評価されるのを防ぐ。

- `data-turbo-method` は、リンクのリクエスト種別をデフォルトの `GET` から変更する。理想的には非 `GET` リクエストはフォームで発火させるべきだが、フォームを使えない場面では `data-turbo-method` が有用となりうる。

- `data-turbo-stream` は、リンクまたはフォームが Turbo Streams レスポンスを受け付けられることを指定する。Turbo は非 `GET` メソッドのフォーム送信については[自動的に stream レスポンスを要求する](../handbook/05_streams.md#streaming-from-http-responses)が、`data-turbo-stream` を使えば `GET` リクエストでも Turbo Streams を利用できるようになる。

- `data-turbo-confirm` は、指定した値で確認ダイアログを表示する。`form` 要素や、`data-turbo-method` を持つリンクで使える。

- `data-turbo-submits-with` は、フォーム送信中に表示するテキストを指定する。`input` または `button` 要素で使える。フォーム送信中は要素のテキストが `data-turbo-submits-with` の値になり、送信後は元のテキストに戻る。処理中に「Saving…」のようなメッセージを表示してユーザーにフィードバックを与えるのに便利である。

- `download` は、リンクがナビゲーションではなくファイルのダウンロード用であるため、Turbo の対象から除外する。

## Automatically Added Attributes

以下の属性は Turbo によって自動的に追加され、任意の時点における Visit の状態を判別するのに役立つ。

- `disabled` は、フォーム送信中のリピート送信を防ぐため、送信元 (submitter) に追加される。

- `data-turbo-preview` は、Visit 中に[プレビュー](../handbook/07_building.md#detecting-when-a-preview-is-visible)を表示している間、`html` 要素に追加される。

- `data-turbo-visit-direction` は、visit 中に `html` 要素に追加され、その方向を示す `forward`、`back`、または `none` の値を持つ。

- `aria-busy` は、次のタイミングで追加される:

  - visit が進行中のあいだ、`html` 要素に

  - 送信が進行中のあいだ、`form` 要素に

  - frame 内で visit またはフォーム送信が進行中のあいだ、`turbo-frame` 要素に

- `busy` は、frame 内でナビゲーションまたはフォーム送信が進行中のあいだ、`turbo-frame` 要素に追加される。

## Meta Tags

`head` に追加する以下の `meta` 要素で、キャッシュや Visit の挙動をカスタマイズできる。

- `<meta name="turbo-cache-control">` は、[キャッシュを無効化する](../handbook/07_building.md#opting-out-of-caching)ために使う。

- `<meta name="turbo-visit-control" content="reload">` は、Turbo がそのページにナビゲートするたびにフルページリロードを行う。リクエストが `<turbo-frame>` に由来する場合も同様である。

- `<meta name="turbo-root">` は、[Turbo Drive を特定のルートロケーションにスコープする](../handbook/02_drive.md#setting-a-root-location)ために使う。

- `<meta name="view-transition" content="same-origin">` は、[View Transition API](https://caniuse.com/view-transitions) をサポートするブラウザで view transition を発火させる。

- `<meta name="turbo-refresh-method" content="morph">` は、[morphing を伴うページリフレッシュ](../handbook/03_page_refreshes.md)を設定する。

- `<meta name="turbo-refresh-scroll" content="preserve">` は、[ページリフレッシュ時のスクロール保持](../handbook/03_page_refreshes.md)を有効化する。

- `<meta name="turbo-prefetch" content="false">` は、[ホバー時のリンク prefetch](../handbook/02_drive.md#prefetching-links-on-hover) を無効化する。
