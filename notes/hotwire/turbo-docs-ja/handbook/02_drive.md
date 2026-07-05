> 原文: https://turbo.hotwired.dev/handbook/drive
> タイトル: Navigate with Turbo Drive
> 最終取得: 2026-07-05

# Navigate with Turbo Drive

Turbo Drive は、Turbo のうちページレベルのナビゲーションを強化する部分です。リンククリックやフォーム送信を監視し、それらをバックグラウンドで実行して、フルリロードなしでページを更新します。これは、以前 [Turbolinks](https://github.com/turbolinks/turbolinks) として知られていたライブラリの進化版です。

## Page Navigation Basics

Turbo Drive はページのナビゲーションを、*action* を伴う *location* (URL) への *visit* としてモデル化します。

visit は、クリックからレンダリングまでの一連のナビゲーションのライフサイクル全体を表します。これにはブラウザ履歴の変更、ネットワークリクエストの発行、キャッシュからのページのコピーの復元、最終的なレスポンスのレンダリング、スクロール位置の更新が含まれます。

レンダリング中、Turbo Drive はリクエスト元ドキュメントの `<body>` の中身をレスポンスドキュメントの `<body>` の中身で置き換え、両者の `<head>` の中身もマージし、必要に応じて `<html>` 要素の `lang` 属性を更新します。`<head>` 要素を置き換えるのではなくマージするのがポイントで、たとえば `<title>` や `<meta>` タグが変わった場合は期待どおり更新される一方、アセットへのリンクが同じであればそれらは触られず、ブラウザが再処理することもありません。

visit には 2 種類あります。ひとつは *application visit* で、*advance* または *replace* の action を持ちます。もうひとつは *restoration visit* で、*restore* の action を持ちます。

## Application Visits

application visit は、Turbo Drive が有効なリンクをクリックするか、[`Turbo.visit(location)`](../reference/drive.md#turbodrivevisit) をプログラムから呼び出すことで開始されます。

application visit は常にネットワークリクエストを発行します。レスポンスが届くと、Turbo Drive はその HTML をレンダリングし、visit を完了します。

可能であれば、Turbo Drive は visit の開始直後にキャッシュからページのプレビューをレンダリングします。これによって、同じページ間の頻繁なナビゲーションで体感速度が向上します。

visit の location にアンカーが含まれている場合、Turbo Drive はアンカー先の要素へのスクロールを試みます。そうでない場合はページの先頭にスクロールします。

application visit はブラウザの履歴を変更します。visit の *action* によってその変更のしかたが決まります。

![Advance visit action](https://s3.amazonaws.com/turbolinks-docs/images/advance.svg)

デフォルトの visit action は *advance* です。advance visit の間、Turbo Drive は [`history.pushState`](https://developer.mozilla.org/en-US/docs/Web/API/History/pushState) を使ってブラウザの履歴スタックに新しいエントリを push します。

Turbo Drive の [iOS アダプタ](https://github.com/hotwired/turbo-ios)を使うアプリケーションでは、通常、新しいビューコントローラをナビゲーションスタックに push することで advance visit を扱います。同様に、[Android アダプタ](https://github.com/hotwired/turbo-android)を使うアプリケーションでは、通常、新しいアクティビティをバックスタックに push します。

![Replace visit action](https://s3.amazonaws.com/turbolinks-docs/images/replace.svg)

新しい履歴エントリをスタックに push せずに、ある location を visit したいこともあるでしょう。*replace* visit action は [`history.replaceState`](https://developer.mozilla.org/en-US/docs/Web/API/History/replaceState) を使って、最上位の履歴エントリを破棄し、新しい location で置き換えます。

リンクをたどったときに replace visit を発火させるように指定するには、リンクに `data-turbo-action="replace"` を付けます。

```html
<a href="/edit" data-turbo-action="replace">Edit</a>
```

プログラムから replace action で location を visit するには、`Turbo.visit` に `action: "replace"` オプションを渡します。

```js
Turbo.visit("/edit", { action: "replace" })
```

Turbo Drive の [iOS アダプタ](https://github.com/hotwired/turbo-ios)を使うアプリケーションでは、通常、最上位のビューコントローラを dismiss してから、アニメーションなしで新しいビューコントローラをナビゲーションスタックに push することで replace visit を扱います。

## Restoration Visits

Turbo Drive は、ブラウザの戻る/進むボタンでナビゲーションしたときに、自動的に restoration visit を開始します。[iOS](https://github.com/hotwired/turbo-ios) や [Android](https://github.com/hotwired/turbo-android) アダプタを使うアプリケーションでは、ナビゲーションスタックを後ろに戻すときに restoration visit が開始されます。

![Restore visit action](https://s3.amazonaws.com/turbolinks-docs/images/restore.svg)

可能であれば、Turbo Drive はリクエストを行わずにキャッシュからページのコピーをレンダリングします。そうでない場合は、ネットワーク越しにページの新しいコピーを取得します。詳しくは [Understanding Caching](07_building.md#understanding-caching) を参照してください。

Turbo Drive は各ページから離れる前にスクロール位置を保存しておき、restoration visit ではその保存された位置に自動的に戻します。

restoration visit の action は *restore* で、Turbo Drive はこれを内部利用のために予約しています。`restore` を action として、リンクに付与したり `Turbo.visit` を呼び出したりしないでください。

## Canceling Visits Before They Start

application visit は、リンククリックによって始まったか [`Turbo.visit`](../reference/drive.md#turbovisit) の呼び出しによって始まったかにかかわらず、開始前にキャンセルできます。

visit が始まろうとしているタイミングで通知を受けるには `turbo:before-visit` イベントを listen し、`event.detail.url` (jQuery を使っている場合は `$event.originalEvent.detail.url`) で visit の location を確認します。そして `event.preventDefault()` を呼び出すことで visit をキャンセルします。

restoration visit はキャンセルできず、`turbo:before-visit` も発火しません。Turbo Drive は、*すでに発生してしまった*履歴ナビゲーション (通常はブラウザの戻る/進むボタン) への応答として restoration visit を発行するためです。

## Custom Rendering

アプリケーションは、ドキュメント全体に対する `turbo:before-render` イベントリスナを追加し、`event.detail.render` プロパティを上書きすることでレンダリング処理をカスタマイズできます。

たとえば、[idiomorph](https://github.com/bigskysoftware/idiomorph) や [morphdom](https://github.com/patrick-steele-idem/morphdom) を使って、レスポンスドキュメントの `<body>` 要素をリクエスト元ドキュメントの `<body>` 要素にマージすることもできます。

```javascript
import { Idiomorph } from "idiomorph"

addEventListener("turbo:before-render", (event) => {
  event.detail.render = (currentElement, newElement) => {
    Idiomorph.morph(currentElement, newElement)
  }
})
```

## Pausing Rendering

アプリケーションはレンダリングを一時停止して、続行前に追加の準備を行うことができます。

レンダリングが始まろうとしているタイミングで通知を受けるには `turbo:before-render` イベントを listen し、`event.preventDefault()` で一時停止します。準備が済んだら `event.detail.resume()` を呼び出してレンダリングを続行します。

ユースケースとしては、visit のための退出アニメーションを追加する例があります。

```javascript
document.addEventListener("turbo:before-render", async (event) => {
  event.preventDefault()

  await animateOut()

  event.detail.resume()
})
```

## Pausing Requests

アプリケーションはリクエストを一時停止して、実行前に追加の準備を行うことができます。

リクエストが始まろうとしているタイミングで通知を受けるには `turbo:before-fetch-request` イベントを listen し、`event.preventDefault()` で一時停止します。準備が済んだら `event.detail.resume()` を呼び出してリクエストを続行します。

ユースケースとしては、リクエストに `Authorization` ヘッダーをセットする例があります。

```javascript
document.addEventListener("turbo:before-fetch-request", async (event) => {
  event.preventDefault()

  const token = await getSessionToken(window.app)
  event.detail.fetchOptions.headers["Authorization"] = `Bearer ${token}`

  event.detail.resume()
})
```

## Performing Visits With a Different Method

デフォルトでは、リンククリックはサーバーに `GET` リクエストを送ります。しかし `data-turbo-method` でこれを変更できます。

```html
<a href="/articles/54" data-turbo-method="delete">Delete the article</a>
```

アクセシビリティの観点から、GET 以外のものについては本物のフォームとボタンを使う方が望ましい、という点は考慮すべきです。

## Requiring Confirmation for a Visit

リンクに `data-turbo-confirm` と `data-turbo-method` の両方を付けると、visit を続行するために確認が必要になります。

```html
<a href="/articles" data-turbo-method="get" data-turbo-confirm="Do you want to leave this page?">Back to articles</a>
<a href="/articles/54" data-turbo-method="delete" data-turbo-confirm="Are you sure you want to delete the article?">Delete the article</a>
```

確認のために呼ばれるメソッドを変更するには `Turbo.config.forms.confirm = confirmMethod` を使います。デフォルトはブラウザ組み込みの `confirm` です。

## Disabling Turbo Drive on Specific Links or Forms

Turbo Drive は、要素自体またはその祖先要素のいずれかに `data-turbo="false"` を付けることで、要素単位で無効化できます。

```html
<a href="/" data-turbo="false">Disabled</a>

<form action="/messages" method="post" data-turbo="false">
  <!-- … -->
</form>

<div data-turbo="false">
  <a href="/">Disabled</a>
  <form action="/messages" method="post">
    <!-- … -->
  </form>
</div>
```

祖先要素でオプトアウトされている状態から再有効化するには `data-turbo="true"` を使います。

```html
<div data-turbo="false">
  <a href="/" data-turbo="true">Enabled</a>
</div>
```

Turbo Drive が無効化されたリンクやフォームは、ブラウザによって通常どおり扱われます。

Drive をオプトアウトではなくオプトインにしたい場合は、`Turbo.session.drive = false` を設定できます。その場合、`data-turbo="true"` を要素ごとに使って Drive を有効にします。JavaScript pack で Turbo をインポートしている場合は、これをグローバルに行うこともできます。

```js
import { Turbo } from "@hotwired/turbo-rails"
Turbo.session.drive = false
```

## View transitions

[View Transition API](https://developer.mozilla.org/en-US/docs/Web/API/View_Transitions_API) を[サポートするブラウザ](https://caniuse.com/?search=View%20Transition%20API)では、Turbo がページ間の遷移時に view transition をトリガーできます。

Turbo は、現在のページと次のページの両方に次の meta タグがある場合に view transition をトリガーします。

```
<meta name="view-transition" content="same-origin" />
```

Turbo はまた、遷移の方向を示すために `<html>` 要素に `data-turbo-visit-direction` 属性を追加します。この属性は次のいずれかの値を取ります。

- advance visit では `forward`

- restoration visit では `back`

- replace visit では `none`

この属性を利用して、遷移中に行われるアニメーションをカスタマイズできます。

```css
html[data-turbo-visit-direction="forward"]::view-transition-old(sidebar):only-child {
  animation: slide-to-right 0.5s ease-out;
}
```

## Displaying Progress

Turbo Drive のナビゲーション中、ブラウザはネイティブのプログレスインジケーターを表示しません。Turbo Drive はリクエスト発行中のフィードバックとして、CSS ベースのプログレスバーをインストールします。

プログレスバーはデフォルトで有効になっており、読み込みに 500ms より長くかかるあらゆるページで自動的に表示されます。(この遅延は [`Turbo.setProgressBarDelay`](../reference/drive.md#turbodrivesetprogressbardelay) メソッドで変更できます。)

プログレスバーは `turbo-progress-bar` というクラス名を持つ `<div>` 要素です。そのデフォルトスタイルはドキュメントの先頭に現れ、あとから来るルールで上書きできます。

たとえば、以下の CSS を使うと太い緑色のプログレスバーになります。

```css
.turbo-progress-bar {
  height: 5px;
  background-color: green;
}
```

プログレスバーを完全に無効化するには、その `visibility` スタイルを `hidden` に設定します。

```css
.turbo-progress-bar {
  visibility: hidden;
}
```

プログレスバーと連動して、Turbo Drive は Visits や Form Submissions から始まったページナビゲーションの間、ページの `<html>` 要素の [`[aria-busy]` 属性](https://www.w3.org/TR/wai-aria/#aria-busy)も切り替えます。Turbo Drive はナビゲーション開始時に `[aria-busy="true"]` をセットし、ナビゲーション完了時に `[aria-busy]` 属性を削除します。

## Reloading When Assets Change

前述のとおり、Turbo Drive は `<head>` 要素の中身をマージします。しかし CSS や JavaScript が変わっている場合、そのマージは既存のものの上にそれらを評価することになります。通常、これは望ましくない衝突につながります。そのような場合には、標準の非 Ajax リクエストで完全に新しいドキュメントを取得する必要があります。

これを実現するには、対象のアセット要素に `data-turbo-track="reload"` を付け、アセット URL にバージョン識別子を含めるだけです。この識別子は数値、最終更新タイムスタンプ、あるいは、より望ましくは以下の例のようにアセット内容のダイジェストにできます。

```html
<head>
  <!-- … -->
  <link rel="stylesheet" href="/application-258e88d.css" data-turbo-track="reload">
  <script src="/application-cbd3cd4.js" data-turbo-track="reload"></script>
</head>
```

## Removing Assets When They Change

前述のとおり、Turbo Drive は `<head>` 要素の中身をマージします。あるページが他のページには存在しない外部アセット (CSS スタイルシートなど) に依存している場合、そのページから離れるときにそれらを取り除くと便利なことがあります。

`<link>` や `<style>` 要素を `[data-turbo-track="dynamic"]` 付きでレンダリングすると、Turbo Drive はその要素がナビゲーション後のレスポンスに含まれていない場合に動的に削除します。これは [`[data-turbo-track="reload"]`](#reloading-when-assets-change) 属性を補完する役割を果たし、スタイルにしか影響しない変更のデプロイでフルリロードが発生するのを避けられます。

```html
<head>
  <!-- … -->
  <link rel="stylesheet" href="/page-specific-styles-258e88d.css" data-turbo-track="dynamic">
  <style data-turbo-track="dynamic">
    .page-specific-styles { /* … */ }
  </style>
</head>
```

`<script>` 要素を `[data-turbo-track="dynamic"]` 付きでレンダリングすると、意図しない副作用があるかもしれない点に注意してください。`<script>` がドキュメントから切り離されても、JavaScript コンテキストは変わりませんし、その要素の既に評価された JavaScript コードがアンロードされたり、何らかの形で変更されたりすることはありません。

## Ensuring Specific Pages Trigger a Full Reload

あるページへの visit で必ずフルリロードをトリガーするようにするには、そのページの `<head>` に `<meta name="turbo-visit-control">` 要素を含めます。

```html
<head>
  <!-- … -->
  <meta name="turbo-visit-control" content="reload">
</head>
```

この設定は、Turbo Drive のページ変更とうまく連携しないサードパーティ JavaScript ライブラリのワークアラウンドとして役立つことがあります。

## Setting a Root Location

Turbo Drive は、現在のドキュメントと同じオリジン (すなわち、同じプロトコル、ドメイン名、ポート) を持つ URL しか読み込みません。他の URL への visit はフルページロードにフォールバックします。

場合によっては、Turbo Drive をさらに同じオリジン内の特定のパスにスコープしたいこともあるでしょう。たとえば、Turbo Drive アプリケーションが `/app` にあり、Turbo Drive を使わないヘルプサイトが `/help` にある場合、アプリからヘルプサイトへのリンクでは Turbo Drive を使うべきではありません。

Turbo Drive を特定のルート location にスコープするには、ページの `<head>` に `<meta name="turbo-root">` 要素を含めます。Turbo Drive は、このパスをプレフィックスとして持つ同一オリジン URL のみを読み込みます。

```html
<head>
  <!-- … -->
  <meta name="turbo-root" content="/app">
</head>
```

## Form Submissions

Turbo Drive はフォーム送信をリンククリックと似たやり方で扱います。大きな違いは、フォーム送信は HTTP POST メソッドを使ってステートフルなリクエストを発行できるのに対し、リンククリックは常にステートレスな HTTP GET リクエストしか発行しない、という点です。

送信の全体を通して、Turbo Drive は `<form>` 要素を対象とし、[ドキュメントを通じてバブルアップする](https://developer.mozilla.org/en-US/docs/Learn/JavaScript/Building_blocks/Events#event_bubbling)一連の[イベント](../reference/events.md)をディスパッチします。

1. `turbo:submit-start`

2. `turbo:before-fetch-request`

3. `turbo:before-fetch-response`

4. `turbo:submit-end`

送信中、Turbo Drive は送信開始時に「submitter」要素の [disabled](https://developer.mozilla.org/en-US/docs/Web/HTML/Attributes/disabled) 属性をセットし、送信終了後にその属性を削除します。`<form>` 要素を送信する際、ブラウザは送信を開始した `<input type="submit">` または `<button>` 要素を [submitter](https://developer.mozilla.org/en-US/docs/Web/API/SubmitEvent/submitter) として扱います。`<form>` 要素をプログラムから送信するには、[HTMLFormElement.requestSubmit()](https://developer.mozilla.org/en-US/docs/Web/API/HTMLFormElement/requestSubmit) メソッドを呼び出し、オプションのパラメータとして `<input type="submit">` または `<button>` 要素を渡します。

`<form>` 送信中に他にも行いたい変更 (たとえば、[送信された `<form>` 内のすべてのフィールド](https://developer.mozilla.org/en-US/docs/Web/API/HTMLFormElement/elements)を disable する、など) があれば、独自のイベントリスナを宣言できます。

```js
addEventListener("turbo:submit-start", ({ target }) => {
  for (const field of target.elements) {
    field.disabled = true
  }
})
```

## Redirecting After a Form Submission

フォーム送信によるステートフルなリクエストのあと、Turbo Drive はサーバーが [HTTP 303 リダイレクトレスポンス](https://en.wikipedia.org/wiki/HTTP_303)を返すことを期待します。Turbo Drive はそのリダイレクトをたどり、リロードなしでページのナビゲーションと更新を行います。

このルールの例外は、レスポンスが 4xx または 5xx ステータスコードでレンダリングされた場合です。これによって、サーバーが `422 Unprocessable Content` で応答することでフォームバリデーションエラーをレンダリングしたり、`500 Internal Server Error` で「Something Went Wrong」画面を表示したりできます。

Turbo が POST リクエストの 200 レスポンスで通常のレンダリングを許可しない理由は、ブラウザには POST visit でのリロードに対する組み込みの挙動があり、「Are you sure you want to submit this form again?」というダイアログを表示するのですが、Turbo ではこれを再現できないからです。代わりに、Turbo はレンダリングを試みるフォーム送信では現在の URL に留まり、URL をフォームの action に変更することはしません。もし変更してしまうと、リロード時にその action URL に対して GET が発行され、そのようなエンドポイントは存在しないかもしれないためです。

フォーム送信が GET リクエストの場合は、フォームに `data-turbo-frame` ターゲットを指定することで、直接レンダリングされたレスポンスをレンダリングできます。レンダリングと同時に URL も更新したい場合は `data-turbo-action` 属性も渡します。

## Streaming After a Form Submission

サーバーはフォーム送信に対して、`Content-Type: text/vnd.turbo-stream.html` ヘッダーとともに、レスポンスボディに 1 つ以上の `<turbo-stream>` 要素を含めた [Turbo Streams](05_streams.md) メッセージで応答することもできます。これによって、ナビゲーションを行うことなくページの複数の部分を更新できます。

## Prefetching Links on Hover

Turbo は、`mouseenter` イベントでリンクを自動的に読み込む (ユーザーがリンクをクリックする前に読み込む) ことで、体感的なリンクナビゲーションのレイテンシを短縮することもできます。これは通常、クリックによるナビゲーションあたり 500〜800ms の速度向上につながります。

リンクのプリフェッチは Turbo v8 以降デフォルトで有効ですが、ページに次の meta タグを追加することで無効化できます。

```html
<meta name="turbo-prefetch" content="false">
```

ユーザーがちょっとリンクにホバーしただけのときにプリフェッチしてしまうのを避けるため、Turbo はリンクにホバーされてから 100ms 待ってからプリフェッチを行います。しかし、レンダリングが重いページへのリンクなど、特定のリンクではプリフェッチの挙動を無効化したいこともあるでしょう。

要素またはその祖先に `data-turbo-prefetch="false"` を付けることで、要素単位でプリフェッチを無効化できます。

```html
<html>
  <head>
    <meta name="turbo-prefetch" content="true">
  </head>
  <body>
    <a href="/articles">Articles</a> <!-- This link is prefetched -->
    <a href="/about" data-turbo-prefetch="false">About</a> <!-- Not prefetched -->
    <div data-turbo-prefetch="false">
      <!-- Links inside this div will not be prefetched -->
    </div>
  </body>
</html>
```

親要素でプリフェッチを無効化した上で、子要素で `data-turbo-prefetch="true"` を使って個別にプリフェッチを許可することもできます。

```html
<html>
  <body data-turbo-prefetch="false">
    <nav id="header" data-turbo-prefetch="true">
      <a href="/articles">Articles</a> <!-- This link is prefetched -->
      <a href="/about">About</a> <!-- This one as well -->
    </nav>
    <div id="body">
      <!-- Links inside this div will not be prefetched -->
    </div>
    <footer id="footer" data-turbo-prefetch="true">
      <!-- Links inside this footer will be prefetched -->
    </footer>
  </body>
</html>
```

プログラム的にプリフェッチを無効化することもできます。`turbo:before-prefetch` イベントをインターセプトして `event.preventDefault()` を呼び出します。

```javascript
document.addEventListener("turbo:before-prefetch", (event) => {
  if (isSavingData() || hasSlowInternet()) {
    event.preventDefault()
  }
})

function isSavingData() {
  return navigator.connection?.saveData
}

function hasSlowInternet() {
  return navigator.connection?.effectiveType === "slow-2g" ||
         navigator.connection?.effectiveType === "2g"
}
```

## Preload Links Into the Cache

[data-turbo-preload](../reference/attributes.md#data-attributes) の真偽値属性を使うと、リンクを Turbo Drive のキャッシュにプリロードできます。

これによって、初回 visit の前でもページのプレビューを提供でき、ページ遷移を電光石火のように感じさせられます。アプリケーションの中で最も重要なページのプリロードに使ってください。必要のないコンテンツを読み込ませてしまうことにつながるため、使いすぎには注意してください。

すべての `<a>` 要素がプリロードできるわけではありません。`[data-turbo-preload]` 属性は、次のようなリンクには効果を持ちません。

- 別ドメインへのナビゲーション

- `<turbo-frame>` 要素を駆動する `[data-turbo-frame]` 属性を持つ

- 祖先の `<turbo-frame>` 要素を駆動する

- `[data-turbo="false"]` 属性を持つ

- `[data-turbo-stream]` 属性を持つ

- `[data-turbo-method]` 属性を持つ

- 祖先要素に `[data-turbo="false"]` 属性を持つものがある

- 祖先要素に `[data-turbo-prefetch="false"]` 属性を持つものがある

これは [Eager-Loading Frames](../reference/frames.md#eager-loaded-frame) や [Lazy-Loading Frames](../reference/frames.md#lazy-loaded-frame) を活用するページとも相性がよいです。ページの構造をプリロードしておくことで、意味のある中身が読み込まれる間、ユーザーに意味のある読み込み中の状態を見せられるためです。

## Ignored Paths

パス/URL の最後のレベルに `.` を含むパスは、`.htm`、`.html`、`.xhtml`、`.php` のファイル拡張子で終わっていない限り Turbo に扱われません。Turbo はこれらのパスをターゲットとするフォームやリンクを無視します。Turbo にこれらのパスを扱わせる最も手っ取り早い方法は、URL の末尾に `/` を追加することです。無視されるフォームの例:

```html
<form action="/messages.67" method="post">
  <!-- ignored -->
</form>

<form action="/messages.php.1" method="post" data-turbo="true">
  <!-- also ignored -->
</form>

<form action="/messages.json" method="post" data-turbo="true">
  <!-- also ignored -->
</form>
```

以下のフォームは扱われます。

```html
<form action="/messages/67" method="post">
  <!-- handled -->
</form>

<form action="/messages.67/action" method="post">
  <!-- also handled -->
</form>

<form action="/messages.php" method="post" data-turbo="true">
  <!-- also handled -->
</form>

<form action="/messages.json/" method="post" data-turbo="true">
  <!-- also handled -->
</form>

<form action="/messages.json/123" method="post" data-turbo="true">
  <!-- also handled -->
</form>
```

`data-turbo` の各種メソッド (`data-turbo="true"` を含む) を設定しても、`.` によって無視されるパスを Turbo に扱わせるように上書き・強制することはできません。

なお、プリロードされた `<a>` 要素は [turbo:before-fetch-request](../reference/events.md) および [turbo:before-fetch-response](../reference/events.md) イベントをディスパッチします。プリロードによって発火した `turbo:before-fetch-request` イベントを、他の仕組みで発火したイベントと区別するには、リクエストの `X-Sec-Purpose` ヘッダー (`event.detail.fetchOptions.headers["X-Sec-Purpose"]` プロパティから読み取れます) が `"prefetch"` になっているかを確認します。

```js
addEventListener("turbo:before-fetch-request", (event) => {
  if (event.detail.fetchOptions.headers["X-Sec-Purpose"] === "prefetch") {
    // do additional preloading setup…
  } else {
    // do something else…
  }
})
```

[次へ: Smooth page refreshes with morphing](03_page_refreshes.md)
