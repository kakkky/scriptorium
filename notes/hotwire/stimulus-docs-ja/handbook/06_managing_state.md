> 原文: https://stimulus.hotwired.dev/handbook/managing-state
> タイトル: Managing State
> 最終取得: 2026-07-19

# Managing State

現代の多くのフレームワークは、状態を常に JavaScript の中に保持することを推奨します。そうしたフレームワークでは、DOM は書き込み専用の描画ターゲットとして扱われ、サーバーから送られてくる JSON を消費するクライアントサイドのテンプレートによって整合が取られます。

Stimulus は別のアプローチを取ります。Stimulus アプリケーションの状態は DOM の属性として存在し、controller 自身はほとんど無状態です。このアプローチのおかげで、初期ドキュメント、Ajax レスポンス、Turbo の visit、あるいは別の JavaScript ライブラリのいずれから届いた HTML であっても扱えて、関連する controller は明示的な初期化ステップを踏まずに自動的に息を吹き返せます。

## Building a Slideshow

前章では、要素にクラス名を追加することで、Stimulus controller が単純な状態をドキュメントの中に保持できることを学びました。しかし単なるフラグではなく、値を保存する必要が出てきたときはどうすればよいでしょうか?

この疑問を、現在選択中のスライドのインデックスを属性に保持するスライドショーの controller を作って探っていきます。

いつも通り HTML から始めます。

```html
<div data-controller="slideshow">
  <button data-action="slideshow#previous"> ← </button>
  <button data-action="slideshow#next"> → </button>

  <div data-slideshow-target="slide">🐵</div>
  <div data-slideshow-target="slide">🙈</div>
  <div data-slideshow-target="slide">🙉</div>
  <div data-slideshow-target="slide">🙊</div>
</div>
```

各 `slide` target はスライドショーの 1 枚のスライドを表します。controller は同時に見えるスライドを 1 枚だけに保つ責務を負います。

controller のドラフトを書いてみましょう。新しいファイル `src/controllers/slideshow_controller.js` を次のように作ります。

```js
// src/controllers/slideshow_controller.js
import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
  static targets = [ "slide" ]

  initialize() {
    this.index = 0
    this.showCurrentSlide()
  }

  next() {
    this.index++
    this.showCurrentSlide()
  }

  previous() {
    this.index--
    this.showCurrentSlide()
  }

  showCurrentSlide() {
    this.slideTargets.forEach((element, index) => {
      element.hidden = index !== this.index
    })
  }
}
```

この controller は `showCurrentSlide()` メソッドを定義し、各 slide target をループしながら、インデックスが一致するかどうかで [`hidden` 属性](https://developer.mozilla.org/en-US/docs/Web/HTML/Global_attributes/hidden) をトグルします。

初期化時に最初のスライドを表示し、`next()` と `previous()` の action メソッドで現在のスライドを進めたり戻したりします。

> ### Lifecycle Callbacks Explained
>
> `initialize()` メソッドは何をしているのでしょう? これまで使ってきた `connect()` メソッドとは何が違うのでしょうか?
>
> これらは Stimulus の *ライフサイクルコールバック* メソッドで、controller がドキュメントに入るときや離れるときに、関連する状態のセットアップ・ティアダウンをするのに便利です。
>
> | Method         | Invoked by Stimulus…                                |
> | -------------- | --------------------------------------------------- |
> | initialize()   | controller が最初にインスタンス化されたときに 1 回だけ |
> | connect()      | controller が DOM に接続されるたびに                |
> | disconnect()   | controller が DOM から切り離されるたびに            |

ページをリロードし、Next ボタンで次のスライドに進めることを確認してください。

## Reading Initial State from the DOM

私たちの controller が現在のスライド (状態) を `this.index` プロパティで追跡していることに注目してください。

さて、あるスライドショーを 1 枚目ではなく 2 枚目から始めたい、としましょう。初期インデックスをマークアップにどうエンコードすればよいでしょうか?

一つのやり方は、HTML の `data` 属性で初期インデックスを持たせることです。たとえば、controller の要素に `data-index` 属性を追加できます。

```html
<div data-controller="slideshow" data-index="1">
```

そして `initialize()` メソッドの中でその属性を読み取り、整数に変換して `showCurrentSlide()` に渡します。

```js
  initialize() {
    this.index = Number(this.element.dataset.index)
    this.showCurrentSlide()
  }
```

これで動くには動きますが、ぎこちないですし、属性の名前をどうするか自分で決めないといけないし、インデックスを再度アクセスしたい・インクリメントして結果を DOM に永続化したい、といった要件に対しても助けにはなりません。

### Using Values

Stimulus controller は、data 属性に自動でマップされる型付きの value プロパティをサポートします。controller クラスの上部に value 定義を追加すると、

```js
  static values = { index: Number }
```

Stimulus は `data-slideshow-index-value` 属性に紐づく `this.indexValue` プロパティを controller に作り出し、数値への変換もこちらで行ってくれます。

実際にやってみましょう。まずは対応する data 属性を HTML に追加します。

```html
<div data-controller="slideshow" data-slideshow-index-value="1">
```

そして controller に `static values` 定義を追加し、`initialize()` メソッドで `this.indexValue` をログ出力するようにします。

```js
export default class extends Controller {
  static values = { index: Number }

  initialize() {
    console.log(this.indexValue)
    console.log(typeof this.indexValue)
  }

  // …
}
```

ページをリロードして、コンソールに `1` と `Number` が出ることを確かめてください。

> ### What's with that `static values` line?
>
> target と同様に、Stimulus controller の value も `values` という名前の静的オブジェクトプロパティで記述して定義します。ここでは `index` という 1 つの数値 value を定義しました。value 定義について詳しくは [reference documentation](../reference/values.md) を参照してください。

続いて、controller の `initialize()` とその他のメソッドを、`this.index` の代わりに `this.indexValue` を使うように書き換えます。書き換え後の controller は次のようになります。

```js
import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
  static targets = [ "slide" ]
  static values = { index: Number }

  initialize() {
    this.showCurrentSlide()
  }

  next() {
    this.indexValue++
    this.showCurrentSlide()
  }

  previous() {
    this.indexValue--
    this.showCurrentSlide()
  }

  showCurrentSlide() {
    this.slideTargets.forEach((element, index) => {
      element.hidden = index !== this.indexValue
    })
  }
}
```

ページをリロードして、スライドを移動するたびに controller 要素の `data-slideshow-index-value` 属性が変化することを Web インスペクタで確認してください。

### Change Callbacks

書き換えた controller はオリジナルより良くなりましたが、`this.showCurrentSlide()` を繰り返し呼び出している部分が目立ちます。controller の初期化時にも、`this.indexValue` を書き換えるすべての箇所でも、手動でドキュメントの状態を更新しなければなりません。

Stimulus の value change コールバックを定義すると、この繰り返しを片付けられ、index value が変化したときに controller がどう反応するかを明示的に指定できます。

まず `initialize()` メソッドを削除し、新しく `indexValueChanged()` メソッドを定義します。そして `next()` と `previous()` から `this.showCurrentSlide()` の呼び出しを取り除きます。

```js
import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
  static targets = [ "slide" ]
  static values = { index: Number }

  next() {
    this.indexValue++
  }

  previous() {
    this.indexValue--
  }

  indexValueChanged() {
    this.showCurrentSlide()
  }

  showCurrentSlide() {
    this.slideTargets.forEach((element, index) => {
      element.hidden = index !== this.indexValue
    })
  }
}
```

ページをリロードして、スライドショーのふるまいが変わっていないことを確認してください。

Stimulus は `indexValueChanged()` メソッドを初期化時に呼び出すのに加え、`data-slideshow-index-value` 属性への任意の変化にも反応して呼び出します。Web インスペクタでこの属性の値をいじると、それに応じて controller がスライドを切り替えるところを見ることさえできます。ぜひ試してみてください!

### Setting Defaults

静的定義の一部としてデフォルト値を指定することもできます。書き方はこうです。

```js
  static values = { index: { type: Number, default: 2 } }
```

controller の要素に `data-slideshow-index-value` 属性が定義されていない場合、この設定によりインデックスは 2 から始まります。他の value がある場合は、デフォルトを持つものと持たないものを混在させて記述できます。

```js
  static values = { index: Number, effect: { type: String, default: "kenburns" } }
```

## Wrap-Up and Next Steps

本章では、スライドショーの controller で現在のインデックスの読み込みと永続化に value を使う方法を見ました。

使いやすさという点では、私たちの controller はまだ不完全です。最初のスライドを見ているときに Previous ボタンを押しても、何も起きないように見えます。内部では `indexValue` が `0` から `-1` になっているだけです。この値がラップアラウンドして *最後* のスライドインデックスへ回り込むようにしてもよいかもしれません。(Next ボタンについても同種の問題があります)

次章では、タイマーや HTTP リクエストといった外部リソースを Stimulus controller の中でどう追跡するかを見ていきます。
