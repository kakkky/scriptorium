goroutine leakを起こすようなコードを構造的になくせるようなpackage 開発を行いたい。

fire-and-forgetで気づかないがちなやつの問題は、親goroutineが子goroutineのlife cycleを掌握できていないことにあると思う。それをかけないようにしたい。
親 goroutine が子の lifetime / cancel / error / panic を掌握する」という structured concurrency の思想を Go で materializeする、というのが大まかなコンセプト。

まずは、goroutine leakについてと、既存対策の理解を経てこれを開発すべきか(現状で事足りるかもしれないし)を判断したい。

# goroutine leakとは
本来はもう必要ないgoroutineが、終了されずに動き続けていること。
chに値を送信できないまま止まってるとか。

https://future-architect.github.io/articles/20260128a/
https://zenn.dev/magicmoment/articles/goroutine-leak-202312



「goroutineの停止タイミングがわからないのは設計上の問題」


# 制御法
groutineの停止タイミングを制御するやり方は現状
- signalを子goroutineに渡す
- 親goroutineで合流する

の主に二つ

### signalを子goroutineに渡す
- context.Context の cancel
- done チャンネル(contextでいいと思う。いまはctx.Doneがあるから実質この択はないか)

### 親goroutineで合流する
- sync.WaitGroup
- <-doneCh パターン
  - doneCh := make(chan struct{}) を子で close(doneCh) → 親で <-doneCh。
- errgroup.Wait()
  - `g, ctx := errgroup.WithContext(parent)`でctxを生成でき、それを共通して子goroutineで参照すれば、伝播して終了する(１点目の要素もある)
  - Waitもできる。


どれもちゃんと書けば、まあleakは防げそう。だが、
- sync.WaitGroup:
  - wg.Add(1) 忘れ
  - defer wg.Done() 忘れ
  - wg.Wait() 忘れ (一番致命的)
  - goroutine 数を動的に増減すると Add 後に goroutine 起動忘れ
- context.Context:
  - 子関数に ctx を渡し忘れ (引数追加の手間)
  - 子の中で ctx.Done() を観測し忘れ
  - 子の中で context.Background() を勝手に作って親を無視
- errgroup:
  - g.Wait() 忘れ
  - g.Go(...) の中で渡された ctx を観測し忘れ
  - g.Go 以外で素の go xxx() を書ける (回避経路がある)

で、個人的に、以下のように構造的に上のは防げる気がする
- Wait/Addは書かない。returnされるまではblockとする。
- contextを強制的に受け取るシグネチャにする
- 
```go
// scope: 書き忘れ構造的に不可能
return scope.Run(parent, func(s *scope.Scope) error {
    s.Go(work1)
    s.Go(work2)
    return nil
}) // Run が return する時点で全員合流済み
```

要はgoroutineがどこまで生き続けるのかが微妙にわかりずらいので、明示的な見やすいscopeを付与する感じ
- 見やすい、守りやすい
    - 見やすい: ブロックを見れば lifetime が一目で分かる (= 認知負荷の軽減)
    - 守りやすい: ブロックを抜けるまで戻ってこないので、 lifetime 違反が物理的に書けない (= misuse-resistance)
    - 読み手と書き手の両方 にメリットがある
- 「scope」 = lexical scope (字句スコープ) = lifetime
  - 変数のスコープは { } の中、 という言語仕様レベルの直感
  - goroutine の lifetime も同じ { } の中、 と扱えるようにする
  - Kotlin / Rust / Python が並行処理 API でこの結論に到達済み、 Go にもまだ無いだけ


goroutine leakを起こしてしまう原因は、開発者がgoroutineのlifetimeを把握できてないから。で。それはコード上の表現のため。変数スコープのように、goroutineの生きるscopeを{}ブロックで決めてやる。

-----
# structred concurrency
上の自分のやろうとしていることは、`structred concurrency`という概念っぽい。
まずはこれについて理解したい。

> "Notes on structured concurrency, or: Go statement considered harmful"

原文：https://vorpus.org/blog/notes-on-structured-concurrency-or-go-statement-considered-harmful/
日本語：https://qiita.com/gotta_dive_into_python/items/6feb3224a5fa572f1e19

本題はここから：
https://vorpus.org/blog/notes-on-structured-concurrency-or-go-statement-considered-harmful/#go-statement-considered-harmful


### Go statements break abstraction
そもそもこの文章は、golangのことについて述べているように見えるが、
- "go statement"  = 「並行タスクを起動してすぐ呼び元に return する構文・関数」
というだけであり、golangに限ったことを言っているわけではない。

- 関数 (function) はプログラミングにおける最も基本的な抽象化
  - 「呼んだら、 仕事をして、 戻ってくる」 という契約
  - 中身を知らなくても、 名前と引数と戻り値だけで使える
  - 戻ってきた時点で「もう仕事は終わっている」 ことが暗黙の前提

という抽象化。

go文があると、関数の実行は終了したが、実は裏側ではgoroutineが生成されていて別スレッドで処理が実行されてているという状態が成立。暗黙的である。
```go
func loadConfig() Config {
     go watchFileChanges()  // ← 裏で goroutine 起動
     return config           // ← return するが、 watchFileChanges は走り続ける
 }

func main() {
    cfg := loadConfig()
    // 戻ってきた。 でも裏では goroutine が動いてる
    // - いつ終わるのか?
    // - メモリは保持され続ける?
    // - ファイル handle は?
    // - error が起きたらどこに行く?
    // 全部 loadConfig のソースを読まないと分からない
}
```

### Go statements break automatic resource cleanup
- クリーンアップ処理は、裏側のgoroutineがそれをまだ参照していることを知らない
的なことを言ってる

言語別クリーンアップ例：
- Python with: ブロック抜けで __exit__ 呼ぶ → close される
- Go defer: 関数抜けで予約された処理を実行
- C++ デストラクタ (RAII): スコープ抜けでオブジェクト破棄
- Java try-with-resources: try ブロック抜けで close される

ブロックを抜ける時点で、 そのリソースを使ってる仕事は全部終わってる..のか？という感じ。
```go
func processFile(path string) error {
     f, err := os.Open(path)
     if err != nil {
         return err
     }
     defer f.Close() // ← この defer が事故の温床

     go backgroundProcess(f) // ← f を渡して goroutine 起動

     return nil
     // ↓ ここで何が起きるか:
     // 1. defer f.Close() が実行される (f closed)
     // 2. backgroundProcess は f を使い続ける
     // 3. backgroundProcess が f.Read() → 「file already closed」 エラー
     // 4. そのエラーは誰にも届かず goroutine が静かに死ぬ
 }
```

### Go statements break error handling.

まず、go文を呼ぶとcallstackから外れる。
で、call stack が切れる → error 伝播が壊れる
```
// 普通は
main()
   └── handler()
        └── doWork()  ← ここで error 発生
            ↑ 戻り値 or panic で main まで「上に向かって」 伝わる
```

エラーを上に伝播していくができねーやんという。

-----

で、これらの点を解決するためにpyでやったのがnurseryという。

> every time our control splits into multiple concurrent paths, we want to make sure that they join up again.

分裂した並行処理は最後に合流すべきだという考え。


- 明示的なオブジェクト (trio.open_nursery() で生成)
- async with ブロックに紐づける (字句スコープ強制)
- 子タスクは nursery.start_soon(func) で起動
- ブロックを抜けるのは全子タスク完了後

```py
async with trio.open_nursery() as nursery:
    nursery.start_soon(work1)
    nursery.start_soon(work2)
# ← ここに来た時点で work1, work2 は必ず完了済み
```

Smith が強調する nursery の 5 性質

1. 字句スコープ束縛 (lexical scoping)

- 子タスクの lifetime は async with ブロックの中 に閉じる
- 「ブロック抜け = 全完了」 が 言語仕様レベルで保証された不変条件
- 関数を信頼して「return = 仕事完了」 が再び成立する

2. error 伝播の回復

- 子タスクが exception を上げる → nursery が受け取る → async with
ブロック外に通常の exception として伝播
- call stack が再接続される
- 既存の例外処理機構 (try/except) で扱える
- traceback / debugger も働く

3. cancel の自動伝播

- 1 つの子が exception を上げた瞬間、 兄弟タスク全員が自動 cancel
- nursery が cancel されたら → ネストした子 nursery も連鎖 cancel
- → fail-fast がデフォルトで保証される


https://en.wikipedia.org/wiki/Structured_concurrency


-----

正直この構造的並行処理は大賛成だが....。

なぜgolangではまだ採用に至っていないのかが気になる。

java
https://adengineer.internet.gmo/2022/12/13/java-19-structured-concurrency/

# なぜgolangでは構造的並行処理が採用されない？
いくつかの点が考えられる。

- `Less is exponentially more`というポリシー
https://commandcenter.blogspot.com/2012/06/less-is-exponentially-more.html

構造的並行処理というのが一個のパラダイムだとすれば、そのパラダイム級の変更をすぐにupdateにのせるのはgolangの思想として合わないんだと思う。



- そもそもgoという予約語をdeprecatedにできるはずもない


他の言語には構造的並行処理が組まれているのにgoにはそれがないことについて驚く人が多いんだと。
ただ、errgroup(or waitgroup), contextを使えば同じことはできるとある。確かにそう。
https://rednafi.com/go/structured-concurrency/

以下でgoにおける構造的並行処理について議論されていたっぽい。
https://github.com/golang/go/issues/29011

提案としては、
- Go 2 で structured concurrency を採用する提案。 具体案は「go func() が
  opaque な値を返し、 その変数がスコープを抜ける時に runtime が goroutine
  完了を待つ」 形 (Trio / Kotlin の lifetime 束縛を Go 風に実装)

クローズ理由は、
- Ian Lance Taylor が (1) 具体案は既存コード全てを破壊するため
  non-starter、 (2) 「structured concurrency の領域には良いアイデアがあると思うが、
  まず幅広いユースケースを列挙する大きな議論が先、 言語変更の具体提案はその後」
  として 保留扱いで close (拒絶ではなく議論不足を理由に差し戻し)
- → 8 年経った今もその「大きな議論」 が正式 proposal として戻ってきていない、
  という構造的停滞の起点になった issue

> また、Trioフォーラムで はGoの状況について議論が交わされている。NJS氏は、構造化並行処理の利点を過大評価することには慎重で、彼らが調査した実際のGoのバグに関する研究では、並行処理バグの約4分の1 （典型的な競合状態）は構造化並行処理では防げなかっただろうと指摘した。しかし、理解するのが最も難しいバグの中には、標準ライブラリモジュールが予期せぬバックグラウンドゴルーチンを生成するものもあったと指摘した。これは、真にスコープ付き並行処理を備えた言語では起こり得ないことだ。また、GoのWaitGroupAPIの使用におけるあらゆるミスは、構造化並行処理によって容易に防げるように思えるとも述べた。


> 同じ構造パターンはすべて Go で実現可能です。 errgroup、WaitGroup、context、 チャネルから自分で構築するだけです。これにより、ゴルーチンのライフタイムをより細かく制御できますが、同時にバグが発生する可能性も高くなります

書けるけど書き方によってはミスりがちというのがgolangの構造的並行処理だな。
>  Go固有の代表的バグパターン
 - 匿名関数によるデータレース: 親goroutineのローカル変数 (例: ループ変数 i) を子goroutineが共有してしまう (Fig.8)
 - WaitGroupのAdd/Wait順序ミス: AddがWaitより先に呼ばれる保証がない (Fig.9)
 - context誤用: タイムアウト分岐で別contextを作ってしまい古いgoroutineに到達不能 (Fig.6)
 - channelの二重close: panic発生 → sync.Onceで囲って修正 (Fig.10)
 - selectの非決定性: 複数caseが同時に有効なときどれが選ばれるか不定 (Fig.11)
 - timeパッケージ起因: time.NewTimerが内部goroutineを起動し、想定外のタイミングでchannelに送信 (Fig.12)

 https://songlh.github.io/paper/go-study.pdf


また、構造的並行処理から微々にもズレるが、ランタイムのスケジューラの違いについても記事で触れられている。
> But as covered earlier, the concurrency paradigms are fundamentally different. Python and Kotlin’s cooperative runtimes can own the cancellation because they own the scheduling. Go’s preemptive scheduler doesn’t know what your goroutine is doing or when it should stop. That’s why cancellation is your job.

協調的スケジューリング：
```py
import asyncio

  async def work():
      print("working...")
      await asyncio.sleep(10)   # ← await が cancel 注入点 (scheduler が自動)
      print("done")              # cancel されたら到達しない

  async def main():
      task = asyncio.create_task(work())
      await asyncio.sleep(1)
      task.cancel()              # フラグ立てるだけ、 残りは scheduler が処理
      try:
          await task
      except asyncio.CancelledError:
          print("cancelled")

  asyncio.run(main())

# 実行結果:
# working...
# cancelled
```

つまり、task自体がawait句で「ここでとまります」をschedulerに教えている。だから scheduler は「今この task は安全に止まってる」 と知っている。
scheduler が握っているのは 「continuation (続きの情報)」:
- 「await の戻り値として何を渡せば task が再開できるか」 を抽象的に知っている
- 言い換えると、 task の 再開インターフェース を持っている
- このインターフェースは 値を渡す / 例外を投げる の 2 通り選べる
- → cancel 時は 例外を投げる方を選ぶ
- task の except 節がそれを受け取って cleanup

先制的スケジューリング(preemptive)：
```go
package main

 import (
     "context"
     "fmt"
     "time"
 )

 func work(ctx context.Context) {
     fmt.Println("working...")
     select {
     case <-time.After(10 * time.Second):
         fmt.Println("done")
     case <-ctx.Done():          // ← 開発者が書かないと cancel は届かない
         return
     }
 }

 func main() {
     ctx, cancel := context.WithCancel(context.Background())
     go work(ctx)
     time.Sleep(1 * time.Second)
     cancel()                     // ctx.Done() を close するだけ
     time.Sleep(100 * time.Millisecond)
 }

```
schedulerは、(表現が怪しいが自律的に)処理が長いやつとか、io待ちに入ってblock状態のタスクをコンテキストスイッチにより入れ替えたりする。
つまり、taskから通達されるわけではなく半ば強制的に処理を止めたりする。
task は「自分が止まる」 ことを 構文で表現していない。 scheduler は task の現在状態を 機械レベル (PC、 register) でしか把握してない。 「コード上のどこか」ではなく「機械語命令のどこか」 で止めている。

- Python: await 地点を scheduler が「安全な中断点」 と認識、 cancel 時に 自動で例外を注入 → 開発者は何も書かない
- Go: scheduler は goroutine の内部状態を知らない、 任意地点での cancel 注入は危険 → runtime は注入しない
  - → Go では開発者が select { case <-ctx.Done() } を書いて「ここなら止めても安全」 を明示する責任を負う。


というかそもそもtaskの性質が異なる。
- 協調型における task = 「再開可能な計算」 (resumable computation)
  - 状態遷移機械として実装
  - 各 await が「次の状態」 への遷移点
  - scheduler は遷移点で介入できる
- 先制型における task (= Go goroutine) = 「通常の関数実行」 (procedural execution)
  - 機械語の連続実行
  - 中断点は機械が決める、 言語は関与しない
  - scheduler は機械状態を保存/復元するだけ


------

# 3rd partyで構造的並行処理をしてるやつ。
以下の二つあたり。concが圧倒的。

### conc 

https://github.com/sourcegraph/conc
https://zenn.dev/kouichi_itagaki/articles/492a55e26c9d86

### nusery
https://github.com/arunsworld/nursery

# 上二つにどう差別化するか。
concに対しては、レキシカルスコープで縛って「見やすい」点が挙げられそう。

arunsworld/nursery は、goroutine taskの中で子taskを生成できなさそう。動的task spawnができなさそう？
```go
type ConcurrentJob func(context.Context, chan error)

RunConcurrently(
    // Job 1
    func(context.Context, chan error) {
        time.Sleep(time.Millisecond * 10)
        log.Println("Job 1 done...")
    },
    // Job 2
    func(context.Context, chan error) {
        time.Sleep(time.Millisecond * 5)
        log.Println("Job 2 done...")
    },
)
```

やるなら、別scopeができる感じになら。
```go
RunConcurrently(
      func(ctx context.Context, errCh chan error) {
          dowork()              // 順次実行
          RunConcurrently(      // dowork 完了後に発火
              func(ctx, errCh) { jobA() }, // jobA と jobB は
              func(ctx, errCh) { jobB() }, // 並行 (兄弟)
          )
          // ← ここでは jobA, jobB 完了済み
      },
  )
```
つまり、親 task は、 自分の子 task と並行に走れない。
親子タスクを並行で処理しつつ、lifetimeを同じscopeで表現する、ができないということ

------
どうやるかはまだ詳細詰めてないが、一旦適当に実装ポイントあげてもらった。

1. 全体構造

- Scope 構造体は ctx, cancel, sync.WaitGroup, errOnce, err を持つ
- Run(ctx, fn) で context.WithCancel から派生 ctx を作って Scope に持たせる
- defer cancel() で確実に ctx を閉じる
- closure (= fn) 実行 → wg.Wait() で子完了待ち → error 集約して return

2. closure 形 + Wait 構造的強制

- Run シグネチャを func(ctx, fn) error で固定
- ユーザーは Wait を直接呼ばない、 Run の return が暗黙の Wait
- 早期 return / panic でも Run の defer 経由で必ず合流
- 「Wait 忘れ」 が API surface 上から消える

3. s.Go の実装

- 内部で wg.Add(1) → go func() { defer wg.Done(); ... }()
- 渡された関数を func(ctx context.Context) error シグネチャで縛る
- ctx は scope の派生 ctx を自動で渡す → 「ctx 渡し忘れ」 不可能

4. panic recovery

- s.Go 内の goroutine に defer recover() 仕込む
- recover 値を fmt.Errorf("panic: %v\n%s", r, debug.Stack()) で error 化
- error 集約パイプラインに流す
- process kill 防止

5. error 集約 + cancel 伝播

- sync.Once で first error を保存
- first error 発生時に cancel() を呼ぶ → 兄弟 goroutine に伝播
- 兄弟は ctx.Done() を観測して停止 (cooperative)

6. 動的 spawn 対応

- s.Go は Run の closure 内 / 子 goroutine 内 / 再帰どこからでも呼べる
- 重要: wg.Add は wg.Wait 開始後に呼んではダメ という race condition
  - 対策: closure 内の最初の s.Go で wg.Add する前に wg.Wait は呼ばれない
  - 親 closure が return するまでは wg.Wait 開始しない設計
  - 「親 closure → wg.Wait 開始」 のタイミングで残ってる子からの追加は保証される

7. nested scope

- Run 内で別の Run を呼べる
- 内側 Run の引数 ctx を外側 scope の ctx にする → 親 cancel が子 scope に伝播
- 内側 scope は独立した cancel/error 境界、 ただし親の影響を受ける
- 内側 Run が return するまで外側の wg.Wait は完了しない

8. 起源 stack 記録 (本命差別化)

- s.Go 内で runtime.Callers(skip=2, pcs[:8]) で PC を 8 frame 取得
- skip=2 で runtime.Callers 自身と s.Go を飛ばし、 ユーザー呼び出し位置を取る
- PC 配列を Scope 内の map に保存 (goroutine id → PC 配列)
- leak 検知時 / error 報告時に runtime.CallersFrames(pcs).Next() で文字列解決

9. error 集約モード

- Run: first-error-cancels (デフォルト、 errgroup 互換)
- RunCollect: 全 error を集めて []error で返す
- RunRace: 1 つでも成功で兄弟 cancel
- 内部の error 集約戦略をモード別に切り替える

10. scope.Defer (cleanup 順序)

- s.Defer(fn) で「子完了後に実行する cleanup」 を登録
- 内部の cleanup スライスに push
- wg.Wait 完了後 → スライス逆順で実行 → Run return
- 通常の defer より「子完了を待つ」 タイミングが後ろにずれる

11. analyzer (別パッケージ)

- golang.org/x/tools/go/analysis で実装
- ルール例: 「scope.Run 内 (= closure body 内) で go 文を使ったら警告」
- AST 走査 + scope の判定 (= ssa の関数境界追跡)
- staticcheck / golangci-lint プラグイン対応

12. テスト戦略

- 同一 scope 内の s.Go × N で全完了確認
- panic → error 化の確認
- 兄弟 cancel 確認 (1 つ失敗で他停止)
- ctx 観測しない子の hang 確認 (limit timeout 付きテスト)
- nested scope の cancel 伝播確認
- 起源 stack の正確性確認 (runtime.Caller で期待 PC と比較)

13. Go 1.26 profile 統合 (将来)

- runtime/pprof の goroutineleak profile を RunTest(t) 内で取得
- profile から leak goroutine の stack を抽出
- scope 内部の起源 map と突き合わせ
- 「この s.Go で生まれた goroutine が leak した」 を t.Fatal に出す

14. ベンチマーク観点

- s.Go × 1M 回 / sec の overhead
- runtime.Callers 取得コスト (起源報告 on/off の比較)
- nested scope の cancel 伝播コスト
- errgroup / conc との 1:1 比較
