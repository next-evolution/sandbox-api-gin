# Lintルール

golangci-lint v2を使用。設定ファイル: `.golangci.yml`

実行: `golangci-lint run ./...` または `make lint`

## 有効なリンター

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
