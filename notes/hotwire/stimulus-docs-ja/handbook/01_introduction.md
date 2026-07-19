> 原文: https://stimulus.hotwired.dev/handbook/introduction
> タイトル: Introduction
> 最終取得: 2026-07-19

# Introduction

## About Stimulus

Stimulus は、ささやかな野心を持つ JavaScript フレームワークです。他のフロントエンドフレームワークとは異なり、Stimulus は *静的な (static)* あるいは *サーバーレンダリングされた (server-rendered)* HTML —「あなたがすでに持っている HTML」— を、シンプルな注釈でページ上の要素と JavaScript オブジェクトを結び付けることによって強化する、という設計思想で作られています。

これらの JavaScript オブジェクトは *controller* と呼ばれます。Stimulus は HTML の `data-controller` 属性が現れるのを常時ページ上で監視しています。各属性について、Stimulus はその値を見て対応する controller クラスを探し出し、そのクラスの新しいインスタンスを生成して要素と結び付けます。

こう考えてみるとよいかもしれません。`class` 属性が HTML と CSS を結ぶ橋であるのと同じように、Stimulus の `data-controller` 属性は HTML と JavaScript を結ぶ橋なのだ、と。

controller のほかに、Stimulus が備えるおもな概念は次の 3 つです。

- *actions*: `data-action` 属性を使って、controller のメソッドを DOM イベントに接続するもの
- *targets*: controller の中で重要な要素を見つけ出すためのもの
- *values*: controller の要素上の data 属性を読み書きし、変化を監視するためのもの

Stimulus が data 属性を用いることは、CSS がコンテンツと見た目を分離するのと同じやり方で、コンテンツと振る舞いを分離するのに役立ちます。さらに Stimulus の規約は、関連するコードを名前で自然にグループ化するように促してくれます。

その結果、Stimulus は小さく再利用しやすい controller を書く手助けをしてくれ、コードが「JavaScript スープ」へと堕ちてしまわないよう、ちょうどよい量の構造を与えてくれるのです。

## About This Book

このハンドブックは、いくつかの完全に動作する controller を書きながら、Stimulus の中核概念を案内していきます。各章は前の章の内容を積み重ねる形で進みます。最初から最後まで通して読むと、次のようなことを学べます。

- テキストフィールドに入力された名前に宛てたあいさつを表示する
- ボタンをクリックしたら、テキストフィールドの内容をシステムのクリップボードにコピーする
- 複数のスライドを持つスライドショーをナビゲーションする
- サーバーから HTML を取得してページ上の要素に自動で埋め込む
- 自分のアプリケーションに Stimulus をセットアップする

ここでの演習を一通り終えたら、Stimulus API の技術的な詳細を理解するために [reference documentation](../reference/controllers.md) を参考にするとよいでしょう。

さあ、始めましょう!
