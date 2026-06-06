### Go GC の基本構造
- mark-sweep GC の 2 フェーズが完全に分離している理由 (mark 中に sweep するとまだ scan してないポインタが指す live
を誤って解放する危険)
- GC は「sweep → off → mark」の 3 フェーズサイクルを常時回す
- mark は live heap だけを scan する。dead は原理的に追跡できない (誰も指してないから)
- sweep は heap 全体を歩くがビットチェックだけで激安なのでコストは無視できる

### コンパイラとランタイムの役割分担
- コンパイラは「決める人」: escape analysis / 割り当て命令の埋め込み / write barrier 差し込み / stack map 生成
- ランタイムは「動かす人」: mallocgc / GC 実行 / pacer / scheduler / OS への返却
- ローカル変数で宣言された値も、コンパイラが関数全体を見て「外に出る」と判断したら最初から heap 割り当て命令を吐く
- 「stack に置いてから heap に動かす」のではなく最初から heap 行きにする

### スタックとヒープ
- スタック = 関数呼び出しに紐づく自動管理領域、GC 対象外
- ヒープ = 動的割り当て領域、GC 対象
- Go の goroutine スタックは動的拡張される
- escape analysis でどちらに置くかが決まる

 ### コストモデル
- GC が消費するリソースは物理メモリと CPU 時間の 2 つだけ
- メモリコスト = live heap (前回) + new heap (今回)
- live heap はアプリの性質、GC が触れない。GC が制御できるのは new heap の許容量だけ
- CPU コストは scan する量 = live + roots に比例

### GC 頻度とトレードオフ
- 頻度が時間と空間のトレードオフの核心のつまみ
- 頻度を下げる = CPU 効率向上 + メモリ使用増、逆も然り
- 「live heap だけ scan する」設計がこのトレードオフを成立させている前提条件
- 1 サイクルの GC コストは頻度に依存しないので、間隔を伸ばすほど単位時間あたりが薄まる

### GOGC の式
- 正しい式: Target heap = Live + (Live + roots) × GOGC/100
- 「work (作業量)」= live + roots (= 次の GC が scan する対象)
- GOGC = 「work に対して何 % の new heap を許すか」
- roots は heap に含まれないので「合計フットプリント」と「work」は別物
- 設計が美しい点: work が消えて GC オーバーヘッドは GOGC で決まり、プログラム規模に依存しない

### 参照切断と回収のタイムラグ
- 参照を切っても heap には残ったまま
- GC は参照変更を逐次追跡していない
- mark フェーズの点呼で初めて live / dead が確定する
- 「論理的に dead だが物理的に free されていない」期間が必ず存在する

### GOMEMLIMIT
- default は off (= 無制限)
- 設定すると GOGC とは別ルートで GC を発火させる上限ガード
- 両方有効なら低い閾値の方が先に発動する
- GOGC=off + GOMEMLIMIT = 最小 GC 頻度で最大効率モード (ただしメモリは常に上限近く)
- 推奨パターン: 環境支配下 + Go 単独 + メモリ予算固定 (= container Web サービス)
- 罠パターン: 他プロセスと共有 / CLI ツール / 既にメモリギリギリ

### 自分の痛点との接続
- ループで heap が下がらない現象 = heap goal が高く設定されすぎて GC が来ない間に heap が積もる
- cgroup memory limit を heap goal より先に踏むと OOM Killer
- GOMEMLIMIT で「heap goal を待ちきれないときに割り込んで GC 発火」させれば回避できる
- 自分の状況はガイド推奨パターン 1 番に完全一致するので、迷いなく GOMEMLIMIT を使ってよい

### ポインタの扱い
- ポインタを返すコストは 3 段重ね: heap 割り当て + mark scan + cache miss
- 「ポインタで返す」を無心で書くべきではない
- 値返却で済むなら stack に乗れる
- mutation あり / 大きい / nil 必要 / interface 経由 → ポインタが妥当

### GC 発火タイミング (主要 3 つ)
1. heap goal 到達 (= live heap ベースライン)
  - Live + (Live + roots) × GOGC/100 を超えたとき
  - GOGC で制御
  - 平常時の主な発火源
2. GOMEMLIMIT に近づいた
  - メモリ上限の安全網
  - 設定されている場合のみ有効 (default off)
  - heap goal より先に来ることがある
3. 強制 GC (2 分間隔)
  - forcegcperiod = 2 * 60 * 1e9 ナノ秒
  - sysmon goroutine が監視
  - 「何も起きないアプリ」での保険
おまけ (+α)
4. 明示的呼び出し
  - runtime.GC() をコードで呼んだとき
  - 主にテスト・バッチ・finalizer 動作確認用
  - 普段は使わない


### Stop-the-world vs Concurrent GC
- Stop-the-world (STW) GC: GC 実行中、アプリは完全停止。シンプルで速い (スループット高い) が、停止時間 = ヒープサイズに比例 →
ヒープが大きいほどレイテンシ悪化
- Concurrent GC (Go の方式): GC のほとんどをアプリと並行に走らせる。停止時間がヒープサイズに比例しないようにしている
