> 原文: https://stimulus.hotwired.dev/reference/outlets
> タイトル: Outlets
> 最終取得: 2026-07-19

# Outlets

*Outlets* を使うと、Stimulus の別の controller の中から、CSS セレクタで Stimulus の *controller インスタンス* とその *controller 要素* を参照できる。

Outlet は、controller 要素にカスタムイベントを dispatch する代わりの手段として、controller 間の通信・協調を助けてくれる。

概念的には [Stimulus Targets](targets.md) に似ているが、Outlet は controller インスタンス *と* それに紐づく controller 要素の両方を参照する、という違いがある。

```html
<div>
  <div class="online-user" data-controller="user-status">...</div>
  <div class="online-user" data-controller="user-status">...</div>
  ...
</div>

...

<div data-controller="chat" data-chat-user-status-outlet=".online-user">
  ...
</div>
```

**target** が自身の controller 要素の **スコープ内** の specifically marked な要素であるのに対して、**outlet** は **ページ上のどこにでも** 置くことができ、必ずしも controller のスコープ内にある必要はない。

## Attributes and Names

`data-chat-user-status-outlet` 属性は *outlet 属性* と呼ばれ、その値は [CSS セレクタ](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_Selectors) である。このセレクタで、*host controller* 側で outlet として利用したい別の controller 要素を参照する。host controller における outlet の identifier は、対象 controller の identifier と一致していなければならない。

```
data-[identifier]-[outlet]-outlet="[selector]"
```

```html
<div data-controller="chat" data-chat-user-status-outlet=".online-user"></div>
```

## Definitions

controller クラスの `static outlets` 配列で controller identifier を定義する。この配列は、この controller で outlet として使える他の controller identifier を宣言する。

```js
// chat_controller.js

export default class extends Controller {
  static outlets = [ "user-status" ]

  connect () {
    this.userStatusOutlets.forEach(status => ...)
  }
}
```

## Properties

`static outlets` 配列で定義した各 outlet について、Stimulus は controller に次の 5 つのプロパティを追加する (`[name]` は outlet の controller identifier に対応)。

| Kind         | Property name                | Return Type          | Effect                                                                                              |
| ------------ | ---------------------------- | -------------------- | --------------------------------------------------------------------------------------------------- |
| Existential  | `has[Name]Outlet`            | `Boolean`            | `[name]` outlet があるかどうかを検査する                                                            |
| Singular     | `[name]Outlet`               | `Controller`         | 最初の `[name]` outlet の `Controller` インスタンスを返す。存在しなければ例外を投げる               |
| Plural       | `[name]Outlets`              | `Array<Controller>`  | すべての `[name]` outlet の `Controller` インスタンスを返す                                          |
| Singular     | `[name]OutletElement`        | `Element`            | 最初の `[name]` outlet の Controller `Element` を返す。存在しなければ例外を投げる                     |
| Plural       | `[name]OutletElements`       | `Array<Element>`     | すべての `[name]` outlet の Controller `Element` を返す                                              |

**Note:** ネストした Stimulus controller のプロパティにアクセスする際、参照する outlet に正しくアクセスするために、名前空間区切り文字は省略する。

```js
// chat_controller.js

export default class extends Controller {
  static outlets = [ "admin--user-status" ]

  selectAll(event) {
    // returns undefined
    this.admin__UserStatusOutlets

    // returns controller reference
    this.adminUserStatusOutlets
  }
}
```

## Accessing Controllers and Elements

`[name]Outlet` および `[name]Outlets` プロパティから返ってくるのは `Controller` インスタンスなので、その controller インスタンスが定義する Values、Classes、Targets、その他のプロパティや関数にもアクセスできる。

```js
this.userStatusOutlet.idValue
this.userStatusOutlet.imageTarget
this.userStatusOutlet.activeClasses
```

outlet controller が定義するどんな関数も呼び出せる。

```js
// user_status_controller.js

export default class extends Controller {
  markAsSelected(event) {
    // ...
  }
}

// chat_controller.js

export default class extends Controller {
  static outlets = [ "user-status" ]

  selectAll(event) {
    this.userStatusOutlets.forEach(status => status.markAsSelected(event))
  }
}
```

Outlet Element も同様で、[`Element`](https://developer.mozilla.org/en-US/docs/Web/API/Element) の任意の関数・プロパティを呼び出せる。

```js
this.userStatusOutletElement.dataset.value
this.userStatusOutletElement.getAttribute("id")
this.userStatusOutletElements.map(status => status.hasAttribute("selected"))
```

## Outlet Callbacks

Outlet コールバックは特別な名前を持った関数で、outlet がページに追加・削除されたときに Stimulus によって呼び出される。

outlet の変化を監視するには、`[name]OutletConnected()` あるいは `[name]OutletDisconnected()` という関数を定義する。

```js
// chat_controller.js

export default class extends Controller {
  static outlets = [ "user-status" ]

  userStatusOutletConnected(outlet, element) {
    // ...
  }

  userStatusOutletDisconnected(outlet, element) {
    // ...
  }
}
```

### Outlets are Assumed to be Present

controller の中で Outlet プロパティにアクセスするということは、対応する Outlet が少なくとも 1 つ存在する、と主張していることになる。宣言に対応する outlet が見つからない場合、Stimulus は例外を投げる。

```
Missing outlet element "user-status" for "chat" controller
```

### Optional outlets

Outlet がオプショナルな場合、あるいは少なくとも 1 つ存在することを担保したい場合は、existential プロパティで存在チェックしてから利用する。

```js
if (this.hasUserStatusOutlet) {
  this.userStatusOutlet.safelyCallSomethingOnTheOutlet()
}
```

### Referencing Non-Controller Elements

対応する `data-controller` と identifier が付いていない要素を outlet として宣言しようとすると、Stimulus は例外を投げる。

```html
<div data-controller="chat" data-chat-user-status-outlet="#user-column"></div>

<div id="user-column"></div>
```

この場合、次のような結果になる。

```
Missing "data-controller=user-status" attribute on outlet element for
"chat" controller`
```
