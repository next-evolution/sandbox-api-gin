# API エンドポイント一覧

ベースパス: `/v1`

すべてのエンドポイントは JWT Middleware → Auth Middleware を通過する。

## 認証 `/v1/auth`

| メソッド | パス | 説明 | Redisセッション |
|---|---|---|---|
| POST | /v1/auth/login | ログイン | 不要（JWT検証のみ） |
| POST | /v1/auth/logout-api | ログアウト | 必要 |

### POST /v1/auth/login

- `Authorization: Bearer <token>` 必須
- リクエストボディの `email`（Base64）と JWT の email が一致しない場合は 401
- DBでユーザーを照会し、blocked なら 401
- 成功時は Redis にセッション保存、UserDto を返す

### POST /v1/auth/logout-api

- `Authorization: Bearer <token>` 必須・Redis セッション必須
- リクエストボディの `userId`（Base64）のセッションを Redis から削除
- エラーは握りつぶす（クライアント側で Cognito ログアウト済みのため）

---

## ユーザー `/v1/user`

| メソッド | パス | 説明 | 備考 |
|---|---|---|---|
| GET | /v1/user | プロフィール取得 | 未承認時は ReturnCodeWarn |
| POST | /v1/user | ユーザー登録 | 重複時は DuplicateError(400) |
| PUT | /v1/user/:userId | ニックネーム更新 | 他ユーザーは ForbiddenError(403) |

### GET /v1/user

- `authUser.Sub` でユーザーを検索
- ユーザー不存在: 404
- `approved = false`: 200 + ReturnCodeWarn + `"利用承認待ちです"`
- `approved = true`: 200 + ReturnCodeOk + UserDto

### POST /v1/user

- `authUser.Sub` / `authUser.Email` / リクエストボディの `nickName` でユーザー新規作成
- `nickName`: required, max=50
- 既存ユーザーの場合: 400 DuplicateError
- 新規作成時は `approved=false` / `admin=false` / `blocked=false`

### PUT /v1/user/:userId

- `:userId` は Base64 エンコード済み（`decodeBase64UserID()` でデコード）
- デコードした userId が `authUser.Sub` と不一致: 403 ForbiddenError
- `nickName`: required, max=50
- ユーザー不存在: 404

---

## レスポンス共通仕様

### ReturnCode

| 値 | 定数 | 意味 |
|---|---|---|
| 0 | ReturnCodeOk | 正常 |
| 1 | ReturnCodeWarn | 警告（処理は成功だが注意あり） |
| 2 | ReturnCodeError | エラー |

### エラーレスポンス

```json
{
  "status": 401,
  "error": "UNAUTHORIZED",
  "message": "..."
}
```

| エラー型 | HTTPステータス | error フィールド |
|---|---|---|
| AuthenticationError | 401 | UNAUTHORIZED |
| ForbiddenError | 403 | FORBIDDEN |
| NotFoundError | 404 | NOT_FOUND |
| DuplicateError / InsertError / UpdateError | 400 | BAD_REQUEST |
| その他 | 500 | INTERNAL_SERVER_ERROR |
