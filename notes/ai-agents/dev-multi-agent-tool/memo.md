マルチエージェントに動くツールをgolangで何か作成したい。
ざっくり、それぞれのエージェントはgoroutineで、また、やり取り（プロトコル自体はおそらくA2Aとかだが）はchennelでいい感じに動くものはできると思う。

また、LLMは一旦無料のオープンモデルを利用することにする。
何にしよう。

- Google AI Studio (Gemini) — 第一推奨
  - 知名度・性能とも最上位。Gemini 3 系が利用可能
  - Gemini 3.1 Flash-Lite: 500 req/day, 250K TPM
  - function calling 対応、コンテキスト長は最大1M トークン級
  - マルチモーダル対応、クレジットカード不要
  - 注意: 無料枠で送ったデータは Google のモデル学習に使われる（商用不可の理由でもあり、商用しないなら問題なし）

- Cerebras — 速度と日次容量の両立が魅力
  - 30 RPM, 60K TPM, 1M tokens/day
  - Llama 3.3 70B / Qwen3 235B / gpt-oss-120b などが無料
  - function calling 対応（OpenAI 互換 API なので base_url 差し替えだけで動く）
  - 推論速度は業界最速クラス（Groq より速い）

- Groq
  - Llama 3.1 8B: 14,400 req/day（リクエスト数では最多）, 6K TPM
  - Llama 3.3 70B: 1,000 req/day, 12K TPM
  - function calling 対応、OpenAI 互換 API
  - TPM が低めなので長文プロンプトはきつい

- OpenRouter（補助枠として）
  - 25+ の無料モデル（Llama 3.3 70B, Hermes 3 405B など）を1つの API キーで切り替えできる
  - 20 req/min, 50 req/day（$10 チャージで 1,000 req/day へ拡張）
  - いろいろなモデルを試したいときに有用

色々みたが、無料が正義なのでGoogle AI Studio (Gemini) を利用。その中でも、Gemini 2.5 Flash-Lite を利用することにする。
https://ai.google.dev/gemini-api/docs/pricing?hl=ja#gemini-2.5-flash-l2.5
- 1,000 リクエスト/日（RPD）
- 15 リクエスト/分（RPM）
- 250,000 トークン/分（TPM）
- コンテキスト長 1M トークン
良さそう。


そして、何を作るかが問題。個人的にはトモコレみたいなやつをマルチAIエージェントにしたら面白そうだと思っている。
聞いてみた。

