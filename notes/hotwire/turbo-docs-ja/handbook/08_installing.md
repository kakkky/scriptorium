> 原文: https://turbo.hotwired.dev/handbook/installing
> タイトル: Installing Turbo in Your Application
> 最終取得: 2026-07-05

# Installing Turbo in Your Application

Turbo は、コンパイル済みの Turbo distributable スクリプトをアプリケーションの `<head>` に直接読み込む形でも、esbuild のようなバンドラーを介して npm 経由でも利用できます。

## In Compiled Form

jsDelivr のような CDN バンドラーを使って、Turbo の最新リリースにそのまま追従できます。アプリケーションの `<head>` に `<script>` タグを含めるだけです。

```html
<head>
  <script type="module" src="https://cdn.jsdelivr.net/npm/@hotwired/turbo@latest/dist/turbo.es2017-esm.min.js"></script>
</head>
```

あるいは [unpkg からコンパイル済みパッケージをダウンロード](https://unpkg.com/browse/@hotwired/turbo@latest/dist/)することもできます。

## As An npm Package

Turbo は `npm` または `yarn` といったパッケージングツール経由で npm からインストールできます。

`Turbo.visit()` のような Turbo の関数を利用する場合は、`Turbo` の関数をコードに import します。

```javascript
import * as Turbo from "@hotwired/turbo"
```

`Turbo.visit()` のような Turbo の関数を利用*しない*場合は、ライブラリだけを import します。こうすることで、一部のバンドラーで発生する tree-shaking や未使用変数の問題を回避できます。MDN の [Import a module for its side effects only](https://developer.mozilla.org/en-US/docs/web/javascript/reference/statements/import#import_a_module_for_its_side_effects_only) を参照してください。

```javascript
import "@hotwired/turbo";
```

## In a Ruby on Rails application

Turbo の JavaScript フレームワークは、アセットパイプラインから直接利用できるように [turbo-rails gem](https://github.com/hotwired/turbo-rails) に含まれています。

[Reference: Drive](../reference/drive.md)
