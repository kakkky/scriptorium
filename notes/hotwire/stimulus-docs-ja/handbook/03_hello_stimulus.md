> 原文: https://stimulus.hotwired.dev/handbook/hello-stimulus
> タイトル: Hello, Stimulus
> 最終取得: 2026-07-19

# Hello, Stimulus

Stimulus の仕組みを学ぶ最良の方法は、シンプルな controller を実際に作ってみることです。本章ではその方法を紹介します。

## Prerequisites

これから手を動かしながら読んでいくには、[`stimulus-starter`](https://github.com/hotwired/stimulus-starter) プロジェクトを動作させておく必要があります。これは Stimulus を試すために事前に設定された空のスターターです。

何もインストールせずにブラウザだけで作業できるように、[`stimulus-starter` を Glitch 上で remix する](https://glitch.com/edit/#!/import/git?url=https://github.com/hotwired/stimulus-starter.git)のがおすすめです。

[![Remix on Glitch](https://cdn.glitch.com/2703baf2-b643-4da7-ab91-7ee2a2d00b5b%2Fremix-button.svg)](https://glitch.com/edit/#!/import/git?url=https://github.com/hotwired/stimulus-starter.git)

自分のテキストエディタで作業したい場合は、`stimulus-starter` をクローンしてセットアップします。

```
$ git clone https://github.com/hotwired/stimulus-starter.git
$ cd stimulus-starter
$ yarn install
$ yarn start
```

そのあとブラウザで http://localhost:9000/ を開きます。

(`stimulus-starter` プロジェクトは依存関係の管理に [Yarn パッケージマネージャ](https://yarnpkg.com/) を使っているので、事前に yarn をインストールしておいてください)

## It All Starts With HTML

まずはテキストフィールドとボタンを使ったシンプルな課題から始めましょう。ボタンをクリックしたら、テキストフィールドの値をコンソールに表示するようにします。

すべての Stimulus プロジェクトは HTML から始まります。`public/index.html` を開き、`<body>` タグの開きすぐ後に次のマークアップを追加します。

```html
<div>
  <input type="text">
  <button>Greet</button>
</div>
```

ブラウザでページをリロードすると、テキストフィールドとボタンが表示されるはずです。

## Controllers Bring HTML to Life

Stimulus の核となる目的は、DOM 要素と JavaScript オブジェクトを自動的に接続することです。それらのオブジェクトを *controller* と呼びます。

最初の controller を、フレームワーク組み込みの `Controller` クラスを継承して作りましょう。`src/controllers/` フォルダに `hello_controller.js` という名前のファイルを新規作成し、次のコードを書き込みます。

```js
// src/controllers/hello_controller.js
import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
}
```

## Identifiers Link Controllers With the DOM

続いて、この controller が HTML とどう結び付けられるべきかを Stimulus に伝える必要があります。そのために、`<div>` の `data-controller` 属性に *identifier* (識別子) を置きます。

```html
<div data-controller="hello">
  <input type="text">
  <button>Greet</button>
</div>
```

identifier は要素と controller をつなぐリンクの役割を果たします。この例では、identifier `hello` によって、`hello_controller.js` の controller クラスのインスタンスを Stimulus が生成するように指示しています。controller が自動的にロードされる仕組みについては [Installation Guide](08_installing.md) で詳しく学べます。

## Is This Thing On?

ブラウザでページをリロードしても、何も変わっていないように見えます。controller がちゃんと動いているかどうか、どうやって確かめればよいでしょうか?

一つの方法は、`connect()` メソッドにログ出力を仕込むことです。Stimulus は controller がドキュメントに接続されるたびにこのメソッドを呼び出します。

`hello_controller.js` に次のように `connect()` メソッドを実装します。

```js
// src/controllers/hello_controller.js
import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
  connect() {
    console.log("Hello, Stimulus!", this.element)
  }
}
```

再度ページをリロードし、開発者コンソールを開いてみてください。`Hello, Stimulus!` に続いて、先ほどの `<div>` の表現が表示されるはずです。

## Actions Respond to DOM Events

次に、コードを書き換えて、「Greet」ボタンをクリックしたときにログメッセージが出るようにしてみましょう。

まず `connect()` を `greet()` にリネームします。

```js
// src/controllers/hello_controller.js
import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
  greet() {
    console.log("Hello, Stimulus!", this.element)
  }
}
```

私たちがやりたいのは、ボタンの `click` イベントが発火したときに `greet()` メソッドを呼び出すことです。Stimulus では、イベントを扱う controller のメソッドを *action メソッド* と呼びます。

action メソッドをボタンの `click` イベントに接続するため、`public/index.html` を開き、ボタンに `data-action` 属性を追加します。

```html
<div data-controller="hello">
  <input type="text">
  <button data-action="click->hello#greet">Greet</button>
</div>
```

> ### Action Descriptors Explained
>
> `data-action` の値 `click->hello#greet` は *action descriptor* と呼ばれます。この記述子はそれぞれ次のような意味を持ちます。
>
> - `click` はイベント名
> - `hello` は controller の identifier
> - `greet` は呼び出すメソッドの名前

ページをブラウザで読み込み、開発者コンソールを開きます。「Greet」ボタンをクリックすると、ログメッセージが表示されるはずです。

## Targets Map Important Elements To Controller Properties

演習を締めくくるため、テキストフィールドに入力された名前を使って挨拶をするように、action を変更していきます。

そのためには、まず controller の中から input 要素への参照を得る必要があります。参照が取れれば、`value` プロパティを読んでその内容を取得できます。

Stimulus では、重要な要素を *target* としてマークすることで、controller 側で対応するプロパティを通じて簡単に参照できるようにできます。`public/index.html` を開き、input 要素に `data-hello-target` 属性を追加します。

```html
<div data-controller="hello">
  <input data-hello-target="name" type="text">
  <button data-action="click->hello#greet">Greet</button>
</div>
```

続いて、controller の target 定義リストに `name` を追加して、target 用のプロパティを作ります。Stimulus は自動的に `this.nameTarget` プロパティを用意してくれます。これは最初にマッチした target 要素を返します。このプロパティを使って要素の `value` を読み取り、挨拶文を組み立てます。

やってみましょう。`hello_controller.js` を次のように更新します。

```js
// src/controllers/hello_controller.js
import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
  static targets = [ "name" ]

  greet() {
    const element = this.nameTarget
    const name = element.value
    console.log(`Hello, ${name}!`)
  }
}
```

ブラウザでページをリロードし、開発者コンソールを開きます。input フィールドに自分の名前を入力して「Greet」ボタンをクリックしてみてください。Hello, world!

## Controllers Simplify Refactoring

Stimulus の controller は、そのメソッドがイベントハンドラとして働ける JavaScript クラスのインスタンスだ、ということが分かりました。

つまり、標準的なリファクタリング技法を武器としてそのまま使えるということです。たとえば `greet()` メソッドをスッキリさせるため、`name` ゲッタを切り出してみましょう。

```js
// src/controllers/hello_controller.js
import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
  static targets = [ "name" ]

  greet() {
    console.log(`Hello, ${this.name}!`)
  }

  get name() {
    return this.nameTarget.value
  }
}
```

## Wrap-Up and Next Steps

おめでとうございます — 最初の Stimulus controller を書き上げました!

ここまでで、フレームワークの最重要概念である controller、action、target を扱いました。次章では、Basecamp から抜き出してきた実際の controller を組み立てるために、これらをどう組み合わせていくかを見ていきます。
