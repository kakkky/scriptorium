> 原文: https://turbo.hotwired.dev/handbook/building
> タイトル: Building Your Turbo Application
> 最終取得: 2026-07-05

# Building Your Turbo Application

Turbo が高速なのは、リンクを辿ったりフォームを送信したりしたときに、ページ全体を再ロードしないからです。あなたのアプリケーションはブラウザ内で永続的に動き続ける長寿命のプロセスになります。そのため、JavaScript の組み立て方を考え直す必要があります。

特に、ナビゲーションのたびに毎回フルページロードによって環境がリセットされる、という前提には頼れなくなります。JavaScript の `window` と `document` オブジェクトはページ遷移をまたいで状態を保持しますし、メモリ上に残しておいたそれ以外のオブジェクトもメモリに残り続けます。

この制約を意識し、少しだけ気を配れば、Turbo に密結合させることなく、アプリケーションをうまくこの制約に対応させて設計できます。

## Working with Script Elements

ブラウザは初回ページロードの際、そのページに存在する `<script>` 要素をすべて自動的にロードし、評価します。

新しいページに遷移すると、Turbo Drive は新しいページの `<head>` の中から、現在のページには存在しない `<script>` 要素を探します。そして、それらを現在の `<head>` に追加し、ブラウザにロード・評価させます。これを利用して、追加の JavaScript ファイルをオンデマンドでロードできます。

Turbo Drive は、ページの `<body>` にある `<script>` 要素を、そのページをレンダリングするたびに評価します。ページごとの JavaScript 状態を設定したり、クライアントサイドのモデルをブートストラップしたりするのに、body 内のインラインスクリプトを使えます。振る舞いをインストールしたり、ページの変化に応じてより複雑な操作を行ったりする場合は、script 要素は避けて、代わりに `turbo:load` イベントを使ってください。

レンダリング後に Turbo が評価してほしくない `<script>` 要素には、`data-turbo-eval="false"` を付けてください。ただし、この注釈は初回ページロード時にブラウザがスクリプトを評価するのを防ぐものではありません。

### Loading Your Application’s JavaScript Bundle

アプリケーションの JavaScript バンドルは、必ずドキュメントの `<head>` にある `<script>` 要素で読み込むようにしてください。そうしないと、Turbo Drive はページ遷移のたびにバンドルを再ロードしてしまいます。

```html
<head>
  ...
  <script src="/application-cbd3cd4.js" defer></script>
</head>
```

さらに、アセットパッケージングシステムを設定して、各スクリプトをフィンガープリント化し、内容が変わったら新しい URL になるようにするのも検討してください。そうすれば、新しい JavaScript バンドルをデプロイしたときに `data-turbo-track` 属性でフルページリロードを強制できます。詳しくは [Reloading When Assets Change](02_drive.md#reloading-when-assets-change) を参照してください。

## Understanding Caching

Turbo Drive は、最近訪れたページのキャッシュを保持します。このキャッシュは 2 つの目的で使われます。restoration visit のときにネットワークにアクセスせずにページを表示すること、そして application visit のときに一時的なプレビューを表示して体感パフォーマンスを向上させることです。

履歴を辿るナビゲーション ([Restoration Visits](02_drive.md#restoration-visits)) では、Turbo Drive は可能な限り、ネットワークから新しいコピーを取得せずにキャッシュからページを復元します。

一方、通常のナビゲーション ([Application Visits](02_drive.md#application-visits)) の場合、Turbo Drive はキャッシュからページを即座に復元してプレビューとして表示すると同時に、ネットワークから新しいコピーを取得します。これにより、よく訪れる場所については、あたかも瞬時にページがロードされたかのような錯覚を与えられます。

Turbo Drive は、新しいページをレンダリングする直前に、現在のページのコピーをキャッシュに保存します。なお、Turbo Drive は [`cloneNode(true)`](https://developer.mozilla.org/en-US/docs/Web/API/Node/cloneNode) を使ってページをコピーするため、アタッチされていたイベントリスナーや関連データはすべて破棄されます。

### Preparing the Page to be Cached

Turbo Drive がドキュメントをキャッシュする前に何らかの準備が必要な場合は、`turbo:before-cache` イベントを購読してください。このイベントを利用して、フォームをリセットしたり、展開されている UI 要素を折りたたんだり、サードパーティ製ウィジェットをティアダウンしたりして、ページが再表示できる状態に整えられます。

```js
document.addEventListener("turbo:before-cache", function() {
  // ...
})
```

ページ要素の中には、フラッシュメッセージやアラートのように、本質的に*一時的*なものがあります。それらをドキュメントと一緒にキャッシュしてしまうと、復元時に再表示されてしまい、たいていの場合これは望ましくありません。そうした要素には `data-turbo-temporary` を付けておくと、Turbo Drive はキャッシュする前にそれらをページから自動的に取り除きます。

```html
<body>
  <div class="flash" data-turbo-temporary>
    Your cart was updated!
  </div>
  ...
</body>
```

### Detecting When a Preview is Visible

Turbo Drive はキャッシュからプレビューを表示するとき、`<html>` 要素に `data-turbo-preview` 属性を追加します。この属性の有無をチェックすることで、プレビュー表示中だけ振る舞いを有効化・無効化するといったことができます。

```js
if (document.documentElement.hasAttribute("data-turbo-preview")) {
  // Turbo Drive is displaying a preview
}
```

### Opting Out of Caching

ページの `<head>` に `<meta name="turbo-cache-control">` 要素を入れてキャッシュディレクティブを宣言することで、ページ単位でキャッシュの挙動を制御できます。

`no-preview` ディレクティブを使うと、application visit のときにそのページのキャッシュ済みバージョンをプレビューとして表示しないよう指定できます。no-preview を指定されたページは restoration visit のときだけ使われます。

ページを一切キャッシュしないようにするには、`no-cache` ディレクティブを使ってください。no-cache を指定されたページは、restoration visit を含め、常にネットワーク越しにフェッチされます。

```html
<head>
  ...
  <meta name="turbo-cache-control" content="no-cache">
</head>
```

アプリケーション全体でキャッシュを完全に無効化したい場合は、すべてのページに no-cache ディレクティブを含めてください。

### Opting Out of Caching from the client-side

`<meta name="turbo-cache-control">` 要素の値は、`Turbo.cache` として公開されているクライアントサイド API からも制御できます。

```js
// Set cache control of current page to `no-cache`
Turbo.cache.exemptPageFromCache()

// Set cache control of current page to `no-preview`
Turbo.cache.exemptPageFromPreview()
```

どちらの関数も、`<meta name="turbo-cache-control">` 要素がまだ存在しなければ `<head>` に作成します。

以前に設定したキャッシュ制御値は、次のようにしてリセットできます。

```js
Turbo.cache.resetCacheControl()
```

## Installing JavaScript Behavior

JavaScript の振る舞いを、`window.onload`、`DOMContentLoaded`、jQuery の `ready` イベントに応答する形でインストールすることに慣れているかもしれません。Turbo では、これらのイベントは初回ページロードのときにしか発火せず、その後のページ遷移では発火しません。以下では、JavaScript の振る舞いを DOM に接続するための 2 つの戦略を比較します。

### Observing Navigation Events

Turbo Drive はナビゲーション中に一連のイベントをトリガーします。その中でも最も重要なのが `turbo:load` イベントで、これは初回ページロード時に一度、そしてその後 Turbo Drive の visit のたびに発火します。

`DOMContentLoaded` の代わりに `turbo:load` イベントを購読することで、ページ遷移のたびに JavaScript の振る舞いをセットアップできます。

```js
document.addEventListener("turbo:load", function() {
  // ...
})
```

このイベントが発火するとき、アプリケーションが常にまっさらな状態であるとは限らないことに注意してください。前のページ向けにインストールした振る舞いをクリーンアップする必要があるかもしれません。

また、Turbo Drive のナビゲーションだけがアプリケーション内でページを更新する唯一の経路とも限らないので、初期化コードを別の関数に切り出し、`turbo:load` からも、DOM を変更する他の場所からも呼べるようにすると良いかもしれません。

可能な限り、ページ本体の要素に直接イベントリスナーを追加するために `turbo:load` を使うのは避けてください。代わりに、[イベント委譲](https://learn.jquery.com/events/event-delegation/) を使って `document` や `window` に一度だけイベントリスナーを登録することを検討してください。

詳しくは [Full List of Events](../reference/events.md) を参照してください。

### Attaching Behavior With Stimulus

frame ナビゲーション、stream メッセージ、クライアントサイドのレンダリング処理などによって、新しい DOM 要素がいつページに現れてもおかしくありません。そうした要素は、多くの場合、まっさらなページロードから来たかのように初期化する必要があります。

Turbo の姉妹フレームワークである [Stimulus](https://stimulus.hotwired.dev) が提供する規約とライフサイクルコールバックを使えば、Turbo Drive のページロードによる更新も含めて、こうした更新を 1 箇所でまとめて扱えます。

Stimulus では、HTML に controller、action、target の属性を付けて注釈を付けられます。

```html
<div data-controller="hello">
  <input data-hello-target="name" type="text">
  <button data-action="click->hello#greet">Greet</button>
</div>
```

対応するコントローラーを実装しておけば、Stimulus が自動で接続してくれます。

```js
// hello_controller.js
import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
  greet() {
    console.log(`Hello, ${this.name}!`)
  }

  get name() {
    return this.targets.find("name").value
  }
}
```

Stimulus は [MutationObserver](https://developer.mozilla.org/en-US/docs/Web/API/MutationObserver) API を使って、ドキュメントが変わるたびにこれらのコントローラーと、関連するイベントハンドラーを接続・切断します。その結果、Turbo Drive のページ遷移も、Turbo Frames のナビゲーションも、Turbo Streams のメッセージも、その他あらゆる DOM 更新と同じように扱えるのです。

## Making Transformations Idempotent

サーバーから受け取った HTML に対して、クライアントサイドで変換をかけたくなることはよくあります。たとえば、ブラウザが把握しているユーザーの現在のタイムゾーンを使って、要素の集まりを日付ごとにグループ化する、といった具合です。

たとえば、要素の作成時刻を UTC で示す `data-timestamp` 属性を、要素群に付けているとしましょう。JavaScript の関数で、そうした要素をすべてドキュメントから探し出し、タイムスタンプをローカルタイムに変換し、日をまたぐ要素の前ごとに日付ヘッダーを差し込む、という処理を書いてあるとします。

この関数を `turbo:load` で走らせるように仕込んだ場合に何が起きるかを考えてみましょう。ページに遷移すると、関数が日付ヘッダーを差し込みます。ページを離れると、Turbo Drive は変換済みのページのコピーをキャッシュに保存します。次に Back ボタンを押すと、Turbo Drive がページを復元し、再び `turbo:load` を発火させ、関数が 2 つ目の日付ヘッダー群を差し込んでしまいます。

この問題を避けるためには、変換関数を*冪等*にしましょう。冪等な変換とは、初回の適用を超えて結果を変えることなく、何度でも安全に適用できる変換のことです。

変換を冪等にする 1 つのやり方は、処理済みの各要素に `data` 属性を立てておいて、既に処理したかどうかを追跡することです。Turbo Drive がキャッシュからページを復元しても、これらの属性はそのまま残っています。変換関数の中でこの属性を検出し、どの要素が既に処理済みかを判定できます。

より堅牢なやり方は、単純に変換そのものを検出することです。上の日付グルーピングの例で言えば、新しい日付ディバイダーを差し込む前に、既に日付ディバイダーが存在するかどうかをチェックする、ということです。このアプローチなら、元の変換で処理されていなかった、後から挿入された要素も自然に扱えます。

## Persisting Elements Across Page Loads

Turbo Drive では、特定の要素を*permanent (永続)* とマークできます。permanent な要素はページロードをまたいで保持されるので、それらの要素に加えた変更をナビゲーション後に再適用する必要はありません。

ショッピングカートを持つ Turbo Drive アプリケーションを考えてみてください。各ページの上部に、現在カートに入っているアイテム数を示すアイコンがあります。このカウンターは、アイテムが追加・削除されるたびに JavaScript で動的に更新されます。

さて、あるユーザーがこのアプリケーションのいくつかのページを訪れているとしましょう。彼女はアイテムを 1 つカートに追加した後、ブラウザの Back ボタンを押します。ナビゲーション時に Turbo Drive はキャッシュから前のページの状態を復元し、その結果、カートのアイテム数の表示が誤って 1 から 0 に変わってしまいます。

この問題は、カウンター要素を permanent としてマークすれば避けられます。permanent な要素にするには、HTML の `id` を付けたうえで `data-turbo-permanent` を付けます。

```html
<div id="cart-counter" data-turbo-permanent>1 item</div>
```

各レンダリングの前に、Turbo Drive はすべての permanent 要素を ID でマッチさせ、元のページから新しいページへと持ち越して、そのデータとイベントリスナーを保持します。

[次へ: Installing Turbo in Your Application](08_installing.md)
