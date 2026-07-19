> 原文: https://stimulus.hotwired.dev/reference/lifecycle-callbacks
> タイトル: Lifecycle Callbacks
> 最終取得: 2026-07-19

# Lifecycle Callbacks

*ライフサイクルコールバック* と呼ばれる特殊なメソッドを定義しておくと、controller や特定の target がドキュメントに接続・切断されるたびに反応できる。

```js
import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
  connect() {
    // …
  }
}
```

## Methods

controller には次のメソッドをいずれも定義できる。

| Method                                        | Invoked by Stimulus…                            |
| --------------------------------------------- | ----------------------------------------------- |
| initialize()                                  | controller が最初にインスタンス化されたときに 1 回だけ |
| [name]TargetConnected(target: Element)        | target が DOM に接続されるたびに                 |
| connect()                                     | controller が DOM に接続されるたびに             |
| [name]TargetDisconnected(target: Element)     | target が DOM から切断されるたびに               |
| disconnect()                                  | controller が DOM から切断されるたびに           |

## Connection

controller は次の 2 つの条件がともに真であるときに、ドキュメントに *接続 (connected)* されている状態となる。

- その要素がドキュメント上に存在する (すなわち `document.documentElement`、つまり `<html>` 要素の子孫である)
- 要素の `data-controller` 属性の中にその identifier が含まれている

controller が接続状態になったとき、Stimulus は `connect()` メソッドを呼び出す。

### Targets

target は次の 2 つの条件がともに真であるときに、ドキュメントに *接続* されている状態となる。

- その要素が、対応する controller の要素の子孫としてドキュメント上に存在する
- 要素の `data-{identifier}-target` 属性の中にその識別名が含まれている

target が接続状態になったとき、Stimulus は controller の `[name]TargetConnected()` メソッドを、その target 要素を引数として呼び出す。`[name]TargetConnected()` のライフサイクルコールバックは、controller の `connect()` コールバックの *前* に発火する。

## Disconnection

接続済みの controller は、上記の条件のいずれかが偽になったとき—たとえば次のいずれかのシナリオで—後から *切断 (disconnected)* 状態になる。

- 要素が `Node#removeChild()` あるいは `ChildNode#remove()` で明示的にドキュメントから取り除かれた
- 要素の親要素のいずれかがドキュメントから取り除かれた
- 要素の親要素のいずれかの中身が `Element#innerHTML=` で置き換えられた
- 要素の `data-controller` 属性が削除もしくは変更された
- Turbo のページ切り替え時など、ドキュメントに新しい `<body>` 要素がインストールされた

controller が切断状態になったとき、Stimulus は `disconnect()` メソッドを呼び出す。

### Targets

接続済みの target も、上記の条件のいずれかが偽になったとき—たとえば次のいずれかのシナリオで—後から *切断* 状態になる。

- 要素が `Node#removeChild()` あるいは `ChildNode#remove()` で明示的にドキュメントから取り除かれた
- 要素の親要素のいずれかがドキュメントから取り除かれた
- 要素の親要素のいずれかの中身が `Element#innerHTML=` で置き換えられた
- 要素の `data-{identifier}-target` 属性が削除もしくは変更された
- Turbo のページ切り替え時など、ドキュメントに新しい `<body>` 要素がインストールされた

target が切断状態になったとき、Stimulus は controller の `[name]TargetDisconnected()` メソッドを、その target 要素を引数として呼び出す。`[name]TargetDisconnected()` のライフサイクルコールバックは、controller の `disconnect()` コールバックの *前* に発火する。

## Reconnection

切断された controller は、後で再度接続状態になることがある。

controller の要素をドキュメントから取り除いた後にまた繋ぎ直したような場合には、Stimulus は要素の以前の controller インスタンスを再利用し、その `connect()` メソッドを複数回呼び出す。

同様に、切断された target も後で再度接続されることがある。Stimulus は controller の `[name]TargetConnected()` メソッドを複数回呼び出す。

## Order and Timing

Stimulus は [DOM の `MutationObserver` API](https://developer.mozilla.org/en-US/docs/Web/API/MutationObserver) を使って、ページの変化を非同期に監視している。

つまり Stimulus は、ドキュメントに変化が発生したあと、その変化ごとの次の [microtask](https://jakearchibald.com/2015/tasks-microtasks-queues-and-schedules/) で、controller のライフサイクルメソッドを非同期に呼び出す。

ライフサイクルメソッドは発生順に実行されるので、controller の `connect()` の 2 回の呼び出しの間には、常に 1 回の `disconnect()` が挟まる。同様に、ある特定の target について controller の `[name]TargetConnected()` が 2 回呼ばれる間には、その同じ target についての `[name]TargetDisconnected()` が必ず 1 回挟まる。
