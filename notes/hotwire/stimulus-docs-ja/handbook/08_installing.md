> 原文: https://stimulus.hotwired.dev/handbook/installing
> タイトル: Installing Stimulus in Your Application
> 最終取得: 2026-07-19

# Installing Stimulus in Your Application

Stimulus をアプリケーションにインストールするには、[`@hotwired/stimulus` npm パッケージ](https://www.npmjs.com/package/@hotwired/stimulus) を JavaScript バンドルに追加します。あるいは、`<script type="module">` タグ内で [`stimulus.js`](https://unpkg.com/@hotwired/stimulus/dist/stimulus.js) を import します。

## Using Stimulus for Rails

[Stimulus for Rails](https://github.com/hotwired/stimulus-rails/) を [import map](https://github.com/rails/importmap-rails) と一緒に使っている場合、この統合によって `app/javascript/controllers` の下のすべての controller ファイルが自動的にロードされます。

### Controller Filenames Map to Identifiers

controller ファイルは `[identifier]_controller.js` という名前を付けてください。`identifier` は、あなたの HTML における各 controller の `data-controller` の identifier に対応します。

Stimulus for Rails では、ファイル名の中で複数の単語を区切るのに慣例としてアンダースコアを用います。controller のファイル名の中のアンダースコア 1 つは、identifier のダッシュ 1 つに変換されます。

サブフォルダを使って controller に名前空間を付けることもできます。名前空間付き controller ファイルのパスに含まれるスラッシュ 1 つは、identifier のダッシュ 2 つに変換されます。

もし好みであれば、controller ファイル名のどこでもアンダースコアの代わりにダッシュを使ってかまいません。Stimulus はそれらを同じものとして扱います。

| If your controller file is named…    | its identifier will be… |
| ------------------------------------ | ----------------------- |
| clipboard_controller.js              | clipboard               |
| date_picker_controller.js            | date-picker             |
| users/list_item_controller.js        | users--list-item        |
| local-time-controller.js             | local-time              |

## Using Webpack Helpers

JavaScript バンドラとして Webpack を使っている場合、[@hotwired/stimulus-webpack-helpers](https://www.npmjs.com/package/@hotwired/stimulus-webpack-helpers) パッケージを使えば Stimulus for Rails と同じ形の自動ロードを行えます。まずパッケージを追加し、次のように使います。

```js
import { Application } from "@hotwired/stimulus"
import { definitionsFromContext } from "@hotwired/stimulus-webpack-helpers"

window.Stimulus = Application.start()
const context = require.context("./controllers", true, /\.js$/)
Stimulus.load(definitionsFromContext(context))
```

## Using Other Build Systems

Stimulus は他のビルドシステムとも動作しますが、controller の自動ロードはサポートされません。その代わりに、controller ファイルを明示的にロードしてアプリケーションインスタンスに登録する必要があります。

```js
// src/application.js
import { Application } from "@hotwired/stimulus"

import HelloController from "./controllers/hello_controller"
import ClipboardController from "./controllers/clipboard_controller"

window.Stimulus = Application.start()
Stimulus.register("hello", HelloController)
Stimulus.register("clipboard", ClipboardController)
```

esbuild のようなビルダと stimulus-rails を組み合わせて使っている場合は、`stimulus:manifest:update` の Rake タスクと `./bin/rails generate stimulus [controller]` ジェネレータを使うことで、`app/javascript/controllers/index.js` にある controller のインデックスファイルを自動更新し続けられます。

## Using Without a Build System

ビルドシステムを使いたくない場合は、Stimulus を `<script type="module">` タグでロードすることもできます。

```html
<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <script type="module">
    import { Application, Controller } from "https://unpkg.com/@hotwired/stimulus/dist/stimulus.js"
    window.Stimulus = Application.start()

    Stimulus.register("hello", class extends Controller {
      static targets = [ "name" ]

      connect() {
      }
    })
  </script>
</head>
<body>
  <div data-controller="hello">
    <input data-hello-target="name" type="text">
    …
  </div>
</body>
</html>
```

## Overriding Attribute Defaults

Stimulus の `data-*` 属性がプロジェクト内の他のライブラリと競合する場合、Stimulus の `Application` を作成する際にそれらをオーバーライドできます。

- `data-controller`
- `data-action`
- `data-target`

これらの Stimulus の中核属性はオーバーライド可能です ([schema.ts](https://github.com/hotwired/stimulus/blob/main/src/core/schema.ts) を参照)。

```js
// src/application.js
import { Application, defaultSchema } from "@hotwired/stimulus"

const customSchema = {
  ...defaultSchema,
  actionAttribute: 'data-stimulus-action'
}

window.Stimulus = Application.start(document.documentElement, customSchema);
```

## Error handling

Stimulus からアプリケーションコードへの呼び出しはすべて `try ... catch` ブロックでラップされています。

もしコードがエラーを投げた場合、そのエラーは Stimulus によってキャッチされ、ブラウザのコンソールに、controller 名や、呼び出されていたイベント・ライフサイクル関数などの詳細情報とともにログ出力されます。`window.onerror` を定義するエラートラッキングシステムを使っている場合、Stimulus はエラーをそちらにも渡します。

`Application#handleError` を定義することで、Stimulus のエラーハンドリングをオーバーライドできます。

```js
// src/application.js
import { Application } from "@hotwired/stimulus"
window.Stimulus = Application.start()

Stimulus.handleError = (error, message, detail) => {
  console.warn(message, detail)
  ErrorTrackingSystem.captureException(error)
}
```

## Debugging

Stimulus アプリケーションを `window.Stimulus` に割り当ててあるなら、コンソールから `Stimulus.debug = true` と実行することで [デバッグモード](https://github.com/hotwired/stimulus/pull/354) を有効化できます。ソースコード内でアプリケーションインスタンスを設定するときに、このフラグを立てておくこともできます。

## Browser Support

Stimulus は、自動更新される主要デスクトップおよびモバイルブラウザすべてを標準でサポートします。Stimulus 3 以降は Internet Explorer 11 をサポートしません (その用途には Stimulus 2 と @stimulus/polyfills の組み合わせを使えます)。
