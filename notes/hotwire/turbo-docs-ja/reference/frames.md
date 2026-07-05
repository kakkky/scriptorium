> 原文: https://turbo.hotwired.dev/reference/frames
> タイトル: Frames
> 最終取得: 2026-07-05

# Frames

## Basic frame

frame 内でのすべてのナビゲーションを閉じ込め、辿ったリンクや送信したフォームのレスポンスに、一致する frame タグが含まれていることを期待する。

```html
<turbo-frame id="messages">
  <a href="/messages/expanded">
    Show all expanded messages in this frame.
  </a>

  <form action="/messages">
    Show response from this form within this frame.
  </form>
</turbo-frame>
```

## Eager-loaded frame

basic frame と同じように動作するが、コンテンツはまずリモートの `src` から読み込まれる。

```html
<turbo-frame id="messages" src="/messages">
  Content will be replaced when /messages has been loaded.
</turbo-frame>
```

## Lazy-loaded frame

eager-loaded frame と同様だが、frame が表示されるまで `src` からのコンテンツ読み込みが行われない。

```html
<turbo-frame id="messages" src="/messages" loading="lazy">
  Content will be replaced when the frame becomes visible and /messages has been loaded.
</turbo-frame>
```

## Frame targeting the whole page by default

```html
<turbo-frame id="messages" target="_top">
  <a href="/messages/1">
    Following link will replace the whole page, not this frame.
  </a>

  <a href="/messages/1" data-turbo-frame="_self">
    Following link will replace just this frame.
  </a>

  <form action="/messages">
    Submitting form will replace the whole page, not this frame.
  </form>
</turbo-frame>
```

## Frame with overwritten navigation targets

```html
<turbo-frame id="messages">
  <a href="/messages/1">
    Following link will replace this frame.
  </a>

  <a href="/messages/1" data-turbo-frame="_top">
    Following link will replace the whole page, not this frame.
  </a>

  <form action="/messages" data-turbo-frame="navigation">
    Submitting form will replace the navigation frame.
  </form>
</turbo-frame>
```

## Frame that promotes navigations to Visits

```html
<turbo-frame id="messages" data-turbo-action="advance">
  <a href="/messages?page=2">Advance history to next page</a>
  <a href="/messages?page=2" data-turbo-action="replace">Replace history with next page</a>
</turbo-frame>
```

## Frame that will get reloaded with morphing during page refreshes & when they are explicitly reloaded with .reload()

```html
<turbo-frame id="my-frame" refresh="morph" src="/my_frame">
</turbo-frame>
```

# Attributes, properties, and functions

`<turbo-frame>` 要素は独自の HTML 属性と JavaScript プロパティを備えた [custom element](https://developer.mozilla.org/en-US/docs/Web/Web_Components/Using_custom_elements) である。

## HTML Attributes

-

`src` は要素のナビゲーションを制御する URL またはパス値を受け取る。

-

`loading` には 2 つの有効な [enumerated](https://www.w3.org/TR/html52/infrastructure.html#keywords-and-enumerated-attributes) 値「eager」と「lazy」がある。`loading="eager"` の場合、`src` 属性の変更は要素を即座にナビゲートさせる。`loading="lazy"` の場合、`src` 属性の変更は要素がビューポート内に表示されるまでナビゲーションを遅延させる。デフォルト値は `eager` である。

-

`busy` は [boolean attribute](https://www.w3.org/TR/html52/infrastructure.html#sec-boolean-attributes) で、`<turbo-frame>` が発行したリクエストの開始時に付与され、リクエストの終了時に外される。

-

`disabled` は [boolean attribute](https://www.w3.org/TR/html52/infrastructure.html#sec-boolean-attributes) で、存在する場合はあらゆるナビゲーションを抑止する。

-

`target` は、子孫の `<a>` がクリックされたときにナビゲートする別の `<turbo-frame>` 要素を ID で参照する。`target="_top"` の場合はウィンドウをナビゲートする。

-

`complete` は boolean attribute であり、その有無は `<turbo-frame>` 要素がナビゲーションを完了したかどうかを示す。

-

`recurse` は、`<turbo-frame>` が別のネストされた frame を含むコンテンツを読み込めるようにする。`recurse` は、frame の読み込まれたコンテンツ内に存在する別の `<turbo-frame>` 要素を ID で参照する。これは、Turbo がターゲットの frame を直接含まないが、そのターゲット frame を読み込める別の frame を含むレスポンスから frame のコンテンツを抽出する必要がある場合に利用できる。

たとえば、`/frames.html` が次の内容を含むとする。

```html
<turbo-frame id="recursive" recurse="composer" src="recursive.html">
</turbo-frame>
```

そして `recursive.html` は次の内容を含むとする。

```html
<turbo-frame id="recursive">
  <turbo-frame id="composer">
    <a href="frames.html">Link</a>
  </turbo-frame>
</turbo-frame>
```

`Link` が `frames.html` にナビゲートし直すとき、Turbo はレスポンス内で ID `composer` の frame を見つけて、現在の `composer` frame を更新する必要がある。しかし `frames.html` 内には直接そのような frame が存在しないため、Turbo は `recurse="composer"` を持つ frame を見つけ、これをアクティブにしてその `src` (`recursive.html`) を読み込ませ、その読み込まれたコンテンツから `composer` frame を検索・抽出して元の frame を更新する。

- `autoscroll` は [boolean attribute](https://www.w3.org/TR/html52/infrastructure.html#sec-boolean-attributes) で、`<turbo-frame>` 要素 (およびその子孫の `<turbo-frame>` 要素) を読み込み後にスクロールして表示するかどうかを制御する。スクロールの垂直方向の位置合わせは、`data-autoscroll-block` 属性に有効な [Element.scrollIntoView({ block: “…” })](https://developer.mozilla.org/en-US/docs/Web/API/Element/scrollIntoView#parameters) 値、すなわち `"end"`, `"start"`, `"center"`, `"nearest"` のいずれかを設定して制御する。`data-autoscroll-block` がない場合のデフォルト値は `"end"` である。スクロールの挙動は、`data-autoscroll-behavior` 属性に有効な [Element.scrollIntoView({ behavior: “…” })](https://developer.mozilla.org/en-US/docs/Web/API/Element/scrollIntoView#parameters) 値、すなわち `"auto"` または `"smooth"` を設定して制御する。`data-autoscroll-behavior` がない場合のデフォルト値は `"auto"` である。

## Properties

すべての `<turbo-frame>` 要素は、`FrameElement` クラスのインスタンスを通じて JavaScript 環境から制御できる。

-

`FrameElement.src` は読み込むパス名または URL を制御する。`src` プロパティを設定すると要素は即座にナビゲートする。`FrameElement.loaded` が `"lazy"` に設定されている場合、`src` プロパティの変更は要素がビューポート内に表示されるまでナビゲーションを遅延させる。

-

`FrameElement.disabled` は boolean プロパティで、要素が読み込みを行うかどうかを制御する。

-

`FrameElement.loading` は、frame がコンテンツを読み込むスタイル (`"eager"` または `"lazy"`) を制御する。

-

`FrameElement.loaded` は、`FrameElement` の現在のナビゲーションが完了した時点で解決する [Promise](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise) インスタンスを参照する。

-

`FrameElement.complete` は読み取り専用の boolean プロパティで、`FrameElement` がナビゲーションを完了した場合は `true`、それ以外は `false` となる。

-

`FrameElement.autoscroll` は、読み込み完了後に要素をスクロールして表示するかどうかを制御する。

-

`FrameElement.isActive` は読み取り専用の boolean プロパティで、frame が読み込まれ操作可能な状態であるかどうかを示す。

-

`FrameElement.isPreview` は読み取り専用の boolean プロパティで、要素を含む `document` が [preview](../handbook/07_building.md#detecting-when-a-preview-is-visible) の場合に `true` を返す。

## Functions

- `FrameElement.reload()` は、frame 要素をその `src` から再読み込みする関数である。
