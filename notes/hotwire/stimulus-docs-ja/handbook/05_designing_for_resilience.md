> 原文: https://stimulus.hotwired.dev/handbook/designing-for-resilience
> タイトル: Designing For Resilience
> 最終取得: 2026-07-19

# Designing For Resilience

クリップボード API は [現在のブラウザで十分サポートされている](https://caniuse.com/#feat=clipboard) とはいえ、それでも古いブラウザを使う少数のユーザーがアプリを利用する可能性は残ります。

さらに、アプリへのアクセスが時折うまくいかない状況も想定すべきです。たとえば断続的なネットワーク障害や CDN の可用性の問題によって、JavaScript の一部あるいは全部が読み込めないことがあり得ます。

古いブラウザへの対応を「手間に見合わない」と切り捨てたり、ネットワーク障害を「リロードすればすぐ直る一時的なトラブル」と片付けたりしたくなる誘惑があります。しかし多くの場合、この種の問題に対して優雅にレジリエンス (回復力) を発揮する形で機能を作ることは、驚くほど簡単です。

このレジリエントなアプローチは一般に *プログレッシブ・エンハンスメント* と呼ばれ、基本機能は HTML と CSS で実装し、その基本体験に対して段階的なアップグレードを、それを支える技術がブラウザに備わっているときに、CSS や JavaScript で少しずつ上乗せしていく Web インターフェースの届け方の実践です。

## Progressively Enhancing the PIN Field

PIN フィールドをプログレッシブに強化して、Copy ボタンをブラウザがサポートしている場合にだけ見えるようにしてみましょう。こうしておけば、動かないボタンを目の前に見せる、という事態を避けられます。

まず CSS で Copy ボタンを隠すところから始めます。次に、Stimulus controller の中でクリップボード API のサポートを *フィーチャーテスト* します。もし API がサポートされていれば、controller の要素にクラス名を付与してボタンを表示させます。

まずは `data-controller` 属性を持つ `div` 要素に `data-clipboard-supported-class="clipboard--supported"` を追加します。

```html
  <div data-controller="clipboard" data-clipboard-supported-class="clipboard--supported">
```

次にボタン要素に `class="clipboard-button"` を追加します。

```html
  <button data-action="clipboard#copy" class="clipboard-button">Copy to Clipboard</button>
```

そして `public/main.css` に次のスタイルを追加します。

```css
.clipboard-button {
  display: none;
}

.clipboard--supported .clipboard-button {
  display: initial;
}
```

まず、この `data-clipboard-supported-class` 属性を controller 内部で静的な class として登録します。

```js
  static classes = [ "supported" ]
```

こうすることで、具体的な CSS クラス名を HTML 側で制御できるようになり、この controller は多様な CSS アプローチにいっそう合わせやすくなります。このやり方で追加された具体的なクラス名は `this.supportedClass` からアクセスできます。

続いて controller に `connect()` メソッドを追加し、クリップボード API がサポートされているかを判定して、controller の要素にクラス名を追加します。

```js
  connect() {
    if ("clipboard" in navigator) {
      this.element.classList.add(this.supportedClass);
    }
  }
```

このメソッドは controller のクラス本体のどこに置いてもかまいません。

もしよければ、ブラウザで JavaScript を無効にしてページをリロードしてみてください。Copy ボタンがもう見えなくなっていることが分かるはずです。

これで PIN フィールドをプログレッシブに強化できました。Copy ボタンのベースラインの状態は非表示、そしてクリップボード API のサポートが JavaScript によって検出されたときにだけ表示される、という設計です。

## Wrap-Up and Next Steps

本章では clipboard controller に穏やかな修正を加え、古いブラウザや不安定なネットワーク条件に対するレジリエンスを持たせました。

次章では、Stimulus controller がどのように状態を管理するかを学んでいきます。
