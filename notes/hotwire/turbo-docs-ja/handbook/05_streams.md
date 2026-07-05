> 原文: https://turbo.hotwired.dev/handbook/streams
> タイトル: Come Alive with Turbo Streams
> 最終取得: 2026-07-05

# Come Alive with Turbo Streams

Turbo Streams は、`<turbo-stream>` 要素で包まれた HTML の断片としてページ変更を届けます。各 stream 要素は、その中の HTML に対して何を行うかを宣言するために、アクションと対象の ID をあわせて指定します。これらの要素は、従来型の HTTP レスポンスとして同期的にブラウザに届けることも、WebSocket や SSE などのトランスポート越しに非同期に届けることもできます。後者の場合、他のユーザーやプロセスによる更新でアプリケーションを「生きている」ように動かすことができます。

Turbo Streams は、たとえばリスト内の要素をページ全体を再ロードせずに取り除くといったユーザー操作後の DOM の外科的な更新に使ったり、他のユーザーから送られてきたメッセージを進行中の会話に追記するといったリアルタイム機能の実装に使ったりできます。

## Stream Messages and Actions

Turbo Streams メッセージは、`<turbo-stream>` 要素からなる HTML の断片です。以下の stream メッセージは、指定できる 8 つの stream アクションを示しています。

```html
<turbo-stream action="append" target="messages">
  <template>
    <div id="message_1">
      This div will be appended to the element with the DOM ID "messages".
    </div>
  </template>
</turbo-stream>

<turbo-stream action="prepend" target="messages">
  <template>
    <div id="message_1">
      This div will be prepended to the element with the DOM ID "messages".
    </div>
  </template>
</turbo-stream>

<turbo-stream action="replace" target="message_1">
  <template>
    <div id="message_1">
      This div will replace the existing element with the DOM ID "message_1".
    </div>
  </template>
</turbo-stream>

<turbo-stream action="replace" method="morph" target="current_step">
  <template>
    <!-- The contents of this template will replace the element with ID "current_step" via morph. -->
    <li>New item</li>
  </template>
</turbo-stream>

<turbo-stream action="update" target="unread_count">
  <template>
    <!-- The contents of this template will replace the
    contents of the element with ID "unread_count" by
    setting innerHtml to "" and then switching in the
    template contents. Any handlers bound to the element
    "unread_count" would be retained. This is to be
    contrasted with the "replace" action above, where
    that action would necessitate the rebuilding of
    handlers. -->
    1
  </template>
</turbo-stream>

<turbo-stream action="update" method="morph" target="current_step">
  <template>
    <!-- The contents of this template will replace the children of the element with ID "current_step" via morph. -->
    <li>New item</li>
  </template>
</turbo-stream>

<turbo-stream action="remove" target="message_1">
  <!-- The element with DOM ID "message_1" will be removed.
  The contents of this stream element are ignored. -->
</turbo-stream>

<turbo-stream action="before" target="current_step">
  <template>
    <!-- The contents of this template will be added before the
    the element with ID "current_step". -->
    <li>New item</li>
  </template>
</turbo-stream>

<turbo-stream action="after" target="current_step">
  <template>
    <!-- The contents of this template will be added after the
    the element with ID "current_step". -->
    <li>New item</li>
  </template>
</turbo-stream>

<turbo-stream action="refresh" request-id="abcd-1234"></turbo-stream>
```

すべての `<turbo-stream>` 要素は、含める HTML を `<template>` 要素で包む必要がある点に注意してください。

Turbo Stream は、[id](https://developer.mozilla.org/en-US/docs/Web/HTML/Global_attributes/id) 属性や [CSS セレクタ](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_Selectors)で解決できるドキュメント内の任意の要素と連携できます (ただし `<template>` 要素や `<iframe>` 要素の中身は例外です)。対象要素を [`<turbo-frame>` 要素](04_frames.md)に変える必要はありません。もし `<turbo-stream>` 要素のためだけに `<turbo-frame>` 要素を使っているのであれば、その `<turbo-frame>` を別の[組み込み要素](https://developer.mozilla.org/en-US/docs/Web/HTML/Element)に変更してください。

1 つの stream メッセージには、WebSocket、SSE、あるいはフォーム送信への応答として、いくつでも stream 要素をレンダリングできます。

また、ページに挿入された `<turbo-stream>` 要素 (たとえばページ全体や frame の読み込みで挿入された場合) は、Turbo によって処理されたあと DOM から取り除かれます。これにより、ページや frame が読み込まれるときに stream アクションを自動的に実行できます。

## Actions With Multiple Targets

アクションは、通常の DOM ID を参照する `target` 属性の代わりに、CSS クエリセレクタを使う `targets` 属性を用いることで、複数の対象に適用できます。例を示します。

```html
<turbo-stream action="remove" targets=".old_records">
  <!-- The element with the class "old_records" will be removed.
  The contents of this stream element are ignored. -->
</turbo-stream>

<turbo-stream action="after" targets="input.invalid_field">
  <template>
    <!-- The contents of this template will be added after the
    all elements that match "inputs.invalid_field". -->
    <span>Incorrect</span>
  </template>
</turbo-stream>
```

## Streaming From HTTP Responses

Turbo は、`text/vnd.turbo-stream.html` の [MIME タイプ](https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types)を宣言する `<form>` 送信への応答として届いた `<turbo-stream>` 要素を、自動的に取り付けることを知っています。[method](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/form#attr-method) 属性が `POST`、`PUT`、`PATCH`、`DELETE` のいずれかに設定された `<form>` 要素を送信するとき、Turbo はリクエストの [Accept](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept) ヘッダのレスポンスフォーマット集合に `text/vnd.turbo-stream.html` を差し込みます。サーバーは [Accept](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept) ヘッダにこの値が含まれるリクエストに応答するとき、Turbo Streams、HTTP リダイレクト、あるいは stream をサポートしない他のクライアント (ネイティブアプリケーションなど) に対して、それぞれ適したレスポンスを返せます。

Rails のコントローラでは、次のようになります。

```ruby
def destroy
  @message = Message.find(params[:id])
  @message.destroy

  respond_to do |format|
    format.turbo_stream { render turbo_stream: turbo_stream.remove(@message) }
    format.html         { redirect_to messages_url }
  end
end
```

デフォルトでは、Turbo はリンクの送信、あるいは method タイプが `GET` のフォームの送信では `text/vnd.turbo-stream.html` MIME タイプを付けません。アプリケーションで `GET` リクエストに Turbo Streams レスポンスを使いたい場合は、リンクやフォームに `data-turbo-stream` 属性を追加することで、Turbo にその MIME タイプを含めるよう指示できます。

## Reusing Server-Side Templates

Turbo Streams の要は、既存のサーバーサイドテンプレートを使い回して、ライブで部分的なページ変更を行えることにあります。リスト内の各メッセージを初回ロード時にレンダリングするための HTML テンプレートは、後から動的に新しいメッセージを 1 件そのリストに追加するときに使うテンプレートと同じです。これこそが HTML-over-the-wire アプローチの本質です。新しいメッセージを JSON にシリアライズし、JavaScript で受け取り、クライアントサイドテンプレートでレンダリングする、といったことは必要ありません。標準的なサーバーサイドテンプレートをそのまま使い回すだけです。

Rails でこれがどう見えるかの、もう 1 つの例を示します。

```erb
<!-- app/views/messages/_message.html.erb -->
<div id="<%= dom_id message %>">
  <%= message.content %>
</div>

<!-- app/views/messages/index.html.erb -->
<h1>All the messages</h1>
<%= render partial: "messages/message", collection: @messages %>
```

```ruby
# app/controllers/messages_controller.rb
class MessagesController < ApplicationController
  def index
    @messages = Message.all
  end

  def create
    message = Message.create!(params.require(:message).permit(:content))

    respond_to do |format|
      format.turbo_stream do
        render turbo_stream: turbo_stream.append(:messages, partial: "messages/message",
          locals: { message: message })
      end

      format.html { redirect_to messages_url }
    end
  end
end
```

新しいメッセージを作るフォームが `MessagesController#create` アクションに送信されると、`MessagesController#index` でメッセージ一覧をレンダリングするのに使ったのとまったく同じ partial テンプレートが、turbo-stream アクションのレンダリングに使われます。これは次のようなレスポンスとして返ってきます。

```html
Content-Type: text/vnd.turbo-stream.html; charset=utf-8

<turbo-stream action="append" target="messages">
  <template>
    <div id="message_1">
      The content of the message.
    </div>
  </template>
</turbo-stream>
```

この `messages/message` の partial テンプレートは、その後、編集・更新操作に伴うメッセージの再レンダリングにも使えます。あるいは、WebSocket や SSE 接続経由で他のユーザーが作成した新しいメッセージを供給するのにも使えます。あらゆる用途にわたって同じテンプレートを再利用できることは非常に強力で、モダンで速いアプリケーションを作るのに必要な作業量を減らす鍵となります。

## Progressively Enhance When Necessary

インタラクション設計は、まず Turbo Streams なしで始めるのがよい習慣です。まずは Turbo Streams が使えなかった場合と同じようにアプリケーション全体を動くようにしてから、そこにレベルアップとして重ねていきます。こうしておけば、Turbo Streams が使えないネイティブアプリケーションや他の環境でも必要となるフローで、更新に頼り切ることがなくなります。

同じことは、WebSocket による更新については特に当てはまります。貧弱な接続やサーバーの問題があるとき、WebSocket は簡単に切れることがあります。アプリケーションを WebSocket なしでも動くように設計しておけば、より頑健になります。

## But What About Running JavaScript?

Turbo Streams は、意図的に append、prepend、(insert) before、(insert) after、replace、update、remove、morph、refresh という 9 つのアクションだけに絞っています。これらのアクションが実行されたときに追加の振る舞いを起こしたいなら、[Stimulus](https://stimulus.hotwired.dev) のコントローラで振る舞いを繋ぎ込むべきです。この制約により、Turbo Streams は HTML をワイヤ越しに届けるという本質的な仕事に集中でき、追加のロジックは専用の JavaScript ファイル内に住まわせておけます。

こうした制約を受け入れることで、個々のレスポンスを、再利用できず、アプリを追いにくくする振る舞いのごちゃ混ぜに変えてしまうことを避けられます。Turbo Streams の核心的な利点は、ページの初回レンダリングからその後のすべての更新まで、同じテンプレートを再利用できることにあります。

## Custom Actions

デフォルトでは、Turbo Streams は [`action` 属性に 8 つの値](../reference/streams.md#the-eight-actions)をサポートしています。もしアプリケーションで他の振る舞いをサポートする必要があるなら、`event.detail.render` 関数を上書きできます。

たとえば、デフォルトのアクションを拡張して `[action="alert"]` や `[action="log"]` を持つ `<turbo-stream>` 要素をサポートしたい場合は、`turbo:before-stream-render` のリスナを宣言してカスタムな振る舞いを提供できます。

```javascript
addEventListener("turbo:before-stream-render", ((event) => {
  const fallbackToDefaultActions = event.detail.render

  event.detail.render = function (streamElement) {
    if (streamElement.action == "alert") {
      // ...
    } else if (streamElement.action == "log") {
      // ...
    } else {
      fallbackToDefaultActions(streamElement)
    }
  }
}))
```

`turbo:before-stream-render` イベントをリッスンするのに加えて、アプリケーションは `StreamActions` のプロパティとして直接アクションを宣言することもできます。

```javascript
import { StreamActions } from "@hotwired/turbo"

// <turbo-stream action="log" message="Hello, world"></turbo-stream>
//
StreamActions.log = function () {
  console.log(this.getAttribute("message"))
}
```

## Integration with Server-Side Frameworks

Turbo に含まれる技術のなかで、バックエンドフレームワークとの緊密な統合による恩恵が最も大きく現れるのが Turbo Streams です。公式の Hotwire スイートの一部として、私たちはそのような統合がどのようなものになりうるかのリファレンス実装を [turbo-rails gem](https://github.com/hotwired/turbo-rails) として作っています。この gem は、Rails が Action Cable と Active Job フレームワークを通じて備えている WebSocket と非同期レンダリングの両方の組み込みサポートに依存しています。

Active Record にミックスインされる [Broadcastable](https://github.com/hotwired/turbo-rails/blob/main/app/models/concerns/turbo/broadcastable.rb) concern を使えば、ドメインモデルから直接 WebSocket の更新をトリガーできます。また、[Turbo::Streams::TagBuilder](https://github.com/hotwired/turbo-rails/blob/main/app/models/turbo/streams/tag_builder.rb) を使えば、インラインのコントローラレスポンスや専用のテンプレート内で `<turbo-stream>` 要素をレンダリングし、8 つのアクションをシンプルな DSL 経由で対応するレンダリングとともに呼び出せます。

とはいえ、Turbo 自体は完全にバックエンドに依存しません。ですので他のエコシステムのフレームワークにも、Rails 向けに提供されたリファレンス実装を見て、独自の緊密な統合を作ってもらえたらと思います。

Turbo の `<turbo-stream-source>` カスタム要素は、その `[src]` 属性を通じて stream ソースに接続します。`ws://` や `wss://` の URL で宣言されると、内部の stream ソースは [WebSocket](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket) インスタンスになります。それ以外の場合、接続は [EventSource](https://developer.mozilla.org/en-US/docs/Web/API/EventSource) 経由になります。

要素がドキュメントに接続されると、stream ソースが接続されます。要素が切断されると、stream も切断されます。

ドキュメントの `<head>` は Turbo のナビゲーションをまたいで永続化されるため、`<turbo-stream-source>` はドキュメントの `<body>` 要素の子孫としてマウントすることが重要です。

Turbo が駆動する典型的なページ全体のナビゲーションでは、`<body>` の中身は破棄され、結果として得られるドキュメントで置き換えられます。streaming を必要とする各ページにこの要素が存在するようにするのは、サーバーの責任です。

もう 1 つの選択肢として、任意のバックエンドアプリケーションを Turbo Streams と統合する分かりやすい方法は、[Mercure プロトコル](https://mercure.rocks)に頼ることです。Mercure は、サーバーアプリケーションが接続中のすべてのクライアントに [Server-Sent Events (SSE)](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events) を通じてページ変更をブロードキャストする便利な方法を定義しています。[Mercure を Turbo Streams と組み合わせて使う方法についてはこちら](https://mercure.rocks/docs/ecosystem/hotwire)。

[次へ: Go Native on iOS & Android](06_native.md)
