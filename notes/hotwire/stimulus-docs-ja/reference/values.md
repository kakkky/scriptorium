> 原文: https://stimulus.hotwired.dev/reference/values
> タイトル: Values
> 最終取得: 2026-07-19

# Values

controller の特別なプロパティを介して、controller 要素上の [HTML data attribute](https://developer.mozilla.org/en-US/docs/Web/HTML/Global_attributes/data-*) を型付きの *value* として読み書きできる。

```html
<div data-controller="loader" data-loader-url-value="/messages">
</div>
```

上の HTML スニペットのように、value の data 属性は `data-controller` 属性と同じ要素に置くこと。

```js
// controllers/loader_controller.js
import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
  static values = {
    url: String
  }

  connect() {
    fetch(this.urlValue).then(/* … */)
  }
}
```

## Definitions

controller の value は `static values` オブジェクトで定義する。左に value の *名前*、右にその *型* を置く。

```js
export default class extends Controller {
  static values = {
    url: String,
    interval: Number,
    params: Object
  }

  // …
}
```

## Types

value の型は `Array`、`Boolean`、`Number`、`Object`、`String` のいずれか。この型が、JavaScript と HTML のあいだで value がどのように変換されるかを決める。

| Type    | Encoded as…               | Decoded as…                                  |
| ------- | ------------------------- | -------------------------------------------- |
| Array   | `JSON.stringify(array)`   | `JSON.parse(value)`                          |
| Boolean | `boolean.toString()`      | `!(value == "0" || value == "false")`        |
| Number  | `number.toString()`       | `Number(value.replace(/_/g, ""))`            |
| Object  | `JSON.stringify(object)`  | `JSON.parse(value)`                          |
| String  | そのまま                  | そのまま                                     |

## Properties and Attributes

Stimulus は、controller に定義された各 value について、getter・setter・existential プロパティを自動生成する。これらのプロパティは controller 要素の data 属性と連動している。

| Kind         | Property name         | Effect                                       |
| ------------ | --------------------- | -------------------------------------------- |
| Getter       | `this.[name]Value`    | `data-[identifier]-[name]-value` を読む      |
| Setter       | `this.[name]Value=`   | `data-[identifier]-[name]-value` に書く      |
| Existential  | `this.has[Name]Value` | `data-[identifier]-[name]-value` の有無を検査 |

### Getters

value の getter は、対応する data 属性をその value の型のインスタンスにデコードする。

controller 要素にその data 属性がない場合、getter は value の型に応じた *デフォルト値* を返す。

| Type    | Default value |
| ------- | ------------- |
| Array   | `[]`          |
| Boolean | `false`       |
| Number  | `0`           |
| Object  | `{}`          |
| String  | `""`          |

### Setters

value の setter は、controller 要素上の対応する data 属性を書き換える。

controller 要素からその data 属性を消したい場合は、value に `undefined` を代入する。

### Existential Properties

value の existential プロパティは、対応する data 属性が controller 要素に存在すれば `true`、なければ `false` になる。

## Change Callbacks

Value の *change callback* を使うと、value の data 属性が変化したときに反応できる。

controller に `[name]ValueChanged` メソッドを定義する (`[name]` は変化を監視したい value の名前)。メソッドは、デコード済みの現在値を第 1 引数、デコード済みの以前の値を第 2 引数として受け取る。

Stimulus は各 change callback を、controller の初期化後に呼び出し、その後は対応する data 属性が変化するたびに呼び出す。これには value の setter への代入による変化も含まれる。

```js
export default class extends Controller {
  static values = { url: String }

  urlValueChanged() {
    fetch(this.urlValue).then(/* … */)
  }
}
```

### Previous Values

`[name]ValueChanged` コールバックのメソッドを 2 引数で定義することで、以前の値にアクセスできる。

```js
export default class extends Controller {
  static values = { url: String }

  urlValueChanged(value, previousValue) {
    /* … */
  }
}
```

2 つの引数の名前は好きに決めてよい。たとえば `urlValueChanged(current, old)` でも構わない。

## Default Values

controller 要素で指定されていない value は、controller 定義の中でデフォルトを指定できる。

```js
export default class extends Controller {
  static values = {
    url: { type: String, default: '/bill' },
    interval: { type: Number, default: 5 },
    clicked: Boolean
  }
}
```

デフォルトを使うときは、`{ type, default }` という展開形式を用いる。この形式と、デフォルトを使わない通常の形式は混在させてよい。

## Naming Conventions

value 名は JavaScript では camelCase、HTML では kebab-case で書く。たとえば `loader` controller の `contentType` という value は、対応する data 属性が `data-loader-content-type-value` になる。
