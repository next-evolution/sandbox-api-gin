# アーキテクチャ詳細

## ディレクトリ構成

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

## 使用ライブラリ

| ライブラリ | 用途 |
|---|---|
| github.com/gin-gonic/gin | HTTPフレームワーク |
| github.com/lestrrat-go/jwx/v2 | JWT検証・JWKS自動取得 |
| github.com/redis/go-redis/v9 | Redisクライアント |
| github.com/jmoiron/sqlx | database/sqlの薄いラッパー（構造体スキャン） |
| github.com/go-sql-driver/mysql | MySQLドライバ |
| github.com/joho/godotenv | .envファイル読み込み |
| github.com/gin-contrib/cors | CORSミドルウェア |

## 認証フロー詳細

### JWT Middleware（`internal/api/middleware/jwt_middleware.go`）
Javaの `JwtAuthFilter` に相当。

1. `Authorization: Bearer <token>` からトークン取得
2. JWKS（CognitoのJWKSエンドポイント）を使ってRS256署名を検証
3. issuer・audience・有効期限をバリデーション
4. Redisからadminフラグ付き `AuthUser` を取得（セッションなければJWT由来のAuthUser使用）
5. Gin Contextに `authUser` キーでセット・セッションTTL更新

### Auth Middleware（`internal/api/middleware/auth_middleware.go`）
Javaの `AuthInterceptor` に相当。

* Context に `authUser` がなければ 401 を返す
* ログイン前でも有効なJWTがあればAuthUserはセットされる（ログインAPIはRedisセッション不要）
