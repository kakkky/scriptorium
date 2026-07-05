> 原文: https://turbo.hotwired.dev/reference/drive
> タイトル: Drive
> 最終取得: 2026-07-05

# Drive

## Turbo.visit

```js
Turbo.visit(location)
Turbo.visit(location, { action: action })
Turbo.visit(location, { frame: frame })
```

指定した *location* (URL またはパスを含む文字列) に対して、指定した *action* (`"advance"` または `"replace"` のいずれかの文字列) で [Application Visit](../handbook/02_drive.md#application-visits) を実行する。

*location* がクロスオリジンの URL である場合、または指定されたルート ([Setting a Root Location](../handbook/02_drive.md#setting-a-root-location) を参照) の外にある場合、Turbo は `window.location` を設定してフルページロードを行う。

*action* が指定されていない場合、Turbo Drive は値として `"advance"` を仮定する。

visit を実行する前に、Turbo Drive は `document` 上で `turbo:before-visit` イベントを発火する。アプリケーションはこのイベントをリッスンし、`event.preventDefault()` によって visit をキャンセルできる ([Canceling Visits Before They Start](../handbook/02_drive.md#canceling-visits-before-they-start) を参照)。

*frame* が指定された場合、その値に一致する `[id]` 属性を持つ `<turbo-frame>` 要素を探し、その frame を指定した *location* にナビゲートする。該当する `<turbo-frame>` が見つからない場合は、ページ全体の [Application Visit](../handbook/02_drive.md#application-visits) を実行する。

## Turbo.cache.clear

```js
Turbo.cache.clear()
```

Turbo Drive のページキャッシュから全エントリを削除する。キャッシュされたページに影響を及ぼす可能性のあるサーバー側の状態変化があったときに呼び出すこと。

**Note:** この関数はかつて `Turbo.clearCache()` として公開されていた。トップレベル関数は新しい `Turbo.cache.clear()` 関数を優先して非推奨となった。

## Turbo.config.drive.progressBarDelay

```js
Turbo.config.drive.progressBarDelay = delayInMilliseconds
```

ナビゲーション中に [progress bar](../handbook/02_drive.md#displaying-progress) が表示されるまでの遅延時間をミリ秒で設定する。progress bar はデフォルトで 500ms 後に表示される。

このメソッドは iOS または Android アダプタと共に使用した場合は効果を持たない。

**Note:** この関数はかつて `Turbo.setProgressBarDelay` として公開されていた。トップレベル関数は新しい `Turbo.config.drive.progressBarDelay = delayInMilliseconds` の記法を優先して非推奨となった。

## Turbo.config.forms.confirm

```js
Turbo.config.forms.confirm = confirmMethod
```

[`data-turbo-confirm`](../handbook/02_drive.md#requiring-confirmation-for-a-visit) が付与されたリンクによって呼び出されるメソッドを設定する。デフォルトはブラウザ組み込みの `confirm` である。このメソッドは、visit を続行すべきかどうかに応じて true または false に解決される `Promise` オブジェクトを返さなければならない。

**Note:** この関数はかつて `Turbo.setConfirmMethod` として公開されていた。トップレベル関数は新しい `Turbo.config.forms.confirm = confirmMethod` の記法を優先して非推奨となった。

## Turbo.session.drive

```js
Turbo.session.drive = false
```

Turbo Drive をデフォルトで無効化する。以降、リンク単位・フォーム単位で `data-turbo="true"` を用いて Turbo Drive をオプトインする必要がある。

## `FetchRequest`

Turbo は [HTTP リクエスト中にさまざまなイベント](events.md#http-requests) をディスパッチする。それらは以下のプロパティを持つ `FetchRequest` オブジェクトを参照する。

Property
Type
Description

`body`
[FormData](https://developer.mozilla.org/en-US/docs/Web/API/FormData) | [URLSearchParams](https://developer.mozilla.org/en-US/docs/Web/API/URLSearchParams)
`"get"` リクエストの場合は [URLSearchParams](https://developer.mozilla.org/en-US/docs/Web/API/URLSearchParams) のインスタンス、それ以外は [FormData](https://developer.mozilla.org/en-US/docs/Web/API/FormData)

`enctype`
`"application/x-www-form-urlencoded" | "multipart/form-data" | "text/plain"`
[HTMLFormElement.enctype](https://developer.mozilla.org/en-US/docs/Web/API/HTMLFormElement/enctype) の値

`fetchOptions`
[RequestInit](https://developer.mozilla.org/en-US/docs/Web/API/Request/Request#options)
リクエストの設定オプション

`headers`
[Headers](https://developer.mozilla.org/en-US/docs/Web/API/Request/Request#headers) | `{ [string]: [any] }`
リクエストの HTTP ヘッダ

`method`
`"get" | "post" | "put" | "patch" | "delete"`
HTTP メソッド

`params`
[URLSearchParams](https://developer.mozilla.org/en-US/docs/Web/API/URLSearchParams)
`"get"` リクエストの場合の [URLSearchParams](https://developer.mozilla.org/en-US/docs/Web/API/URLSearchParams) インスタンス

`target`
[HTMLFormElement](https://developer.mozilla.org/en-US/docs/Web/API/HTMLFormElement) | [HTMLAnchorElement](https://developer.mozilla.org/en-US/docs/Web/API/HTMLAnchorElement) | `FrameElement` | `null`
リクエストの発生源となった要素

`url`
[URL](https://developer.mozilla.org/en-US/docs/Web/API/URL)
リクエストの [URL](https://developer.mozilla.org/en-US/docs/Web/API/URL)

## `FetchResponse`

Turbo は [HTTP リクエスト中にさまざまなイベント](events.md#http-requests) をディスパッチする。それらは以下のプロパティを持つ `FetchResponse` オブジェクトを参照する。

Property
Type
Description

`clientError`
`boolean`
ステータスが 400 から 499 の範囲にあれば `true`、そうでなければ `false`

`contentType`
`string` | `null`
[Content-Type](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Type) ヘッダの値

`failed`
`boolean`
レスポンスが成功しなかった場合は `true`、そうでなければ `false`

`isHTML`
`boolean`
コンテンツタイプが HTML の場合は `true`、そうでなければ `false`

`location`
[URL](https://developer.mozilla.org/en-US/docs/Web/API/URL)
[Response.url](https://developer.mozilla.org/en-US/docs/Web/API/Response/url) の値

`redirected`
`boolean`
[Response.redirected](https://developer.mozilla.org/en-US/docs/Web/API/Response/redirected) の値

`responseHTML`
`Promise<string>`
`Response` が HTML の場合にクローンし、[Response.text()](https://developer.mozilla.org/en-US/docs/Web/API/Response/text) を呼び出す

`responseText`
`Promise<string>`
`Response` をクローンし、[Response.text()](https://developer.mozilla.org/en-US/docs/Web/API/Response/text) を呼び出す

`response`
[Response](https://developer.mozilla.org/en-US/docs/Web/API/Response)
`Response` のインスタンス

`serverError`
`boolean`
ステータスが 500 から 599 の範囲にあれば `true`、そうでなければ `false`

`statusCode`
`number`
[Response.status](https://developer.mozilla.org/en-US/docs/Web/API/Response/status) の値

`succeeded`
`boolean`
[Response.ok](https://developer.mozilla.org/en-US/docs/Web/API/Response/ok) の場合は `true`、そうでなければ `false`

## `FormSubmission`

Turbo は [フォーム送信中にさまざまなイベント](events.md#forms) をディスパッチする。それらは以下のプロパティを持つ `FormSubmission` オブジェクトを参照する。

Property
Type
Description

`action`
`string`
`<form>` 要素の送信先

`body`
[FormData](https://developer.mozilla.org/en-US/docs/Web/API/FormData) | [URLSearchParams](https://developer.mozilla.org/en-US/docs/Web/API/URLSearchParams)
背後にある [Request](https://developer.mozilla.org/en-US/docs/Web/API/Request) のペイロード

`enctype`
`"application/x-www-form-urlencoded" | "multipart/form-data" | "text/plain"`
[HTMLFormElement.enctype](https://developer.mozilla.org/en-US/docs/Web/API/HTMLFormElement/enctype)

`fetchRequest`
[FetchRequest](#fetchrequest)
背後にある [FetchRequest](#fetchrequest) インスタンス

`formElement`
[HTMLFormElement](https://developer.mozilla.org/en-US/docs/Web/API/HTMLFormElement)
送信中の `<form>` 要素

`isSafe`
`boolean`
`method` が `"get"` の場合は `true`、そうでなければ `false`

`location`
[URL](https://developer.mozilla.org/en-US/docs/Web/API/URL)
`action` 文字列を [URL](https://developer.mozilla.org/en-US/docs/Web/API/URL) インスタンスへ変換したもの

`method`
`"get" | "post" | "put" | "patch" | "delete"`
HTTP メソッド

`submitter`
[HTMLButtonElement](https://developer.mozilla.org/en-US/docs/Web/API/HTMLButtonElement) | [HTMLInputElement](https://developer.mozilla.org/en-US/docs/Web/API/HTMLInputElement) | `undefined`
`<form>` 要素の送信を発生させた要素
