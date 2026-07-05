> 原文: https://turbo.hotwired.dev/reference/streams
> タイトル: Streams
> 最終取得: 2026-07-05

# Streams

## The eight actions

### Append

template タグの内容を、target dom id で指定されたコンテナに append する。

```html
<turbo-stream action="append" target="dom_id">
  <template>
    Content to append to container designated with the dom_id.
  </template>
</turbo-stream>
```

template の最初の要素が、dom_id で指定されたコンテナの直下の子要素と同じ id を既に使っている場合、append ではなく置換される。

### Prepend

template タグの内容を、target dom id で指定されたコンテナに prepend する。

```html
<turbo-stream action="prepend" target="dom_id">
  <template>
    Content to prepend to container designated with the dom_id.
  </template>
</turbo-stream>
```

template の最初の要素が、dom_id で指定されたコンテナの直下の子要素と同じ id を既に使っている場合、prepend ではなく置換される。

### Replace

target dom id で指定された要素を置換する。

```html
<turbo-stream action="replace" target="dom_id">
  <template>
    Content to replace the element designated with the dom_id.
  </template>
</turbo-stream>
```

`turbo-stream` 要素に `[method="morph"]` 属性を追加すると、target dom id で指定された要素を morph によって置換できる。

```html
<turbo-stream action="replace" method="morph" target="dom_id">
  <template>
    Content to replace the element.
  </template>
</turbo-stream>
```

### Update

template タグの内容を、target dom id で指定されたコンテナに反映して更新する。

```html
<turbo-stream action="update" target="dom_id">
  <template>
    Content to update to container designated with the dom_id.
  </template>
</turbo-stream>
```

`turbo-stream` 要素に `[method="morph"]` 属性を追加すると、target dom id で指定された要素の子要素のみを morph できる。

```html
<turbo-stream action="update" method="morph" target="dom_id">
  <template>
    Content to replace the element.
  </template>
</turbo-stream>
```

### Remove

target dom id で指定された要素を削除する。

```html
<turbo-stream action="remove" target="dom_id">
</turbo-stream>
```

### Before

template タグの内容を、target dom id で指定された要素の前に挿入する。

```html
<turbo-stream action="before" target="dom_id">
  <template>
    Content to place before the element designated with the dom_id.
  </template>
</turbo-stream>
```

### After

template タグの内容を、target dom id で指定された要素の後に挿入する。

```html
<turbo-stream action="after" target="dom_id">
  <template>
    Content to place after the element designated with the dom_id.
  </template>
</turbo-stream>
```

### Refresh

新しいコンテンツを描画するために [Page Refresh](../handbook/03_page_refreshes.md) を発火させる。

```html
<!-- without `[request-id]` -->
<turbo-stream action="refresh"></turbo-stream>

<!-- debounced with `[request-id]` -->
<turbo-stream action="refresh" request-id="abcd-1234"></turbo-stream>

<!-- refresh behavior with `[method]` and `[scroll]` -->
<turbo-stream action="refresh" method="morph" scroll="preserve"></turbo-stream>
```

## Targeting Multiple Elements

1 つのアクションで複数の要素を対象にするには、`target` 属性の代わりに、CSS クエリセレクタを指定した `targets` 属性を使う。

```html
<turbo-stream action="remove" targets=".elementsWithClass">
</turbo-stream>

<turbo-stream action="after" targets=".elementsWithClass">
  <template>
    Content to place after the elements designated with the css query.
  </template>
</turbo-stream>
```

## Processing Stream Elements

Turbo はあらゆる形式のストリームに接続して stream アクションを受信・処理できる。stream ソースは、そのイベントの `data` 属性に stream アクションの HTML を含んだ [MessageEvent](https://developer.mozilla.org/en-US/docs/Web/API/MessageEvent) メッセージをディスパッチする必要がある。接続は `Turbo.session.connectStreamSource(source)` で行い、切断は `Turbo.session.disconnectStreamSource(source)` で行う。`MessageEvent` を生成するもの以外のソースから stream アクションを処理する必要があるなら、`Turbo.renderStreamMessage(streamActionHTML)` を使うとよい。

これらを一つにまとめる良い方法は、turbo-rails が [TurboCableStreamSourceElement](https://github.com/hotwired/turbo-rails/blob/main/app/javascript/turbo/cable_stream_source_element.js) で行っているように、カスタム要素を使うことである。

## Stream Elements inside HTML

Turbo streams は[カスタム HTML 要素](https://developer.mozilla.org/en-US/docs/Web/API/Web_components/Using_custom_elements)として実装されている。要素は、ページの DOM に接続された際にブラウザが呼び出す `connectedCallback` 関数の中で解釈される。

つまり、DOM にレンダリングされた stream 要素はすべて解釈される、ということである。解釈が終わったあと、Turbo はその要素を DOM から取り除く。より具体的に言えば、stream アクションをページや frame のコンテンツ HTML の中にレンダリングすれば、それらは実行される。これはメインのコンテンツ読み込みに加えて追加の副作用を実行するのに使える。
