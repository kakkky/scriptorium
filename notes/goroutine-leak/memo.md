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
