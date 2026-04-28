---
name: zenn-markdown
description: scriptorium の `articles/` `books/` 配下で Zenn 記事・本を新規作成・編集するとき、および既存ファイルの format/syntax チェック (frontmatter 必須項目・Zenn 独自記法の閉じ忘れ・slug 規則・chapters 整合 等) を行うときに使う。Zenn 公式の markdown 記法と frontmatter 仕様を真とする。
---

# Zenn markdown ガイド (執筆 + format/syntax チェック)

公式: https://zenn.dev/zenn/articles/markdown-guide / https://zenn.dev/zenn/articles/zenn-cli-guide

## 適用対象
- `articles/*.md` (記事)
- `books/<slug>/config.yaml` (本のメタ)
- `books/<slug>/*.md` (チャプター)

## 新規作成
ファイルは手書きせず CLI を使う (slug が自動採番される)。
- 記事: `make zenn.article.create` (= `npx zenn new:article`)
- 本: `make zenn.book.create` (= `npx zenn new:book`)
- プレビュー: `make zenn.preview` (= `npx zenn preview`)

## frontmatter

### 記事 (`articles/<slug>.md`)
```yaml
---
title: ""              # 必須
emoji: "🌊"            # 必須・絵文字 1 文字
type: "tech"           # 必須・tech | idea
topics: []             # 任意・最大 5 個
published: false       # 必須・true で公開、false で下書き
published_at: ""       # 任意・"YYYY-MM-DD" or "YYYY-MM-DD hh:mm" (JST)。一度設定したら変更不可
publication_name: ""   # 任意・publication 掲載時のみ
---
```

### 本のメタ (`books/<slug>/config.yaml`)
```yaml
title: ""              # 必須
summary: ""            # 必須・有料本でも公開される紹介文
topics: []             # 任意・最大 5 個
published: true        # 必須
price: 0               # 任意・0 は無料、有料は 200〜5000 の 100 円単位
toc_depth: 2           # 任意・0〜3、デフォルト 2
chapters:              # 必須・slug を表示順に並べる (最大 100)
  - "intro"
  - "chapter1"
```

### チャプター (`books/<slug>/<chapter-slug>.md`)
```yaml
---
title: ""              # 必須
free: false            # 任意・有料本で一部無料公開する場合 true
---
```

slug の規則: `[a-z0-9_-]` で 12〜50 字 (記事・本)。チャプターは 1〜50 字。

## 共通 markdown 記法 (公式と同じ)
- 見出し: `## 見出し2` から始める (アクセシビリティ)
- リスト・番号付きリスト・引用 (`>`)・区切り線 (`-----`)
- インライン: `*イタリック*` `**太字**` `~~打ち消し~~` `` `code` ``
- HTML コメントは `<!-- ... -->` (出力されない)
- リンク: `[text](URL)`

## Zenn 固有記法

### メッセージ
```
:::message
通常の補足
:::

:::message alert
警告
:::
```

### アコーディオン (details)
```
:::details タイトル
中身
:::
```
ネストするときは `::` を 1 段増やす (`::::details`)。

### リンクカード
URL を行単独で書くと自動でカード化。明示するなら:
```
@[card](https://example.com)
```

### 画像
```
![Alt](URL)
![Alt](URL =250x)         # 幅指定 (px)
*キャプション*             # 画像直下に書くとキャプション扱い
[![Alt](画像URL)](リンクURL)  # 画像をリンクにする
```

### コードブロック
- 言語指定: ` ```js `
- ファイル名表示: ` ```js:src/foo.js `
- diff: ` ```diff js ` (`+` `-` 行が色付け)

### 数式 (KaTeX)
- ブロック: `$$ ... $$` (前後に空行必須)
- インライン: `$a \ne 0$`

### 脚注
```
本文[^1] と書く。

[^1]: 脚注の内容
```
インライン形式 `^[ここに脚注]` も可。

### Mermaid
` ```mermaid ` で図示。1 ブロック 2000 文字以内、chain 10 以下。

### 埋め込み
URL を行単独で書くと自動展開されるもの:
- Twitter / X: `https://twitter.com/<user>/status/<id>`
- YouTube: `https://www.youtube.com/watch?v=...`
- GitHub の単一ファイル / Permalink: `https://github.com/<o>/<r>/blob/<sha>/<path>` (行範囲は `#L1-L3`)

`@[...]` 構文で埋めるもの:
```
@[gist](https://gist.github.com/...)
@[codepen](ページURL)
@[codesandbox](embed 用 URL)
@[stackblitz](embed 用 URL)
@[jsfiddle](ページURL)
@[slideshare](key)
@[speakerdeck](ID)
@[figma](ファイルURL)
@[blueprintue](ページURL)
@[docswell](slide URL)
```

### 絵文字補完
本文中で `:` + 文字を入力すると候補表示 (Zenn エディタ機能)。

## format / syntax チェック観点

既存ファイルをチェックするときも、新規作成時も同じ観点で見る。指摘するときは「ファイルパス + 行番号 + どの規則違反か」をセットで報告する。

### frontmatter
- [ ] 必須項目が全部埋まっている
  - 記事: `title` / `emoji` / `type` / `published`
  - 本 (`config.yaml`): `title` / `summary` / `published` / `chapters`
  - チャプター: `title`
- [ ] `emoji` は 1 文字 (絵文字 1 個)。文字列複数や空文字は NG
- [ ] `type` は `tech` か `idea` のいずれか
- [ ] `topics` の要素数 ≤ 5
- [ ] `published_at` の形式が `YYYY-MM-DD` または `YYYY-MM-DD hh:mm` (JST 前提)
- [ ] `published_at` は **公開後に書き換えていない** (Zenn 側で不変)
- [ ] 本の `price` は 0、または 200〜5000 の 100 円単位
- [ ] 本の `toc_depth` は 0〜3
- [ ] YAML として壊れていない (タブ混入・クォート抜け・インデント崩れがない)

### slug / ファイル名
- [ ] 記事・本の slug: `[a-z0-9_-]` で 12〜50 字
- [ ] チャプター slug: `[a-z0-9_-]` で 1〜50 字
- [ ] 本の `chapters:` 配列の各 slug が、`books/<book-slug>/<chapter-slug>.md` として実在する
- [ ] `chapters:` に存在しないチャプター slug が残っていない (削除漏れ)
- [ ] 逆に、ディレクトリにあるが `chapters:` に載っていない `.md` がない (載せ忘れ)

### markdown 本文
- [ ] 見出しは `##` から始まる (`#` を本文で使っていない)
- [ ] `:::message` / `:::message alert` / `:::details` がすべて `:::` で閉じられている
- [ ] ネストした details で `:` の数が一致している (外側 `::::details` には外側 `::::`)
- [ ] コードブロックの ``` ``` `` ` ``` が対になっている
- [ ] コードブロックに言語指定が付いている (`plain` / `text` でも可)
- [ ] KaTeX ブロック `$$` の前後に空行がある
- [ ] 画像の alt が空でない (アクセシビリティ)
- [ ] 外部 URL でカード化したいものは行単独 or `@[card](URL)`
- [ ] 脚注 `[^x]` の参照と定義が対応している (孤立した参照や未参照定義がない)
- [ ] Mermaid ブロックは 2000 文字以内、chain 10 以下

## してはいけないこと
- `published_at` を **公開後に変更** する
- frontmatter の YAML を壊すような自動整形
- slug 規則外の文字 (大文字・記号) でファイル名を作る
- 本の `chapters` に存在しないチャプター slug を残す
