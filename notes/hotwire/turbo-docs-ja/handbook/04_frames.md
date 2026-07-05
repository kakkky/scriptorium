> 原文: https://turbo.hotwired.dev/handbook/frames
> タイトル: Decompose with Turbo Frames
> 最終取得: 2026-07-05

# Decompose with Turbo Frames

Turbo Frames を使うと、ページの決められた部分をリクエストに応じて更新できます。frame の内側にあるリンクやフォームはすべて捕捉され、レスポンスを受け取ったあとに frame の中身が自動で更新されます。サーバーがドキュメント全体を返してきても、更新済みの frame を含む断片だけを返してきても、レスポンスからはその frame だけが抜き出され、既存の中身が置き換えられます。

frame は、ページのあるセグメントを `<turbo-frame>` 要素で囲むことで作ります。各要素には一意な ID が必要で、この ID はサーバーから新しいページをリクエストする際に置き換え対象を照合するために使われます。1 つのページには複数の frame を置くことができ、それぞれが独自のコンテキストを確立します。

```html
<body>
  <div id="navigation">Links targeting the entire page</div>

  <turbo-frame id="message_1">
    <h1>My message title</h1>
    <p>My message content</p>
    <a href="/messages/1/edit">Edit this message</a>
  </turbo-frame>

  <turbo-frame id="comments">
    <div id="comment_1">One comment</div>
    <div id="comment_2">Two comments</div>

    <form action="/messages/comments">...</form>
  </turbo-frame>
</body>
```

このページには 2 つの frame があります。1 つはメッセージ本体を表示し、その編集リンクを含むもの。もう 1 つはコメント一覧を表示し、新しいコメントを追加するフォームを含むものです。それぞれが独自のナビゲーションコンテキストを作り、リンクとフォーム送信の両方を捕捉します。

メッセージ編集リンクがクリックされると、`/messages/1/edit` が返すレスポンスから `<turbo-frame id="message_1">` セグメントが抜き出され、その中身がクリック元の frame を置き換えます。編集レスポンスはたとえば次のようになります。

```html
<body>
  <h1>Editing message</h1>

  <turbo-frame id="message_1">
    <form action="/messages/1">
      <input name="message[name]" type="text" value="My message title">
      <textarea name="message[content]">My message content</textarea>
      <input type="submit">
    </form>
  </turbo-frame>
</body>
```

`<h1>` が `<turbo-frame>` の外にあることに注目してください。これは、編集時にメッセージ表示がフォームに置き換わっても、この `<h1>` はそのまま残るということです。frame が更新される際には、対応する `<turbo-frame>` の内側にあるコンテンツだけが使われます。

このおかげで、ページは 2 つの役割を簡単に兼ねられます。frame 内でインラインに編集することもできれば、frame の外側で、ページ全体をその編集操作に充てるかたちで編集することもできます。

frame には特定の目的があります。ドキュメントの一部分について、コンテンツとナビゲーションを区画化することです。frame の存在は、その子要素に含まれる `<a>` 要素や `<form>` 要素すべてに影響を及ぼすため、必要もないのに導入すべきではありません。Turbo Frames は [Turbo Stream](05_streams.md) の利用をサポートするためのものではありません。もし `<turbo-stream>` 要素のために `<turbo-frame>` 要素を使っているのであれば、`<turbo-frame>` を他の[組み込み要素](https://developer.mozilla.org/en-US/docs/Web/HTML/Element)に変えるべきです。

## Eager-Loading Frames

frame は、それを含むページがロードされたときに必ず埋められている必要はありません。`turbo-frame` タグに `src` 属性が付いていれば、そのタグがページに現れた時点で、参照先の URL が自動的にロードされます。

```html
<body>
  <h1>Imbox</h1>

  <div id="emails">
    ...
  </div>

  <turbo-frame id="set_aside_tray" src="/emails/set_aside">
  </turbo-frame>

  <turbo-frame id="reply_later_tray" src="/emails/reply_later">
  </turbo-frame>
</body>
```

このページはロードと同時に [imbox](http://itsnotatypo.com) にある全メールを表示しますが、そのあとに 2 つのリクエストを追加で発行し、ページ下部に「set aside」しておいたメール用と「後で返信」待ちのメール用の小さなトレイを表示します。これらのトレイは、`src` に指定された URL に対する別々の HTTP リクエストから生成されます。

上の例ではトレイは空の状態で始まりますが、eager-loading frame に初期コンテンツを入れておき、`src` から取得したコンテンツで上書きさせることもできます。

```html
<turbo-frame id="set_aside_tray" src="/emails/set_aside">
  <img src="/icons/spinner.gif">
</turbo-frame>
```

imbox ページをロードすると、set-aside トレイは `/emails/set_aside` からロードされます。そのレスポンスには、最初の例と同じように、対応する `<turbo-frame id="set_aside_tray">` 要素が含まれていなければなりません。

```html
<body>
  <h1>Set Aside Emails</h1>

  <p>These are emails you've set aside</p>

  <turbo-frame id="set_aside_tray">
    <div id="emails">
      <div id="email_1">
        <a href="/emails/1">My important email</a>
      </div>
    </div>
  </turbo-frame>
</body>
```

このページは、imbox ページのトレイ frame に個々のメール入りの `div` だけがロードされる縮小版としても機能しますし、ヘッダーと説明が付いた直接の遷移先としても機能します。ちょうど、メッセージ編集フォームの例と同じです。

`/emails/set_aside` 側の `<turbo-frame>` には `src` 属性が付いていないことに注意してください。この属性は、コンテンツを遅延ロードしたい側の frame にだけ付けます。コンテンツを提供する側でレンダリングされた frame には付けません。

ナビゲーション中、Frame は新しいコンテンツを取得している間、`<turbo-frame>` 要素に `[aria-busy="true"]` を設定します。ナビゲーションが完了すると、Frame は `[aria-busy]` 属性を削除します。`<form>` 送信を通じて `<turbo-frame>` をナビゲートする場合、Turbo はそのフォームの `[aria-busy="true"]` 属性も frame のものと連動して切り替えます。

ナビゲーションが終わると、Frame は `<turbo-frame>` 要素に `[complete]` 属性を設定します。

## Lazy-Loading Frames

ページの初回ロード時に表示されていない frame には `loading="lazy"` を付けておくことで、表示されるまでロードを始めないようにできます。これは `img` 要素の `loading="lazy"` 属性とまったく同じ挙動です。`summary`/`detail` ペアの中やモーダルの中など、最初は隠れていて後から現れるものに含まれる frame のロードを遅延させるのにうってつけです。

## Cache Benefits to Loading Frames

ページのセグメントを frame に切り出すことは、ページの実装をシンプルにするのに役立ちますが、それと同じくらい重要な理由がキャッシュ効率の向上です。多くのセグメントを含む複雑なページは、特に多くのユーザーで共有されるコンテンツと個々のユーザー向けに特化されたコンテンツが混ざっている場合、効率的にキャッシュするのが難しくなります。セグメントが増えれば増えるほど、キャッシュ引きに必要な依存キーは増え、キャッシュはより頻繁に無効化されます。

frame は、変化するタイムスケールや対象オーディエンスが異なるセグメントを分離するのに理想的です。ユーザーごとに変わる部分を frame にすれば、ページの残りの大半を全ユーザーで簡単に共有できることもあります。逆に、大量にパーソナライズされたページで、1 つだけある共有可能なセグメントを frame にして、共有キャッシュから配信することが合理的な場合もあります。

とはいえ、loading frame の取得にかかるオーバーヘッド自体は一般に非常に小さいものですが、それでもいくつロードするかは慎重に判断すべきです。特に、それらの frame がページ上に読み込みジッターを引き起こすような場合には注意が必要です。一方、ロード直後には見えないコンテンツ、たとえばモーダルの背後にあるものやスクロールしないと見えない位置にあるものについては、frame は実質タダのようなものです。

## Targeting Navigation Into or Out of a Frame

デフォルトでは、frame 内でのナビゲーションはその frame だけを対象にします。これはリンクを辿る場合もフォームを送信する場合も同様です。しかし、`target` を `_top` に設定すれば、frame ではなくページ全体を駆動するナビゲーションにできます。あるいは、`target` に別の frame の ID を設定して、名前付きの他の frame を駆動させることもできます。

先ほどの set-aside トレイの例では、トレイ内のリンクは個々のメールを指しています。これらのリンクには、`set_aside_tray` という ID の frame タグを探してほしくはありません。そのメールに直接ナビゲートしてほしいのです。これは、トレイ frame に `target` 属性を付けることで実現します。

```html
<body>
  <h1>Imbox</h1>
  ...
  <turbo-frame id="set_aside_tray" src="/emails/set_aside" target="_top">
  </turbo-frame>
</body>

<body>
  <h1>Set Aside Emails</h1>
  ...
  <turbo-frame id="set_aside_tray" target="_top">
    ...
  </turbo-frame>
</body>
```

大半のリンクは frame のコンテキスト内で動作させたいが、一部だけはそうしたくない、ということもあるでしょう。フォームでも同じです。frame でない要素に `data-turbo-frame` 属性を付けることで、これを制御できます。

```html
<body>
  <turbo-frame id="message_1">
    ...
    <a href="/messages/1/edit">
      Edit this message (within the current frame)
    </a>

    <a href="/messages/1/permission" data-turbo-frame="_top">
      Change permissions (replace the whole page)
    </a>
  </turbo-frame>

  <form action="/messages/1/delete" data-turbo-frame="message_1">
    <a href="/messages/1/warning" data-turbo-frame="_self">
      Load warning within current frame
    </a>

    <input type="submit" value="Delete this message">
    (with a confirmation shown in a specific frame)
  </form>
</body>
```

## Promoting a Frame Navigation to a Page Visit

Frame のナビゲーションは、ドキュメントの残りの状態 (たとえば現在のスクロール位置やフォーカス中の要素) を保ったままページの一部を変更する機会をアプリケーションに提供します。とはいえ、Frame の変更をブラウザの[history](https://developer.mozilla.org/en-US/docs/Web/API/History)にも反映させたいときがあります。

Frame ナビゲーションを Visit に格上げするには、要素をレンダリングする際に `[data-turbo-action]` 属性を付けます。この属性はすべての [Visit](02_drive.md#page-navigation-basics) 値をサポートし、次の要素に指定できます。

- `<turbo-frame>` 要素そのもの

- `<turbo-frame>` をナビゲートする `<a>` 要素

- `<turbo-frame>` をナビゲートする `<form>` 要素

- `<turbo-frame>` をナビゲートする `<form>` 要素に含まれる `<input type="submit">` や `<button>` 要素

たとえば、記事のページネーション付き一覧をレンダリングし、ナビゲーションを [“advance” Action](02_drive.md#application-visits) に変換する Frame を考えてみます。

```html
<turbo-frame id="articles" data-turbo-action="advance">
  <a href="/articles?page=2" rel="next">Next page</a>
</turbo-frame>
```

この `<a rel="next">` 要素をクリックすると、`<turbo-frame>` 要素の `[src]` 属性 *と* ブラウザのパスの *両方* が `/articles?page=2` に設定されます。

**Note:** ブラウザをリフレッシュしたあとにページをレンダリングする際、URL パスや検索パラメータから導かれる状態と一緒に、2 ページ目の記事一覧を描画するのは *アプリケーション側の* 責任です。

## “Breaking out” from a Frame

多くの場合、`<turbo-frame>` を起点としたリクエストは、その frame (もしくは `target` や `data-turbo-frame` 属性の使い方によっては別の frame) のためのコンテンツを取得することが期待されます。つまり、レスポンスには常に期待される `<turbo-frame>` 要素が含まれているべきです。もし Turbo が期待している `<turbo-frame>` 要素がレスポンスに含まれていない場合はエラーとみなされ、その frame に説明的なメッセージを書き込んだ上で例外が投げられます。

一方で特定のケースでは、`<turbo-frame>` リクエストへのレスポンスを frame に閉じ込めず、新しいフルページのナビゲーションとして扱わせたいことがあります。いわゆる「frame から抜け出す」動きです。典型的な例は、セッションが失効したか失われたためにアプリケーションがログインページにリダイレクトする場合です。このとき、Turbo にはエラーとして扱ってもらうよりも、ログインページをそのまま表示してもらったほうが望ましいでしょう。

これを実現する一番シンプルな方法は、[`turbo-visit-control`](../reference/attributes.md#meta-tags) meta タグを含めることで、そのログインページはフルページの再読み込みを要求すると宣言することです。

```html
<head>
  <meta name="turbo-visit-control" content="reload">
  ...
</head>
```

Turbo Rails を使っているのであれば、同じことを `turbo_page_requires_reload` ヘルパーで実現できます。

`turbo-visit-control` が `reload` に指定されたページは、リクエストが frame の内側から発生したものであっても、常にフルページのナビゲーションになります。

missing frame をアプリケーション独自のやり方で扱いたい場合は、[`turbo:frame-missing`](../reference/events.md) イベントを横取りして、たとえばレスポンスを変換したり、別の場所への visit を行ったりできます。

## Anti-Forgery Support (CSRF)

Turbo は、`name` が `csrf-param` または `csrf-token` の `<meta>` タグが DOM に存在するかを確認することで [CSRF](https://en.wikipedia.org/wiki/Cross-site_request_forgery) 対策を提供します。たとえば次のようになります。

```html
<meta name="csrf-token" content="[your-token]">
```

フォーム送信時、トークンは自動的にリクエストヘッダーに `X-CSRF-TOKEN` として追加されます。`data-turbo="false"` を付けて行うリクエストでは、ヘッダーへのトークン追加はスキップされます。

## Custom Rendering

Turbo のデフォルトの `<turbo-frame>` レンダリング処理は、リクエスト元の `<turbo-frame>` 要素の中身を、レスポンス内で対応する `<turbo-frame>` 要素の中身に置き換えます。実質的には、`<turbo-frame>` 要素の中身は [`<turbo-stream action="update">`](../reference/streams.md#update) 要素で処理されるかのようにレンダリングされます。内部のレンダラーは、レスポンス内の `<turbo-frame>` の中身を抜き出し、それを使ってリクエスト元の `<turbo-frame>` 要素の中身を置き換えます。`<turbo-frame>` 要素自体は変更されず、要素のリクエスト・レスポンスライフサイクル全体を通じて [Turbo Drive が管理する `[src]`、`[busy]`、`[complete]` 属性](../reference/frames.md#html-attributes)だけが変化します。

アプリケーションは、`turbo:before-frame-render` イベントリスナーを追加して `event.detail.render` プロパティを上書きすることで、`<turbo-frame>` のレンダリング処理をカスタマイズできます。

たとえば、[morphdom](https://github.com/patrick-steele-idem/morphdom) を使って、レスポンスの `<turbo-frame>` 要素をリクエスト元の `<turbo-frame>` 要素にマージすることもできます。

```javascript
import morphdom from "morphdom"

addEventListener("turbo:before-frame-render", (event) => {
  event.detail.render = (currentElement, newElement) => {
    morphdom(currentElement, newElement, { childrenOnly: true })
  }
})
```

`turbo:before-frame-render` イベントはドキュメントに向かってバブルアップするため、イベントリスナーを個々の `<turbo-frame>` 要素に直接付ければその要素のレンダリングだけを上書きでき、`document` に付ければすべての `<turbo-frame>` 要素のレンダリングを上書きできます。

## Pausing Rendering

アプリケーションはレンダリングを一時停止し、追加の準備を行ってから再開できます。

`turbo:before-frame-render` イベントを購読すると、レンダリングが始まりそうになったタイミングで通知を受け取れます。`event.preventDefault()` を呼べばそこで一時停止できます。準備が終わったら `event.detail.resume()` を呼んでレンダリングを再開します。

具体例としては、退場アニメーションを付けるようなユースケースがあります。

```javascript
document.addEventListener("turbo:before-frame-render", async (event) => {
  event.preventDefault()

  await animateOut()

  event.detail.resume()
})
```

[次へ: Come Alive with Turbo Streams](05_streams.md)
