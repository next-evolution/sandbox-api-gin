# Tips

開発中に気づいた、重要ではないが知っておくと効率的な注意点をまとめる。

---

## Go + MySQL: NULL許容カラムとGoの型不一致

### 症状

```
sql: Scan error on column index N, name "xxx": converting NULL to string is unsupported
```

実行時にのみ発生し、コンパイルでは検出できない。

### 原因

DBのNULL許容カラムをGoの `string` 型でスキャンしようとすると上記エラーになる。

```go
// DBカラムがNULL許容でも、Goのstructはコンパイルが通る
type record struct {
    Description string `db:"description"` // NULLが来たら実行時エラー
}
```

### 対処パターン

**NULLを空文字として扱ってよい場合（推奨）** → SQLで `COALESCE` を使う

```sql
-- NULLなら空文字、値があればそのまま返す
COALESCE(t.description, '') AS description
```

Go側のstructはそのまま `string` で済み、変更が最小限になる。

**NULLと空文字を区別したい場合** → `sql.NullString` を使う

```go
type record struct {
    Description sql.NullString `db:"description"`
}
// NULL判定: rec.Description.Valid
// 値取得:   rec.Description.String
```

### 発生しやすいケース

1. **任意入力のテキストカラム** — `description`, `memo`, `sub_title` など、登録時に省略できるカラムはNULL許容になりやすい。
2. **後から追加したカラム** — ALTER TABLEで追加した場合、既存行が `DEFAULT NULL` になる。
3. **JOINした側のカラム** — `LEFT JOIN` でマッチしない行は結合先の全カラムがNULLになる。`INNER JOIN` でも結合先テーブル自体のカラムがNULL許容なら同様。

### 予防策

- struct作成時にテーブル定義を確認し、NULL許容カラム（`DEFAULT NULL` or `NULL`）を洗い出す。
- JOINするクエリでは、結合先テーブルのカラムを意識して確認する。
