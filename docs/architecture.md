# アーキテクチャ詳細

## ディレクトリ構成

```
cmd/
  main.go                          # エントリポイント・依存性注入
internal/
  config/                          # 環境変数読み込み
  domain/
    apperror/                      # カスタムエラー型
    model/                         # ドメインモデル（AuthUser, User, KeyValue）
      fx/                          # FXドメインモデル（TradeEntry, TradePosition, PriceInfo 他）
    repository/                    # リポジトリインターフェース
      fx/                          # FXリポジトリIF（TradeSimulation, Symbol, Country, EconomicIndicator）
    service/
      fx/                          # FXドメインサービス（FxTradeCalculator）
  application/
    command/                       # コマンドオブジェクト（入力値の集約）
      fx/                          # FXコマンド（TradeSimulationCommand）
    dto/                           # データ転送オブジェクト（UserDto）
      fx/                          # FX DTO（SymbolDto）
    usecase/                       # ユースケース（ビジネスロジック）
      fx/                          # FXユースケース（TradeSimulationUseCase, GetMasterUseCase）フラット
        country/                   # 国マスタ CRUD UseCase
        symbol/                    # シンボルマスタ CRUD UseCase
        summertime/                # サマータイム CRUD UseCase
        bardata/                   # バーデータ検索・ステータス UseCase
        economicindicator/         # 経済指標 CRUD UseCase
        economicindicatordata/     # 経済指標データ CRUD + インポート UseCase
  infrastructure/
    infraredis/                    # Redisセッション実装
    infradb/                       # MySQL実装（User）
      fx/                          # FX MySQL実装（TradeSimulation, Symbol, Country, EconomicIndicator）
    external/                      # 外部サービス（GaitameRateService）
  security/
    jwt_provider.go                # JWT検証（JWKS自動取得・RS256）
  api/
    middleware/                    # jwt_middleware.go / auth_middleware.go
    controller/                    # HTTPハンドラ（auth_controller.go, user_controller.go, fx_*.go）
    dto/
      request/                     # リクエストDTO
        fx/                        # FXリクエストDTO（TradeSimulationRequest, SymbolRequest）
      response/                    # レスポンスDTO（ErrorResponse, ApiResponse など共通型）
        fx/                        # FXレスポンスDTO（BarDataSearchResponse, TradeSimulationResponse 他）
    router/                        # ルーティング設定
```

### パッケージ命名規約（FX機能）

ドメイン境界ごとにサブパッケージを切り、競合を避けるためにプレフィックスを付ける。

| ディレクトリ | パッケージ名 |
|---|---|
| `internal/domain/model/fx/` | `fx` |
| `internal/domain/repository/fx/` | `fxrepository` |
| `internal/domain/service/fx/` | `fxservice` |
| `internal/application/command/fx/` | `fxcommand` |
| `internal/application/usecase/fx/` | `fxusecase`（TradeSimulationUseCase, GetMasterUseCase のみ） |
| `internal/application/usecase/fx/country/` | `country` |
| `internal/application/usecase/fx/symbol/` | `symbol` |
| `internal/application/usecase/fx/summertime/` | `summertime` |
| `internal/application/usecase/fx/bardata/` | `bardata` |
| `internal/application/usecase/fx/economicindicator/` | `economicindicator` |
| `internal/application/usecase/fx/economicindicatordata/` | `economicindicatordata` |
| `internal/infrastructure/infradb/fx/` | `infradbfx` |
| `internal/infrastructure/external/` | `external` |
| `internal/application/dto/fx/` | `fxdto` |
| `internal/api/dto/request/fx/` | `fxrequest` |
| `internal/api/dto/response/fx/` | `fxresponse` |

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
