# Stimulus ドキュメント日本語訳

[Stimulus](https://stimulus.hotwired.dev/) の `handbook/` 配下と `reference/` 配下を日本語訳した Markdown 集です。

- 原文取得日: 2026-07-19
- 翻訳方針: 地の文のみ日本語化。Stimulus 用語 / API 名 / 属性名 / イベント名 / 要素名は原文の英語のまま。コード例は改変せず。
- 内部リンク: 訳ファイル同士の相対パスに書き換え済み。外部リンク (MDN, GitHub 等) はそのまま維持。
- Handbook は「です・ます調」、Reference は「である調」で統一。

## Handbook

1. [Introduction](handbook/01_introduction.md) — Stimulus の全体像
2. [The Origin of Stimulus](handbook/02_origin.md) — Stimulus 誕生の経緯と設計思想
3. [Hello, Stimulus](handbook/03_hello_stimulus.md) — 最初の controller を書く
4. [Building Something Real](handbook/04_building_something_real.md) — クリップボードコピー機能で実践
5. [Designing For Resilience](handbook/05_designing_for_resilience.md) — プログレッシブ・エンハンスメント
6. [Managing State](handbook/06_managing_state.md) — value で状態を扱う
7. [Working With External Resources](handbook/07_working_with_external_resources.md) — 外部リソースの扱い
8. [Installing Stimulus in Your Application](handbook/08_installing.md) — Stimulus のインストール

## Reference

- [Controllers](reference/controllers.md) — controller クラスの基本、identifier、スコープ、登録、通信
- [Lifecycle Callbacks](reference/lifecycle-callbacks.md) — `initialize` / `connect` / `disconnect` などのライフサイクル
- [Actions](reference/actions.md) — `data-action` の記述子・オプション・パラメータ
- [Targets](reference/targets.md) — `data-*-target` の定義とアクセスプロパティ
- [Outlets](reference/outlets.md) — 別 controller のインスタンスを CSS セレクタで参照する仕組み
- [Values](reference/values.md) — 型付き data 属性の読み書きと change callback
- [CSS Classes](reference/css-classes.md) — 論理名で CSS クラスを差し替える仕組み
- [Using Typescript](reference/using-typescript.md) — TypeScript で型付けする際のパターン

## 原文サイト

- Handbook: https://stimulus.hotwired.dev/handbook/introduction
- Reference: https://stimulus.hotwired.dev/reference/controllers
