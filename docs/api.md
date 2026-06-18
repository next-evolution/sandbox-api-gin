# API エンドポイント一覧

ベースパス: `/api/v1`

認証必須エンドポイントは JWT Middleware → Auth Middleware を通過する。
認証不要エンドポイント（`@PublicApi` 相当）はミドルウェアを通過しない。

## 認証 `/api/v1/auth`

| メソッド | パス | 説明 | Redisセッション |
|---|---|---|---|
| POST | /api/v1/auth/login | ログイン | 不要（JWT検証のみ） |
| POST | /api/v1/auth/logout-api | ログアウト | 必要 |

### POST /api/v1/auth/login

- `Authorization: Bearer <token>` 必須
- リクエストボディの `email`（Base64）と JWT の email が一致しない場合は 401
- DBでユーザーを照会し、blocked なら 401
- 成功時は Redis にセッション保存、UserDto を返す

### POST /api/v1/auth/logout-api

- `Authorization: Bearer <token>` 必須・Redis セッション必須
- リクエストボディの `userId`（Base64）のセッションを Redis から削除
- エラーは握りつぶす（クライアント側で Cognito ログアウト済みのため）

---

## ユーザー `/api/v1/user`

| メソッド | パス | 説明 | 備考 |
|---|---|---|---|
| GET | /api/v1/user | プロフィール取得 | 未承認時は ReturnCodeWarn |
| POST | /api/v1/user | ユーザー登録 | 重複時は DuplicateError(400) |
| PUT | /api/v1/user/:userId | ニックネーム更新 | 他ユーザーは ForbiddenError(403) |

### GET /api/v1/user

- `authUser.Sub` でユーザーを検索
- ユーザー不存在: 404
- `approved = false`: 200 + ReturnCodeWarn + `"利用承認待ちです"`
- `approved = true`: 200 + ReturnCodeOk + UserDto

### POST /api/v1/user

- `authUser.Sub` / `authUser.Email` / リクエストボディの `nickName` でユーザー新規作成
- `nickName`: required, max=50
- 既存ユーザーの場合: 400 DuplicateError
- 新規作成時は `approved=false` / `admin=false` / `blocked=false`

### PUT /api/v1/user/:userId

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

---

## FX `/api/v1/fx`

### マスターリスト（認証不要）

| メソッド | パス | 説明 |
|---|---|---|
| GET | /api/v1/fx/master-list/symbol/:symbolType | シンボル一覧取得 |
| GET | /api/v1/fx/master-list/country | 国一覧取得 |
| GET | /api/v1/fx/master-list/currency-pair | 通貨ペア一覧取得（symbolType=Trade固定） |
| GET | /api/v1/fx/master-list/currency-index | 通貨インデックス一覧取得（symbolType=Analyze固定） |
| GET | /api/v1/fx/master-list/economic-indicator/:countryCode | 経済指標一覧取得 |

すべて認証不要（JWT・Authミドルウェアを通過しない）。

#### GET /api/v1/fx/master-list/symbol/:symbolType

- `:symbolType`: `"Trade"` または `"Analyze"`（それ以外は 400）
- `fx_symbol` テーブルから `symbol_type` が一致するレコードを返す

#### GET /api/v1/fx/master-list/country

- `fx_country` テーブルから全国一覧を返す（key=code、value=name_short）

#### GET /api/v1/fx/master-list/currency-pair

- `fx_symbol` テーブルの `symbol_type='Trade'` を返す（`/symbol/Trade` のショートカット）

#### GET /api/v1/fx/master-list/currency-index

- `fx_symbol` テーブルの `symbol_type='Analyze'` を返す（`/symbol/Analyze` のショートカット）

#### GET /api/v1/fx/master-list/economic-indicator/:countryCode

- `:countryCode`: 国コード（例: `"JP"`）または `"ALL"`（全件取得）
- `fx_economic_indicator` テーブルから取得（key=id文字列、value=name）

#### レスポンス形式（全エンドポイント共通）

ApiResponse ラッパーなし。`KeyValue` の配列を直接返す。

```json
[
  {"key": "USDJPY", "value": "ドル円"},
  {"key": "EURUSD", "value": "ユーロドル"}
]
```

---

### シンボル管理（認証必須）

| メソッド | パス | 説明 |
|---|---|---|
| GET | /api/v1/fx/symbol/currency-pair-list | 通貨ペア一覧（SymbolDto配列） |
| GET | /api/v1/fx/symbol/currency-index-list | 通貨インデックス一覧（SymbolDto配列） |
| POST | /api/v1/fx/symbol/search | シンボル検索（ページネーション） |
| POST | /api/v1/fx/symbol | シンボル追加 |
| GET | /api/v1/fx/symbol/:symbol | シンボル取得 |
| PUT | /api/v1/fx/symbol/:symbol | シンボル更新 |

#### GET /api/v1/fx/symbol/currency-pair-list

- `symbolType=Trade` のシンボルを最大500件返す
- レスポンス: `SymbolDto` の配列（ApiResponseラッパーなし）

#### GET /api/v1/fx/symbol/currency-index-list

- `symbolType=Analyze` のシンボルを最大500件返す
- レスポンス: `SymbolDto` の配列（ApiResponseラッパーなし）

#### POST /api/v1/fx/symbol/search

リクエスト:
```json
{"symbolType": "Trade", "page": 1, "size": 20}
```
- `symbolType`: `"Trade"` または `"Analyze"`（必須、それ以外は400）
- `page`: 1始まり（必須、min=1）
- `size`: 取得件数（必須、min=1）

レスポンス:
```json
{
  "returnCode": 0,
  "totalCount": 10,
  "searchCount": 10,
  "totalPage": 1,
  "list": [{"symbol":"USDJPY","symbolType":"Trade","name":"ドル円","validScale":3,"targetVolatility":0.005,"sortOrder":100}]
}
```

#### POST /api/v1/fx/symbol

リクエスト:
```json
{"symbol": {"symbol":"USDJPY","symbolType":"Trade","name":"ドル円","validScale":3,"targetVolatility":0.005,"sortOrder":100}}
```
- `symbol.symbol`: 必須
- `symbol.symbolType`: 必須（`"Trade"` または `"Analyze"`）
- `symbol.name`: 必須
- 重複シンボル: 400 DuplicateError
- 成功: 200（ボディなし）

#### GET /api/v1/fx/symbol/:symbol

- `:symbol`: シンボルコード（例: `USDJPY`）
- 存在しない場合: 404 NotFoundError
- レスポンス: `SymbolDto`（ApiResponseラッパーなし）

#### PUT /api/v1/fx/symbol/:symbol

- `:symbol`: 更新対象のシンボルコード（`baseSymbol`）
- リクエストボディは POST /api/v1/fx/symbol と同じ構造
- `baseSymbol == req.symbol.symbol` の場合: 同一シンボルの属性更新
- `baseSymbol != req.symbol.symbol` の場合: シンボルコードの変更（新コードが既存の場合は400 DuplicateError）
- 対象不存在: 400 UpdateError
- 成功: 200（ボディなし）

#### SymbolDto

```json
{
  "symbol": "USDJPY",
  "symbolType": "Trade",
  "name": "ドル円",
  "validScale": 3,
  "targetVolatility": 0.005,
  "sortOrder": 100
}
```

---

### 国管理（認証必須）

| メソッド | パス | 説明 |
|---|---|---|
| POST | /api/v1/fx/country/search | 国一覧検索（ページネーション） |
| POST | /api/v1/fx/country | 国追加 |
| GET | /api/v1/fx/country/:code | 国取得 |
| PUT | /api/v1/fx/country/:code | 国更新 |

#### POST /api/v1/fx/country/search

リクエスト:
```json
{"page": 1, "size": 20}
```
- `page`: 1始まり（必須、min=1）
- `size`: 取得件数（必須、min=1）

レスポンス:
```json
{
  "returnCode": 0,
  "totalCount": 10,
  "searchCount": 10,
  "totalPage": 1,
  "list": [{"code":"JP","name":"日本","currencyCode":"JPY","nameEn":"Japan","nameShort":"日本","sortOrder":1}]
}
```

#### POST /api/v1/fx/country

リクエスト:
```json
{"country": {"code":"JP","name":"日本","currencyCode":"JPY","nameEn":"Japan","nameShort":"日本","sortOrder":1}}
```
- `country.code`: 必須
- `country.name`: 必須
- `country.currencyCode`: 必須
- `country.nameEn`: 必須
- `country.nameShort`: 必須
- 重複コード: 400 DuplicateError
- 成功: 200（ボディなし）

#### GET /api/v1/fx/country/:code

- `:code`: ISO 3166-1 alpha-2（例: `"JP"`）
- 存在しない場合: 404 NotFoundError
- レスポンス: `CountryDto`（ApiResponseラッパーなし）

#### PUT /api/v1/fx/country/:code

- `:code`: 更新対象の国コード（`baseCode`）
- リクエストボディは POST /api/v1/fx/country と同じ構造
- `baseCode == req.country.code` の場合: 属性更新（不存在は400 UpdateError）
- `baseCode != req.country.code` の場合: コード変更（新コードが既存の場合は400 DuplicateError）
- 成功: 200（ボディなし）

#### CountryDto

```json
{
  "code": "JP",
  "name": "日本",
  "currencyCode": "JPY",
  "nameEn": "Japan",
  "nameShort": "日本",
  "sortOrder": 1
}
```

---

### 夏時間管理（認証必須）

| メソッド | パス | 説明 |
|---|---|---|
| POST | /api/v1/fx/summer-time/search | 夏時間一覧検索（ページネーション） |
| POST | /api/v1/fx/summer-time | 夏時間追加 |
| GET | /api/v1/fx/summer-time/:targetYear | 夏時間取得 |
| PUT | /api/v1/fx/summer-time/:targetYear | 夏時間更新 |

#### POST /api/v1/fx/summer-time/search

リクエスト:
```json
{"page": 1, "size": 20}
```

レスポンス:
```json
{
  "returnCode": 0,
  "totalCount": 5,
  "searchCount": 5,
  "totalPage": 1,
  "list": [{"targetYear": 2024, "applyStart": "2024-03-10", "applyEnd": "2024-11-03"}]
}
```

#### POST /api/v1/fx/summer-time

リクエスト:
```json
{"summerTime": {"targetYear": 2024, "applyStart": "2024-03-10", "applyEnd": "2024-11-03"}}
```
- `summerTime.targetYear`: 必須
- 重複年: 400 DuplicateError
- 成功: 200（ボディなし）

#### GET /api/v1/fx/summer-time/:targetYear

- `:targetYear`: 対象年（整数）
- 存在しない場合: 404 NotFoundError
- レスポンス: `SummerTimeDto`（ApiResponseラッパーなし）

#### PUT /api/v1/fx/summer-time/:targetYear

- `:targetYear`: 更新対象の年（`baseYear`）
- リクエストボディは POST /api/v1/fx/summer-time と同じ構造
- `baseYear == req.summerTime.targetYear` の場合: 属性更新（不存在は400 UpdateError）
- `baseYear != req.summerTime.targetYear` の場合: 年変更（新年が既存の場合は400 DuplicateError）
- 成功: 200（ボディなし）

#### SummerTimeDto

```json
{
  "targetYear": 2024,
  "applyStart": "2024-03-10",
  "applyEnd": "2024-11-03"
}
```

---

### バーデータ管理（認証必須）

| メソッド | パス | 説明 |
|---|---|---|
| POST | /api/v1/fx/bar-data | バーデータ検索（ページネーション） |
| GET | /api/v1/fx/bar-data/:symbolType/:barType | バーデータ件数ステータス |
| POST | /api/v1/fx/bar-data/import-csv/:symbol/:barType/:skipLatest | バーデータCSVインポート |

#### POST /api/v1/fx/bar-data

リクエスト:
```json
{"barType": "15M", "symbol": "USDJPY", "barDateFrom": "20260101", "barDateTo": "20260131", "sortAsc": false, "page": 1, "size": 100}
```
- `barType`: 必須（`"15M"`, `"1H"`, `"4H"`, `"1D"` のいずれか）
- `symbol`: 必須
- `barDateFrom` / `barDateTo`: 任意（`yyyyMMdd` 形式）
- `sortAsc`: 任意（デフォルト: false = 降順）
- `page`: 1始まり（必須、min=1）
- `size`: 取得件数（必須、min=1）

レスポンス:
```json
{
  "returnCode": 0,
  "totalCount": 500,
  "searchCount": 500,
  "totalPage": 5,
  "list": [{"symbol":"USDJPY","barDateTime":"2026-01-31 23:45:00","openPrice":154.123,"highPrice":154.200,...}]
}
```

#### GET /api/v1/fx/bar-data/:symbolType/:barType

- `:symbolType`: `"Trade"` または `"Analyze"`
- `:barType`: `"15M"`, `"1H"`, `"4H"`, `"1D"` のいずれか
- シンボルごとの件数とデータ期間を返す

レスポンス（ApiResponseラッパーなし）:
```json
[
  {"symbol":"USDJPY","barDateTime":null,"existsCount":1000,"message":"2025-01-01 00:00~2026-01-31 23:45"},
  {"symbol":"EURUSD","barDateTime":null,"existsCount":0,"message":null}
]
```

#### POST /api/v1/fx/bar-data/import-csv/:symbol/:barType/:skipLatest

multipart/form-data でCSVファイルをアップロードし、バーデータをインポートする。

- `:symbol`: シンボル（例: `"USDJPY"`）
- `:barType`: `"15M"`, `"1H"`, `"4H"`, `"1D"` のいずれか
- `:skipLatest`: `"true"` で最新1件をスキップ（未確定足を除外する場合）
- フォームフィールド `uploadFile`: CSVファイル

CSVファイル名に `{symbol}_{keyword}`（例: `USDJPY_240`）が含まれていない場合はエラー。

| barType | keyword |
|---|---|
| 15M | 15 |
| 1H | 60 |
| 4H | 240 |
| 1D | 1D |

レスポンス（ApiResponseラッパーなし）:
```json
{
  "symbol": "USDJPY",
  "barDateTime": "2026-03-01 00:00:00",
  "fileName": "USDJPY_240_20260301.csv",
  "fileSize": 102400,
  "resultStatus": "OK",
  "readCount": 1000,
  "existsCount": 950,
  "insertCount": 50,
  "differenceCount": 0
}
```
- `resultStatus`: `"OK"`（新規挿入あり）/ `"SKIP"`（全件既存）/ `"ERROR"`（ファイル名不正・整合性チェックNG）

---

### トレードシミュレーション（認証必須）

| メソッド | パス | 説明 |
|---|---|---|
| POST | /api/v1/fx/trade/simulation | トレードシミュレーション実行 |

### POST /api/v1/fx/trade/simulation

JWT Middleware → Auth Middleware を通過する（認証必須）。

#### リクエスト

```json
{
  "riskAmount": 10000,
  "firstLotRatio": 30,
  "entry": {
    "id": null,
    "tradeVersion": "v1",
    "entryType": "F3",
    "symbol": "USDJPY",
    "tradeType": "L",
    "contractAt": "2026-01-02T11:22:33+09:00",
    "fibonacciType": "382",
    "fibonacciBar": "4H",
    "contractPrice": 149.500,
    "lossPrice": 149.200,
    "positionRatio": 0,
    "priceJpy": 149.500,
    "lot": null,
    "settlementAmount": 0,
    "lossPips": 0,
    "settlementRatio": null,
    "comment": null,
    "imagePath": null
  },
  "positionList": [
    {
      "id": null,
      "positionNumber": 1,
      "settlementPrice": 0,
      "settlementPips": 0,
      "settlementRatio": null,
      "lot": null,
      "profitAmount": 0,
      "lossAmount": 0
    }
  ]
}
```

**フィールド仕様:**

| フィールド | 必須 | 説明 |
|---|---|---|
| `riskAmount` | — | リスク金額（整数、最大6桁）。0の場合はデフォルト10,000円 |
| `firstLotRatio` | ✓ (>0) | 第1ポジションのロット比率（%単位、例: 30 = 30%）。0の場合はデフォルト30% |
| `entry.tradeVersion` | ✓ | トレードバージョン |
| `entry.symbol` | ✓ | 通貨ペア（例: USDJPY、EURUSD）。JPY末尾でドル円フラグが変わる |
| `entry.fibonacciType` | ✓ | フィボナッチタイプ |
| `entry.fibonacciBar` | ✓ | フィボナッチバー |
| `entry.contractAt` | — | 約定日時（RFC 3339形式）。15分足に切り捨てて価格取得に使用 |
| `entry.contractPrice` | — | 約定価格。0の場合はDBから取得した価格を使用 |
| `entry.lossPrice` | — | 損切価格。0の場合は自動計算（ドル円: ±0.003、円建て: ±0.3） |
| `entry.priceJpy` | — | JPY換算レート。0の場合はDBから取得 |
| `entry.lot` / `entry.settlementRatio` | — | null許容。null時は0として扱う |
| `positionList[].positionNumber` | ✓ (>0) | ポジション番号 |
| `positionList[].settlementPrice` | — | 決済価格。リスト先頭が0の場合は全ポジションにデフォルト計算値を設定 |

**positionList の挙動:**

- 先頭要素の `settlementPrice == 0` の場合: 全ポジションに `contractPrice` からの自動計算値を設定（1件: +0.6pip、2件: +0.6/+0.9pip、3件: +0.6/+0.9/+1.2pip）
- 先頭要素の `settlementPrice > 0` の場合: `settlementPrice > 0` のポジションのみ有効とし、`positionNumber` を振り直す

#### レスポンス

```json
{
  "returnCode": 0,
  "entry": {
    "id": -1,
    "tradeVersion": "v1",
    "entryType": "F3",
    "symbol": "USDJPY",
    "tradeType": "L",
    "contractAt": "2026-01-02 11:15:00",
    "fibonacciType": "382",
    "fibonacciBar": "4H",
    "contractPrice": 149.500,
    "lossPrice": 149.200,
    "positionRatio": 0,
    "priceJpy": 149.500,
    "lot": 0.33,
    "settlementAmount": 5940,
    "lossPips": 300,
    "settlementRatio": 0.60,
    "comment": null,
    "imagePath": null
  },
  "positionList": [
    {
      "id": null,
      "positionNumber": 1,
      "settlementPrice": 149.600,
      "settlementPips": 100,
      "settlementRatio": 0.33,
      "lot": 0.33,
      "profitAmount": 3300,
      "lossAmount": 9900
    }
  ]
}
```

**レスポンスの計算内容:**

| フィールド | 計算内容 |
|---|---|
| `entry.contractAt` | 15分足に切り捨てた時刻（形式: `"yyyy-MM-dd HH:mm:ss"`） |
| `entry.id` | リクエストのidがnull/0の場合は -1 を設定 |
| `entry.lot` | `riskAmount / lossPips / 100,000` で計算（ドル建ては `/ priceJpy` も除算） |
| `entry.lossPips` | `|contractPrice - lossPrice|` のpips値 |
| `entry.settlementAmount` | 全ポジションの `profitAmount` 合計 |
| `entry.settlementRatio` | `settlementAmount / lossAmount(totalLot)` |
| `position.lot` | ポジション数に応じてロットを分割（1件: totalLot、2〜3件: `firstLotRatio` で按分） |
| `position.settlementPips` | `|settlementPrice - contractPrice|` のpips値 |
| `position.profitAmount` | 損益金額（損失時は負値） |
| `position.lossAmount` | そのポジションのロットに対する最大損失額 |
| `position.settlementRatio` | `profitAmount / lossAmount` |

---

### ZigZag（認証必須）

| メソッド | パス | 説明 |
|---|---|---|
| POST | /api/v1/fx/zigzag | ZigZag検索 |
| POST | /api/v1/fx/zigzag/status | ZigZagステータス一覧取得 |
| POST | /api/v1/fx/zigzag/generate | ZigZag生成 |
| POST | /api/v1/fx/zigzag/bar-data | ZigZagバーデータ取得 |

#### POST /api/v1/fx/zigzag（ZigZag検索）

リクエスト:

| フィールド | 型 | 必須 | 説明 |
|---|---|---|---|
| `barType` | string | ○ | `"15M"/"1H"/"4H"/"1D"` |
| `symbol` | string | ○ | シンボル |
| `depth` | int16 | ○ | ZigZag深さ（≥1） |
| `barDateTimeMin` | datetime | ○ | 検索期間（開始） |
| `barDateTimeMax` | datetime | ○ | 検索期間（終了） |
| `wave` | int | - | Waveフィルタ |
| `previousWave` | int | - | 前waveフィルタ |
| `nextWave` | int | - | 次waveフィルタ |
| `next2Wave` | int | - | 2つ先waveフィルタ |
| `direction4h200` / `direction4h75` / `direction4h20` | int | - | 4H SMA方向フィルタ |
| `direction1h200` / `direction15m200` | int | - | 1H/15M SMA方向フィルタ |
| `wave4h` | int | - | 4H waveフィルタ |
| `directionTarget4h200` | int | - | Target4h SMA200方向フィルタ |
| `page` | int | ○ | ページ（≥1） |
| `size` | int | ○ | ページサイズ（≥1） |

direction フィルタ値: `999`=全件、`0`=ニュートラル、`1`/`2`=上昇、`-1`/`-2`=下降

レスポンス: `ZigZagSearchResponse`（returnCode/totalCount/searchCount/totalPage/list）

#### POST /api/v1/fx/zigzag/status

リクエスト: `symbolType`（Trade/Analyze）、`barType`、`depth`

#### POST /api/v1/fx/zigzag/generate

リクエスト: `symbol`、`barType`、`depth`、`barDateTime`、`loadSize`（処理件数上限）

warn=true の場合 returnCode=1（Warn）、message にエラー理由を返す。

#### POST /api/v1/fx/zigzag/bar-data

リクエスト: `barType`、`symbol`、`depth`、`waveStart`、`wave`

---

### 経済指標管理（認証必須）

| メソッド | パス | 説明 |
|---|---|---|
| POST | /api/v1/fx/economic-indicator/search | 経済指標検索（ページネーション） |
| POST | /api/v1/fx/economic-indicator | 経済指標追加 |
| GET | /api/v1/fx/economic-indicator/:countryCode/:code | 経済指標取得 |
| PUT | /api/v1/fx/economic-indicator/:countryCode/:code | 経済指標更新 |

#### POST /api/v1/fx/economic-indicator/search

リクエスト:
```json
{"page": 1, "size": 20, "countryCode": "JP", "importance": "A", "name": "GDP"}
```
- `page`: 1始まり（必須、min=1）
- `size`: 取得件数（必須、min=1）
- `countryCode` / `importance` / `name`: 任意（フィルタ）

レスポンス:
```json
{
  "returnCode": 0,
  "totalCount": 10,
  "searchCount": 10,
  "totalPage": 1,
  "list": [{"code":"GDP","countryCode":"JP","name":"国内総生産","importance":"A","description":"","unitOfValue":"%","countryName":"日本","countryNameShort":"日本"}]
}
```

#### POST /api/v1/fx/economic-indicator

リクエスト:
```json
{"indicator": {"code":"GDP","countryCode":"JP","name":"国内総生産","importance":"A","description":"","unitOfValue":"%"}}
```
- `indicator.code`: 必須
- `indicator.countryCode`: 必須
- `indicator.name`: 必須
- `indicator.importance`: 必須
- 重複（countryCode + code）: 400 DuplicateError
- 成功: 200（ボディなし）

#### GET /api/v1/fx/economic-indicator/:countryCode/:code

- 存在しない場合: 404 NotFoundError
- レスポンス: `EconomicIndicatorDto`（ApiResponseラッパーなし）

#### PUT /api/v1/fx/economic-indicator/:countryCode/:code

- `:countryCode` + `:code` が更新対象の現在キー
- リクエストボディは POST と同じ構造
- 対象不存在: 400 UpdateError
- 成功: 200（ボディなし）

#### EconomicIndicatorDto

```json
{
  "code": "GDP",
  "countryCode": "JP",
  "name": "国内総生産",
  "importance": "A",
  "description": "",
  "unitOfValue": "%",
  "countryName": "日本",
  "countryNameShort": "日本"
}
```

---

### 経済指標データ管理（認証必須）

| メソッド | パス | 説明 |
|---|---|---|
| POST | /api/v1/fx/economic-indicator-data/search | 経済指標データ検索（ページネーション） |
| POST | /api/v1/fx/economic-indicator-data | 経済指標データ追加 |
| GET | /api/v1/fx/economic-indicator-data/:countryCode/:code/:publication | 経済指標データ取得 |
| PUT | /api/v1/fx/economic-indicator-data/:countryCode/:code/:publication | 経済指標データ更新 |
| POST | /api/v1/fx/economic-indicator-data/import-text | 経済指標データテキストインポート |

#### POST /api/v1/fx/economic-indicator-data/search

リクエスト:
```json
{"page": 1, "size": 20, "code": "GDP", "importance": "A", "countryCode": "JP", "publicationBaseDate": "2026-01-01", "sortAsc": false}
```
- `page`: 1始まり（必須、min=1）
- `size`: 取得件数（必須、min=1）
- `code` / `importance` / `countryCode` / `publicationBaseDate`: 任意フィルタ
- `sortAsc`: 任意（デフォルト: false = 降順）

レスポンス:
```json
{
  "returnCode": 0,
  "totalCount": 100,
  "searchCount": 100,
  "totalPage": 5,
  "list": [{"code":"GDP","countryCode":"JP","name":"国内総生産","importance":"A","description":"","publication":"2026-01-15 08:50:00","publicationDate":"2026-01-15","publicationTime":"08:50","dayOfWeek":4,"subTitle":"","resultValue":"0.5","forecastValue":"0.4","previousValue":"0.3","unitOfValue":"%","memo":"","countryName":"日本","countryNameShort":"日本"}]
}
```

#### POST /api/v1/fx/economic-indicator-data

リクエスト:
```json
{"data": {"code":"GDP","countryCode":"JP","publication":"2026-01-15 08:50:00","resultValue":"0.5","forecastValue":"0.4","previousValue":"0.3"}}
```
- 重複（countryCode + code + publication）: 400 DuplicateError
- 成功: 200（ボディなし）

#### GET /api/v1/fx/economic-indicator-data/:countryCode/:code/:publication

- `:publication`: `yyyy-MM-dd HH:mm:ss` 形式
- 存在しない場合: 404 NotFoundError
- レスポンス: `EconomicIndicatorDataDto`（ApiResponseラッパーなし）

#### PUT /api/v1/fx/economic-indicator-data/:countryCode/:code/:publication

- `:publication`: `yyyy-MM-dd HH:mm:ss` 形式
- 対象不存在: 400 UpdateError
- 成功: 200（ボディなし）

#### POST /api/v1/fx/economic-indicator-data/import-text

multipart/form-data で経済指標スケジュールテキストを複数ファイル一括インポートする。

- フォームフィールド `uploadFileList`: テキストファイル（複数可）
- ファイル名形式: `{year}_{countryName}_{importance}.txt`（例: `2026_JP_A.txt`）
- ファイル名からインポート年・国・重要度を判別して DB に登録

レスポンス（ApiResponseラッパーなし、ファイルごとの結果配列）:
```json
[
  {"fileName":"2026_JP_A.txt","fileSize":10240,"readCount":150,"resultStatus":"OK"}
]
```

#### EconomicIndicatorDataDto

```json
{
  "code": "GDP",
  "countryCode": "JP",
  "name": "国内総生産",
  "importance": "A",
  "description": "",
  "publication": "2026-01-15 08:50:00",
  "publicationDate": "2026-01-15",
  "publicationTime": "08:50",
  "dayOfWeek": 4,
  "subTitle": "",
  "resultValue": "0.5",
  "forecastValue": "0.4",
  "previousValue": "0.3",
  "unitOfValue": "%",
  "memo": "",
  "countryName": "日本",
  "countryNameShort": "日本"
}
```

---

## Admin 管理 API

全エンドポイントは JWT + Auth ミドルウェア必須。`authUser.Admin == true` でなければ 403。

### マスターキャッシュ

#### GET /api/v1/admin/master-refresh

Redisの `master*` キーの件数ステータスを返す。

レスポンス: `ApiResponse`（returnCode=0、message=ステータス文字列）

#### PUT /api/v1/admin/master-refresh

国・シンボル・経済指標のキャッシュをリフレッシュし、`price*` キーを削除してステータスを返す。

レスポンス: `ApiResponse`（returnCode=0、message=更新後ステータス文字列）

### ユーザー管理

#### POST /api/v1/admin/users

リクエスト:

| フィールド | 型 | 必須 | 説明 |
|---|---|---|---|
| `page` | int | ○ | ページ（≥1） |
| `size` | int | ○ | ページサイズ（≥1） |
| `emailAddress` | string | - | メールアドレス部分一致 |
| `approved` | bool | - | 承認フラグフィルタ |

レスポンス: `UserSearchResponse`（returnCode/totalCount/searchCount/totalPage/list）

#### PUT /api/v1/admin/users/approved/:userId

パスパラメータ `:userId` は Base64 エンコード済みの userId。承認フラグを `true` にセット。

承認済みの場合は 400（DuplicateError）。

レスポンス: `UserResponse`（returnCode=0、user）

#### PUT /api/v1/admin/users/block/:userId

リクエスト: `blocked`（bool、必須）

Block状態が変わらない場合は 400（DuplicateError）。

レスポンス: `UserResponse`（returnCode=0、user）

#### PUT /api/v1/admin/users/admin/:userId

リクエスト: `admin`（bool、必須）

管理者権限が変わらない場合は 400（DuplicateError）。

レスポンス: `UserResponse`（returnCode=0、user）
