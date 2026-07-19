> 原文: https://stimulus.hotwired.dev/reference/actions
> タイトル: Actions
> 最終取得: 2026-07-19

# Actions

*Actions* は、controller の中で DOM イベントを扱うための仕組みである。

```html
<div data-controller="gallery">
  <button data-action="click->gallery#next">…</button>
</div>
```

```js
// controllers/gallery_controller.js
import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
  next(event) {
    // …
  }
}
```

action は次の 3 つのあいだの接続である。

- controller のメソッド
- controller の要素
- DOM のイベントリスナー

## Descriptors

`data-action` の値 `click->gallery#next` は *action descriptor* と呼ばれる。この記述子の各要素は次の意味を持つ。

- `click` は待ち受ける DOM イベントの名前
- `gallery` は controller の identifier
- `next` は呼び出すメソッドの名前

### Event Shorthand

Stimulus では、上の button/click のようによく使われる要素/イベントの組について、イベント名を省略して action descriptor を短く書ける。

```html
<button data-action="gallery#next">…</button>
```

こうしたショートハンドの組は次のとおり。

| Element             | Default Event |
| ------------------- | ------------- |
| a                   | click         |
| button              | click         |
| details             | toggle        |
| form                | submit        |
| input               | input         |
| input type=submit   | click         |
| select              | change        |
| textarea            | input         |

## KeyboardEvent Filter

[KeyboardEvent](https://developer.mozilla.org/en-US/docs/Web/API/KeyboardEvent) の Action について、特定のキーが押されたときにだけ controller メソッドを呼びたい場合がある。

`Escape` キーにだけ反応するイベントリスナーは、次の例のように action descriptor のイベント名に `.esc` を付けて登録できる。

```html
<div data-controller="modal"
     data-action="keydown.esc->modal#close" tabindex="0">
</div>
```

これは発火するイベントがキーボードイベントである場合にのみ機能する。

これらのフィルターとキーの対応は下表の通り。

| Filter     | Key Name    |
| ---------- | ----------- |
| enter      | Enter       |
| tab        | Tab         |
| esc        | Escape      |
| space      | " "         |
| up         | ArrowUp     |
| down       | ArrowDown   |
| left       | ArrowLeft   |
| right      | ArrowRight  |
| home       | Home        |
| end        | End         |
| page_up    | PageUp      |
| page_down  | PageDown    |
| [a-z]      | [a-z]       |
| [0-9]      | [0-9]       |

他のキーをサポートしたい場合は、カスタムスキーマでモディファイアを拡張できる。

```javascript
import { Application, defaultSchema } from "@hotwired/stimulus"

const customSchema = {
  ...defaultSchema,
  keyMappings: { ...defaultSchema.keyMappings, at: "@" },
}

const app = Application.start(document.documentElement, customSchema)
```

モディファイアキーとの組み合わせフィルタを購読したい場合は、`ctrl+a` のように書ける。

```html
<div data-action="keydown.ctrl+a->listbox#selectAll" role="option" tabindex="0">...</div>
```

サポートされているモディファイアキーは下表の通り。

| Modifier   | Notes                    |
| ---------- | ------------------------ |
| `alt`      | MacOS では `option`      |
| `ctrl`     |                          |
| `meta`     | MacOS では Command キー  |
| `shift`    |                          |

### Global Events

グローバルな `window` や `document` オブジェクトに向けて dispatch されるイベントを controller で待ち受けたい場合がある。

イベントリスナーを `window` あるいは `document` に登録するには、action descriptor のイベント名 (フィルタモディファイアがあればそれを含めた形) に `@window` または `@document` を付ける。

```html
<div data-controller="gallery"
     data-action="resize@window->gallery#layout">
</div>
```

### Options

[DOM イベントリスナーのオプション](https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/addEventListener#Parameters) を指定したい場合は、action descriptor の末尾に *action option* を 1 つ以上付け加えられる。

```html
<div data-controller="gallery"
     data-action="scroll->gallery#layout:!passive">
  <img data-action="click->gallery#open:capture">
```

Stimulus は次の action option をサポートする。

| Action option | DOM event listener option  |
| ------------- | -------------------------- |
| `:capture`    | `{ capture: true }`        |
| `:once`       | `{ once: true }`           |
| `:passive`    | `{ passive: true }`        |
| `:!passive`   | `{ passive: false }`       |

これらに加えて、DOM のイベントリスナーオプションにはネイティブに存在しない、Stimulus 独自の action option もある。

| Custom action option | Description                                                             |
| -------------------- | ----------------------------------------------------------------------- |
| `:stop`              | メソッド呼び出しの前にイベントの `.stopPropagation()` を呼ぶ            |
| `:prevent`           | メソッド呼び出しの前にイベントの `.preventDefault()` を呼ぶ             |
| `:self`              | イベントが要素自身から発火した場合にのみメソッドを呼ぶ                  |

`Application.registerActionOption` メソッドで独自の action option を登録できる。

たとえば、`<details>` 要素は開閉されるたびに [toggle](https://developer.mozilla.org/en-US/docs/Web/API/HTMLDetailsElement/toggle_event) イベントを dispatch する。独自の `:open` action option を用意すれば、要素が *open* にトグルされたときだけイベントをルーティングできる。

```javascript
import { Application } from "@hotwired/stimulus"

const application = Application.start()

application.registerActionOption("open", ({ event }) => {
  if (event.type == "toggle") {
    return event.target.open == true
  } else {
    return true
  }
})
```

同様に、独自の `:!open` option を使えば、要素が *close* にトグルされたときにイベントをルーティングできる。action descriptor のオプションを `!` プレフィックス付きで宣言すると、コールバックの `value` 引数には `false` が渡される。

```javascript
import { Application } from "@hotwired/stimulus"

const application = Application.start()

application.registerActionOption("open", ({ event, value }) => {
  if (event.type == "toggle") {
    return event.target.open == value
  } else {
    return true
  }
})
```

イベントが controller の action にルーティングされないようにするには、`registerActionOption` のコールバック関数が `false` を返す必要がある。ルーティングしたい場合は `true` を返す。

コールバックは次のキーを持つオブジェクトを 1 つの引数として受け取る。

| Name        | Description                                                                                                              |
| ----------- | ------------------------------------------------------------------------------------------------------------------------ |
| name        | String: option の名前 (上の例では `"open"`)                                                                              |
| value       | Boolean: option の値 (`:open` なら `true`、`:!open` なら `false`)                                                        |
| event       | [Event](https://developer.mozilla.org/en-US/docs/web/api/event): イベントインスタンス。submitter 要素の action parameter を含む `params` 付き |
| element     | [Element](https://developer.mozilla.org/en-US/docs/Web/API/element): action descriptor が宣言された要素                   |
| controller  | メソッド呼び出しを受け取ることになる `Controller` インスタンス                                                            |

## Event Objects

*action メソッド* は、action のイベントリスナーとして機能する controller のメソッドである。

action メソッドの第 1 引数は DOM の *event object* である。イベントにアクセスしたい理由はいろいろある。

- キーボードイベントからキーコードを読む
- マウスイベントから座標を読む
- input イベントからデータを読む
- action の submitter 要素から params を読む
- ブラウザのデフォルト挙動をキャンセルする
- イベントがバブリングしてこの action にたどり着く前に、どの要素から発火したかを知る

すべてのイベントに共通する基本的なプロパティは次の通り。

| Event Property         | Value                                                                                                   |
| ---------------------- | ------------------------------------------------------------------------------------------------------- |
| event.type             | イベントの名前 (例: `"click"`)                                                                          |
| event.target           | イベントを dispatch した対象 (すなわちクリックされた最内側の要素)                                       |
| event.currentTarget    | イベントリスナーがインストールされている対象 (`data-action` 属性を持つ要素、あるいは `document` や `window`) |
| event.params           | action の submitter 要素が渡した action params                                                          |

次のイベントメソッドを使うと、イベントの扱いをより細かく制御できる。

| Event Method                | Result                                                                              |
| --------------------------- | ----------------------------------------------------------------------------------- |
| event.preventDefault()      | イベントのデフォルト挙動 (リンクを辿る、フォームを送信する 等) をキャンセルする    |
| event.stopPropagation()     | イベントが親要素の他のリスナーへバブリングする前に止める                            |

## Multiple Actions

`data-action` 属性の値は、スペース区切りの action descriptor のリストである。

1 つの要素に複数の action が乗るのはよくあることだ。たとえば次の input 要素は、フォーカスを得たときに `field` controller の `highlight()` を、value が変わるたびに `search` controller の `update()` を呼び出す。

```html
<input type="text" data-action="focus->field#highlight input->search#update">
```

同じイベントに対して複数の action が定義されているとき、Stimulus は descriptor が並んでいる順序 (左から右) で action を呼び出す。

action の連鎖は、action の中で `Event#stopImmediatePropagation()` を呼び出すことで任意のタイミングで止められる。右側にある追加の action は無視される。

```javascript
highlight(event) {
  event.stopImmediatePropagation()
  // ...
}
```

## Naming Conventions

action 名は controller のメソッドに直接マップされるので、常に camelCase を使う。

`click`、`onClick`、`handleClick` のように、単にイベント名を繰り返しただけの action 名は避ける。

```html
<button data-action="click->profile#click">Don't</button>
```

代わりに、呼び出されたときに *何が起きるか* に基づいて action メソッドを命名する。

```html
<button data-action="click->profile#showDialog">Do</button>
```

こうすることで、controller のソースを見なくても、HTML ブロックの挙動を推測しやすくなる。

## Action Parameters

action は、submitter 要素から渡されるパラメータを持てる。書式は `data-[identifier]-[param-name]-param`。パラメータは、それを渡したい action と同じ要素上に指定する必要がある。

すべてのパラメータは、値から推論されて `Number`、`String`、`Object`、`Boolean` のいずれかに自動的にキャストされる。

| Data attribute                                     | Param                          | Type    |
| -------------------------------------------------- | ------------------------------ | ------- |
| `data-item-id-param="12345"`                       | `12345`                        | Number  |
| `data-item-url-param="/votes"`                     | `"/votes"`                     | String  |
| `data-item-payload-param='{"value":"1234567"}'`    | `{ value: 1234567 }`           | Object  |
| `data-item-active-param="true"`                    | `true`                         | Boolean |

以下のセットアップを考える。

```html
<div data-controller="item spinner">
  <button data-action="item#upvote spinner#start" 
    data-item-id-param="12345" 
    data-item-url-param="/votes"
    data-item-payload-param='{"value":"1234567"}' 
    data-item-active-param="true">…</button>
</div>
```

これは `ItemController#upvote` と `SpinnerController#start` の両方を呼ぶが、パラメータが渡されるのは前者だけである。

```js
// ItemController
upvote(event) {
  // { id: 12345, url: "/votes", active: true, payload: { value: 1234567 } }
  console.log(event.params) 
}

// SpinnerController
start(event) {
  // {}
  console.log(event.params) 
}
```

event の他の情報が要らない場合は、params を分割代入できる。

```js
upvote({ params }) {
  // { id: 12345, url: "/votes", active: true, payload: { value: 1234567 } }
  console.log(params) 
}
```

同じ controller の複数 action が同じ submitter 要素を共有している場合など、必要な params だけを分割代入することもできる。

```js
upvote({ params: { id, url } }) {
  console.log(id) // 12345
  console.log(url) // "/votes"
}
```
