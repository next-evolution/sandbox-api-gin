# docs-check

`docs/api.md` と `docs/architecture.md` が実装と乖離していないかチェックする。

## 手順

以下を順番に実施し、差異をまとめてレポートする。

### 1. ルーター vs `docs/api.md`

- `internal/api/router/router.go` を読む
- 各エンドポイント（HTTPメソッド・パス・認証必須/不要）を抽出し、`docs/api.md` の記載と照合する
- 認証不要グループ（`v1Public`）と認証必須グループ（`v1`）の分類が `docs/api.md` と一致しているか確認する
- レスポンス形式（`ApiResponse` でラップされているか、直接返却か）を代表コントローラーで確認し、`docs/api.md` の共通仕様と照合する

### 2. ミドルウェア vs `docs/architecture.md`

- `internal/api/middleware/jwt_middleware.go` を読む
- `internal/api/middleware/auth_middleware.go` を読む
- 実際のフロー（トークン取得・JWT検証・Redisセッション取得・`authUser` セット）と `docs/architecture.md` の「認証フロー詳細」セクションを照合する

### 3. Config（環境変数） vs `CLAUDE.md`

- `internal/config/config.go` を読む
- 全フィールドと対応する環境変数名・デフォルト値を抽出する
- `CLAUDE.md` の「必要な環境変数」テーブルと照合する（変数名の追加・削除・変更がないか）

### 4. ディレクトリ構成 vs `docs/architecture.md`

- `internal/` 配下のディレクトリ一覧を取得する
- 実際の構成と `docs/architecture.md` の「ディレクトリ構成」セクションを照合する（新規追加・削除されたパッケージがないか）

### 5. パッケージ命名 vs `docs/architecture.md`

- `internal/application/usecase/fx/` 配下・`internal/infrastructure/infradb/fx/` 配下のパッケージ名を確認する
- `docs/architecture.md` の「パッケージ命名規約（FX機能）」テーブルと照合する

## レポート形式

差異がある場合:

```
## 差異あり

### [確認項目名]
- **項目**: （ドキュメントの記載）
- **実装**: （実際の実装）
- **修正案**: （推奨される修正内容）
```

差異がない場合:

```
## 差異なし
docs/*.md と実装の間に乖離は見つかりませんでした。
```
