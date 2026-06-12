# sandbox-api-gin

![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-1.12-00ADD8?logo=go&logoColor=white)
![MySQL](https://img.shields.io/badge/MySQL-8.4-4479A1?logo=mysql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-8.0-DC382D?logo=redis&logoColor=white)
![AWS Cognito](https://img.shields.io/badge/AWS_Cognito-FF9900?logo=amazonaws&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?logo=docker&logoColor=white)

FX トレード支援を目的としたバックエンド REST API。  
Go 1.26 / Gin で構築し、**DDD（ドメイン駆動設計）** に基づくレイヤー構造を採用。  
認証は AWS Cognito（RS256 JWT）、セッション管理は Redis で実装。  
[sandbox-api-springboot](../sandbox-api-springboot) を Go へ移植したプロジェクト。

---

## アーキテクチャ

### パッケージ構成

| パッケージ | 役割 |
|---|---|
| `internal/api` | コントローラー、ミドルウェア、ルーティング |
| `internal/application` | ユースケース、コマンド、DTO |
| `internal/domain` | ドメインモデル（値オブジェクト・集約）、リポジトリインターフェース、カスタムエラー型。フレームワーク非依存。 |
| `internal/infrastructure` | リポジトリ実装（sqlx / MySQL）、Redis セッション |
| `internal/security` | JWT 検証（RS256）、JWKS キャッシュ |
| `internal/config` | 環境変数読み込み・設定管理 |

### 依存関係

```
api → application → domain
  ↓                    ↑
infrastructure ─────────┘
security → domain
```

詳細は [docs/architecture.md](./docs/architecture.md) を参照。

---

## 主な機能

| ドメイン | 主なエンドポイント |
|---|---|
| **認証 / Auth** | ログイン・ログアウト（AWS Cognito JWT + Redis セッション） |
| **ユーザー / User** | プロフィール取得・ユーザー登録・情報更新 |
| **管理者 / Admin** | ユーザー検索・承認・ブロック・管理者権限付与、Redis キャッシュ管理 |
| **FX マスター** | 通貨シンボル・国・通貨ペア・経済指標（公開 API） |
| **FX バーデータ** | OHLC バーデータ検索・CSV 一括インポート |
| **ZigZag 分析** | ZigZag 生成・検索・ステータス取得・バーデータ取得 |
| **トレードシミュレーション** | リスク額・ロット比率・エントリーに基づくシミュレーション |

エンドポイント詳細は [docs/api.md](./docs/api.md) を参照。

---

## Getting Started

### 1. ローカルインフラ起動

```bash
# MySQL（43306）+ Redis（46379）を Docker で起動
cd docs/base_src && docker compose up -d
```

### 2. アプリケーション環境変数

`.env.local` ファイルを作成して設定：

| 変数名 | 説明 | 例 |
|---|---|---|
| `DB_HOST` | MySQL ホスト | `localhost` |
| `DB_PORT` | MySQL ポート | `43306` |
| `DB_SCHEMA` | データベース名 | `sandbox_local` |
| `DB_USER` | DB ユーザー | `sandbox_app` |
| `DB_PASSWORD` | DB パスワード | `s4ndb0x_app` |
| `REDIS_HOST` | Redis ホスト | `localhost` |
| `REDIS_PORT` | Redis ポート | `46379` |
| `JWT_ISSUER1` | Cognito URL | `https://cognito-idp.ap-northeast-1.amazonaws.com/...` |
| `JWT_AUDIENCE1/2/3` | Cognito App Client ID | — |
| `JWT_ORIGIN1/2` | 許可オリジン | `http://localhost:3000` |
| `FX_RATE_URL` | Gaitame レート API ベース URL | `https://api.gaitame.com` |

### 3. ビルド & 起動

```bash
# ビルド
go build ./...
# または
make build

# API サーバー起動（APP_ENV で環境を指定）
APP_ENV=local go run ./cmd/main.go   # .env.local を読み込む
APP_ENV=docker go run ./cmd/main.go  # .env.docker を読み込む
# または
APP_ENV=local make run

# Lint チェック（golangci-lint v2）
golangci-lint run ./...
# または
make lint
```

---

## API Documentation

- API Docs: [docs/api.md](./docs/api.md)
- アーキテクチャ詳細: [docs/architecture.md](./docs/architecture.md)
