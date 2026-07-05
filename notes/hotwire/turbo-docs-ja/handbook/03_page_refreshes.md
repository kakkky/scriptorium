> 原文: https://turbo.hotwired.dev/handbook/page_refreshes
> タイトル: Smooth page refreshes with morphing
> 最終取得: 2026-07-05

# Smooth page refreshes with morphing

[Turbo Drive](02_drive.md) はページ全体のリロードを避けることでナビゲーションを高速化します。しかし Turbo がさらに体験の質を引き上げられるシナリオがあります。同じページをもう一度読み込む場合 (ページリフレッシュ) です。

ページリフレッシュの典型的なシナリオは、フォームを送信して元のページにリダイレクトされて戻ってくる場合です。こうしたシナリオでは、ページの `<body>` を丸ごと置き換えるのではなく、変更された内容だけを更新できると体感が大幅に向上します。Turbo は morphing とスクロール保持によって、これを代わりにやってくれます。

## Morphing

Turbo がページリフレッシュをどう扱うかは、ページの head 内の `<meta name="turbo-refresh-method">` で設定できます。

```html
<head>
  ...
  <meta name="turbo-refresh-method" content="morph">
</head>
```

指定できる値は `morph` または `replace` (デフォルト) です。`morph` にすると、ページリフレッシュが発生したときに、Turbo はページの `<body>` の中身を置き換えるのではなく、変更のあった DOM 要素だけを更新し、それ以外はそのまま残します。このアプローチは画面の状態を保つため、より良い体感をもたらします。

内部的には、Turbo は素晴らしい [idiomorph ライブラリ](https://github.com/bigskysoftware/idiomorph)を利用しています。

## Scroll preservation

Turbo がスクロールをどう扱うかは、ページの head 内の `<meta name="turbo-refresh-scroll">` で設定できます。

```html
<head>
  ...
  <meta name="turbo-refresh-scroll" content="preserve">
</head>
```

指定できる値は `preserve` または `reset` (デフォルト) です。`preserve` にすると、ページリフレッシュが発生したときに、Turbo はページの垂直方向と水平方向のスクロール位置を保持します。

## Exclude sections from morphing

morphing の際に特定の要素を無視したい場合があります。たとえば、ページがリフレッシュされても開いたままにしておきたいポップオーバーがあるかもしれません。そうした要素には `data-turbo-permanent` を付けておけば、Turbo は morph しようとしません。

```html
<div data-turbo-permanent>...</div>
```

## Turbo frames

[turbo frames](04_frames.md) を使うと、ページリフレッシュ時に morphing でリロードされる領域を画面内に定義できます。そうするには、対象の frame に `refresh="morph"` を付けます。

```html
<turbo-frame id="my-frame" refresh="morph" src="/my_frame">
</turbo-frame>
```

この仕組みを使えば、初回ページロードでは届かなかった追加コンテンツ (たとえばページネーション) をロードできます。ページリフレッシュが起こっても、Turbo は frame の中身を削除せず、代わりにその turbo frame をリロードし、中身を morphing でレンダリングします。

## Broadcasting page refreshes

新しい [turbo stream アクション](05_streams.md)として `refresh` が追加されており、これによってページリフレッシュを起こせます。

```html
<turbo-stream action="refresh"></turbo-stream>
```

リフレッシュの挙動は `method` と `scroll` 属性で指定できます。

```html
<turbo-stream action="refresh" method="morph" scroll="preserve"></turbo-stream>
```

`method` 属性は `morph` または `replace`、`scroll` 属性は `preserve` または `reset` を指定できます。

サーバーサイドフレームワークはこの stream を活用して、シンプルながら強力なブロードキャストモデルを提供できます。サーバーは単一の汎用シグナルをブロードキャストするだけで、ページは morphing によってなめらかにリフレッシュされます。

Rails 向けに [`turbo-rails`](https://github.com/hotwired/turbo-rails) gem がどう実現しているかを見てみてください。

```ruby
# In the model
class Calendar < ApplicationRecord
  broadcasts_refreshes
end

# View
turbo_stream_from @calendar
```

[次へ: Decompose with Turbo Frames](04_frames.md)
