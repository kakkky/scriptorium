> 原文: https://stimulus.hotwired.dev/handbook/working-with-external-resources
> タイトル: Working With External Resources
> 最終取得: 2026-07-19

# Working With External Resources

前章では、value を使って controller の内部状態を読み込み・永続化する方法を学びました。

controller が扱う状態は、外部リソースの状態であることもあります。ここで言う *外部* とは、DOM の中でも Stimulus の一部でもない何か、という意味です。たとえば HTTP リクエストを発行してその状態の変化に応答したい場合もあるでしょうし、タイマーを開始して、controller が接続されなくなったときに停止したい、といった場面もあるでしょう。本章では、そうした 2 つのケースをどう扱うかを見ていきます。

## Asynchronously Loading HTML

リモートの HTML 断片を非同期に読み込んで挿入することで、ページの一部を動的に埋めていく方法を学びます。私たちは Basecamp でこの手法を使って、初期ページの読み込みを速く保ちつつ、view からユーザー固有のコンテンツを切り離して、キャッシュをより効かせやすくしています。

汎用的な content-loader controller を作り、要素の中身をサーバーから取得した HTML で埋めるようにします。そのうえで、メールの受信箱のような未読メッセージ一覧をロードするのに使ってみます。

まずは受信箱のスケッチを `public/index.html` に書きます。

```html
<div data-controller="content-loader"
     data-content-loader-url-value="/messages.html"></div>
```

そして、メッセージ一覧の HTML を持つ `public/messages.html` を新規作成します。

```html
<ol>
  <li>New Message: Stimulus Launch Party</li>
  <li>Overdue: Finish Stimulus 1.0</li>
</ol>
```

(実アプリケーションではこの HTML はサーバー側で動的に生成しますが、ここではデモのため静的ファイルを使います)

これで controller を実装できます。

```js
// src/controllers/content_loader_controller.js
import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
  static values = { url: String }

  connect() {
    this.load()
  }

  load() {
    fetch(this.urlValue)
      .then(response => response.text())
      .then(html => this.element.innerHTML = html)
  }
}
```

controller が connect すると、要素の `data-content-loader-url-value` 属性に指定された URL に対して [Fetch](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch) リクエストを発火します。そして返ってきた HTML を要素の `innerHTML` プロパティに代入して読み込ませます。

ブラウザの開発者コンソールでネットワークタブを開き、ページをリロードしてください。最初のページロードを表すリクエストに続いて、controller が `messages.html` に発行するリクエストが表示されるはずです。

## Refreshing Automatically With a Timer

controller を改良して、受信箱を定期的にリフレッシュして常に最新に保つようにしてみましょう。

`data-content-loader-refresh-interval-value` 属性で、controller が中身を再読み込みする間隔をミリ秒で指定します。

```html
<div data-controller="content-loader"
     data-content-loader-url-value="/messages.html"
     data-content-loader-refresh-interval-value="5000"></div>
```

controller 側でも、この interval を確認し、指定があればリフレッシュ用のタイマーを開始するようにします。

controller に `static values` 定義を追加し、新しく `startRefreshing()` メソッドを定義します。

```js
export default class extends Controller {
  static values = { url: String, refreshInterval: Number }

  startRefreshing() {
    setInterval(() => {
      this.load()
    }, this.refreshIntervalValue)
  }

  // …
}
```

次に、`connect()` メソッドを更新して、interval value が存在すれば `startRefreshing()` を呼ぶようにします。

```js
  connect() {
    this.load()

    if (this.hasRefreshIntervalValue) {
      this.startRefreshing()
    }
  }
```

ページをリロードし、開発者コンソールで 5 秒ごとに新しいリクエストが発生することを確認しましょう。次に `public/messages.html` を書き換えて、その変更が受信箱に現れるのを待ってみてください。

## Releasing Tracked Resources

タイマーは controller が connect したときに開始していますが、停止する処理がまだありません。つまり、もし controller の要素が消えたとしても、controller はバックグラウンドで HTTP リクエストを発行し続けてしまいます。

この問題は、`startRefreshing()` メソッドを修正してタイマーへの参照を保持するようにすれば直せます。

```js
  startRefreshing() {
    this.refreshTimer = setInterval(() => {
      this.load()
    }, this.refreshIntervalValue)
  }
```

そして、タイマーをキャンセルするための `stopRefreshing()` メソッドを併せて用意します。

```js
  stopRefreshing() {
    if (this.refreshTimer) {
      clearInterval(this.refreshTimer)
    }
  }
```

最後に、controller が disconnect したときにタイマーをキャンセルするよう Stimulus に指示するため、`disconnect()` メソッドを追加します。

```js
  disconnect() {
    this.stopRefreshing()
  }
```

これで、content-loader controller は DOM に接続されているときにだけリクエストを発行する、と安心して言えるようになりました。

最終形の controller クラスを見てみましょう。

```js
// src/controllers/content_loader_controller.js
import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
  static values = { url: String, refreshInterval: Number }

  connect() {
    this.load()

    if (this.hasRefreshIntervalValue) {
      this.startRefreshing()
    }
  }

  disconnect() {
    this.stopRefreshing()
  }

  load() {
    fetch(this.urlValue)
      .then(response => response.text())
      .then(html => this.element.innerHTML = html)
  }

  startRefreshing() {
    this.refreshTimer = setInterval(() => {
      this.load()
    }, this.refreshIntervalValue)
  }

  stopRefreshing() {
    if (this.refreshTimer) {
      clearInterval(this.refreshTimer)
    }
  }
}
```

## Using action parameters

もし loader を複数の異なるソースで動かしたいのなら、action parameter を使って実現できます。次の HTML を見てみましょう。

```html
<div data-controller="content-loader">
  <a href="#" data-content-loader-url-param="/messages.html" data-action="content-loader#load">Messages</a>
  <a href="#" data-content-loader-url-param="/comments.html" data-action="content-loader#load">Comments</a>
</div>
```

そうすると、`load` action の中でこれらのパラメータを使えます。

```js
import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
  load({ params }) {
    fetch(params.url)
      .then(response => response.text())
      .then(html => this.element.innerHTML = html)
  }
}
```

さらに、params を分割代入して URL パラメータだけを取り出すこともできます。

```js
  load({ params: { url } }) {
    fetch(url)
      .then(response => response.text())
      .then(html => this.element.innerHTML = html)
  }
```

## Wrap-Up and Next Steps

本章では、Stimulus のライフサイクルコールバックを使って外部リソースを取得・解放する方法を見ました。

次章では、自分のアプリケーションに Stimulus をインストールして設定する方法を見ていきます。
