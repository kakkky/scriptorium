> 原文: https://stimulus.hotwired.dev/reference/targets
> タイトル: Targets
> 最終取得: 2026-07-19

# Targets

*Targets* を使うと、重要な要素を名前で参照できる。

```html
<div data-controller="search">
  <input type="text" data-search-target="query">
  <div data-search-target="errorMessage"></div>
  <div data-search-target="results"></div>
</div>
```

## Attributes and Names

`data-search-target` 属性は *target 属性* と呼ばれ、その値は *target 名* のスペース区切りリストである。`search` controller の中でその要素を参照するのに使える。

```html
<div data-controller="search">
  <div data-search-target="results"></div>
</div>
```

## Definitions

controller クラスの `static targets` 配列で target 名を定義する。

```js
// controllers/search_controller.js
import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
  static targets = [ "query", "errorMessage", "results" ]
  // …
}
```

## Properties

`static targets` 配列で定義した各 target 名について、Stimulus は controller に次のプロパティを追加する (`[name]` は target 名に対応)。

| Kind         | Name                    | Value                                                          |
| ------------ | ----------------------- | -------------------------------------------------------------- |
| Singular     | `this.[name]Target`     | スコープ内で最初にマッチした target                            |
| Plural       | `this.[name]Targets`    | スコープ内のマッチするすべての target の配列                   |
| Existential  | `this.has[Name]Target`  | スコープ内にマッチする target があるかどうかを示す真偽値       |

**Note:** マッチする要素が存在しない状態で singular の target プロパティにアクセスするとエラーが投げられる。

## Shared Targets

要素は複数の target 属性を持てるし、target が複数の controller で共有されるのもよくあることだ。

```html
<form data-controller="search checkbox">
  <input type="checkbox" data-search-target="projects" data-checkbox-target="input">
  <input type="checkbox" data-search-target="messages" data-checkbox-target="input">
  …
</form>
```

上の例では、`search` controller の中からはこれらのチェックボックスがそれぞれ `this.projectsTarget`、`this.messagesTarget` としてアクセスできる。

`checkbox` controller の中からは、`this.inputTargets` が両方のチェックボックスを含む配列を返す。

## Optional Targets

controller が扱う target が存在するかどうか分からない場合は、existential target プロパティの値でコードを分岐させる。

```js
if (this.hasResultsTarget) {
  this.resultsTarget.innerHTML = "…"
}
```

## Connected and Disconnected Callbacks

target の *element callback* を使うと、target 要素が controller の要素の中で追加・削除されたときに反応できる。

controller に `[name]TargetConnected` あるいは `[name]TargetDisconnected` メソッドを定義する (`[name]` は追加・削除を監視したい target の名前)。メソッドは第 1 引数として要素を受け取る。

Stimulus は各 element callback を、`connect()` の後・`disconnect()` の前に、target 要素が追加・削除されるたびに呼び出す。

```js
export default class extends Controller {
  static targets = [ "item" ]

  itemTargetConnected(element) {
    this.sortElements(this.itemTargets)
  }

  itemTargetDisconnected(element) {
    this.sortElements(this.itemTargets)
  }

  // Private
  sortElements(itemTargets) { /* ... */ }
}
```

**Note:** `[name]TargetConnected` と `[name]TargetDisconnected` のコールバック実行中は、その裏で動いている `MutationObserver` インスタンスが一時停止される。つまりコールバックの中で対応する名前の target を追加・削除しても、対応するコールバックはさらに呼び出されない。

## Naming Conventions

target 名は controller のプロパティに直接マップされるので、常に camelCase を使う。

```html
<span data-search-target="camelCase"></span>
<span data-search-target="do-not-do-this"></span>
```

```js
export default class extends Controller {
  static targets = [ "camelCase" ]  
}
```
