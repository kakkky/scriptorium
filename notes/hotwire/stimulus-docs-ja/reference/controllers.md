> 原文: https://stimulus.hotwired.dev/reference/controllers
> タイトル: Controllers
> 最終取得: 2026-07-19

# Controllers

*controller* は、Stimulus アプリケーションにおける組織化の基本単位である。

```js
import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
  // …
}
```

controller は、アプリケーション側で定義した JavaScript クラスのインスタンスである。各 controller クラスは、`@hotwired/stimulus` モジュールが export する `Controller` 基底クラスを継承する。

## Properties

すべての controller は、いずれかの Stimulus `Application` インスタンスに属し、ある HTML 要素に紐づいている。controller クラスの中からは、次の情報にアクセスできる。

- application: `this.application` プロパティ経由
- HTML 要素: `this.element` プロパティ経由
- identifier: `this.identifier` プロパティ経由

## Modules

controller クラスは JavaScript モジュール内に、1 ファイル 1 クラスで定義する。各 controller クラスは、上記の例のように、モジュールのデフォルトオブジェクトとして export する。

これらのモジュールは `controllers/` ディレクトリに置く。ファイル名は `[identifier]_controller.js` とし、`[identifier]` は各 controller の identifier に対応させる。

## Identifiers

*identifier* は、HTML の中で controller クラスを参照するために使う名前である。

要素に `data-controller` 属性を追加すると、Stimulus はその属性値から identifier を読み取り、対応する controller クラスの新しいインスタンスを生成する。

たとえば、次の要素は `controllers/reference_controller.js` に定義されたクラスのインスタンスの controller を持つ。

```html
<div data-controller="reference"></div>
```

以下は、Stimulus が require コンテキストの controller にどのように identifier を生成するかの例である。

| If your controller file is named…     | its identifier will be… |
| ------------------------------------- | ----------------------- |
| clipboard_controller.js               | clipboard               |
| date_picker_controller.js             | date-picker             |
| users/list_item_controller.js         | users--list-item        |
| local-time-controller.js              | local-time              |

## Scopes

Stimulus が controller を要素に接続すると、その要素とすべての子孫要素がその controller の *scope* を構成する。

たとえば、次の `<div>` と `<h1>` は controller のスコープに含まれるが、それらを取り囲む `<main>` 要素はスコープには含まれない。

```html
<main>
  <div data-controller="reference">
    <h1>Reference</h1>
  </div>
</main>
```

## Nested Scopes

ネストしている場合、各 controller は自身のスコープしか意識しない。ネストした内側の controller のスコープは除外される。

たとえば、下の `#parent` の controller は、自身のスコープ内の直接の `item` target しか認識せず、`#child` の controller の target には触れない。

```html
<ul id="parent" data-controller="list">
  <li data-list-target="item">One</li>
  <li data-list-target="item">Two</li>
  <li>
    <ul id="child" data-controller="list">
      <li data-list-target="item">I am</li>
      <li data-list-target="item">a nested list</li>
    </ul>
  </li>
</ul>
```

## Multiple Controllers

`data-controller` 属性の値は、スペース区切りの identifier のリストである。

```html
<div data-controller="clipboard list-item"></div>
```

ページ上の 1 つの要素に複数の controller が乗ることは珍しくない。上の例では、`<div>` に `clipboard` と `list-item` の 2 つの controller が接続されている。

同様に、ページ上の複数の要素が同じ controller クラスを参照することも珍しくない。

```html
<ul>
  <li data-controller="list-item">One</li>
  <li data-controller="list-item">Two</li>
  <li data-controller="list-item">Three</li>
</ul>
```

ここでは、各 `<li>` が `list-item` controller の独自のインスタンスを持つ。

## Naming Conventions

controller クラスのメソッド名・プロパティ名には常に camelCase を使うこと。

identifier が複数の単語から成る場合、単語同士は kebab-case (ダッシュ区切り、たとえば `date-picker`、`list-item`) で記述する。

ファイル名で複数の単語を区切るときは、アンダースコアまたはダッシュ (snake_case または kebab-case、たとえば `controllers/date_picker_controller.js`、`controllers/list-item-controller.js`) を使うこと。

## Registration

Stimulus for Rails を import map と組み合わせて使うか、`@hotwired/stimulus-webpack-helpers` パッケージを Webpack と組み合わせて使う場合、アプリケーションは上記の規約に従って controller クラスを自動的にロード・登録する。

そうでない場合は、アプリケーション側で各 controller クラスを手動でロード・登録する必要がある。

### Registering Controllers Manually

controller クラスを identifier とともに手動で登録するには、まずクラスを import し、その後アプリケーションオブジェクトの `Application#register` メソッドを呼び出す。

```js
import ReferenceController from "./controllers/reference_controller"

application.register("reference", ReferenceController)
```

モジュールから import する代わりに、controller クラスをインラインで登録することもできる。

```js
import { Controller } from "@hotwired/stimulus"

application.register("reference", class extends Controller {
  // …
})
```

### Preventing Registration Based On Environmental Factors

特定の環境要因 (たとえば特定のユーザーエージェント) が満たされた場合にだけ controller を登録・ロードしたい場合は、静的な `shouldLoad` メソッドを上書きすればよい。

```js
class UnloadableController extends ApplicationController {
  static get shouldLoad() {
    return false
  }
}

// This controller will not be loaded
application.register("unloadable", UnloadableController)
```

### Trigger Behaviour When A Controller Is Registered

controller の登録が完了した時点で何らかの振る舞いを走らせたい場合は、静的な `afterLoad` メソッドを追加できる。

```js
class SpinnerButton extends Controller {
  static legacySelector = ".legacy-spinner-button"

  static afterLoad(identifier, application) {
    // use the application instance to read the configured 'data-controller' attribute
    const { controllerAttribute } = application.schema

    // update any legacy buttons with the controller's registered identifier
    const updateLegacySpinners = () => {
      document.querySelector(this.legacySelector).forEach((element) => {
        element.setAttribute(controllerAttribute, identifier)
      })
    }

    // called as soon as registered so DOM may not have loaded yet
    if (document.readyState == "loading") {
      document.addEventListener("DOMContentLoaded", updateLegacySpinners)
    } else {
      updateLegacySpinners()
    }
  }
}

// This controller will update any legacy spinner buttons to use the controller
application.register("spinner-button", SpinnerButton)
```

`afterLoad` メソッドは、コントロールする対象要素が DOM 上に存在しない場合であっても、controller が登録された直後に呼び出される。この関数はもとの controller コンストラクタに `this` を束縛された状態で、2 つの引数—`identifier` (controller の登録に使われた identifier) と Stimulus のアプリケーションインスタンス—とともに呼び出される。

## Cross-Controller Coordination With Events

controller 同士が通信する必要がある場合はイベントを使うべきである。`Controller` クラスにはこれを容易にする `dispatch` という便利メソッドがある。第 1 引数の `eventName` は、コロンを挟んで controller 名を自動的にプレフィックスとして付ける。ペイロードは `detail` に保持される。使い方は次のとおり。

```js
class ClipboardController extends Controller {
  static targets = [ "source" ]

  copy() {
    this.dispatch("copy", { detail: { content: this.sourceTarget.value } })
    navigator.clipboard.writeText(this.sourceTarget.value)
  }
}
```

このイベントは、別の controller の action にルーティングできる。

```html
<div data-controller="clipboard effects" data-action="clipboard:copy->effects#flash">
  PIN: <input data-clipboard-target="source" type="text" value="1234" readonly>
  <button data-action="clipboard#copy">Copy to Clipboard</button>
</div>
```

`Clipboard#copy` action が呼び出されると、`Effects#flash` action も呼び出される。

```js
class EffectsController extends Controller {
  flash({ detail: { content } }) {
    console.log(content) // 1234
  }
}
```

2 つの controller が同じ HTML 要素に属していない場合は、`data-action` 属性を *受信側の* controller の要素に追加する必要がある。受信側 controller の要素が送信側 controller の要素の親 (あるいは同じ要素) でない場合は、イベントに `@window` を付ける必要がある。

```html
<div data-action="clipboard:copy@window->effects#flash">
```

`dispatch` は、第 2 引数として次の追加オプションを受け付ける。

| option       | default                | notes                                                                                                                                                                                       |
| ------------ | ---------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `detail`     | `{}` 空オブジェクト     | [CustomEvent.detail](https://developer.mozilla.org/en-US/docs/Web/API/CustomEvent/detail) を参照                                                                                              |
| `target`     | `this.element`         | [Event.target](https://developer.mozilla.org/en-US/docs/Web/API/Event/target) を参照                                                                                                          |
| `prefix`     | `this.identifier`      | prefix が falsy (たとえば `null` や `false`) の場合、`eventName` のみが使われる。文字列を指定した場合は、その文字列とコロンが `eventName` の前に付与される                                       |
| `bubbles`    | `true`                 | [Event.bubbles](https://developer.mozilla.org/en-US/docs/Web/API/Event/bubbles) を参照                                                                                                        |
| `cancelable` | `true`                 | [Event.cancelable](https://developer.mozilla.org/en-US/docs/Web/API/Event/cancelable) を参照                                                                                                  |

`dispatch` は生成した [`CustomEvent`](https://developer.mozilla.org/en-US/docs/Web/API/CustomEvent) を返す。これを使うことで、他のリスナーがイベントをキャンセルできるようにできる。

```js
class ClipboardController extends Controller {
  static targets = [ "source" ]

  copy() {
    const event = this.dispatch("copy", { cancelable: true })
    if (event.defaultPrevented) return
    navigator.clipboard.writeText(this.sourceTarget.value)
  }
}
```

```js
class EffectsController extends Controller {
  flash(event) {
    // this will prevent the default behaviour as determined by the dispatched event
    event.preventDefault()
  }
}
```

## Directly Invoking Other Controllers

なんらかの理由で controller 間の通信にイベントを使えない場合、application の `getControllerForElementAndIdentifier` メソッドから controller インスタンスに直接到達できる。この方法は、より一般的なイベントの利用では解決できない特殊な問題を抱えている場合にのみ用いるべきである。それでも必要なら、次のようにする。

```js
class MyController extends Controller {
  static targets = [ "other" ]

  copy() {
    const otherController = this.application.getControllerForElementAndIdentifier(this.otherTarget, 'other')
    otherController.otherMethod()
  }
}
```
