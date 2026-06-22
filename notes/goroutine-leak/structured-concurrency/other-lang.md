# Structured Concurrency: Kotlin & Swift

## Kotlin Coroutines

### スコープ

**`coroutineScope`**
- 子が1つでも失敗 → 兄弟全員キャンセル → 親に伝播
- キャンセルは双方向（子→親、親→子）

**`supervisorScope`**
- 子が失敗しても兄弟は継続、親に伝播しない
- キャンセルは一方向のみ（親→子）
- ブロック本体自体が例外で終了した場合は全子がキャンセルされる

**`SupervisorJob()`**
- 子の失敗が親をキャンセルしない `Job`
- 落とし穴: `launch(SupervisorJob())` のように子のコンテキストに直接付けると親子関係が切れる。`supervisorScope` の方が安全

### エラー伝播

- `coroutineScope` で複数の子が同時失敗した場合、最初の例外が「主例外」になり残りは `suppressed` として付加される
- `CoroutineExceptionHandler`: `launch` で作ったルートコルーチン専用のlast-resort handler。`async` や子コルーチンには無効

### キャンセル（協調的）

- `job.cancel()` → 次のsuspension pointで `CancellationException` がスロー
- `CancellationException` は `CoroutineExceptionHandler` で無視される（特殊扱い）
- `catch` で捕捉したら必ず再スローしないとキャンセルが止まる
- `withContext(NonCancellable) { }` でキャンセル不可ブロック

明示的チェックAPI:
- `isActive` — Boolチェック
- `ensureActive()` — キャンセル済みなら即スロー
- `yield()` — 他コルーチンに譲りつつキャンセルチェック

### Timeout

- `withTimeout(duration)` — タイムアウト時に `TimeoutCancellationException`（`CancellationException` のサブクラス）をスロー
- `withTimeoutOrNull(duration)` — タイムアウト時は `null` 返却、例外なし

---

## Swift Concurrency

### スコープ

**`withThrowingTaskGroup`（structured, throwing）**
- 子がthrow → 他の子が暗黙的にキャンセル → 親に伝播
- Kotlin の `coroutineScope` に近い

**`withTaskGroup`（structured, non-throwing）**
- 子の失敗が兄弟を自動キャンセルしない
- `group.next()` で個別にハンドリングする
- Kotlin の `supervisorScope` に近い

**`async let`（structured）**
- 変数宣言と同時に子タスク起動、`await` で結果を受け取る
- `await` より前にスコープを抜けると自動キャンセル

**`Task { }` / `Task.detached { }`（unstructured）**

| | `Task { }` | `Task.detached { }` |
|--|--|--|
| アクター継承 | あり | なし |
| 優先度継承 | あり | なし |
| 親キャンセル伝播 | なし | なし |

### キャンセル（協調的）

- `task.cancel()` でフラグが立つだけ（強制停止なし）
- `Task.isCancelled` — Boolチェック
- `Task.checkCancellation()` — キャンセル済みなら `CancellationError` をスロー
- `Task.sleep(for:)` はsuspension pointでキャンセルを受け取り `CancellationError`

**`withTaskCancellationHandler`**
```swift
await withTaskCancellationHandler {
    // 通常処理
} onCancel: {
    urlSession.invalidateAndCancel() // 同期的に即呼ばれる
}
```
- URLSession など、ネイティブキャンセルサポートがないAPIとのブリッジに使う
- すでにキャンセル済みなら operation 開始前に `onCancel` が呼ばれる

### Timeout

ビルトインなし。定石パターン:

```swift
try await withThrowingTaskGroup(of: T.self) { group in
    group.addTask { try await work() }
    group.addTask {
        try await Task.sleep(for: timeout)
        throw TimeoutError()
    }
    defer { group.cancelAll() }
    return try await group.next()!
}
```

---

## 対比

| 機能 | Kotlin | Swift |
|--|--|--|
| 子失敗→兄弟キャンセル | `coroutineScope` | `withThrowingTaskGroup` |
| 子失敗→兄弟継続 | `supervisorScope` | `withTaskGroup` |
| タイムアウト | `withTimeout` / `withTimeoutOrNull`（ビルトイン） | なし（自前実装） |
| キャンセル例外型 | `CancellationException` | `CancellationError` |
| グローバル例外ハンドラ | `CoroutineExceptionHandler` | なし |
| 非構造化 | `GlobalScope.launch` | `Task.detached` |
