# API 設計

## 概要

RPC スタイルの HTTP API

POST メソッドのみを使用し、JSON 形式で通信する。

## 命名規則

全て POST で統一し、パスは以下の命名規則に従う。

```
POST /<serviceName>/<version>/<useCaseName>/<methodName>

例: POST /core/v1/task/get
```

- 単語は `lowerCamelCase` で統一
- `<serviceName>` はマイクロサービス化を想定して設定、モノリス時は`core`というサービス名を使う。
- `<useCaseName>`, `<methodName>` は`usecase`ディレクトリ配下の構造体名、メソッド名に合わせる。
- `get`, `list` から始まる method は安全性・冪等性を担保する。

## API スキーマ定義

OpenAPI などを使うとスキーマファイルの管理コストがかかるため用意しない。

代わりに、ソースコードを綺麗に保ち API のパス、リクエスト、レスポンスを読みやすくする。

`cmd/httpserver/main.go`を見れば全てのエンドポイントが分かる。

```go
// cmd/httpserver/main.go
mux.HandleFunc("GET /healthz", handler.HandleGetHealthz) // ヘルスチェックのみGETを使用

mux.HandleFunc("POST /core/v1/task/get", taskHandler.HandleGetV1)
mux.HandleFunc("POST /core/v1/task/list", taskHandler.HandleListV1)
mux.HandleFunc("POST /core/v1/task/create", taskHandler.HandleCreateV1)
mux.HandleFunc("POST /core/v1/task/update", taskHandler.HandleUpdateV1)
mux.HandleFunc("POST /core/v1/task/delete", taskHandler.HandleDeleteV1)
mux.HandleFunc("POST /core/v1/task/done", taskHandler.HandleDoneV1)
```

handler のコードを追えばリクエスト、レスポンスの内容が分かる。

```go
// interface/inbound/http/handler/task.go

// POST /core/v1/task/get
func (h *TaskHandler) HandleGetV1(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ID string `json:"id"`
	} // リクエストスキーマ
	var successResponse response.Task // 成功時のレスポンススキーマ
	var errorResponse response.Error  // エラーレスポンススキーマ

	ctx := r.Context()

	// 400
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		slog.WarnContext(ctx, "handler.TaskHandler.HandleGetV1", "err", err)
		errorResponse = response.MapError(response.ErrInvalidRequestBody)
		response.RenderJson(ctx, w, http.StatusBadRequest, errorResponse)
		return
	}

	t, err := h.taskUseCase.Get(ctx, request.ID)

	// 200
	if err == nil {
		successResponse = response.MapTask(t)
		response.RenderJson(ctx, w, http.StatusOK, successResponse)
		return
	}

	// 400
	if errors.Is(err, task.ErrInvalidID) {
		slog.WarnContext(ctx, "handler.TaskHandler.HandleGetV1", "err", err)
		errorResponse = response.MapError(response.ErrInvalidRequestBody)
		response.RenderJson(ctx, w, http.StatusBadRequest, errorResponse)
		return
	}

	// 404
	if errors.Is(err, task.ErrNotFound) {
		slog.WarnContext(ctx, "handler.TaskHandler.HandleGetV1", "err", err)
		errorResponse = response.MapError(response.ErrNotFound)
		response.RenderJson(ctx, w, http.StatusNotFound, errorResponse)
		return
	}

	// 500
	slog.ErrorContext(ctx, "handler.TaskHandler.HandleGetV1", "err", err)
	errorResponse = response.MapError(response.ErrInternalServerError)
	response.RenderJson(ctx, w, http.StatusInternalServerError, errorResponse)
}
```

## リクエスト

- POST メソッドを利用する
- クエリパラメータは使用しない

## レスポンス

- `Content-Type`レスポンスヘッダーに`application/json; charset=utf-8`を設定する
- [HTTP response status codes - HTTP | MDN](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status)を参考に適切なステータスコードを設定する。

## その他の API 形式との比較

それぞれのメリット、デメリットはこのリポジトリの設計との比較。

**RESTful**

メリット

- メソッドから安全性、冪等性、キャッシュ可否が明確 ( [HTTP request methods - HTTP | MDN](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Methods) )
- 採用実績が多く、チームメンバー間での認識齟齬が生まれにくい

デメリット

- 同じリソースに対して多様な操作を実装する場合、9 つのメソッドだけでは表現しきれない。

**gRPC**

メリット

- protocol buffer というバイナリ形式で通信するため、データサイズを圧縮できる。

デメリット

- proto ファイルを元に生成したクライアントコードの使用が前提になるため、JavaScript の fetch API や、curl コマンドなどが使用できず、開発体験が悪い。

**GraphQL**

メリット

- クライアントから必要なデータを指定できるため、オーバーフェッチ、アンダーフェッチを回避できる。

デメリット

- Apollo などのライブラリを利用する前提になるため、依存が増える。
- JavaScript の fetch API を直接使用しないことになる。

**JSON-RPC**

このリポジトリの設計に 1 番近い。[JSON-RPC](https://www.jsonrpc.org/)では単一のパスにしてリクエストボディでメソッドも指定している。

メリット

- コミュニティで仕様を規定しているため、独自ルールよりはチームメンバー間での認識齟齬が生まれにくい

デメリット

- パスにメソッドを含まれないのでログの検索性が下がる

**結論**

- gRPC, GraphQL のように依存が増えるものは開発体験が下がるので採用しない。
- RESTful は優れた設計だが、複雑なアプリケーションではメソッドの拡張性に課題がある。
- JSON-RPC は 1 番近いが、ログの検索性を優先してこれには準拠しない。
