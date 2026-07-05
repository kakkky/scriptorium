> 原文: https://turbo.hotwired.dev/reference/events
> タイトル: Events
> 最終取得: 2026-07-05

# Events

Turbo はさまざまな [Custom Events](https://developer.mozilla.org/en-US/docs/Web/API/CustomEvent) を発火する。以下のソースから dispatch される。

- [Document](#document)

- [Page Refreshes](#page-refreshes)

- [Forms](#forms)

- [Frames](#frames)

- [Streams](#streams)

- [HTTP Requests](#http-requests)

jQuery を使う場合、イベントの data は `$event.originalEvent.detail` としてアクセスする必要がある。

## Document

Turbo Drive は、ナビゲーションのライフサイクルを追跡し、ページの読み込みに応答するためのイベントを発火する。特に注記のない限り、以下のイベントは [document.documentElement](https://developer.mozilla.org/en-US/docs/Web/API/Document/documentElement) オブジェクト (すなわち `<html>` 要素) 上で発火する。

### `turbo:click`

Turbo が有効なリンクをクリックしたときに発火する。クリックされた要素が [event.target](https://developer.mozilla.org/en-US/docs/Web/API/Event/target) となる。リクエスト先のロケーションには `event.detail.url` でアクセスできる。このイベントをキャンセルすると、クリックは通常のブラウザナビゲーションとしてそのまま通される。

`event.detail` property
Type
Description

`url`
`string`
the requested location

`originalEvent`
[MouseEvent](https://developer.mozilla.org/en-US/docs/Web/API/MouseEvent)
the original [`click` event](https://developer.mozilla.org/en-US/docs/Web/API/Element/click_event)

### `turbo:before-visit`

ロケーションを訪問する直前に発火する。ただし履歴によるナビゲーションのときは発火しない。リクエスト先のロケーションには `event.detail.url` でアクセスできる。このイベントをキャンセルするとナビゲーションを止められる。

`event.detail` property
Type
Description

`url`
`string`
the requested location

### `turbo:visit`

visit が開始された直後に発火する。リクエスト先のロケーションには `event.detail.url` で、action には `event.detail.action` でアクセスできる。

`event.detail` property
Type
Description

`url`
`string`
the requested location

`action`
`"advance" | "replace" | "restore"`
the visit's [Action](../handbook/02_drive.md#page-navigation-basics)

### `turbo:before-cache`

Turbo が現在のページをキャッシュに保存する直前に発火する。

`turbo:before-cache` イベントのインスタンスには `event.detail` プロパティが無い。

### `turbo:before-render`

ページをレンダリングする直前に発火する。新しい `<body>` 要素には `event.detail.newBody` でアクセスできる。レンダリングは `event.detail.resume` によってキャンセルおよび再開できる ([Pausing Rendering](../handbook/02_drive.md#pausing-rendering) を参照)。Turbo Drive がレスポンスをレンダリングする方法をカスタマイズするには、`event.detail.render` 関数をオーバーライドする ([Custom Rendering](../handbook/02_drive.md#custom-rendering) を参照)。

`event.detail` property
Type
Description

`renderMethod`
`"replace" | "morph"`
the strategy that will be used to render the new content

`newBody`
[HTMLBodyElement](https://developer.mozilla.org/en-US/docs/Web/API/HTMLBodyElement)
the new `<body>` element that will replace the document's current `<body>` element

`resume`
`(value?: any) => void`
called when [Pausing Requests](../handbook/02_drive.md#pausing-requests)

`render`
`(currentBody, newBody) => void`
override to [Customize Rendering](../handbook/02_drive.md#custom-rendering)

### `turbo:render`

Turbo がページをレンダリングした後に発火する。application visit でキャッシュされたロケーションを訪れる場合、このイベントは 2 回発火する。1 回目はキャッシュ版をレンダリングした後、2 回目は fresh 版をレンダリングした後である。

`event.detail` property
Type
Description

`renderMethod`
`"replace" | "morph"`
the strategy used to render the new content

### `turbo:load`

初回のページロード後に一度、そしてその後の Turbo visit ごとに発火する。

`event.detail` property
Type
Description

`url`
`string`
the requested location

`timing.visitStart`
`number`
timestamp at the start of the Visit

`timing.requestStart`
`number`
timestamp at the start of the HTTP request for the next page

`timing.requestEnd`
`number`
timestamp at the end of the HTTP request for the next page

`timing.visitEnd`
`number`
timestamp at the end of the Visit

### `turbo:reload`

Turbo がページ全体のリロードを行う直前に発火する。

`event.detail` property
Type
Description

`reason`
`string`
the reason for the reload (e.g. "turbo_disabled", "turbo_visit_control_is_reload", "tracked_element_mismatch", "request_failed")

## Page Refreshes

Turbo Drive はページの内容を morph している最中にイベントを発火する。

### `turbo:morph`

Turbo がページを morph した後に発火する。

`event.detail` property
Type
Description

`currentElement`
[Element](https://developer.mozilla.org/en-US/docs/Web/API/Element)
the original [Element](https://developer.mozilla.org/en-US/docs/Web/API/Element) that remains connected after the morph (most commonly `document.body`)

`newElement`
[Element](https://developer.mozilla.org/en-US/docs/Web/API/Element)
the [Element](https://developer.mozilla.org/en-US/docs/Web/API/Element) with the new attributes and children that is not connected after the morph

### `turbo:before-morph-element`

Turbo が要素を morph する直前に発火する。[event.target](https://developer.mozilla.org/en-US/docs/Web/API/Event/target) は、morph 後もドキュメントに接続され続ける元の要素を指す。`event.preventDefault()` を呼び出してこのイベントをキャンセルすると、morph をスキップして元の要素・その属性・子要素を保持できる。

`event.detail` property
Type
Description

`newElement`
[Element](https://developer.mozilla.org/en-US/docs/Web/API/Element)
the [Element](https://developer.mozilla.org/en-US/docs/Web/API/Element) with the new attributes and children that is not connected after the morph

### `turbo:before-morph-attribute`

Turbo が要素の属性を morph する直前に発火する。[event.target](https://developer.mozilla.org/en-US/docs/Web/API/Event/target) は、morph 後もドキュメントに接続され続ける元の要素を指す。`event.preventDefault()` を呼び出してこのイベントをキャンセルすると、morph をスキップして元の属性値を保持できる。

`event.detail` property
Type
Description

`attributeName`
`string`
the name of the attribute to be mutated

`mutationType`
`"update" | "remove"`
how the attribute will be mutated

### `turbo:morph-element`

Turbo が要素を morph した後に発火する。[event.target](https://developer.mozilla.org/en-US/docs/Web/API/Event/target) は、morph 後も接続され続ける、morph された要素を指す。

`event.detail` property
Type
Description

`newElement`
[Element](https://developer.mozilla.org/en-US/docs/Web/API/Element)
the [Element](https://developer.mozilla.org/en-US/docs/Web/API/Element) with the new attributes and children that is not connected after the morph

## Forms

Turbo Drive はフォームの送信、リダイレクト、送信失敗の際にイベントを発火する。以下のイベントは送信中の `<form>` 要素上で発火する。

### `turbo:submit-start`

フォーム送信中に発火する。[FormSubmission](drive.md#formsubmission) オブジェクトには `event.detail.formSubmission` でアクセスできる。バリデーション失敗後などにフォーム送信を中止するには `event.detail.formSubmission.stop()` を使う。jQuery を使う場合は `event.originalEvent.detail.formSubmission.stop()` を使う。

`event.detail` property
Type
Description

`formSubmission`
[FormSubmission](drive.md#formsubmission)
the `<form>` element submission

### `turbo:submit-end`

フォーム送信によって開始されたネットワークリクエストが完了した後に発火する。[FormSubmission](drive.md#formsubmission) オブジェクトには `event.detail.formSubmission` でアクセスでき、`event.detail` に含まれるプロパティも利用できる。

`event.detail` property
Type
Description

`formSubmission`
[FormSubmission](drive.md#formsubmission)
the `<form>` element submission

`success`
`boolean`
a `boolean` representing the request's success

`fetchResponse`
[FetchResponse](drive.md#fetchresponse) | `undefined`
present when a response is received, even if `success: false`. `undefined` if the request errored before a response was received

`error`
[Error](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Errors) | `undefined`
`undefined` unless an actual fetch error occurs (e.g., network issues)

## Frames

Turbo Frames はナビゲーションのライフサイクル中にイベントを発火する。以下のイベントは `<turbo-frame>` 要素上で発火する。

### `turbo:before-frame-render`

`<turbo-frame>` 要素をレンダリングする直前に発火する。新しい `<turbo-frame>` 要素には `event.detail.newFrame` でアクセスできる。レンダリングは `event.detail.resume` によってキャンセルおよび再開できる ([Pausing Rendering](../handbook/04_frames.md#pausing-rendering) を参照)。Turbo Drive がレスポンスをレンダリングする方法をカスタマイズするには、`event.detail.render` 関数をオーバーライドする ([Custom Rendering](../handbook/04_frames.md#custom-rendering) を参照)。

`event.detail` property
Type
Description

`newFrame`
`FrameElement`
the new `<turbo-frame>` element that will replace the current `<turbo-frame>` element

`resume`
`(value?: any) => void`
called when [Pausing Requests](../handbook/02_drive.md#pausing-requests)

`render`
`(currentFrame, newFrame) => void`
override to [Customize Rendering](../handbook/02_drive.md#custom-rendering)

### `turbo:frame-render`

`<turbo-frame>` 要素がビューをレンダリングした直後に発火する。該当の `<turbo-frame>` 要素が [event.target](https://developer.mozilla.org/en-US/docs/Web/API/Event/target) となる。[FetchResponse](drive.md#fetchresponse) オブジェクトには `event.detail.fetchResponse` プロパティでアクセスできる。

`event.detail` property
Type
Description

`fetchResponse`
[FetchResponse](drive.md#fetchresponse)
the HTTP request's response

### `turbo:frame-load`

`<turbo-frame>` 要素がナビゲートされ、読み込みを終えたときに発火する (`turbo:frame-render` の後に発火する)。該当の `<turbo-frame>` 要素が [event.target](https://developer.mozilla.org/en-US/docs/Web/API/Event/target) となる。

`turbo:frame-load` イベントのインスタンスには `event.detail` プロパティが無い。

### `turbo:frame-missing`

`<turbo-frame>` 要素のリクエストに対するレスポンスに、対応する `<turbo-frame>` 要素が含まれていないときに発火する。デフォルトでは、Turbo は frame に案内のメッセージを書き込み、例外を投げる。このイベントをキャンセルすればこの挙動を上書きできる。[Response](https://developer.mozilla.org/en-US/docs/Web/API/Response) インスタンスには `event.detail.response` でアクセスでき、`event.detail.visit(location, visitOptions)` を呼び出せば visit を実行できる (`VisitOptions` については [Turbo.visit](drive.md#turbo.visit) を参照)。

`event.detail` property
Type
Description

`response`
[Response](https://developer.mozilla.org/en-US/docs/Web/API/Response)
the HTTP response for the request initiated by a `<turbo-frame>` element

`visit`
`async (location: string | URL, visitOptions: VisitOptions) => void`
a convenience function to initiate a page-wide navigation

## Streams

Turbo Streams はライフサイクル中にイベントを発火する。以下のイベントは `<turbo-stream>` 要素上で発火する。

### `turbo:before-stream-render`

Turbo Stream によるページ更新をレンダリングする直前に発火する。新しい `<turbo-stream>` 要素には `event.detail.newStream` でアクセスできる。要素の振る舞いをカスタマイズするには、`event.detail.render` 関数をオーバーライドする ([Custom Actions](../handbook/05_streams.md#custom-actions) を参照)。

`event.detail` property
Type
Description

`newStream`
`StreamElement`
the new `<turbo-stream>` element whose action will be executed

`render`
`async (currentElement) => void`
override to define [Custom Actions](../handbook/05_streams.md#custom-actions)

## HTTP Requests

Turbo は HTTP 越しにコンテンツを取得する際にイベントを発火する。リクエストを起こした主体に応じて、イベントは以下のいずれかで発火する。

- ナビゲーション中の `<turbo-frame>`

- 送信中の `<form>`

- ページ全体の Turbo Visit 中の `<html>` 要素

### `turbo:before-fetch-request`

Turbo がネットワークリクエストを発行する直前 (ページ取得、フォーム送信、リンクの preload など) に発火する。リクエスト先のロケーションには `event.detail.url` で、fetch の options オブジェクトには `event.detail.fetchOptions` でアクセスできる。このイベントは、それを引き起こした対応する要素 (`<turbo-frame>` または `<form>` 要素) 上で発火し、[event.target](https://developer.mozilla.org/en-US/docs/Web/API/Event/target) プロパティでアクセスできる。リクエストは `event.detail.resume` によってキャンセルおよび再開できる ([Pausing Requests](../handbook/02_drive.md#pausing-requests) を参照)。

`event.detail` property
Type
Description

`fetchOptions`
[RequestInit](https://developer.mozilla.org/en-US/docs/Web/API/Request/Request#options)
the `options` used to construct the [Request](https://developer.mozilla.org/en-US/docs/Web/API/Request/Request)

`url`
[URL](https://developer.mozilla.org/en-US/docs/Web/API/URL)
the request's location

`resume`
`(value?: any) => void` callback
called when [Pausing Requests](../handbook/02_drive.md#pausing-requests)

### `turbo:before-fetch-response`

ネットワークリクエストが完了した後に発火する。fetch の options オブジェクトには `event.detail` でアクセスできる。このイベントは、それを引き起こした対応する要素 (`<turbo-frame>` または `<form>` 要素) 上で発火し、[event.target](https://developer.mozilla.org/en-US/docs/Web/API/Event/target) プロパティでアクセスできる。

`event.detail` property
Type
Description

`fetchResponse`
[FetchResponse](drive.md#fetchresponse)
the HTTP request's response

### `turbo:before-prefetch`

Turbo がリンクを prefetch する直前に発火する。リンクが `event.target` となる。このイベントをキャンセルすると prefetch を止められる。

### `turbo:fetch-request-error`

フォームまたは frame の fetch リクエストがネットワークエラーで失敗したときに発火する。このイベントは、それを引き起こした対応する要素 (`<turbo-frame>` または `<form>` 要素) 上で発火し、[event.target](https://developer.mozilla.org/en-US/docs/Web/API/Event/target) プロパティでアクセスできる。このイベントはキャンセル可能である。

`event.detail` property
Type
Description

`request`
[FetchRequest](drive.md#fetchrequest)
The HTTP request that failed

`error`
[Error](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Errors)
provides the cause of the failure
