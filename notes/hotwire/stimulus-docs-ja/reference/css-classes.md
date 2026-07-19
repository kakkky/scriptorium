> 原文: https://stimulus.hotwired.dev/reference/css-classes
> タイトル: CSS Classes
> 最終取得: 2026-07-19

# CSS Classes

HTML において、*CSS class* は `class` 属性を使って要素に適用できるスタイルの集合を定義する。

CSS class は、プログラム的にスタイルを切り替えたりアニメーションを再生したりするための便利な道具である。たとえば、Stimulus の controller は、バックグラウンドで処理を実行しているあいだ要素に "loading" クラスを追加し、CSS 側でそのクラスにプログレスインジケータを表示するスタイルを当てられる。

```html
<form data-controller="search" class="search--busy">
```

```css
.search--busy {
  background-image: url(throbber.svg) no-repeat;
}
```

JavaScript の文字列でクラス名をハードコードする代わりに、Stimulus では data 属性と controller のプロパティを組み合わせて、CSS class を *論理名* で参照できる。

## Definitions

controller の `static classes` 配列で、CSS class を論理名で定義する。

```js
// controllers/search_controller.js
import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
  static classes = [ "loading" ]

  // …
}
```

## Attributes

controller の `static classes` 配列で定義した論理名は、controller 要素上の *CSS class 属性* に対応する。

```html
<form data-controller="search"
      data-search-loading-class="search--busy">
  <input data-action="search#loadResults">
</form>
```

CSS class 属性は、controller の identifier と論理名を組み合わせて `data-[identifier]-[logical-name]-class` の形式で構築する。属性値は 1 つの CSS クラス名でも、複数のクラス名のリストでもよい。

**Note:** CSS class 属性は、`data-controller` 属性と同じ要素に指定しなければならない。

1 つの論理名に複数の CSS クラスを指定したい場合は、クラス同士をスペースで区切る。

```html
<form data-controller="search"
      data-search-loading-class="bg-gray-500 animate-spinner cursor-busy">
  <input data-action="search#loadResults">
</form>
```

## Properties

`static classes` 配列で定義した各論理名について、Stimulus は controller に次の *CSS class プロパティ* を追加する。

| Kind         | Name                          | Value                                                                             |
| ------------ | ----------------------------- | --------------------------------------------------------------------------------- |
| Singular     | `this.[logicalName]Class`     | `logicalName` に対応する CSS class 属性の値                                        |
| Plural       | `this.[logicalName]Classes`   | 対応する CSS class 属性のすべてのクラスを、スペースで分割した配列                  |
| Existential  | `this.has[LogicalName]Class`  | 対応する CSS class 属性が存在するかどうかを示す真偽値                              |

これらのプロパティを使って、[DOM `classList` API](https://developer.mozilla.org/en-US/docs/Web/API/Element/classList) の `add()` や `remove()` メソッドで要素に CSS class を適用する。

たとえば、`search` controller の要素で結果取得の前にローディングインジケータを表示するには、`loadResults` action を次のように実装できる。

```js
export default class extends Controller {
  static classes = [ "loading" ]

  loadResults() {
    this.element.classList.add(this.loadingClass)

    fetch(/* … */)
  }
}
```

CSS class 属性がクラス名のリストを含んでいる場合、singular の CSS class プロパティはそのリストの最初のクラスを返す。

すべてのクラス名を配列としてアクセスするには plural の CSS class プロパティを使う。これを [スプレッド構文](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/Spread_syntax) と組み合わせれば、複数のクラスを一度に適用できる。

```js
export default class extends Controller {
  static classes = [ "loading" ]

  loadResults() {
    this.element.classList.add(...this.loadingClasses)

    fetch(/* … */)
  }
}
```

**Note:** 対応する CSS class 属性が存在しないときに CSS class プロパティにアクセスしようとすると、Stimulus はエラーを投げる。

## Naming Conventions

CSS class 定義の論理名は camelCase で指定する。論理名は camelCase の CSS class プロパティにマップされる。

```js
export default class extends Controller {
  static classes = [ "loading", "noResults" ]

  loadResults() {
    // …
    if (results.length == 0) {
      this.element.classList.add(this.noResultsClass)
    }
  }
}
```

HTML では CSS class 属性を kebab-case で書く。

```html
<form data-controller="search"
      data-search-loading-class="search--busy"
      data-search-no-results-class="search--empty">
```

CSS class 属性を構築する際は、[Controllers: Naming Conventions](controllers.md#naming-conventions) に記された identifier の規約に従うこと。
