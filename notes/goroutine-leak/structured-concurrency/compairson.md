
golangで構造的並行処理をやろうとしているpackageはconc, nursery の二つくらい（star 0は除く）。
以下、scopeが自分の作ったものとかのbiasは抜きにして比較してもらった。


# AIによるGo Structured Concurrency パッケージ比較

対象:
- [sourcegraph/conc](https://github.com/sourcegraph/conc)
- [arunsworld/nursery](https://github.com/arunsworld/nursery)
- [kakkky/scope](https://github.com/kakkky/scope)

## 書き方

### sourcegraph/conc

ビルダーパターンで用途別のプリミティブを使い分ける。最もシンプルなAPIは `conc.WaitGroup`。

```go
// WaitGroup: sync.WaitGroup のパニックセーフラッパー。最もシンプル
var wg conc.WaitGroup
wg.Go(func() { /* ... */ })
wg.Wait()

// Pool: 並行数制限などオプションが必要な場合
p := pool.New().WithMaxGoroutines(5)
p.Go(func() { /* ... */ })
p.Wait()

// ErrorPool: エラーを集約するプール
p := pool.New().WithErrors()
p.Go(func() error { return nil })
err := p.Wait()

// ResultPool: 結果を返すプール
p := pool.NewWithResults[int]().WithErrors()
p.Go(func() (int, error) { return 1, nil })
results, err := p.Wait()

// iter: コレクション並列処理
iter.ForEach(items, func(item *Item) { /* ... */ })
```

### arunsworld/nursery

`ConcurrentJob = func(context.Context, chan error)` のスライスを渡す関数型API。エラーはチャネルで送信する。

```go
nursery.RunConcurrently(
    func(ctx context.Context, errCh chan error) {
        // job 1
    },
    func(ctx context.Context, errCh chan error) {
        if err := doSomething(); err != nil {
            errCh <- err
        }
    },
)
```

### kakkky/scope

`Run` の閉じ括弧がgoroutineのライフタイム境界になるlexical scope。

```go
err := scope.Run(context.Background(), func(s *scope.Scope) error {
    s.Go(func(ctx context.Context) error {
        return nil
    })
    s.Go(func(ctx context.Context) error {
        return nil
    })
    return nil
}) // ここで全goroutineが終了していることが保証される
```

ネストしたスコープも作れる（同期的にブロック）。

```go
scope.Run(ctx, func(s *scope.Scope) error {
    s.Scope(func(s *scope.Scope) error {
        s.Go(func(ctx context.Context) error { return nil })
        return nil
    })
    return nil
})
```

---

## 機能比較

| 機能 | conc | nursery | scope |
|---|---|---|---|
| goroutineリーク防止 | △ Wait()必須 | ○ 自動 | ◎ lexical構文で不可能 |
| panicリカバリ | ○ 自動 | × なし | ○ スタックトレース付き |
| 複数エラー集約 | ◎ ErrorPool | × 最初のみ | × 最初のみ |
| 並行数制限 | ○ WithMaxGoroutines | × | × |
| コンテキスト連携 | △ opt-in (WithContext) | ○ 全jobに伝播 | ◎ 自動でcancel/伝播 |
| ネストスコープ | × | △ 非同期ネスト | ○ 同期ネスト |
| Supervisor mode | × | × | ○ WithSupervisor |
| 並列イテレーション | ○ iter.ForEach | × | × |
| First-completion semantics | × | ○ RunUntilFirstCompletion | × |
| 結果の返却 | ○ ResultPool | × | × |
| 動的goroutine追加 | ○ | × 事前定義のみ | ○ goroutine内からGo()可 |

---

## 評価

### sourcegraph/conc

機能が最も豊富。

- `WithMaxGoroutines()` で並行数を制限できる唯一のパッケージ
- `iter.ForEach` / `stream` でコレクション並列処理が簡潔に書ける
- `ResultPool` で結果をまとめて返せる
- ただし `Wait()` 呼び忘れでgoroutineがリークするリスクが残る。構造化並行性の哲学からは外れる

### arunsworld/nursery

APIはシンプルで学習コストが低い。

- `RunUntilFirstCompletion` は他にない独自機能（最初に完了したjobでキャンセル）
- goroutineのライフタイムは自動保証される
- **panicリカバリがない**のが最大の欠点。運用コードでは自前でrecoverを書く必要がある
- errChへの送信強制でコードが冗長になりがち
- errChのバッファが10固定で、大量エラー時にブロックリスクがある

### kakkky/scope

構造化並行性の哲学に最も忠実。

- `Run` の閉じ括弧でgoroutineリークが**構文的に不可能**
- panicをerrorに自動変換（スタックトレース付き）
- `WithSupervisor()` でsibling goroutineへのキャンセル伝播を抑制できる独自設計
- misuse検出: `Run` 返却後に `Go()` を呼ぶとpanicで通知
- イテレーション・結果返却・並行数制限などの高レベルAPIはない

---

## まとめ

| | conc | nursery | scope |
|---|---|---|---|
| 設計思想 | ユーティリティ集合 | Nursery (Python Trio) | Structured Concurrency |
| 向いているケース | 並行数制限・結果集約が必要な場面 | シンプルなfan-outパターン | goroutineリーク撲滅・安全性重視 |



-------

# memo

### concのpoolのWithMaxGoroutineについて

- スロットが空いている → 直接goroutineをspawn
- スロットが埋まっている → 既存workerのtasksチャネルに積む（新しいgoroutineはspawnしない）

```go
type limiter chan struct{}


select {
case p.limiter <- struct{}{}: // スロット取得できたら新しいgoroutineをspawn
    p.handle.Go(...)
case p.tasks <- f:            // 埋まっていたら既存workerのタスクキューに積む
    return
}

// goroutineが終わったら defer <-p.limiter でスロットを解放。
```
