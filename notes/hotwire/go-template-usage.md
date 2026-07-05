# sample-todo で使っている Go `html/template` 技術まとめ

`sample-todo/` の実装で使っている `html/template` の要素と設計判断を整理する。Hotwire 固有ではなく **Go 標準の話**。

---

## 1. Parse 系 (テンプレートを読み込む)

### `template.ParseFS(fs, files...)`

```go
template.ParseFS(viewsFS, "views/layout.html", "views/index.html", "views/_item.html")
```

- 埋め込み FS (`embed.FS`) から複数ファイルを parse
- **すべての `{{define}}` が 1 つの template set にマージされる**
- ファイル名がテンプレート名になるわけではない (define 名が主)

### `template.Must(...)`

```go
template.Must(template.ParseFS(...))
```

- エラーを返す関数を「エラー時に panic する」形にラップ
- パッケージレベル変数の初期化で使う定番パターン

### `//go:embed views/*.html`

```go
//go:embed views/*.html
var viewsFS embed.FS
```

- Go 1.16+ の embed 機能
- テンプレートを **バイナリに含めてしまう** ので配布が単一実行ファイルで済む
- 実行時にファイルシステムに依存しない

---

## 2. 定義系 (テンプレートを名前で登録)

### `{{define "name"}}...{{end}}`

```html
{{define "todo_item"}}
<turbo-frame id="todo_{{.ID}}">...</turbo-frame>
{{end}}
```

- 名前付きテンプレートを **1 つ登録** する
- 同じ set 内で `{{template "todo_item" .}}` で呼び出せる
- 同じ set 内で同名 define が複数あると **後勝ち** で上書きされる (だから set を分ける)

### `{{block "name" .}}...{{end}}`

```html
{{block "content" .}}{{end}}
```

- `{{define}}` + `{{template}}` の糖衣構文
- 「default 実装を持つ差し替え可能ポイント」を作る
- レイアウトの「ここに page 固有の内容が入る」を表現するのに使う
- 同じ set 内で別ファイルが `{{define "content"}}` すると **上書きされる**

---

## 3. 呼び出し系 (テンプレートを埋め込む)

### `{{template "name" データ}}`

```html
{{template "todo_item" .}}
```

- 別 define を **その場に展開**
- 第 2 引数は呼び出し先の `.` になる
- **第 1 引数はリテラル文字列のみ**。変数は使えない (これが `yield` 相当が要る場面の理由)

---

## 4. データ参照系

### `.` (ドット)

- **現在のスコープのデータ**
- `ExecuteTemplate` に渡した第 3 引数から始まる
- `range` の中では 1 要素に切り替わる
- `template` 呼び出しで別スコープに引き継げる

### `.Field`

```html
<span>{{.Title}}</span>
```

- 構造体フィールドや map キーへのアクセス
- `map[string]any{"Todos": ...}` の場合も `.Todos` で取れる

### `{{range 集合}}...{{end}}`

```html
{{range .Todos}}
  {{template "todo_item" .}}
{{end}}
```

- slice / array / map / channel をループ
- **ループの中では `.` が 1 要素に変わる**
- 空の集合ではブロックが実行されない (`{{else}}` で空時の表示を書ける)

### `{{if 条件}}...{{end}}`

```html
<div class="todo {{if .Done}}done{{end}}">
```

- 条件分岐
- 値が「zero value」なら false 扱い (Go の慣習)
- `{{else}}` / `{{else if}}` あり

---

## 5. Execute 系 (レスポンスに書き出す)

### `t.ExecuteTemplate(w, name, data)`

```go
indexTmpl.ExecuteTemplate(w, "layout", map[string]any{"Todos": store.List()})
```

- **set の中から指定した名前の define を発火**
- そこから `{{block}}` / `{{template}}` で他 define へ連鎖する
- `w` は `io.Writer` (レスポンスライタでよい)

### `ExecuteTemplate` vs `Execute`

- `Execute`: template 自身の名前を実行 (単一ファイル parse で使う)
- `ExecuteTemplate`: **名前を指定して実行** (set に複数 define がある時はこっち)
- サンプルは全部 `ExecuteTemplate`

---

## 6. Set 分割の設計判断

サンプルで **7 個の `*template.Template` を分けている理由**:

```go
indexTmpl        = mustParse("views/layout.html", "views/index.html", "views/_item.html")
aboutTmpl        = mustParse("views/layout.html", "views/about.html")
itemFrameTmpl    = mustParse("views/_item.html")
editFrameTmpl    = mustParse("views/_edit.html")
streamCreateTmpl = mustParse("views/stream_create.html", "views/_item.html")
streamToggleTmpl = mustParse("views/stream_toggle.html", "views/_item.html")
streamDeleteTmpl = mustParse("views/stream_delete.html")
```

- **`content` という名前が `index.html` と `about.html` で衝突する**ので同じ set に入れられない (後勝ちで上書きされる)
- レイアウトが要る用途 (index / about) と要らない用途 (Frame 応答 / Stream 応答) で分ける
- Frame / Stream 応答は **layout を parse しない** = レイアウト HTML が出力されない
- **共通 partial (`_item.html`) は必要な set にそれぞれ parse** して使い回す

### Set 分割の代替案

ページが増えたら `Clone()` を使う方が DRY:

```go
base := template.Must(template.ParseFS(viewsFS, "views/layout.html", "views/_item.html"))

indexTmpl := template.Must(template.Must(base.Clone()).ParseFS(viewsFS, "views/index.html"))
aboutTmpl := template.Must(template.Must(base.Clone()).ParseFS(viewsFS, "views/about.html"))
```

- 共通 layout + partial を 1 回だけ parse
- 各ページは base のコピーに自分の `content` define を足すだけ
- ページが 10 個くらいに増えたときに効いてくる

現サンプルはページ 2 個なので `ParseFS(files...)` を並べる方式で十分。

---

## 7. HTML エスケープ (暗黙で効いている)

- `html/template` は **文脈依存の自動エスケープ** を行う
- `<span>{{.Title}}</span>` — HTML エスケープ
- `<a href="{{.URL}}">` — URL エスケープ
- `<script>var x = {{.Data}}</script>` — JS エスケープ
- サンプルでは意識していないが、`text/template` ではなく `html/template` を使う理由がこれ
- XSS 対策の第一防衛線

---

## 8. サンプルで使っていない機能

参考までに、`html/template` にはあるが今回使っていないもの:

- **FuncMap**: `template.New("").Funcs(...)` でカスタム関数を登録 (今回不要)
- **`with`**: 特定フィールドを `.` にしてブロック実行 (`range` で代用できるので今回未使用)
- **パイプ**: `{{.Name | upper}}` みたいな関数チェーン (Funcs 定義が要る)
- **`Clone`**: base template を複製して各ページで別の define を足す (ページが増えたら使う)
- **`Delims`**: 区切り文字を `{{ }}` から変える (Vue.js 等との衝突回避)

---

## 全体像

```
[埋め込み]
  //go:embed views/*.html → embed.FS

[parse]
  ParseFS(fs, files...) → *template.Template (set)
    ├─ 内部に {{define}} 群を保持
    └─ Must でエラー時 panic

[定義]
  {{define "layout"}}     ← 名前付きテンプレート
  {{block "content" .}}   ← 上書き可能な差し替え点

[呼び出し]
  {{template "foo" .}}    ← partial 呼び出し (リテラルのみ)

[データ制御]
  {{.Field}} {{range}} {{if}} {{else}}

[出力]
  ExecuteTemplate(w, "layout", data) → 連鎖して HTML 完成

[Set 分割]
  用途ごとに *template.Template を分ける
  = define 名の衝突を回避 & 責務を分離
```

---

## Hotwire 側との関係

- `<turbo-frame>` / `<turbo-stream>` はただの HTML 属性 / タグ
- Go template から見ると **文字列を出力しているだけ**
- Hotwire 固有の知識は不要
- **`{{template}}` で partial を使い回す** ことで「一覧の中の 1 行 = Stream で追加する 1 行 = Frame で置換する 1 行」が同じ HTML になる
- これが Hotwire × Go template の相性の良さ (partial 共有で DRY)

---

## 参考

- Go 公式 `html/template`: <https://pkg.go.dev/html/template>
- Go 公式 `text/template` (アクション構文はこちら): <https://pkg.go.dev/text/template>
- 実装: `sample-todo/templates.go`, `sample-todo/views/*.html`
