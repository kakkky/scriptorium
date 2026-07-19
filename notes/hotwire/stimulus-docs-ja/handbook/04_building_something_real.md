> 原文: https://stimulus.hotwired.dev/handbook/building-something-real
> タイトル: Building Something Real
> 最終取得: 2026-07-19

# Building Something Real

最初の controller を実装し、Stimulus が HTML と JavaScript をどうつなぐかを学びました。次は、Basecamp から抜き出した controller を再現しながら、実アプリケーションでも使えるものを見ていきましょう。

## Wrapping the DOM Clipboard API

Basecamp の UI のあちこちに、次のようなボタンが散りばめられています。

![Screenshot showing a text field with an email address inside and a ”Copy to clipboard“ button to the right](https://stimulus.hotwired.dev/assets/bc3-clipboard-ui.png)

これらのボタンをクリックすると、Basecamp は URL やメールアドレスといったちょっとしたテキストをクリップボードにコピーしてくれます。

Web プラットフォームには [システムクリップボードにアクセスするための API](https://www.w3.org/TR/clipboard-apis/) が用意されていますが、私たちがやりたいことを実現してくれる HTML 要素はありません。「Copy to clipboard」ボタンを実装するには、JavaScript の力が必要です。

## Implementing a Copy Button

たとえば、他の人にアクセス権を与えるために PIN を発行できるアプリを考えてみましょう。発行した PIN と一緒に、クリップボードにコピーするボタンを並べて表示できたら、共有もぐっと楽になります。

`public/index.html` を開き、`<body>` の中身をボタンのラフスケッチで置き換えます。

```html
<div>
  PIN: <input type="text" value="1234" readonly>
  <button>Copy to Clipboard</button>
</div>
```

## Setting Up the Controller

次に、`src/controllers/clipboard_controller.js` を作成し、空の `copy()` メソッドを持たせます。

```js
// src/controllers/clipboard_controller.js
import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
  copy() {
  }
}
```

そして外側の `<div>` に `data-controller="clipboard"` を追加します。この属性が要素に現れると、Stimulus はいつでも私たちの controller のインスタンスを結び付けてくれます。

```html
<div data-controller="clipboard">
```

## Defining the Target

クリップボード API を呼び出す前に、テキストフィールドの中身を選択するためにテキストフィールドへの参照が必要です。テキストフィールドに `data-clipboard-target="source"` を追加します。

```html
  PIN: <input data-clipboard-target="source" type="text" value="1234" readonly>
```

続いて controller に target 定義を追加し、`this.sourceTarget` としてテキストフィールド要素にアクセスできるようにします。

```js
export default class extends Controller {
  static targets = [ "source" ]

  // ...
}
```

> ### What's With That `static targets` Line?
>
> Stimulus が controller クラスを読み込むとき、`targets` という名前の静的配列に含まれる target 名の文字列を探します。配列内の各 target 名について、Stimulus は controller に 3 つのプロパティを追加します。ここでは `"source"` という target 名が次のプロパティになります。
>
> - `this.sourceTarget` は controller のスコープ内で最初に見つかった `source` target を返します。`source` target が存在しない場合、このプロパティへのアクセスはエラーを投げます。
> - `this.sourceTargets` は controller のスコープ内のすべての `source` target の配列を返します。
> - `this.hasSourceTarget` は `source` target が存在する場合は `true` を、存在しない場合は `false` を返します。
>
> target についてもっと詳しくは [reference documentation](../reference/targets.md) を参照してください。

## Connecting the Action

これで Copy ボタンをつなぐ準備ができました。

ボタンをクリックしたら controller の `copy()` メソッドを呼び出したいので、`data-action="clipboard#copy"` を追加します。

```html
  <button data-action="clipboard#copy">Copy to Clipboard</button>
```

> ### Common Events Have a Shorthand Action Notation
>
> お気づきかもしれませんが、action descriptor から `click->` を省略しています。これは Stimulus が `<button>` 要素の action におけるデフォルトイベントとして `click` を定義しているためです。
>
> 他の一部の要素にもデフォルトイベントがあります。全リストは次の通りです。
>
> | Element           | Default Event |
> | ----------------- | ------------- |
> | a                 | click         |
> | button            | click         |
> | details           | toggle        |
> | form              | submit        |
> | input             | input         |
> | input type=submit | click         |
> | select            | change        |
> | textarea          | input         |

最後に、`copy()` メソッドの中で input フィールドの内容を選択し、クリップボード API を呼び出します。

```js
  copy() {
    navigator.clipboard.writeText(this.sourceTarget.value)
  }
```

ブラウザでページを読み込み、Copy ボタンをクリックしてください。テキストエディタに戻ってペーストすると、PIN `1234` が貼り付けられているはずです。

## Stimulus Controllers are Reusable

ここまでは、ページ上に controller のインスタンスが 1 つだけあるケースを見てきました。

とはいえ、1 つのページに同じ controller のインスタンスが複数ある、というのはよくある話です。たとえば PIN のリストを表示していて、各行にそれぞれ Copy ボタンが付いている、といった具合です。

私たちの controller は再利用可能です。テキストをクリップボードにコピーする手段を提供したいときはいつでも、適切な注釈の付いたマークアップをページに置くだけで済みます。

もう 1 つ PIN をページに追加してみましょう。`<div>` をコピー&ペーストして PIN フィールドが 2 つ並ぶようにし、2 つ目の `value` 属性を書き換えます。

```html
<div data-controller="clipboard">
  PIN: <input data-clipboard-target="source" type="text" value="3737" readonly>
  <button data-action="clipboard#copy">Copy to Clipboard</button>
</div>
```

ページをリロードして、両方のボタンが動作することを確認しましょう。

## Actions and Targets Can Go on Any Kind of Element

さらにもう 1 つ PIN フィールドを追加してみましょう。今度はボタンではなく Copy *リンク* を使います。

```html
<div data-controller="clipboard">
  PIN: <input data-clipboard-target="source" type="text" value="3737" readonly>
  <a href="#" data-action="clipboard#copy">Copy to Clipboard</a>
</div>
```

Stimulus は、適切な `data-action` 属性がついている限り、どんな種類の要素でも使わせてくれます。

なお、この場合はリンクをクリックすると、ブラウザは `href` に従ってリンク先へ遷移しようともします。この既定の挙動は、action の中で `event.preventDefault()` を呼び出して打ち消せます。

```js
  copy(event) {
    event.preventDefault()
    navigator.clipboard.writeText(this.sourceTarget.value)
  }
```

同様に、`source` target は `<input type="text">` である必要はありません。controller が期待するのは、その要素に `value` プロパティと `copy()` メソッド (原文ママ) があること、それだけです。つまり `<textarea>` を代わりに使うこともできます。

```html
  PIN: <textarea data-clipboard-target="source" readonly>3737</textarea>
```

## Wrap-Up and Next Steps

本章では、ブラウザ API を Stimulus controller でラップする実例を見ました。同じ controller のインスタンスがページ上に複数同時に登場できること、そして action と target のおかげで HTML と JavaScript が緩やかに結び付けられていることを確認できました。

次は、controller の設計に小さな変更を加えるだけで、より堅牢な実装へと導けることを見ていきましょう。
