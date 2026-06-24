# kakkky/scope 機能性・有用性評価

リポジトリ: https://github.com/kakkky/scope  
バージョン: v0.7.0（未安定、v1未到達）  
Go バージョン要件: 1.25以上  
ライセンス: MIT

---

## 何をするパッケージか

goroutine のライフタイムを lexical scope（`Run` の `{}` ブロック）に束縛する。
ブロックを抜けた時点で、`s.Go` で起動した全 goroutine が完了済みであることをランタイムレベルで保証する。

```go
err := scope.Run(ctx, func(s *scope.Scope) error {
    s.Go(func(ctx context.Context) error { return work1(ctx) })
    s.Go(func(ctx context.Context) error { return work2(ctx) })
    return nil
}) // この行を通過した時点で work1, work2 は必ず終了済み
```

---

## API 一覧

### `Run` — エントリーポイント

```go
func Run(ctx context.Context, body func(s *Scope) error, opts ...Option) error
```

- `body` を実行し、`body` 内で起動した全 goroutine が完了するまでブロック
- 最初の non-nil error でスコープの context をキャンセル（デフォルト）
- panic を recover してスタックトレース付き error に変換

### `Scope.Go` — goroutine 起動

```go
func (s *Scope) Go(fn func(ctx context.Context) error)
```

- スコープに束縛した goroutine を起動
- `fn` の引数 `ctx` はスコープの context が自動で渡る（渡し忘れ不可能）
- goroutine 内から再帰的に呼び出し可能（動的 spawn 対応）
- `Run` が返却した後に `Go` を呼ぶと panic

### `Scope.Scope` — 子スコープ

```go
func (s *Scope) Scope(body func(child *Scope) error, opts ...Option) error
```

- 親スコープから派生した子スコープを**同期的に**実行（ブロッキング）
- 親のキャンセルは子に伝播し、子のエラーは親に伝播してキャンセルを引き起こす
- goroutine 内からは呼び出し不可（body または別の Scope から呼ぶ）

### `GoFuture` — 値を返す goroutine

```go
func GoFuture[T any](s *Scope, fn func(ctx context.Context) (T, error)) Future[T]
func (f Future[T]) Wait() (T, error)
```

- 値を返す goroutine を起動し、`Wait()` で結果を遅延取得
- DAG 形式の依存計算（A と B を並列、C は A の結果を使う）に対応

```go
aF := scope.GoFuture(s, func(ctx context.Context) (int, error) { return heavyCalcA(ctx) })
bF := scope.GoFuture(s, func(ctx context.Context) (int, error) { return heavyCalcB(ctx) })
// aF と bF は並列実行
a, _ := aF.Wait()
b, _ := bF.Wait()
```

### Option

| Option | 動作 |
|---|---|
| `WithSupervisor()` | 1つの goroutine の失敗が兄弟の context をキャンセルしない |
| `WithErrAggregation()` | 全エラーを収集して `errors.Join()` で返す |
| `WithMaxConcurrency(n)` | 同時実行 goroutine 数を制限（上限に達すると `Go()` がブロック） |
| `WithTimeout(d)` | 指定期間後に context をキャンセル（`WithSupervisor` 関係なく常に発火） |

なお、Option は子スコープに継承されない。各 `Run`/`Scope` で明示的に指定する。

---

## キャンセル原因の特定

`context.Cause(ctx)` で「なぜキャンセルされたか」を取得できる。

```go
scope.Run(ctx, func(s *scope.Scope) error {
    s.Go(func(ctx context.Context) error {
        return errors.New("network failure")
    })
    s.Go(func(ctx context.Context) error {
        <-ctx.Done()
        fmt.Println(context.Cause(ctx)) // "network failure"
        return ctx.Err()
    })
    return nil
})
```

---

## 機能性の評価

### 強み

**goroutine リーク防止が構文的に不可能**  
`wg.Wait()` 忘れ・`defer wg.Done()` 忘れ・`go func()` の fire-and-forget といった典型的なミスが API surface から消える。`Run` の閉じ括弧が物理的な lifetime の境界になる。

**ctx 渡し忘れが不可能**  
`s.Go` の受け取る関数は `func(ctx context.Context) error` 固定。goroutine 内で `context.Background()` を勝手に作って親の context を無視するパターンを設計レベルで防ぐ。

**panic がプロセスを落とさない**  
goroutine 内の panic を自動で recover し、スタックトレース付きの error として Run に返す。

**動的 spawn に対応**  
goroutine 内から `s.Go` を呼べる。fan-out の再帰パターンや、前の処理結果を見てから goroutine 数を決めるようなケースに対応できる。

**`WithSupervisor` による独立した障害分離**  
他のパッケージにない機能。1つの goroutine の失敗を他に波及させたくない（監視タスク群や独立したバックグラウンドジョブ）場合に使える。

### 弱み

- 並列イテレーション（`iter.ForEach` 相当）がない
- v1 未到達のため API が変わる可能性がある
- コミュニティサポートが薄い（star 数・採用実績が少ない）

---

## 向いているプロジェクト・コード

### 向いている

**goroutine を多用するバックエンドサービス**  
HTTP ハンドラや gRPC サービスで複数の外部 API や DB を並行呼び出しするケース。`Run` のブロックでリクエストスコープに goroutine を束縛でき、リクエスト完了時の確実なリソース解放が保証される。

```go
func (h *Handler) GetDashboard(ctx context.Context, id string) (*Dashboard, error) {
    var profile *Profile
    var orders []*Order
    err := scope.Run(ctx, func(s *scope.Scope) error {
        pF := scope.GoFuture(s, func(ctx context.Context) (*Profile, error) {
            return h.profileRepo.Get(ctx, id)
        })
        oF := scope.GoFuture(s, func(ctx context.Context) ([]*Order, error) {
            return h.orderRepo.List(ctx, id)
        })
        var err error
        profile, err = pF.Wait()
        if err != nil {
            return err
        }
        orders, err = oF.Wait()
        return err
    })
    // ...
}
```

**バッチ処理・ジョブワーカー**  
複数ジョブを並行処理しつつ、1つが失敗したら他を止める fail-fast パターン。`WithErrAggregation` で全エラーを収集するモードにも切り替えられる。

**ライブラリ内の内部並行処理**  
ライブラリ関数内で goroutine を起動して return すると「裏で goroutine が動き続ける」暗黙の状態が生まれ、利用者が察知できない。`scope.Run` で包めばライブラリの呼び出し元にその状態が漏れない。

**goroutine リークを構造的に禁止したいチーム・プロジェクト**  
「goroutine は必ず `scope.Run` の中で起動する」をルールにすることで、コードレビューで素の `go func()` を機械的に弾ける。

**DAG 形式の並列計算**  
`GoFuture` の `Wait()` を使えば、依存関係のある並列タスク（A と B を並行、C は A の完了後に開始）を自然に記述できる。

### 向いていない

**スライスの全要素を並行処理したいだけ**  
`for` ループ内で `s.Go` を呼べば代替できるが、`conc` の `iter.ForEach` の方がシンプル。

**first-completion セマンティクス**  
「最初に完了したタスクで他をキャンセルする」パターンは `scope` が直接サポートしていない。`WithSupervisor` + 手動チャネルで近いことはできるが素直でない。

**高スループットなホットパスで最軽量の並行プリミティブが欲しい場合**  
`Run` はスコープ管理のオーバーヘッドを持つ。純粋な速度が必要なら `errgroup` や `conc` の方が選択肢として自然。

---

## 他パッケージとの比較

| 観点 | conc | nursery | scope |
|---|---|---|---|
| リーク防止の強さ | △ Wait()必須 | ○ 自動 | ◎ 構文的に不可能 |
| panic リカバリ | ○ | × | ○ スタックトレース付き |
| ctx 自動伝播 | △ opt-in | ○ | ◎ 強制 |
| 動的 spawn | ○ | × | ○ |
| 並行数制限 | ○ | × | ○ |
| 並列イテレーション | ○ | × | × |
| 値の返却 | ○ ResultPool | × | △ GoFuture |
| Supervisor mode | × | × | ○ |
| 設計思想への忠実度 | △ | ○ | ◎ |

---

## 総評

`kakkky/scope` は Structured Concurrency の思想に最も忠実な Go 実装。
goroutine のライフタイムを lexical scope と一致させるという設計方針が一貫しており、`wg.Wait()` 忘れ・`ctx` 渡し忘れ・fire-and-forget という典型的な Go の並行バグを構造レベルで排除できる。

機能の網羅性では `conc` に劣るが、**goroutine リーク撲滅・コードの可読性・安全性を最優先にするプロジェクト**では有力な選択肢。
v1 未到達のため本番導入には API 固まり待ちが無難だが、パッケージの思想と設計はそのまま採用する価値がある。
