# TODO: 未実装エンドポイント

Java SpringBoot（`docs/ref_src`）と Go Gin の実装差異。

## ZigZag（4エンドポイント）

| メソッド | パス |
|---|---|
| POST | `/v1/fx/zigzag` |
| POST | `/v1/fx/zigzag/status` |
| POST | `/v1/fx/zigzag/generate` |
| POST | `/v1/fx/zigzag/bar-data` |

## BarData CSV インポート（1エンドポイント）

| メソッド | パス |
|---|---|
| POST | `/v1/fx/bar-data/import-csv/{symbol}/{barType}/{skipLatest}` |

## Admin 管理（6エンドポイント）

| メソッド | パス |
|---|---|
| GET | `/v1/admin/master-refresh` |
| PUT | `/v1/admin/master-refresh` |
| POST | `/v1/admin/users` |
| PUT | `/v1/admin/users/approved/{userId}` |
| PUT | `/v1/admin/users/block/{userId}` |
| PUT | `/v1/admin/users/admin/{userId}` |
