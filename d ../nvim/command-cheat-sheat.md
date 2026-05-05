# LazyVim 操作メモ

LazyVim デフォルト + 主要プラグイン (neo-tree / snacks.nvim / which-key / flash.nvim / blink.cmp / gitsigns / trouble) のキーバインドまとめ。

## 表記の凡例

- **キー** 列: 実際に押すキー (`<kbd>` 表記)
- **表記** 列: nvim ドキュメント標準のバッククォート記法
- `<leader>` = <kbd>Space</kbd>
- `<C-x>` = <kbd>Ctrl</kbd>+<kbd>x</kbd>、`<S-x>` = <kbd>Shift</kbd>+<kbd>x</kbd>、`<D-x>` = <kbd>Cmd</kbd>+<kbd>x</kbd>
- `<CR>` = <kbd>Enter</kbd>、`<BS>` = <kbd>Backspace</kbd>、`<Esc>` = <kbd>Esc</kbd>

> 困ったら **<kbd>Space</kbd> を押して止まる** → which-key が候補を表示してくれる。
> 全 keymap は `<leader>sk` (<kbd>Space</kbd>+<kbd>s</kbd>+<kbd>k</kbd>) で検索可能。

---

## 1. Filer (neo-tree)

### 開く / 閉じる

| キー | 表記 | 機能 |
|---|---|---|
| <kbd>Space</kbd>+<kbd>e</kbd> | `<leader>e` | neo-tree を開く（root ベース、toggle） |
| <kbd>q</kbd> | `q` | tree を閉じる |

### tree 内のナビゲーション

| キー | 表記 | 機能 |
|---|---|---|
| <kbd>Enter</kbd> | `<CR>` | ファイルを開く / ディレクトリ展開 |
| <kbd>l</kbd> | `l` | 同上 |
| <kbd>h</kbd> | `h` | ディレクトリを閉じる |
| <kbd>P</kbd> | `P` | プレビュー toggle |
| <kbd>Backspace</kbd> | `<BS>` | 親ディレクトリへ |
| <kbd>.</kbd> | `.` | カーソル位置を root に設定 |

### ファイル操作

| キー | 表記 | 機能 |
|---|---|---|
| <kbd>a</kbd> | `a` | 新規作成（末尾 `/` でディレクトリ） |
| <kbd>A</kbd> | `A` | ディレクトリを新規作成 |
| <kbd>d</kbd> | `d` | 削除 |
| <kbd>r</kbd> | `r` | リネーム |
| <kbd>c</kbd> | `c` | コピー（パス指定プロンプト） |
| <kbd>m</kbd> | `m` | 移動（パス指定プロンプト） |
| <kbd>y</kbd> | `y` | clipboard に copy（neo-tree 内部） |
| <kbd>x</kbd> | `x` | clipboard に cut（neo-tree 内部） |
| <kbd>p</kbd> | `p` | clipboard から paste |

### 表示・検索

| キー | 表記 | 機能 |
|---|---|---|
| <kbd>H</kbd> | `H` | 隠しファイル表示 toggle |
| <kbd>R</kbd> | `R` | refresh |
| <kbd>/</kbd> | `/` | tree 内 filter |
| <kbd>f</kbd> | `f` | filter 確定 |
| <kbd>F</kbd> | `F` | filter 解除 |
| <kbd>?</kbd> | `?` | このキー一覧を表示 |

---

## 2. Yank / Paste / Register

> `y` 単独はオペレータ。押しっぱなしで止まると which-key が **モーション候補** を出す。これがいわゆる "yank メニュー"。`yy` で行単位 yank、`yi"` で `"..."` の中身 yank、のように **モーション/テキストオブジェクトと組み合わせる** のが本質。

### Yank

| キー | 表記 | 機能 |
|---|---|---|
| <kbd>y</kbd>+<kbd>y</kbd> | `yy` | 行を yank |
| <kbd>Y</kbd> | `Y` | 行末まで yank |
| <kbd>y</kbd>+<kbd>w</kbd> | `yw` | 単語末まで yank |
| <kbd>y</kbd>+<kbd>i</kbd>+<kbd>w</kbd> | `yiw` | 単語の内側 yank |
| <kbd>y</kbd>+<kbd>i</kbd>+<kbd>"</kbd> | `yi"` | `"..."` の中身 yank |
| <kbd>y</kbd>+<kbd>a</kbd>+<kbd>"</kbd> | `ya"` | `"..."` 全体 yank |
| <kbd>y</kbd>+<kbd>i</kbd>+<kbd>p</kbd> | `yip` | 段落の内側 yank |
| <kbd>y</kbd>+<kbd>$</kbd> | `y$` | 行末まで yank |

### Paste / Register

| キー | 表記 | 機能 |
|---|---|---|
| <kbd>p</kbd> | `p` | カーソル後に paste |
| <kbd>P</kbd> | `P` | カーソル前に paste |
| <kbd>"</kbd>+<kbd>+</kbd>+<kbd>y</kbd> | `"+y` | システムクリップボードへ yank |
| <kbd>"</kbd>+<kbd>+</kbd>+<kbd>p</kbd> | `"+p` | システムクリップボードから paste |
| <kbd>Space</kbd>+<kbd>s</kbd>+<kbd>"</kbd> | `<leader>s"` | レジスタ picker (snacks)、yank 履歴の確認に便利 |

---

## 3. Search / Picker (snacks.picker)

| キー | 表記 | 機能 |
|---|---|---|
| <kbd>Space</kbd>+<kbd>Space</kbd> | `<leader><space>` | ファイル検索 (root) |
| <kbd>Space</kbd>+<kbd>f</kbd>+<kbd>f</kbd> | `<leader>ff` | ファイル検索 |
| <kbd>Space</kbd>+<kbd>f</kbd>+<kbd>r</kbd> | `<leader>fr` | 最近開いたファイル |
| <kbd>Space</kbd>+<kbd>f</kbd>+<kbd>b</kbd> | `<leader>fb` | バッファ検索 |
| <kbd>Space</kbd>+<kbd>f</kbd>+<kbd>c</kbd> | `<leader>fc` | nvim config を開く |
| <kbd>Space</kbd>+<kbd>/</kbd> | `<leader>/` | プロジェクト全文検索 (live grep) |
| <kbd>Space</kbd>+<kbd>s</kbd>+<kbd>g</kbd> | `<leader>sg` | grep |
| <kbd>Space</kbd>+<kbd>s</kbd>+<kbd>w</kbd> | `<leader>sw` | カーソル下の単語を grep |
| <kbd>Space</kbd>+<kbd>s</kbd>+<kbd>h</kbd> | `<leader>sh` | help タグ |
| <kbd>Space</kbd>+<kbd>s</kbd>+<kbd>k</kbd> | `<leader>sk` | keymap 検索 |
| <kbd>Space</kbd>+<kbd>:</kbd> | `<leader>:` | command 履歴 |
| <kbd>Space</kbd>+<kbd>s</kbd>+<kbd>R</kbd> | `<leader>sR` | resume (前回の picker を再開) |

---

## 4. LSP / Code

| キー | 表記 | 機能 |
|---|---|---|
| <kbd>g</kbd>+<kbd>d</kbd> | `gd` | 定義ジャンプ |
| <kbd>g</kbd>+<kbd>D</kbd> | `gD` | 宣言ジャンプ |
| <kbd>g</kbd>+<kbd>r</kbd> | `gr` | references |
| <kbd>g</kbd>+<kbd>I</kbd> | `gI` | 実装 |
| <kbd>g</kbd>+<kbd>y</kbd> | `gy` | 型定義 |
| <kbd>K</kbd> | `K` | hover docs |
| <kbd>Space</kbd>+<kbd>c</kbd>+<kbd>a</kbd> | `<leader>ca` | code action |
| <kbd>Space</kbd>+<kbd>c</kbd>+<kbd>r</kbd> | `<leader>cr` | rename |
| <kbd>Space</kbd>+<kbd>c</kbd>+<kbd>f</kbd> | `<leader>cf` | format |
| <kbd>Space</kbd>+<kbd>c</kbd>+<kbd>d</kbd> | `<leader>cd` | line diagnostic |
| <kbd>]</kbd>+<kbd>d</kbd> / <kbd>[</kbd>+<kbd>d</kbd> | `]d` / `[d` | 次/前の diagnostic |
| <kbd>]</kbd>+<kbd>e</kbd> / <kbd>[</kbd>+<kbd>e</kbd> | `]e` / `[e` | 次/前の error |
| <kbd>]</kbd>+<kbd>w</kbd> / <kbd>[</kbd>+<kbd>w</kbd> | `]w` / `[w` | 次/前の warning |

---

## 5. Trouble (診断/シンボル一覧)

| キー | 表記 | 機能 |
|---|---|---|
| <kbd>Space</kbd>+<kbd>x</kbd>+<kbd>x</kbd> | `<leader>xx` | diagnostics (workspace) |
| <kbd>Space</kbd>+<kbd>x</kbd>+<kbd>X</kbd> | `<leader>xX` | diagnostics (buffer) |
| <kbd>Space</kbd>+<kbd>x</kbd>+<kbd>L</kbd> | `<leader>xL` | location list |
| <kbd>Space</kbd>+<kbd>x</kbd>+<kbd>Q</kbd> | `<leader>xQ` | quickfix list |
| <kbd>Space</kbd>+<kbd>c</kbd>+<kbd>s</kbd> | `<leader>cs` | document symbols |
| <kbd>Space</kbd>+<kbd>c</kbd>+<kbd>S</kbd> | `<leader>cS` | LSP references/definitions |

---

## 6. Buffer / Window / Tab

### 用語整理 (混乱しがち)

| 名前 | 何か | 表示 | VSCode で言うと |
|---|---|---|---|
| **buffer** | 開いてるファイルの中身 (メモリ上) | 上部の `\|1\|2\|3\|` バー (bufferline) | 「タブ」 |
| **window** | 画面の分割エリア | `<C-w>v` などで作れる枠 | 「split editor」 |
| **tab** | window と buffer を束ねた作業空間 | bufferline 右端 / タブバー | 「ウィンドウ」 |

→ 上部の `|1|2|3|` は **buffer** (= ファイル)。VSCode の「タブ」と思えば OK。

### Buffer (= 開いてるファイルの切替)

| キー | 表記 | 機能 |
|---|---|---|
| <kbd>Shift</kbd>+<kbd>l</kbd> | `<S-l>` | 次の buffer |
| <kbd>Shift</kbd>+<kbd>h</kbd> | `<S-h>` | 前の buffer |
| <kbd>[</kbd>+<kbd>b</kbd> / <kbd>]</kbd>+<kbd>b</kbd> | `[b` / `]b` | 同上 |
| <kbd>Space</kbd>+<kbd>b</kbd>+<kbd>j</kbd> | `<leader>bj` | **buffer pick** (ラベル文字でジャンプ) |
| <kbd>Space</kbd>+<kbd>`</kbd> | `` <leader>` `` | 直前の buffer に戻る (alt-tab 的) |
| <kbd>Space</kbd>+<kbd>b</kbd>+<kbd>d</kbd> | `<leader>bd` | buffer を閉じる |
| <kbd>Space</kbd>+<kbd>b</kbd>+<kbd>o</kbd> | `<leader>bo` | 他の buffer を全閉じ |
| <kbd>Space</kbd>+<kbd>b</kbd>+<kbd>r</kbd> | `<leader>br` | 右側の buffer を閉じる |
| <kbd>Space</kbd>+<kbd>b</kbd>+<kbd>l</kbd> | `<leader>bl` | 左側の buffer を閉じる |
| <kbd>Space</kbd>+<kbd>b</kbd>+<kbd>p</kbd> | `<leader>bp` | buffer を pin (固定) |
| <kbd>Space</kbd>+<kbd>b</kbd>+<kbd>P</kbd> | `<leader>bP` | pin していない buffer を全閉じ |

> **番号で直接ジャンプ** (`<leader>1`〜`<leader>9`) は LazyVim デフォルトでは未定義。
> 代わりに **`<leader>bj`** で各 buffer に出るラベル文字を押すのが標準的な使い方。

### Window (= 画面分割)

| キー | 表記 | 機能 |
|---|---|---|
| <kbd>Ctrl</kbd>+<kbd>w</kbd>+<kbd>s</kbd> | `<C-w>s` | 水平 split |
| <kbd>Ctrl</kbd>+<kbd>w</kbd>+<kbd>v</kbd> | `<C-w>v` | 垂直 split |
| <kbd>Ctrl</kbd>+<kbd>h/j/k/l</kbd> | `<C-h/j/k/l>` | window 移動 |
| <kbd>Ctrl</kbd>+<kbd>↑/↓/←/→</kbd> | `<C-Up/Down/Left/Right>` | window resize |
| <kbd>Space</kbd>+<kbd>w</kbd>+<kbd>d</kbd> | `<leader>wd` | window 削除 |

### Tab (= 別の作業空間、滅多に使わない)

| キー | 表記 | 機能 |
|---|---|---|
| <kbd>Space</kbd>+<kbd>Tab</kbd>+<kbd>Tab</kbd> | `<leader><tab><tab>` | 新規 tab |
| <kbd>Space</kbd>+<kbd>Tab</kbd>+<kbd>d</kbd> | `<leader><tab>d` | tab 閉じる |
| <kbd>Space</kbd>+<kbd>Tab</kbd>+<kbd>]</kbd> | `<leader><tab>]` | 次の tab |
| <kbd>Space</kbd>+<kbd>Tab</kbd>+<kbd>[</kbd> | `<leader><tab>[` | 前の tab |

---

## 7. Git (LazyGit + snacks)

> 役割の住み分け:
> - **LazyGit** → commit / push / rebase などフルの git 操作
> - **snacks.picker** → log / status / diff / blame などを picker で見る

### よく使う

| キー | 表記 | 機能 |
|---|---|---|
| <kbd>Space</kbd>+<kbd>g</kbd>+<kbd>g</kbd> | `<leader>gg` | LazyGit (root dir) |
| <kbd>Space</kbd>+<kbd>g</kbd>+<kbd>l</kbd> | `<leader>gl` | Git Log |
| <kbd>Space</kbd>+<kbd>g</kbd>+<kbd>f</kbd> | `<leader>gf` | 現ファイルの履歴 |
| <kbd>Space</kbd>+<kbd>g</kbd>+<kbd>b</kbd> | `<leader>gb` | line blame |
| <kbd>Space</kbd>+<kbd>g</kbd>+<kbd>s</kbd> | `<leader>gs` | Git Status |
| <kbd>Space</kbd>+<kbd>g</kbd>+<kbd>p</kbd> | `<leader>gp` | GitHub PRs (open) |
| <kbd>Space</kbd>+<kbd>g</kbd>+<kbd>P</kbd> | `<leader>gP` | GitHub PRs (all) |

### diff / log の別バリエーション

| キー | 表記 | 機能 |
|---|---|---|
| <kbd>Space</kbd>+<kbd>g</kbd>+<kbd>L</kbd> | `<leader>gL` | Git Log (cwd) |
| <kbd>Space</kbd>+<kbd>g</kbd>+<kbd>d</kbd> | `<leader>gd` | Git Diff (hunks) |
| <kbd>Space</kbd>+<kbd>g</kbd>+<kbd>D</kbd> | `<leader>gD` | Git Diff (origin) |
| <kbd>Space</kbd>+<kbd>g</kbd>+<kbd>S</kbd> | `<leader>gS` | Git Stash |
| <kbd>Space</kbd>+<kbd>g</kbd>+<kbd>G</kbd> | `<leader>gG` | LazyGit (cwd) |

### browse (GitHub URL)

| キー | 表記 | 機能 |
|---|---|---|
| <kbd>Space</kbd>+<kbd>g</kbd>+<kbd>B</kbd> | `<leader>gB` | Git Browse (ブラウザで開く) |
| <kbd>Space</kbd>+<kbd>g</kbd>+<kbd>Y</kbd> | `<leader>gY` | Git Browse (URL を copy) |

### GitHub (issue)

| キー | 表記 | 機能 |
|---|---|---|
| <kbd>Space</kbd>+<kbd>g</kbd>+<kbd>i</kbd> | `<leader>gi` | GitHub Issues (open) |
| <kbd>Space</kbd>+<kbd>g</kbd>+<kbd>I</kbd> | `<leader>gI` | GitHub Issues (all) |

---

## 8. Flash (motion)

| キー | 表記 | 機能 |
|---|---|---|
| <kbd>s</kbd> | `s` | flash jump (画面内の任意位置へ) |
| <kbd>S</kbd> | `S` | flash treesitter (構文要素へジャンプ) |
| <kbd>r</kbd> | `r` | remote flash (operator-pending) |
| <kbd>R</kbd> | `R` | treesitter search (operator-pending) |

---

## 9. Terminal (snacks.terminal)

| キー | 表記 | 機能 |
|---|---|---|
| <kbd>Ctrl</kbd>+<kbd>/</kbd> | `<C-/>` | terminal toggle |
| <kbd>Space</kbd>+<kbd>f</kbd>+<kbd>t</kbd> | `<leader>ft` | terminal (cwd) |
| <kbd>Space</kbd>+<kbd>f</kbd>+<kbd>T</kbd> | `<leader>fT` | terminal (root) |
| <kbd>Ctrl</kbd>+<kbd>/</kbd> (term 内) | `<C-/>` | terminal を閉じる |

---

## 10. UI / その他

| キー | 表記 | 機能 |
|---|---|---|
| <kbd>Space</kbd>+<kbd>u</kbd>+<kbd>w</kbd> | `<leader>uw` | word wrap toggle |
| <kbd>Space</kbd>+<kbd>u</kbd>+<kbd>s</kbd> | `<leader>us` | spell toggle |
| <kbd>Space</kbd>+<kbd>u</kbd>+<kbd>n</kbd> | `<leader>un` | 全 notification dismiss |
| <kbd>Space</kbd>+<kbd>n</kbd> | `<leader>n` | notification 履歴 |
| <kbd>Space</kbd>+<kbd>l</kbd> | `<leader>l` | Lazy plugin manager |
| <kbd>Space</kbd>+<kbd>c</kbd>+<kbd>m</kbd> | `<leader>cm` | Mason (LSP/formatter manager) |
| <kbd>Space</kbd>+<kbd>q</kbd>+<kbd>q</kbd> | `<leader>qq` | 終了 |

---

## Tips

- **覚えていない時は which-key**: `<leader>` を押して 1 秒待つと候補メニューが出る。`g` や `]`, `[` も同様にプレフィックスとして候補を表示してくれる。
- **正確な keymap を知りたい時**: `<leader>sk` で picker が起動 → 部分一致で検索できる。自分のメモより信頼できる一次情報。
- **プラグインのバージョン**: `lazy-lock.json` 起因で動作が変わるので、挙動が違う時はまず `<leader>l` (Lazy) で update 履歴を確認。
