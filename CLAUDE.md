# Golang Gin RestAPI

Java SpringBoot で構築された RestAPI をベースに Go Gin の RestAPI を構築。

* ベースとなるソース: `./docs/base_src`（**読み取りのみ・書き込み禁止**）

---

## このファイルの管理方針

**CLAUDE.md は「Claudeの行動を変える指示書」** であり、ドキュメントではない。
毎回コンテキストに全文読み込まれるため、肥大化させない。

| 種類 | 置き場所 |
|---|---|
| コーディング規約・禁止事項 | **CLAUDE.md** |
| アーキテクチャの制約（依存方向など） | **CLAUDE.md** |
| ビルド・実行コマンド | **CLAUDE.md** |
| 重要な落とし穴（bit(1)、エラー型など） | **CLAUDE.md** |
| エンドポイント一覧 | `docs/api.md` |
| タスク指示（step系） | **プロンプトで渡す** |

---

## ドキュメント参照先

| 内容 | ファイル |
|---|---|
| APIエンドポイント一覧・レスポンス仕様 | [docs/api.md](docs/api.md) |

---

## Build & Run

```bash
# ビルド
go build ./...
# または
make build

# 実行（APP_ENVで環境を指定する）
APP_ENV=local go run ./cmd/main.go   # .env.local を読み込む
APP_ENV=docker go run ./cmd/main.go  # .env.docker を読み込む
# または
APP_ENV=local make run

# Lintチェック（golangci-lint v2）
golangci-lint run ./...
# または
make lint

# go.mod整理
go mod tidy
# または
make tidy
```

### 必要な環境変数

| 変数名 | 例 | 説明 |
|---|---|---|
| DB_HOST | localhost | MySQLホスト |
| DB_PORT | 43306 | MySQLポート |
| DB_SCHEMA | sandbox_local | DBスキーマ名 |
| DB_USER | sandbox_app | DBユーザー |
| DB_PASSWORD | s4ndb0x_app | DBパスワード |
| REDIS_HOST | localhost | Redisホスト |
| REDIS_PORT | 46379 | Redisポート |
| JWT_ISSUER1 | https://cognito-idp.ap-northeast-1.amazonaws.com/... | CognitoのIssuer URL |
| JWT_AUDIENCE1/2/3 | Cognito App Client ID | 許可するaudience（複数可） |
| JWT_ORIGIN1/2 | http://localhost:3000 | CORS許可オリジン（未設定時はCORS無効） |
| SESSION_TTL | 3600 | セッションTTL（秒）デフォルト3600 |
| SERVER_PORT | 8080 | サーバーポート デフォルト8080 |
| GIN_MODE | debug | Ginモード（debug / release / test） |

### 環境設定ファイル

`APP_ENV` 環境変数で読み込むファイルを切り替える。ファイルはすべて `.gitignore` 対象。

| APP_ENV | 読み込むファイル | 用途 |
|---|---|---|
| 未設定 | `.env` | production |
| `local` | `.env.local` | ローカル起動 |
| `docker` | `.env.docker` | Docker起動 |

### ローカルインフラ起動

```bash
# MySQL + Redis を Docker で起動（./docs/base_srcのdocker-compose.ymlを使用）
cd docs/base_src && docker compose up -d
```

---

## アーキテクチャ

DDD（ドメイン駆動設計）に基づくレイヤー構造。

### ディレクトリ構成

```
cmd/
  main.go                          # エントリポイント・依存性注入
internal/
  config/                          # 環境変数読み込み
  domain/
    apperror/                      # カスタムエラー型
    model/                         # ドメインモデル（AuthUser, User）
    repository/                    # リポジトリインターフェース
  application/
    command/                       # コマンドオブジェクト（入力値の集約）
    dto/                           # データ転送オブジェクト（UserDto）
    usecase/                       # ユースケース（ビジネスロジック）
  infrastructure/
    infraredis/                    # Redisセッション実装
    infradb/                       # MySQL実装
  security/
    jwt_provider.go                # JWT検証（JWKS自動取得・RS256）
  api/
    middleware/                    # jwt_middleware.go / auth_middleware.go
    controller/                    # HTTPハンドラ
    dto/request/ & response/       # リクエスト/レスポンスDTO
    router/                        # ルーティング設定
```

### レイヤー依存関係

```
api → application → domain
  ↓                    ↑
infrastructure ─────────┘
security → domain
```

---

## 認証フロー

### JWT Middleware（`internal/api/middleware/jwt_middleware.go`）
Javaの `JwtAuthFilter` に相当。

1. `Authorization: Bearer <token>` からトークン取得
2. JWKS（CognitoのJWKSエンドポイント）を使ってRS256署名を検証
3. issuer・audience・有効期限をバリデーション
4. Redisからadminフラグ付き `AuthUser` を取得（セッションなければJWT由来のAuthUser使用）
5. Gin Contextに `authUser` キーでセット・セッションTTL更新

### Auth Middleware（`internal/api/middleware/auth_middleware.go`）
Javaの `AuthInterceptor` に相当。

* Context に `authUser` がなければ401を返す
* ログイン前でも有効なJWTがあればAuthUserはセットされる（ログインAPIはRedisセッション不要）

---

## 実装規約

### エラー型
`internal/domain/apperror` のカスタムエラー型を使用。

| エラー型 | HTTPステータス |
|---|---|
| AuthenticationError | 401 |
| ForbiddenError | 403 |
| NotFoundError | 404 |
| DuplicateError / InsertError / UpdateError | 400 |

### Base64デコード
`encoding/base64.StdEncoding` を使用（JavaのBase64.getDecoder()に相当）。
パディングなしの場合は `RawStdEncoding` にフォールバック。

### MySQL `bit(1)` カラム
Goのdatabase/sqlドライバは `bit(1)` を `[]byte` で返す。
SQLクエリで `(col+0) AS col` にキャストし、Go側では `uint8` でスキャンして `!= 0` でbool変換する。

### 日時フォーマット
`UserDto` の日時フィールドはJavaに合わせて `"yyyy-MM-dd HH:mm:ss"` 形式。
`internal/application/dto/user_dto.go` の `DateTime` 型でカスタムMarshal。

### 使用ライブラリ
| ライブラリ | 用途 |
|---|---|
| github.com/gin-gonic/gin | HTTPフレームワーク |
| github.com/lestrrat-go/jwx/v2 | JWT検証・JWKS自動取得 |
| github.com/redis/go-redis/v9 | Redisクライアント |
| github.com/jmoiron/sqlx | database/sqlの薄いラッパー（構造体スキャン） |
| github.com/go-sql-driver/mysql | MySQLドライバ |
| github.com/joho/godotenv | .envファイル読み込み |
| github.com/gin-contrib/cors | CORSミドルウェア |

### Lintルール（`.golangci.yml`）
golangci-lint v2を使用。有効なリンター：

| リンター | 目的 |
|---|---|
| errcheck | エラーの戻り値を無視していないかチェック |
| govet | go vetによる静的解析 |
| ineffassign | 無効な代入の検出 |
| staticcheck | 高度な静的解析（gosimpleも統合済み） |
| unused | 未使用コードの検出 |
| misspell | スペルミスの検出 |
| gocritic | バグ・パフォーマンス・スタイルのレビュー |
| noctx | contextなしのDB/HTTPリクエスト検出 |
| bodyclose | HTTPレスポンスbodyのclose漏れ検出 |
| goimports（formatter） | importの順序・不要importのチェック |

#### importグループの順序（goimports）
```go
import (
    // 1. 標準ライブラリ
    "context"
    "fmt"

    // 2. サードパーティ
    "github.com/gin-gonic/gin"

    // 3. ローカル（sandbox-api-gin/...）
    "sandbox-api-gin/internal/..."
)
```

### DB操作の規約

**sqlxの `GetContext` / `SelectContext` を使って構造体に直接スキャンする**。

```go
// 1件取得
var rec sandboxUserRecord
err := db.GetContext(ctx, &rec, query, args...)
if err == sql.ErrNoRows { ... }

// 複数件取得
var recs []sandboxUserRecord
err := db.SelectContext(ctx, &recs, query, args...)
```

**DB操作は必ず `Context` を渡す**（`noctx` リンターで検出される）。

```go
// NG
row := db.QueryRow(query, args...)

// OK
row := db.QueryRowContext(context.Background(), query, args...)
```

**INSERT / UPDATE は `ExecContext` を使い、`RowsAffected` で件数を検証する**。INSERT後は `LastInsertId()` でIDを取得してドメインモデルにセットする。

```go
result, err := db.ExecContext(ctx, query, args...)
if err != nil { return err }
rows, err := result.RowsAffected()
if err != nil { return err }
if rows != 1 { return apperror.NewInsertError("...") }
// INSERT のみ: 自動採番IDをドメインモデルに反映
id, err := result.LastInsertId()
if err != nil { return err }
user.ID = id
```

**リソースの Close は `defer func(){}()` でエラーハンドリングする**（`errcheck` リンターで検出される）。

```go
// NG
defer db.Close()

// OK
defer func() {
    if err := db.Close(); err != nil {
        slog.Error("切断エラー", "error", err)
    }
}()
```

### main.goのパターン
`run()` 関数にロジックを分離し、`defer` が確実に実行されるようにしている。
`log.Fatal` は `main()` からのみ呼び出す（`gocritic: exitAfterDefer` で検出される）。

```go
func main() {
    if err := run(); err != nil {
        log.Fatal(err)
    }
}

func run() error {
    // deferを使って安全にリソース解放
    // エラーはreturnで伝播させる
}
```

### CORS設定

`JWT_ORIGIN1` / `JWT_ORIGIN2` が設定されている場合のみ CORS ミドルウェアを有効化。
未設定の場合は CORS なし（同一オリジンのみ許可）。

```go
if len(cfg.JWTOrigins) > 0 {
    engine.Use(cors.New(cors.Config{
        AllowOrigins:     cfg.JWTOrigins,
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))
}
```

### Graceful Shutdown

SIGINT / SIGTERM を受け取ったら10秒のタイムアウトで安全に終了。

```go
srv := &http.Server{Addr: ":" + cfg.ServerPort, Handler: engine}
go func() {
    srv.ListenAndServe()
}()
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
srv.Shutdown(shutdownCtx)
```

### コンテキスト伝播

HTTP リクエストのコンテキスト（タイムアウト・キャンセル）をすべての層に伝播させる。

```
Controller: ctx := c.Request.Context()
              ↓
UseCase:    Execute(ctx context.Context, cmd)
              ↓
Repository: Save(ctx context.Context, ...) / FindBySub(ctx, ...) etc.
              ↓
DB / Redis: db.GetContext(ctx, ...) / redisClient.Set(ctx, ...)
```

### リクエストバリデーション

Gin の `binding` タグで入力バリデーションを行う。`ShouldBindJSON` でエラーが返った場合は 400 を返す。

```go
type UserRegistrationRequest struct {
    NickName string `json:"nickName" binding:"required,max=50"`
}

var req UserRegistrationRequest
if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, response.ErrorResponse{...})
    return
}
```

### パスパラメータの userId デコード

パスに `:userId` を含むエンドポイントでは Base64 エンコード済みで渡される。
`decodeBase64UserID()` でデコードし、`authUser.Sub` と一致しなければ `ForbiddenError`。

Javaの `UserId.decodeUserIdValue()` に相当。`StdEncoding` → `RawStdEncoding` フォールバックは Base64 デコードと同じパターン。

### パッケージ命名規約
Goのパッケージ名は短い識別子。ディレクトリ名と一致させる。
サードパーティと競合する場合はprefixを付ける。

```
internal/infrastructure/infraredis/  → package infraredis  （"redis"だとgo-redisと競合）
internal/infrastructure/infradb/     → package infradb
```
