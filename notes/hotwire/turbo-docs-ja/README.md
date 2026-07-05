# Turbo ドキュメント日本語訳

[Turbo](https://turbo.hotwired.dev/) の `handbook/` 配下と `reference/` 配下を日本語訳した Markdown 集です。

- 原文取得日: 2026-07-05
- 翻訳方針: 地の文のみ日本語化。Turbo 用語 / API 名 / 属性名 / イベント名 / 要素名は原文の英語のまま。コード例は改変せず。
- 内部リンク: 訳ファイル同士の相対パスに書き換え済み。外部リンク (MDN, GitHub 等) はそのまま維持。
- Handbook は「です・ます調」、Reference は「である調」で統一。

## Handbook

1. [Introduction](handbook/01_introduction.md) — Turbo の全体像
2. [Navigate with Turbo Drive](handbook/02_drive.md) — ページ遷移の高速化
3. [Smooth page refreshes with morphing](handbook/03_page_refreshes.md) — morphing によるページリフレッシュ
4. [Decompose with Turbo Frames](handbook/04_frames.md) — Frame でのページ分割
5. [Come Alive with Turbo Streams](handbook/05_streams.md) — Stream によるライブ更新
6. [Go Native on iOS & Android](handbook/06_native.md) — Turbo Native
7. [Building Your Turbo Application](handbook/07_building.md) — Turbo アプリの構築
8. [Installing Turbo in Your Application](handbook/08_installing.md) — Turbo のインストール

## Reference

- [Drive](reference/drive.md) — `Turbo.visit`、`Turbo.session`、`Turbo.config` などの API リファレンス
- [Frames](reference/frames.md) — `<turbo-frame>` 要素と属性
- [Streams](reference/streams.md) — `<turbo-stream>` 要素とアクション
- [Events](reference/events.md) — `turbo:*` イベント一覧
- [Attributes and Meta Tags](reference/attributes.md) — `data-turbo-*` 属性と `<meta name="turbo-*">` タグ

## 原文サイト

- Handbook: https://turbo.hotwired.dev/handbook/introduction
- Reference: https://turbo.hotwired.dev/reference/drive
