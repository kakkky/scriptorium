---
title: "第11章 アダプター層 ー プレゼンテーション(プレゼンター・ハンドラー)の実装"
free: false
---

# この章で行うこと
- プレゼンテーションの実装
    - プレゼンター
    - ハンドラー
- swaggoを用いたSwaggerドキュメントの生成

# プレゼンテーションとは
プレゼンテーションは、エンドユーザーとのやり取りを担う層です。  
外部からの入力（リクエスト）をアプリケーションが処理できる形に変換し、アプリケーションの処理結果（レスポンス）を外部に返却する役割を果たします。

クリーンアーキテクチャでいうとインターフェースアダプター層に位置しています。「内側」であるユースケースと、リクエストを受け付けている「外側」であるWebサーバーとの橋渡しを行なっていると言えます。

プレゼンテーション層は **ハンドラー** と **プレゼンター** の二つに分けられ、それぞれ以下のような役割を持ちます。

#### ハンドラー
  - 外部からのリクエストを受け取り、必要な情報を抽出してアプリケーション層へ伝える
  - アプリケーション層の処理結果を受け取り、プレゼンターに渡す

#### プレゼンター
  - ハンドラーから受け取った処理結果を、適切なフォーマット（JSON、HTML など）に変換する
  - ユーザーや外部システムが受け取りやすい形でレスポンスを返す

![](https://storage.googleapis.com/zenn-user-upload/5781bc18a89e-20250129.png)

# 実装していくハンドラー
エンドポイントとともに、対応するハンドラをまとめました。

| HTTPメソッド | エンドポイント       | 説明                        | 対応のハンドラ |
|-------------|----------------------|-----------------------------|----------------|
| POST        | /users                | 新しいユーザーを作成する    | `PostUserHandler` |
| DELETE      | /users/me                | ユーザーを削除する         | `DeleteUserHandler` |
| GET         | /users               | 全てのユーザーを取得する    | `GetUsersHandler` |
| PATCH       | /users/me                | ユーザーのプロフィールを更新する | `UpdateUserHandler` |
| POST        | /tasks               | 新しいタスクを作成する     | `PostTaskHandler` |
| DELETE      | /tasks/{id}          | 指定したIDのタスクを削除する | `DeleteTaskHandler` |
| PATCH       | /tasks/{id}/state          | 指定したIDのタスクの状態を更新する | `UpdateTaskStateHandler` |
| GET         | /tasks/{id}          | 指定したIDのタスクを取得する | `GetTaskHandler` |
| GET         | /tasks               | 全てのタスクを取得する     | `GetTasksHandler` |
| GET         | /users/me/tasks          | ログインしているユーザーに紐づくタスクを取得する | `GetUserTasksHandler` |
| POST        | /login               | ログインを行う             | `LoginHandler` |
| DELETE      | /logout              | ログアウトを行う           | `LogoutHandler` |

ハンドラには、HTTPメソッドとリソースを対応させた命名をしてみました。


それでは、プレゼンター用の関数を実装してから、それを利用する形でハンドラを実装していきます。
なお、プレゼンテーションのテストは、第14章にて「エンドポイントの統合テスト」として扱います。

# プレゼンターの実装
レスポンスを、整形して返す役割を果たすのがプレゼンターの役割です。今回は、JSON形式でレスポンスを返すようにします。

`app/adapter/presentation/presenter`ディレクトリを作成してください。
::::details プレゼンターの実装
処理成功の場合に返す用
```go:app/adapter/presentation/presenter/success.go
package presenter

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse[T any] struct {
	Status int `json:"status"`
	Data   T   `json:"data"`
}

func RespondOK[T any](rw http.ResponseWriter, respBody T) {
	respondJsonSuccess(rw, http.StatusOK, respBody)
}

func RespondCreated[T any](rw http.ResponseWriter, respBody T) {
	respondJsonSuccess(rw, http.StatusCreated, respBody)
}

func RespondNoContent(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusNoContent)
}

func respondJsonSuccess[T any](rw http.ResponseWriter, statusCode int, respBody T) {
	rw.Header().Set("Content-Type", "application/json;charset=utf-8")
	rw.WriteHeader(statusCode)
	jsonResp := SuccessResponse[T]{
		Status: statusCode,
		Data:   respBody,
	}
	if err := json.NewEncoder(rw).Encode(jsonResp); err != nil {
		RespondInternalServerError(rw, err.Error())
	}
}
```
処理失敗（エラー）の場合に返す用
```go:app/adapter/presentation/presenter/failure.go
package presenter

import (
	"encoding/json"
	"net/http"
)

type FailureResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func RespondBadRequest(rw http.ResponseWriter, message string) {
	respondJsonFailure(rw, http.StatusBadRequest, message)
}

func RespondInternalServerError(rw http.ResponseWriter, message string) {
	respondJsonFailure(rw, http.StatusInternalServerError, message)
}

func RespondUnAuthorized(rw http.ResponseWriter, message string) {
	respondJsonFailure(rw, http.StatusUnauthorized, message)
}

func RespondForbidden(rw http.ResponseWriter, message string) {
	respondJsonFailure(rw, http.StatusForbidden, message)
}

func respondJsonFailure(rw http.ResponseWriter, statusCode int, message string) {
	rw.Header().Set("Content-Type", "application/json;charset=utf-8")
	rw.WriteHeader(statusCode)
	jsonResp := FailureResponse{
		Status:  statusCode,
		Message: message,
	}
	json.NewEncoder(rw).Encode(jsonResp)
}
```
::::

ポイントは以下の通りです。

- 成功時と失敗時に分けて考えている
- ジェネリクスを使用している
- ステータスコード別に関数を用意している

### 成功時と失敗時に分けて考えている
成功時と失敗時には返したい内容が大きく異なると考えました。

成功時のレスポンスモデル
```go
type SuccessResponse[T any] struct {
	Status int `json:"status"`
	Data   T   `json:"data"`
}
```
失敗時のレスポンスモデル
```go
type FailureResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
```

成功時・失敗時のどちらにもステータスコードを含めるようにしています。

成功時には、`Data`フィールドを含みます。ここには、各ハンドラが返したいデータモデルの型が入ります。ジェネリクスを使用することによってそれを表現していますが、ここについては後述します。s

失敗時には、成功時に含んだ`Data`のように返したいデータは特にありません。しかしながら、デバッグのためには、「何が原因でコケたのか」を判断するためにエラーメッセージを含めたいです。
そのため、`Message`フィールドを持たせるようにしました。

成功時と失敗時で、実際には以下のようなフォーマットのJSONを返すようにします。
```json
// 成功時
{
    "status": 201,
    "data": {
        "id": "01JK7W0ZRTA88X2C5SCPS0PZS8",
        "email": "fuga@fuga.com",
        "name": "fuga"
    }
}

// 失敗時
{
    "status": 400,
    "message": "errType : ErrAlreadyRegisterd , Message : you have been already registerd"
}

```


### ジェネリクスの使用

#### ジェネリクスについて（ざっくり）
ジェネリクスを使用すると、多様なデータ型で動作する関数や構造体を定義することができます。

以下は、簡単な例です。
```go
package main

import "fmt"


func PrintValue[T any](value T) {
	fmt.Println(value)
}

func main() {
	// 型推論（明示的な型指定なし）
	PrintValue(42)        // int
	PrintValue("Hello")   // string
	PrintValue(3.14)      // float64

	// 明示的な型指定
	PrintValue[string]("setbs")  // 明示的に string を指定
	PrintValue[int](10)  // 明示的に int を指定
	PrintValue[float64](3.14)    // 明示的に float64 を指定
}
```
`T`は型パラメータです。この場合、`[T any]`とあるように、型パラメータ`T`を`any`としているので引数に入る型を全て受け入れます。
`main`関数で呼び出していますが、型を明示的に指定しても書けますし、型推論を使用して書くこともできます。

https://recursionist.io/learn/languages/go/oop/generics
https://speakerdeck.com/syumai/gonogenericswohuo-yong-suru
https://speakerdeck.com/syumai/gonogenericswohuo-yong-suru


#### 今回ジェネリクスを使ったワケ
成功時のプレゼンター関数にのみジェネリクスを使用しています。
成功時に返したいデータモデルはハンドラによって異なってくるため、`Data`フィールドをどんな型でも受け入れられるようにジェネリクスを使いました。

```go
type SuccessResponse[T any] struct {
	Status int `json:"status"`
	Data   T   `json:"data"` //　型パラメータとして宣言したT型
}
```

関数にも型パラメータを持たせています。引数`respBody`に入れるデータ型によって、`SuccessResponse`型の`Data`フィールドの型が決定されます。

```go
func respondJsonSuccess[T any](rw http.ResponseWriter, statusCode int, respBody T) {
    // 省略
	jsonResp := SuccessResponse[T]{
		Status: statusCode,
		Data:   respBody,
	}
	// 省略
}
```

### ステータスコード別に関数を用意している
失敗時の関数を見てみると、プライベート関数`respondJsonFailure()`が定義されていて、それを他のパブリック関数が呼び出しています。
この`respondJsonFailure()`には、各プレゼンター関数における共通処理を記述しています。
```go
func respondJsonFailure(rw http.ResponseWriter, statusCode int, message string) {
	rw.Header().Set("Content-Type", "application/json;charset=utf-8")
	rw.WriteHeader(statusCode)
	jsonResp := FailureResponse{
		Status:  statusCode,
		Message: message,
	}
	json.NewEncoder(rw).Encode(jsonResp)
}
```
レスポンスヘッダにステータスコードや`Content-Type`等を付加した後、失敗時のデータモデル`FailureResponse`に詰め替えたデータを、JSONにエンコードしてレスポンスに書き込みます。

`RespondBadRequest()`や`RespondUnAuthorized()`は、内部で`respondJsonFailure()`を呼び出しています。
```go
func RespondBadRequest(rw http.ResponseWriter, message string) {
	respondJsonFailure(rw, http.StatusBadRequest, message)
}
```

このように用意した`RespondXxxx()`関数を使い分けることで、レスポンスに含めたいステータスコードを直感的に操作できるようにしています。
成功時のプレゼンター関数も概ね同じことを行なっています。

### 今回扱うステータスコード一覧
今回扱うステータスコードを以下にまとめました。

| ステータスコード       | 意味                             |
|----------------------|----------------------------------|
| 200 OK               | 正常に処理され、結果が返された   |
| 201 Created          | 新しいリソースが作成された       |
| 204 No Content       | コンテンツなしで正常終了         |
| 400 Bad Request      | リクエストが不正                 |
| 401 Unauthorized     | 認証が必要または無効             |
| 403 Forbidden        | アクセス禁止                     |
| 500 Internal Server Error | サーバー内部エラー             |

https://developer.mozilla.org/ja/docs/Web/HTTP/Status


# ハンドラーの実装
定義したプレゼンターを利用し、ハンドラを実装していきます。
ハンドラは、ユースケースの処理を呼び出してレスポンスを返す役割を果たします。

第5章にて、ユースケースでは、入出力用のDTO(Data Transfer Object)を定義してドメイン知識が漏れないようにしていました。
そのため、ユースケースとのやりとりはユースケースに定義したDTOを使用します。
![](https://storage.googleapis.com/zenn-user-upload/c6a519f81f5b-20250130.png)

## `net/http`におけるハンドラー周りの解説
今回は、Webフレームワークを使用しません。そのため、ハンドラー定義は標準パッケージである`net/http`に沿ってやっていきます。

`net/http`におけるハンドラー定義の仕方やルーティングへの登録の仕方は決まっているので、それについてここで述べておきます。

なお、ハンドラーのルーティングへの登録は内容として13章の話にはなりますが、併せて取り上げる方が整理しやすいと考えたのでここでまとめています。

#### `http.Handler`インターフェース
`net/http`では、ハンドラーは以下のインターフェースを満たす必要があります。
```go
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```
この`ServeHTTP()`メソッドは、HTTPリクエストを処理し、レスポンスを返します。
HTTPリクエストを表すのが`*Request`型であり、`ResponseWriter`型にレスポンス内容が書き込まれます。
つまり、`ServeHTTP()`メソッドはリクエストを元にレスポンスを返す、まさにハンドラーの挙動を表すメソッドであることがわかります。

第13章においてハンドラをルーティングに登録するには、このインターフェースを満たさなければなりません。
#### `http.Handle()`関数によるハンドラの登録
`http.Handle()`関数によって、`http.Handler`を満たす構造体をルーティングに登録できます。
```go
http.Handle("GET /xxx",XxxHandler)
```
内部では、`register()`メソッドを呼び出して実際に登録しているようです。
```go
func Handle(pattern string, handler Handler) {
    DefaultServeMux.register(pattern, handler)
}
```
`registger()`メソッドのシグネチャを以下に示しておきます。
```go
func (mux *ServeMux) register(pattern string, handler Handler)
```

#### `http.HandlerFunc`関数型
ハンドラーには上で述べたように、`http.Handler`インターフェースを満たす必要があるので、構造体を用意して、`ServeHTTP()`メソッドを定義して....というふうにしなければなりませんでした。
しかし、ハンドラーの定義の仕方には他にもあります。

その一つは、`http.HandlerFunc`関数型によるハンドラー定義です。
```go
// 関数型
type HandlerFunc func(ResponseWriter, *Request)

// この関数型は、ServeHTTPメソッドを持つ
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r)
}
```
`http.HandlerFunc`は、関数`func(ResponseWriter, *Request)`のシグネチャをとる関数型のようです。
そして、`http.HandlerFunc`型には`ServeHTTP()`メソッドが定義されています。

`ServeHTTP()`メソッドを持っていることが、`http.Handler`インターフェースを満たす、つまりハンドラーとして定義できる条件でした。
つまり、`http.HandlerFunc`関数型を利用することによって、`func(ResponseWriter, *Request)`のシグネチャをとる関数をハンドラとして登録できるようになるのです。

以下に例を示します。
```go
// func(ResponseWriter, *Request) のシグネチャをとる関数
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello,world")
}

// ルーティングに登録する
// HelloWorld関数をhttp.HandlerFuncでキャストしてhttp.Handlerを満たすようにしている
http.Handle("GET /helloworld", http.HandlerFunc(HelloWorld))
```

結局、`http.HandlerFunc`関数型は`http.Handler`インターフェースを満たすようにするためのものです。ハンドラーに登録するためには`http.Handler`インターフェースが必要なのですね。

#### `http.HandleFunc()`によるハンドラの登録
`http.Handler`インターフェースを満たすハンドラをルーティング登録するには、`http.Handle()`関数を使用するのでした。

それとは別に、`http.HandleFunc()`関数でハンドラの登録を行うこともできます。

第二引数の`func(ResponseWriter, *Request)`に注目してください。
```go
func HandleFunc(pattern string, handler func(ResponseWriter, *Request))
```
`HandleFunc()`関数内部の処理を見てみます。
```go
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
     DefaultServeMux.register(pattern, HandlerFunc(handler))
}
```

先ほど見たような形ではありませんか？

結局は同じように、`register(pattern, HandlerFunc(handler))`というふうに、`func(ResponseWriter, *Request)`のシグネチャをとる関数を`http.HandlerFunc`関数型でキャストして、`http.Handler`インターフェースを満たすようにしてから、ルーティングに登録しているのがわかります。

先に示した例を以下のように書き換えることができます。
```go
// func(ResponseWriter, *Request) のシグネチャをとる関数
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello,world")
}
// ルーティングに登録
http.HandleFunc("GET /helloworld",HelloWorld())
```

似たような名前がたくさん出てきてややこしいと思いますが、**ハンドラーの基本は`http.Handler`インターフェース**であり、**`http.Handle()`関数によってハンドラを登録**します。
ハンドラの定義を若干簡素化するために、`http.HandlerFunc`関数型や、`http.HandleFunc()`関数があるということを覚えておくといいと思います。

![](https://storage.googleapis.com/zenn-user-upload/cf1e1f0f8971-20250202.png)

以下にまとめておきます。

| 用語                 | 説明 |
|----------------------|--------------------------------------------|
| `http.Handler`      | `ServeHTTP(ResponseWriter, *Request)` メソッドを持つインターフェース。ハンドラーはこれを満たす必要がある。 |
| `http.HandlerFunc`  | `func(ResponseWriter, *Request)` のシグネチャを持つ関数型。`ServeHTTP()` メソッドが定義されており、`http.Handler` を満たす。 |
| `http.Handle()`      | `http.Handler` を満たす構造体をルーティングに登録する。 |
| `http.HandleFunc()` | `func(ResponseWriter, *Request))`のシグネチャをとる関数を`HandlerFunc(handler)` とキャストして `http.Handler` を満たす形で登録する。 |

ハンドラを登録する`http.Handle()` `http.HandleFunc()`については、13章にて再度取り上げます。



参考
https://zenn.dev/hsaki/books/golang-httpserver-internal

::::message
#### ハンドラー構造体とハンドラー関数の使い分け(私見)
ハンドラーを登録するには、大きく２パターンです。

(1). `ServeHTTP()`を持つことで`http.Handler`インターフェースを満たす**ハンドラー構造体**を定義し、`http.Handle()`関数でルーティングに登録する

(2). `func(ResponseWriter, *Request)`のシグネチャをとる**ハンドラー関数**を定義し、`http.HandleFunc()`関数でルーティングに登録する


「`http.HandlerFunc`関数型でキャストして、`http.Handle()`で登録」という方法も取れますが、それなら、その処理を内部で行なってくれている`http.HandleFunc()`関数を使うのがシンプルでしょう。

さて、挙げた２つの使い分けですが、

- **他のコンポーネント（構造体）を利用して複雑なロジックを表現したいなら(1)の方法**
- **シンプルなロジックなら(2)の方法**

というように考えています。

**ハンドラー構造体**であれば、構造体の中にフィールドとして、いくつかの構造体を埋め込むことができます。これなら、クリーンアーキテクチャなどに則り、責務わけされた別コンポーネントを利用する形で処理を記述できます。依存を持たせることができるのですね。

逆に、それ以外であればハンドラー関数を使うでいいのかなと考えます。関数内にシンプルな処理を書くだけなら、(2)の方法でいいのではないでしょうか。

あくまでこれは筆者の私見ですが、いくつかあるハンドラの登録方法について考えてみました。
::::

## 実装
`net/http`パッケージにおけるハンドラー関連の話をなんとなく理解したところで、実装を始めていきます。

### JSONバリデーション
ハンドラーの実装を始める前に、JSONの内容が欠落していたりしないかを検証するバリデーションを実装します。
汎用的処理として、`pkg`ディレクトリにて実装していきます。

以下のパッケージをインストールしてください。
```
make get-pkg name="github.com/go-playground/validator/v10"
```
https://github.com/go-playground/validator

以下のように実装します。
```go:pkg/validation/validator.go
package validation

import "github.com/go-playground/validator/v10"

// シングルトンインスタンスとして用意
var validate *validator.Validate

func NewValidator() *validator.Validate {
	if validate != nil {
		return validate
	}
	// 非ポインタ構造体でもrequiredgタグが有効になる
    // v11からデフォルトになるのでオプションをつけている
	validate = validator.New(validator.WithRequiredStructEnabled())
	return validate
}
```
これにより、以下のように`validator`パッケージに用意された`Struct()`メソッドを利用できます。
```go
validation.NewValidator().Struct(構造体)
```
`Struct()`メソッドは、構造体のフィールドについた`validate:"required"`タグを認識します。そのフィールドが`nil`だった場合は、エラーを返します。

リクエストモデルには`validate:"required"`タグをつけるようにして、不正なリクエストを弾けるようにしましょう。


次からハンドラの実装を始めます。
`app/adapter/presentation/handler`ディレクトリを作成してください。

### ヘルスチェックハンドラ
エンドポイントが通っているかを簡単に確かめるためだけのハンドラを実装しておきます。

::::details ヘルスチェックハンドラのコード
```go:app/adapter/presentation/handler/health/handler.go
package health

import (
	"net/http"

	"github.com/kakkky/app/adapter/presentation/presenter"
)

// @Summary     apiのヘルスチェックを行う
// @Description apiのヘルスチェックを行う。ルーティングが正常に登録されているかを確かめる。
// @Tags        HealthCheck
// @Success     200 {object} presenter.SuccessResponse[healthResponse] "Health check message""
// @Router      /health [get]
func HealthCheckHandler(rw http.ResponseWriter, r *http.Request) {
	resp := healthResponse{
		HealthCheck: "ok",
	}
	presenter.RespondOK(rw, resp)
}
```
レスポンスモデル
```go:app/adapter/presentation/handler/health/health_models.go
package health

type healthResponse struct {
	HealthCheck string `json:"health_check"`
}
```
::::

ただレスポンスを返すだけの処理なので、ハンドラー関数を定義しておきました。
`presenter.RespondOK(rw, resp)`はプレゼンター用の関数です。引数`resp`は`healthResponse`型なので、型推論により型パラメータ`T`は`healthResponse`型となり、実質的に
```go
func RespondOK(rw http.ResponseWriter, respBody health.healthResponse)
```
となっています。

`// @Success`といったアノテーション部分は後で触れます。

### User系
::::details PostUserHandlerの実装
```go:app/adapter/presentation/handler/user/post_user_handler.go
package user

import (
	"encoding/json"
	"net/http"

	"github.com/kakkky/app/adapter/presentation/presenter"
	"github.com/kakkky/app/application/usecase/user"
	"github.com/kakkky/app/domain/errors"
	"github.com/kakkky/pkg/validation"
)

type PostUserHandler struct {
	registerUsecase *user.RegisterUsecase
}

func NewPostUserHandler(registerUsecase *user.RegisterUsecase) *PostUserHandler {
	return &PostUserHandler{
		registerUsecase: registerUsecase,
	}
}

// @Summary     ユーザーの登録
// @Description 新しいユーザーを登録する
// @Tags        User
// @Accept      json
// @Produce     json
// @Param       request body     PostUserRequest                             true "ユーザー登録のための情報"
// @Success     201     {object} presenter.SuccessResponse[PostUserResponse] "登録されたユーザーの情報"
// @Failure     400     {object} presenter.FailureResponse                   "不正なリクエスト"
// @Failure     500     {object} presenter.FailureResponse                   "内部サーバーエラー"
// @Router      /users [post]
func (puh *PostUserHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// jsonをデコード
	var params PostUserRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		presenter.RespondBadRequest(rw, err.Error())
		return
	}
	if err := validation.NewValidator().Struct(&params); err != nil {
		presenter.RespondBadRequest(rw, err.Error())
		return
	}
	// DTOに詰め替える
	input := user.RegisterUsecaseInputDTO{
		Name:     params.Name,
		Email:    params.Email,
		Password: params.Password,
	}
	// ユースケースに渡して実行
	ctx := r.Context()
	output, err := puh.registerUsecase.Run(ctx, input)
	if (err != nil) && errors.IsDomainErr(err) {
		presenter.RespondBadRequest(rw, err.Error())
		return
	}
	if err != nil {
		presenter.RespondInternalServerError(rw, err.Error())
		return
	}
	resp := PostUserResponse{
		ID:    output.ID,
		Email: output.Email,
		Name:  output.Name,
	}
	presenter.RespondCreated(rw, resp)
}

```
リクエストモデル、レスポンスモデルの定義
```go:app/adapter/presentation/handler/user/post_user_models.go
package user

type PostUserRequest struct {
	Email    string `json:"email" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type PostUserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
```
::::

::::details DeleteUserHandlerの実装
```go:app/adapter/presentation/handler/user/delete_user_handler.go
package user

import (
	"net/http"

	"github.com/kakkky/app/adapter/presentation/middleware"
	"github.com/kakkky/app/adapter/presentation/presenter"
	"github.com/kakkky/app/application/usecase/user"
	"github.com/kakkky/app/domain/errors"
)

type DeleteUserHandler struct {
	unregisterUsecase *user.UnregisterUsecase
}

func NewDeleteUserHandler(unregisterUsecase *user.UnregisterUsecase) *DeleteUserHandler {
	return &DeleteUserHandler{
		unregisterUsecase: unregisterUsecase,
	}
}

// @Summary     ユーザーの退会
// @Description ユーザーを退会させ、ユーザー情報を削除する
// @Tags        User
// @Produce     json
// @Security    BearerAuth
// @Success     204
// @Failure     400 {object} presenter.FailureResponse "不正なリクエスト"
// @Failure     500 {object} presenter.FailureResponse "内部サーバーエラー"
// @Router      /users/me [delete]
func (duh *DeleteUserHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// リクエストスコープのコンテキストからuserIdを取得
	userID := middleware.GetUserID(r.Context())
	input := user.UnregisterUsecaseInputDTO{
		ID: userID,
	}
	ctx := r.Context()
	err := duh.unregisterUsecase.Run(ctx, input)
	if (err != nil) && errors.IsDomainErr(err) {
		presenter.RespondBadRequest(rw, err.Error())
		return
	}
	if err != nil {
		presenter.RespondInternalServerError(rw, err.Error())
		return
	}
	presenter.RespondNoContent(rw)
}
```
::::

::::details GetUsersHandlerの実装
```go:app/adapter/presentation/handler/user/get_users_handler.go
package user

import (
	"net/http"

	"github.com/kakkky/app/adapter/presentation/presenter"
	"github.com/kakkky/app/application/usecase/user"
	"github.com/kakkky/app/domain/errors"
)

type GetUsersHandler struct {
	fetchAllUserUsecase *user.FetchUsersUsecase
}

func NewGetUsersHandler(fetchAllUserUsecase *user.FetchUsersUsecase) *GetUsersHandler {
	return &GetUsersHandler{
		fetchAllUserUsecase: fetchAllUserUsecase,
	}
}

// @Summary     全ユーザーを取得する
// @Description 全てのユーザーのID・名前をリストで取得する
// @Tags        User
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} presenter.SuccessResponse[[]GetUsersResponse] "登録されたユーザーの情報"
// @Failure     400 {object} presenter.FailureResponse                     "不正なリクエスト"
// @Failure     500 {object} presenter.FailureResponse                     "内部サーバーエラー"
// @Router      /users [get]
func (guh *GetUsersHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	outputs, err := guh.fetchAllUserUsecase.Run(ctx)
	if (err != nil) && errors.IsDomainErr(err) {
		presenter.RespondBadRequest(rw, err.Error())
		return
	}
	if err != nil {
		presenter.RespondInternalServerError(rw, err.Error())
		return
	}
	resp := make([]GetUsersResponse, 0, len(outputs))
	for _, output := range outputs {
		resp = append(resp, GetUsersResponse{
			ID:   output.ID,
			Name: output.Name,
		})
	}
	presenter.RespondOK(rw, resp)

}
```
レスポンスモデルを定義
```go:app/adapter/presentation/handler/user/get_users_models.go
package user

type GetUsersResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

```
::::

::::details UpdateUserHandlerの実装
```go:app/adapter/presentation/handler/user/update_user_handler.go
package user

import (
	"encoding/json"
	"net/http"

	"github.com/kakkky/app/adapter/presentation/middleware"
	"github.com/kakkky/app/adapter/presentation/presenter"

	"github.com/kakkky/app/application/usecase/user"
	"github.com/kakkky/app/domain/errors"
	"github.com/kakkky/pkg/validation"
)

type UpdateUserHandler struct {
	updateUserUsecase *user.UpdateProfileUsecase
}

func NewUpdateUserHandler(updateUserUsecase *user.UpdateProfileUsecase) *UpdateUserHandler {
	return &UpdateUserHandler{
		updateUserUsecase: updateUserUsecase,
	}
}

// @Summary     ユーザーの更新
// @Description ユーザー情報（名前・メールアドレス）を更新する
// @Tags        User
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body     UpdateUserRequest                             true "ユーザー更新のための情報"
// @Success     200     {object} presenter.SuccessResponse[UpdateUserResponse] "登録されたユーザーの情報"
// @Failure     400     {object} presenter.FailureResponse                     "不正なリクエスト"
// @Failure     500     {object} presenter.FailureResponse                     "内部サーバーエラー"
// @Router      /users/me [patch]
func (uuh *UpdateUserHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	var params UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		presenter.RespondBadRequest(rw, err.Error())
		return
	}
	if err := validation.NewValidator().Struct(&params); err != nil {
		presenter.RespondBadRequest(rw, err.Error())
		return
	}
	input := user.UpdateProfileUsecaseInputDTO{
		ID:    userID,
		Email: params.Email,
		Name:  params.Name,
	}
	ctx := r.Context()
	output, err := uuh.updateUserUsecase.Run(ctx, input)
	if (err != nil) && errors.IsDomainErr(err) {
		presenter.RespondBadRequest(rw, err.Error())
		return
	}
	if err != nil {
		presenter.RespondInternalServerError(rw, err.Error())
		return
	}
	resp := UpdateUserResponse{
		ID:    output.ID,
		Email: output.Email,
		Name:  output.Name,
	}
	presenter.RespondOK(rw, resp)
}
```
リクエストモデル、レスポンスモデルを定義
```go:app/adapter/presentation/handler/user/update_user_models.go
package user

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
```
::::

#### ハンドラ内処理の流れ
大体のハンドラは、以下のような流れです。

1. リクエストボディ（`r.Body`）をデコードしてJSONからリクエストモデルを表現したGoの構造体に変換
   - エラーだとステータスコード400とエラーメッセージを返す(`presenter.RespondBadRequest(rw, err.Error())
`)
2. リクエストモデルにバリデーションをかけて正しいリクエストか検証
   - エラーだとステータスコード400とエラーメッセージを返す(`presenter.RespondBadRequest(rw, err.Error())
`)
3. リクエストモデルからユースケースで定義したInputDTOにデータを詰め替える
4. ユースケースの処理を呼び出す
   - ユースケースの処理から帰ってきたエラーによってレスポンスが異なる
       - ドメインエラー:`presenter.RespondBadRequest(rw, err.Error())`
       - その他のエラー（つまりシステム側のエラー）:`presenter.RespondInternalServerError(rw, err.Error())`
5. ユースケースからの出力（`OutputDTO`）をレスポンスモデルに詰め替えて成功時のプレゼンターに渡す（`presenter.RespondCreated(rw, resp)`）


また、`DeleteUserHandler`と`UpdateUserHandler`については、
```go
// リクエストスコープのコンテキストからuserIdを取得
userID := middleware.GetUserID(r.Context())
```
の記述をしています。まだミドルウェアを実装していないのでエラーになっているかと思いますが、次章で行いますのでそのままで大丈夫です。
この記述に関してもそこで触れることにします。

### Task系
::::details PostTaskHandlerの実装
```go:app/adapter/presentation/handler/task/post_task_handler.go
package task

import (
	"encoding/json"
	"net/http"

	"github.com/kakkky/app/adapter/presentation/middleware"
	"github.com/kakkky/app/adapter/presentation/presenter"
	"github.com/kakkky/app/application/usecase/task"
	"github.com/kakkky/app/domain/errors"
	"github.com/kakkky/pkg/validation"
)

type PostTaskHandler struct {
	createPostUsecase *task.CreateTaskUsecase
}

func NewPostTaskHandler(createPostUsecase *task.CreateTaskUsecase) *PostTaskHandler {
	return &PostTaskHandler{
		createPostUsecase: createPostUsecase,
	}
}

// @Summary     タスクを作成する
// @Description 内容、タスク状態からユーザーに紐づくタスクを作成する
// @Tags        Task
// @Produce     json
// @Param       request body PostTaskRequest true "タスク作成のための情報"
// @Security    BearerAuth
// @Success     201 {object} presenter.SuccessResponse[PostTaskResponse] "作成したタスクの情報"
// @Failure     400 {object} presenter.FailureResponse                   "不正なリクエスト"
// @Failure     500 {object} presenter.FailureResponse                   "内部サーバーエラー"
// @Router      /tasks [post]
func (pth *PostTaskHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// jsonをデコード
	var params PostTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		presenter.RespondBadRequest(rw, err.Error())
		return
	}
	if err := validation.NewValidator().Struct(&params); err != nil {
		presenter.RespondBadRequest(rw, err.Error())
		return
	}
	// contextからuserIdを取得
	userID := middleware.GetUserID(r.Context())
	// inputDTOに詰め替える
	input := task.CreateTaskUsecaseInputDTO{
		UserId:  userID,
		Content: params.Content,
		State:   params.State,
	}
	output, err := pth.createPostUsecase.Run(r.Context(), input)
	if err != nil && errors.IsDomainErr(err) {
		presenter.RespondBadRequest(rw, err.Error())
		return
	}
	if err != nil {
		presenter.RespondInternalServerError(rw, err.Error())
		return
	}
	resp := PostTaskResponse{
		ID:      output.ID,
		UserId:  output.UserId,
		Content: output.Content,
		State:   output.State,
	}
	presenter.RespondCreated(rw, resp)
}
```
リクエストモデル・レスポンスモデルを定義
```go:app/adapter/presentation/handler/task/post_task_models.go
package task

type PostTaskRequest struct {
	Content string `json:"content" validate:"required"`
	State   string `json:"state" validate:"required"`
}

type PostTaskResponse struct {
	ID      string `json:"id"`
	UserId  string `json:"user_id" `
	Content string `json:"content"`
	State   string `json:"state"`
}
```
::::

::::details DeleteTaskHandlerの実装
```go:app/adapter/presentation/handler/task/delete_task_handler.go
package task

import (
	"net/http"

	"github.com/kakkky/app/adapter/presentation/middleware"
	"github.com/kakkky/app/adapter/presentation/presenter"
	"github.com/kakkky/app/application/usecase/task"
	"github.com/kakkky/app/domain/errors"
)

type DeleteTaskHandler struct {
	deleteTaskUsease *task.DeleteTaskUsecase
}

func NewDeleteTaskHandler(deleteTaskUsease *task.DeleteTaskUsecase) *DeleteTaskHandler {
	return &DeleteTaskHandler{
		deleteTaskUsease: deleteTaskUsease,
	}
}

// @Summary     タスクを削除する
// @Description 指定したidのタスクを削除する
// @Tags        Task
// @Produce     json
// @Security    BearerAuth
// @Success     204
// @Failure     400 {object} presenter.FailureResponse "不正なリクエスト"
// @Failure     403 {object} presenter.FailureResponse "権限エラー"
// @Failure     500 {object} presenter.FailureResponse "内部サーバーエラー"
// @Router      /tasks/{id} [delete]
func (dth *DeleteTaskHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// パスパラメータから取得
	id := r.PathValue("id")
	// contextからuserIdを取得
	userID := middleware.GetUserID(r.Context())
	// inputDTOに詰め替える
	input := task.DeleteTaskUsecaseInputDTO{
		ID:     id,
		UserId: userID,
	}
	err := dth.deleteTaskUsease.Run(r.Context(), input)
	// タスクを削除する権限がない（ログインしているユーザーのタスクでない）場合
	if err != nil && errors.Is(err, errors.ErrForbiddenTaskOperation) {
		presenter.RespondForbidden(rw, err.Error())
		return
	}
	// ドメインエラー
	if err != nil && errors.IsDomainErr(err) {
		presenter.RespondBadRequest(rw, err.Error())
		return
	}
	// その他エラー
	if err != nil {
		presenter.RespondInternalServerError(rw, err.Error())
		return
	}
	presenter.RespondNoContent(rw)
}
```
::::

::::details UpdateTaskStateHandlerの実装
```go:app/adapter/presentation/handler/task/update_task_state_handler.go
package task

import (
	"encoding/json"
	"net/http"

	"github.com/kakkky/app/adapter/presentation/middleware"
	"github.com/kakkky/app/adapter/presentation/presenter"
	"github.com/kakkky/app/application/usecase/task"
	"github.com/kakkky/app/domain/errors"
	"github.com/kakkky/pkg/validation"
)

type UpdateTaskStateHandler struct {
	updateTaskStateUsecase *task.UpdateTaskStateUsecase
}

func NewUpdateTaskStateHandler(updateTaskStateUsecase *task.UpdateTaskStateUsecase) *UpdateTaskStateHandler {
	return &UpdateTaskStateHandler{
		updateTaskStateUsecase: updateTaskStateUsecase,
	}
}

// @Summary     タスク状態を更新する
// @Description タスクの状態(todo/doing/done)を 指定して更新する
// @Tags        Task
// @Produce     json
// @Param       request body UpdateTaskStateRequest true "タスク更新のための情報"
// @Security    BearerAuth
// @Success     201 {object} presenter.SuccessResponse[UpdateTaskStateResponse] "更新したタスクの情報"
// @Failure     400 {object} presenter.FailureResponse                          "不正なリクエスト"
// @Failure     403 {object} presenter.FailureResponse                          "権限エラー"
// @Failure     500 {object} presenter.FailureResponse                          "内部サーバーエラー"
// @Router      /tasks/{id} [patch]
func (utsh *UpdateTaskStateHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// jsonをデコード
	var params UpdateTaskStateRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		presenter.RespondBadRequest(rw, err.Error())
		return
	}
	if err := validation.NewValidator().Struct(&params); err != nil {
		presenter.RespondBadRequest(rw, err.Error())
		return
	}
	// idをパスパラメータから取得
	id := r.PathValue("id")
	// contextからuserIdを取得
	userID := middleware.GetUserID(r.Context())
	// inputDTOに詰め替える
	input := task.UpdateTaskStateUsecaseInputDTO{
		ID:     id,
		UserId: userID,
		State:  params.State,
	}
	output, err := utsh.updateTaskStateUsecase.Run(r.Context(), input)
	// タスクを削除する権限がない（ログインしているユーザーのタスクでない）場合
	if err != nil && errors.Is(err, errors.ErrForbiddenTaskOperation) {
		presenter.RespondForbidden(rw, err.Error())
		return
	}
	// ドメインエラー
	if err != nil && errors.IsDomainErr(err) {
		presenter.RespondBadRequest(rw, err.Error())
		return
	}
	// その他エラー
	if err != nil {
		presenter.RespondInternalServerError(rw, err.Error())
		return
	}
	resp := UpdateTaskStateResponse{
		ID:      output.ID,
		UserId:  output.UserId,
		Content: output.Content,
		State:   output.State,
	}
	presenter.RespondOK(rw, resp)
}
```
リクエストモデル・レスポンスモデルを定義
```go:app/adapter/presentation/handler/task/update_task_state_models.go
package task

type UpdateTaskStateRequest struct {
	State string `json:"state" validate:"required"`
}
type UpdateTaskStateResponse struct {
	ID      string `json:"id"`
	UserId  string `json:"user_id"`
	Content string `json:"content"`
	State   string `json:"state"`
}
```
::::

::::details GetTaskHandlerの実装
```go:app/adapter/presentation/handler/task/get_task_handler.go
package task

import (
	"net/http"

	"github.com/kakkky/app/adapter/presentation/presenter"
	"github.com/kakkky/app/application/usecase/task"
	"github.com/kakkky/app/domain/errors"
)

type GetTaskHandler struct {
	fetchTaskUsecase *task.FetchTaskUsease
}

func NewGetTaskHandler(fetchTaskUsecase *task.FetchTaskUsease) *GetTaskHandler {
	return &GetTaskHandler{
		fetchTaskUsecase: fetchTaskUsecase,
	}
}

// @Summary     タスクを表示する
// @Description idを指定してタスクを表示する
// @Tags        Task
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} presenter.SuccessResponse[GetTaskResponse] 　　"タスクの情報"
// @Failure     400 {object} presenter.FailureResponse                  "不正なリクエスト"
// @Failure     500 {object} presenter.FailureResponse                  "内部サーバーエラー"
// @Router      /tasks/{id} [get]
func (gth *GetTaskHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// パスパラメータから取得
	id := r.PathValue("id")
	// inputDTOに詰め替える
	input := task.FetchTaskUsecaseInputDTO{
		ID: id,
	}
	output, err := gth.fetchTaskUsecase.Run(r.Context(), input)
	if err != nil && errors.IsDomainErr(err) {
		presenter.RespondBadRequest(rw, err.Error())
		return
	}
	if err != nil {
		presenter.RespondInternalServerError(rw, err.Error())
		return
	}
	resp := GetTaskResponse{
		ID:       output.ID,
		UserId:   output.UserId,
		UserName: output.UserName,
		Content:  output.Content,
		State:    output.State,
	}
	presenter.RespondOK(rw, resp)
}
```
レスポンスモデルを定義
```go:app/adapter/presentation/handler/task/get_task_handler.go
package task

type GetTaskResponse struct {
	ID       string `json:"id"`
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
	Content  string `json:"content"`
	State    string `json:"state"`
}
```
::::

::::details GetTasksHandlerの実装
```go:app/adapter/presentation/handler/task/get_tasks_handler.go
package task

import (
	"net/http"

	"github.com/kakkky/app/adapter/presentation/presenter"
	"github.com/kakkky/app/application/usecase/task"
	"github.com/kakkky/app/domain/errors"
)

type GetTasksHandler struct {
	fetchTasksUsecase *task.FetchTasksUsease
}

func NewGetTasksHandler(fetchTasksUsecase *task.FetchTasksUsease) *GetTasksHandler {
	return &GetTasksHandler{
		fetchTasksUsecase: fetchTasksUsecase,
	}
}

// @Summary     全てのタスクを表示する
// @Description 全ユーザーのタスクを全て表示する
// @Tags        Task
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} presenter.SuccessResponse[[]GetTaskResponse] "タスクの情報"
// @Failure     400 {object} presenter.FailureResponse                    "不正なリクエスト"
// @Failure     500 {object} presenter.FailureResponse                    "内部サーバーエラー"
// @Router      /tasks [get]
func (gth *GetTasksHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	outputs, err := gth.fetchTasksUsecase.Run(r.Context())
	if err != nil && errors.IsDomainErr(err) {
		presenter.RespondBadRequest(rw, err.Error())
		return
	}
	if err != nil {
		presenter.RespondInternalServerError(rw, err.Error())
		return
	}
	resp := make([]GetTaskResponse, 0, len(outputs))
	for _, output := range outputs {
		resp = append(resp, GetTaskResponse{
			ID:       output.ID,
			UserId:   output.UserId,
			UserName: output.UserName,
			Content:  output.Content,
			State:    output.State,
		})
	}
	presenter.RespondOK(rw, resp)
}
```
::::

::::details GetUserTasksHandlerの実装
```go:app/adapter/presentation/handler/task/get_user_tasks_handler.go
package task

import (
	"net/http"

	"github.com/kakkky/app/adapter/presentation/middleware"
	"github.com/kakkky/app/adapter/presentation/presenter"
	"github.com/kakkky/app/application/usecase/task"
	"github.com/kakkky/app/domain/errors"
)

type GetUserTasksHandler struct {
	fetchUserTasksUsecase *task.FetchUserTasksUsecase
}

func NewGetUserTasksHandler(fetchUserTasksUsecase *task.FetchUserTasksUsecase) *GetUserTasksHandler {
	return &GetUserTasksHandler{
		fetchUserTasksUsecase: fetchUserTasksUsecase,
	}
}

// @Summary ユーザーが持つ全てのタスクを表示する
// @Description　ログインしているユーザーのタスクを全て表示する
// @Tags     Task
// @Produce  json
// @Security BearerAuth
// @Success  200 {object} presenter.SuccessResponse[[]GetTaskResponse] "タスクの情報"
// @Failure  400 {object} presenter.FailureResponse                    "不正なリクエスト"
// @Failure  500 {object} presenter.FailureResponse                    "内部サーバーエラー"
// @Router   /users/me/tasks [get]
func (guth *GetUserTasksHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	input := task.FetchUserTasksUsecaseInputDTO{
		UserId: userID,
	}
	outputs, err := guth.fetchUserTasksUsecase.Run(r.Context(), input)
	if err != nil && errors.IsDomainErr(err) {
		presenter.RespondBadRequest(rw, err.Error())
		return
	}
	if err != nil {
		presenter.RespondInternalServerError(rw, err.Error())
		return
	}
	resp := make([]GetUserTaskResponse, 0, len(outputs))
	for _, output := range outputs {
		resp = append(resp, GetUserTaskResponse{
			ID:      output.ID,
			Content: output.Content,
			State:   output.State,
		})
	}
	presenter.RespondOK(rw, resp)
}
```
レスポンスモデルを定義
```go:app/adapter/presentation/handler/task/get_user_tasks_models.go
package task

type GetUserTaskResponse struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	State   string `json:"state"`
}
```
::::
#### パスパラメータの取得
Task系ハンドラには、このような記述がいくつかあります。
```go
// idをパスパラメータから取得
id := r.PathValue("id")
```
これによって、パスパラメータ`id`を取得します。タスクのidを指定して操作を行うためにこうしています。
```go
// PathValue は、リクエストにマッチした [ServeMux] のパターン内にある指定されたパスワイルドカードの値を返します。
// リクエストがパターンにマッチしなかった場合、またはパターン内にそのようなワイルドカードが存在しない場合は、空の文字列を返します。
func (r *Request) PathValue(name string) string {
	if i := r.patIndex(name); i >= 0 {
		return r.matches[i]
	}
	return r.otherValues[name]
}
```
Go1.22のアプデにより、パスパラメータの取得のためのメソッドが追加されました。それまではなかったので標準パッケージは苦しい部分がありましたが、これにより`net/http`が強化されたと言えます。
https://tip.golang.org/doc/go1.22#enhanced_routing_patterns


### 認証系
::::details LoginHandlerの実装
```go:app/adapter/presentation/handler/auth/login_handler.go
package auth

import (
	"encoding/json"
	"net/http"

	"github.com/kakkky/app/adapter/presentation/presenter"
	"github.com/kakkky/app/application/usecase/auth"
	"github.com/kakkky/app/domain/errors"
	"github.com/kakkky/pkg/validation"
)

type LoginHandler struct {
	loginUsecase *auth.LoginUsecase
}

func NewLoginHandler(loginUsecase *auth.LoginUsecase) *LoginHandler {
	return &LoginHandler{
		loginUsecase: loginUsecase,
	}
}

// @Summary     ユーザーのログイン
// @Description メールアドレス・パスワードで認証し、署名されたトークンを返す
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       request body     LoginRequest                             true "認証に必要な情報"
// @Success     200     {object} presenter.SuccessResponse[LoginResponse] "署名されたトークンを含む情報"
// @Failure     400     {object} presenter.FailureResponse                "不正なリクエスト"
// @Failure     401     {object} presenter.FailureResponse                "パスワードが不一致"
// @Failure     500     {object} presenter.FailureResponse                "内部サーバーエラー"
// @Router      /login [post]
func (lh *LoginHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var params LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		presenter.RespondBadRequest(rw, err.Error())
		return
	}
	if err := validation.NewValidator().Struct(&params); err != nil {
		presenter.RespondBadRequest(rw, err.Error())
		return
	}
	input := auth.LoginUsecaseInputDTO{
		Email:    params.Email,
		Password: params.Password,
	}
	output, err := lh.loginUsecase.Run(r.Context(), input)
	if (err != nil) && errors.Is(err, errors.ErrPasswordMismatch) || errors.IsDomainErr(err) {
		presenter.RespondUnAuthorized(rw, "email or password invalid")
		return
	}
	if err != nil {
		presenter.RespondInternalServerError(rw, err.Error())
		return
	}
	resp := LoginResponse{
		JwtToken: output.JwtToken,
	}
	presenter.RespondOK(rw, resp)
}
```
リクエストモデル・レスポンスモデルを定義
```go:app/adapter/presentation/handler/auth/login_models.go
package auth

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	JwtToken string `json:"jwt_token"`
}
```
::::

::::details LogoutHandlerの実装
```go:app/adapter/presentation/handler/auth/logout_handler.go
package auth

import (
	"net/http"

	"github.com/kakkky/app/adapter/presentation/middleware"
	"github.com/kakkky/app/adapter/presentation/presenter"
	"github.com/kakkky/app/application/usecase/auth"
)

type LogoutHandler struct {
	logoutUsecase *auth.LogoutUsecase
}

func NewLogoutHandler(logoutUsecase *auth.LogoutUsecase) *LogoutHandler {
	return &LogoutHandler{
		logoutUsecase: logoutUsecase,
	}
}

// @Summary     ユーザーのログアウト
// @Description メールアドレス・パスワードで認証し、署名されたトークンを返す
// @Tags        Auth
// @Produce     json
// @Security    BearerAuth
// @Success     204
// @Failure     500 {object} presenter.FailureResponse "内部サーバーエラー"
// @Router      /logout [delete]
func (lh *LogoutHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// 認可制御のミドルウェアでコンテキストに値が付加されている
	// リクエストスコープのコンテキストからuserIDを取得する
	userID := middleware.GetUserID(r.Context())
	ctx := r.Context()
	if err := lh.logoutUsecase.Run(ctx, auth.LogoutUsecaseInputDTO{UserID: userID}); err != nil {
		presenter.RespondInternalServerError(rw, err.Error())
	}
	presenter.RespondNoContent(rw)
}

```
::::

特に、User系やTask系のハンドラの実装と変わりありません。

# swaggoによるドキュメント生成
Swaggerは、APIの設計、文書化、テストをサポートするツールセットであり、OpenAPI仕様に基づいています。
OpenAPIとは、RESTful APIを記述するための標準的な仕様のことです。「どのような項目をどのような形式で記載すればよいのか」といった体裁を定義しており、ドキュメントに統一性を持たせられるため理解しやすくなります。
https://qiita.com/MaSi1031/items/ada0854f3d1d93ee66f8
https://swagger.io/
https://qiita.com/yuya_sega/items/0b87e8e7d494f6fa3d69

SwaggerをGoで扱いやすくしたのがswaggoというツールです。
https://github.com/swaggo/swag
https://qiita.com/pei0804/items/3a0b481d1e47e5a72078

swaggoを操作するためのCLIは、dockerファイルにてインストールしていました。
```Dockerfile
RUN go install github.com/swaggo/swag/cmd/swag@latest
```

それに加え、swaggoを`net/http`で扱うためのパッケージをインストールしておきます。
```
make get-app name="github.com/swaggo/http-swagger"
```

https://github.com/swaggo/http-swagger

swaggoは、アノテーションを解析してドキュメントを自動で生成してくれるとても便利なツールです。

以下のようなアノテーションを、定義したハンドラの上に記述していました。
```
// @Summary     タスクを作成する
// @Description 内容、タスク状態からユーザーに紐づくタスクを作成する
// @Tags        Task
// @Produce     json
// @Param       request body PostTaskRequest true "タスク作成のための情報"
// @Security    BearerAuth
// @Success     201 {object} presenter.SuccessResponse[PostTaskResponse] "作成したタスクの情報"
// @Failure     400 {object} presenter.FailureResponse                   "不正なリクエスト"
// @Failure     500 {object} presenter.FailureResponse                   "内部サーバーエラー"
// @Router      /tasks [post]
```
### アノテーションの意味合い
今回扱っているアノテーションの意味合いをまとめています。

| 記法                     | 記述例                                                      | 説明                                      |
|--------------------------|-----------------------------------------------------------|-------------------------------------------|
| `@Summary`              | `@Summary タスクを作成する`                                | API の概要を記述                          |
| `@Description`          | `@Description 内容、タスク状態からユーザーに紐づくタスクを作成する` | API の詳細な説明を記述                    |
| `@Tags`                 | `@Tags Task`                                              | API のカテゴリ（グループ）を指定          |
| `@Produce`              | `@Produce json`                                           | レスポンスのフォーマットを指定            |
| `@Param` (ボディ)       | `@Param request body PostTaskRequest true "タスク作成のための情報"` | リクエストボディのパラメータを定義        |
| `@Param` (パス)         | `@Param id path int true "タスクID"`                      | パスパラメータの定義                      |
| `@Security`             | `@Security BearerAuth`                                    | 認証方式を指定                            |
| `@Success` (オブジェクト) | `@Success 201 {object} presenter.SuccessResponse[PostTaskResponse] "作成したタスクの情報"` | 成功時のレスポンスの型を指定              |
| `@Failure` (400)        | `@Failure 400 {object} presenter.FailureResponse "不正なリクエスト"` | 失敗時のレスポンス（400 Bad Request）     |
| `@Failure` (500)        | `@Failure 500 {object} presenter.FailureResponse "内部サーバーエラー"` | 失敗時のレスポンス（500 Internal Server Error） |
| `@Router`               | `@Router /tasks [post]`                                  | API のエンドポイントと HTTP メソッドを指定 |

`@Success`を取り上げてみます。
```
// @Success     201 {object} presenter.SuccessResponse[PostTaskResponse] "作成したタスクの情報"
```
`{object}`とは一つのオブジェクトを表し、そのオブジェクトが`presenter.SuccessResponse[PostTaskResponse]`型であることを示しています。つまり、成功時にはその型のレスポンスが返るとということですね。

swaggoは型へのジェネリクスにも対応しています。
https://github.com/swaggo/swag?tab=readme-ov-file#how-to-use-generics

### APIの概要も書く
各エンドポイントにアノテーションを書きましたが、本APIの概要もドキュメントに載せたいと思います。
`main.go`の`main()`関数上部に、以下を追加してください。
```go:app/main.go
// @title                      TODO API
// @version                    1.0
// @description                This is TODO API by golang.
// @host                       localhost:8080
// @BasePath                   /
// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
```

使用したアノテーションの意味合いは以下の通りです。

| アノテーション                     | 説明 |
|----------------------------------|------|
| `@title`                        | API のタイトル |
| `@version`                      | API のバージョン |
| `@description`                  | API の概要説明 |
| `@host`                         | API のホスト名（例: `localhost:8080`） |
| `@BasePath`                     | API のベースパス（例: `/`） |
| `@securityDefinitions.apikey`   | 認証の種類（APIキー方式の認証を定義） |
| `@in`                           | 認証情報の取得場所（`header` は HTTP ヘッダー） |
| `@name`                         | 認証情報のキー名（例: `Authorization`） |

### アノテーションからドキュメントを生成する
アノテーションを追加しただけでドキュメントは生成されません。

以下のようにタスクを定義しています。
```Makefile
# アノテーションをパースしてドキュメント生成
swag:
	@echo "Generating document by swagger..."
	cd ./app && swag fmt && swag init
```

コマンドを実行してみてください。
```
make swag
```

ドキュメントの自動生成に成功すれば、以下のようにログが出るはずです。
```
Generating document by swagger...
cd ./app && swag fmt && swag init
// 省略
2025/01/30 15:39:49 Generating user.UpdateUserResponse
2025/01/30 15:39:49 create docs.go at  docs/docs.go
2025/01/30 15:39:49 create swagger.json at  docs/swagger.json
2025/01/30 15:39:49 create swagger.yaml at  docs/swagger.yaml
```
`app/docs`ディレクトリが作成され、いくつかファイルが生成されていることも確認してください。

さらに、生成された`docs`パッケージを`app/main.go`内でブランクインポートしておく必要があります。以下のように追加してください。
```go:app/main.go
import (
	_ "github.com/あなたのプロジェクト名/app/docs"	 //追加
)
```

生成したドキュメントを閲覧するには、エンドポイントを通して、ページにアクセスする必要があります。
その作業については13章にて行うこととします。

# 現段階のディレクトリ
ファイルも示したいところなのですが、文字数超過(5万字)がきてしまったのでディレクトリのみに絞って載せておきます。
```
.
├── app
│   ├── adapter
│   │   ├── presentation
│   │   │   ├── handler
│   │   │   │   ├── auth
│   │   │   │   ├── health
│   │   │   │   ├── task
│   │   │   │   └── user
│   │   │   └── presenter
│   │   ├── queryservice
│   │   └── repository
│   ├── application
│   │   └── usecase
│   │       ├── auth
│   │       ├── task
│   │       └── user
│   ├── config
│   ├── docs
│   ├── domain
│   │   ├── errors
│   │   ├── task
│   │   └── user
│   ├── infrastructure
│   │   ├── auth
│   │   │   └── certificate
│   │   ├── db
│   │   │   ├── container
│   │   │   ├── migrations
│   │   │   ├── mysql
│   │   │   │   └── conf.d
│   │   │   ├── sqlc
│   │   │   │   └── queries
│   │   │   └── testhelper
│   │   ├── kvs
│   │   ├── queryservice_test
│   │   │   └── testdata
│   │   │       └── fixtures
│   │   └── repository_test
│   │       └── testdata
│   │           └── fixtures
│   └── tmp
└── pkg
    ├── hash
    ├── ulid
    └── validation
```
