> 原文: https://stimulus.hotwired.dev/reference/using-typescript
> タイトル: Using Typescript
> 最終取得: 2026-07-19

# Using Typescript

Stimulus 自体は [TypeScript](https://www.typescriptlang.org/) で書かれており、パッケージから直接型定義が提供されている。以下では、Stimulus のプロパティに対して型をどう定義するかを示す。

## Define Controller Element Type

デフォルトでは controller の `element` は `Element` 型になる。[Generic Type](https://www.typescriptlang.org/docs/handbook/2/generics.html) として型を指定することで、controller 要素の型を上書きできる。たとえば、要素の型が `HTMLFormElement` であることが期待される場合は次のようになる。

```ts
import { Controller } from "@hotwired/stimulus"

export default class MyController extends Controller<HTMLFormElement> {
  submit() {
    new FormData(this.element)
  }
}
```

## Define Value Properties

TypeScript の `declare` キーワードを使って、設定した value のプロパティを定義できる。controller 内で参照する value についてだけ、プロパティを定義しておけばよい。

```ts
import { Controller } from "@hotwired/stimulus"

export default class MyController extends Controller {
  static values = {
    code: String
  }

  declare codeValue: string
  declare readonly hasCodeValue: boolean
}
```

> `declare` キーワードは既存の Stimulus プロパティを上書きしないようにしつつ、TypeScript のための型定義だけを与える。

## Define Target Properties

TypeScript の `declare` キーワードを使って、設定した target のプロパティを定義できる。controller 内で参照する target についてだけ、プロパティを定義しておけばよい。

`[name]Target` と `[name]Targets` プロパティの戻り値の型は、`Element` 型を継承したものであればどれでもよい。目的に最も合う型を選ぶこと。汎用的な HTML 要素として定義したい場合は `Element` か `HTMLElement` を選ぶ。

```ts
import { Controller } from "@hotwired/stimulus"

export default class MyController extends Controller {
  static targets = [ "input" ]

  declare readonly hasInputTarget: boolean
  declare readonly inputTarget: HTMLInputElement
  declare readonly inputTargets: HTMLInputElement[]
}
```

> `declare` キーワードは既存の Stimulus プロパティを上書きしないようにしつつ、TypeScript のための型定義だけを与える。

## Custom properties and methods

controller クラスには、他の任意のカスタムプロパティを TypeScript 流に定義できる。

```ts
import { Controller } from "@hotwired/stimulus"

export default class MyController extends Controller {
  container: HTMLElement
}
```

詳しくは [TypeScript Documentation](https://www.typescriptlang.org/docs/handbook/intro.html) を参照。
