package infradbfx

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	fxmodel "sandbox-api-gin/internal/domain/model/fx"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type MySQLBarDataRepository struct {
	db *sqlx.DB
}

func NewMySQLBarDataRepository(db *sqlx.DB) fxrepository.BarDataRepository {
	return &MySQLBarDataRepository{db: db}
}

type fxBarDataStatusRecord struct {
	Symbol          string  `db:"symbol"`
	BarDateTimeMinS *string `db:"barDateTimeMinS"`
	BarDateTimeMaxS *string `db:"barDateTimeMaxS"`
	Count           int     `db:"count"`
}

type fxBarDataRecord struct {
	Symbol      string    `db:"symbol"`
	BarDateTime time.Time `db:"barDateTime"`
	OpenPrice   float64   `db:"openPrice"`
	HighPrice   float64   `db:"highPrice"`
	LowPrice    float64   `db:"lowPrice"`
	ClosePrice  float64   `db:"closePrice"`
	Volume      int       `db:"volume"`
	HighProfit  float64   `db:"highProfit"`
	LowProfit   float64   `db:"lowProfit"`
	CloseProfit float64   `db:"closeProfit"`
	RangeProfit float64   `db:"rangeProfit"`
	RsiValue    float64   `db:"rsiValue"`
	RsiMa       float64   `db:"rsiMa"`
}

type fxBarLoadDataRecord struct {
	Symbol      string    `db:"symbol"`
	BarDateTime time.Time `db:"barDateTime"`
	OpenPrice   float64   `db:"openPrice"`
	HighPrice   float64   `db:"highPrice"`
	LowPrice    float64   `db:"lowPrice"`
	ClosePrice  float64   `db:"closePrice"`
	Volume      int       `db:"volume"`
}

type fxBarLoadSmaRecord struct {
	Symbol      string    `db:"symbol"`
	BarDateTime time.Time `db:"barDateTime"`
	SmaRange    int       `db:"smaRange"`
	SmaPrice    *float64  `db:"smaPrice"`
	SmaCross    uint8     `db:"smaCross"`
}

type fxBarLoadRsiRecord struct {
	Symbol      string    `db:"symbol"`
	BarDateTime time.Time `db:"barDateTime"`
	RsiRange    int       `db:"rsiRange"`
	RsiValue    *float64  `db:"rsiValue"`
	RsiMa       *float64  `db:"rsiMa"`
}

type fxBarCsvImportCheckRecord struct {
	ExistsCount int `db:"existsCount"`
	DiffCount   int `db:"diffCount"`
}

func (r *MySQLBarDataRepository) StatusList(ctx context.Context, symbolType string, barType fxmodel.BarType) ([]fxmodel.BarDataStatus, error) {
	query := fmt.Sprintf(`
		SELECT
			c.symbol                                              AS symbol,
			DATE_FORMAT(MIN(b.bar_date_time), '%%Y-%%m-%%d %%H:%%i') AS barDateTimeMinS,
			DATE_FORMAT(MAX(b.bar_date_time), '%%Y-%%m-%%d %%H:%%i') AS barDateTimeMaxS,
			COUNT(b.symbol)                                       AS `+"`count`"+`
		FROM fx_symbol c
		LEFT JOIN fx_bar_%s b ON b.symbol = c.symbol
		WHERE (c.deleted+0) = 0 AND c.symbol_type = ?
		GROUP BY c.symbol
		ORDER BY c.sort_order`, barType.Suffix())

	var recs []fxBarDataStatusRecord
	if err := r.db.SelectContext(ctx, &recs, query, symbolType); err != nil {
		return nil, err
	}
	result := make([]fxmodel.BarDataStatus, len(recs))
	for i, rec := range recs {
		result[i] = fxmodel.BarDataStatus{
			Symbol:          rec.Symbol,
			BarDateTimeMinS: rec.BarDateTimeMinS,
			BarDateTimeMaxS: rec.BarDateTimeMaxS,
			Count:           rec.Count,
		}
	}
	return result, nil
}

func (r *MySQLBarDataRepository) SearchCount(ctx context.Context, symbol string, barType fxmodel.BarType, barDateFrom, barDateTo string) (int, error) {
	suffix := barType.Suffix()
	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM fx_bar_%s B
		INNER JOIN fx_symbol C ON C.symbol = B.symbol
		WHERE B.symbol = ?`, suffix)
	args := []any{symbol}
	query, args = appendDateRange(query, args, barDateFrom, barDateTo)
	query += ` AND (C.deleted+0) = 0`

	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	return count, err
}

func (r *MySQLBarDataRepository) Search(ctx context.Context, symbol string, barType fxmodel.BarType, barDateFrom, barDateTo string, sortAsc bool, page, size int) ([]fxmodel.BarData, error) {
	suffix := barType.Suffix()
	query := fmt.Sprintf(`
		SELECT
			B.symbol                                   AS symbol,
			B.bar_date_time                            AS barDateTime,
			ROUND(B.open_price,   C.valid_scale)       AS openPrice,
			ROUND(B.high_price,   C.valid_scale)       AS highPrice,
			ROUND(B.low_price,    C.valid_scale)       AS lowPrice,
			ROUND(B.close_price,  C.valid_scale)       AS closePrice,
			B.volume                                   AS volume,
			ROUND(B.high_profit,  C.valid_scale)       AS highProfit,
			ROUND(B.low_profit,   C.valid_scale)       AS lowProfit,
			ROUND(B.close_profit, C.valid_scale)       AS closeProfit,
			ROUND(B.range_profit, C.valid_scale)       AS rangeProfit,
			COALESCE(R.rsi_value, 0)                   AS rsiValue,
			COALESCE(R.rsi_ma,    0)                   AS rsiMa
		FROM fx_bar_%s B
		INNER JOIN fx_symbol C ON C.symbol = B.symbol
		LEFT JOIN fx_bar_%s_rsi R
			ON R.symbol = B.symbol AND R.bar_date_time = B.bar_date_time AND R.rsi_range = 14
		WHERE B.symbol = ?`, suffix, suffix)
	args := []any{symbol}
	query, args = appendDateRange(query, args, barDateFrom, barDateTo)
	query += ` AND (C.deleted+0) = 0`

	if sortAsc {
		query += ` ORDER BY B.bar_date_time`
	} else {
		query += ` ORDER BY B.bar_date_time DESC`
	}

	offset := (page - 1) * size
	if offset > 0 {
		query += ` LIMIT ? OFFSET ?`
		args = append(args, size, offset)
	} else {
		query += ` LIMIT ?`
		args = append(args, size)
	}

	var recs []fxBarDataRecord
	if err := r.db.SelectContext(ctx, &recs, query, args...); err != nil {
		return nil, err
	}
	result := make([]fxmodel.BarData, len(recs))
	for i, rec := range recs {
		result[i] = fxmodel.BarData{
			Symbol:      rec.Symbol,
			BarDateTime: fxmodel.LocalDateTime{Time: rec.BarDateTime},
			OpenPrice:   rec.OpenPrice,
			HighPrice:   rec.HighPrice,
			LowPrice:    rec.LowPrice,
			ClosePrice:  rec.ClosePrice,
			Volume:      rec.Volume,
			HighProfit:  rec.HighProfit,
			LowProfit:   rec.LowProfit,
			CloseProfit: rec.CloseProfit,
			RangeProfit: rec.RangeProfit,
			RsiValue:    rec.RsiValue,
			RsiMa:       rec.RsiMa,
		}
	}
	return result, nil
}

func (r *MySQLBarDataRepository) DeleteLoad(ctx context.Context, symbol string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM fx_bar_load WHERE symbol = ?`, symbol)
	return err
}

func (r *MySQLBarDataRepository) DeleteLoadSma(ctx context.Context, symbol string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM fx_bar_load_sma WHERE symbol = ?`, symbol)
	return err
}

func (r *MySQLBarDataRepository) DeleteLoadRsi(ctx context.Context, symbol string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM fx_bar_load_rsi WHERE symbol = ?`, symbol)
	return err
}

func (r *MySQLBarDataRepository) BulkLoad(ctx context.Context, list []fxmodel.BarLoadData) error {
	if len(list) == 0 {
		return nil
	}
	placeholders := make([]string, len(list))
	args := make([]interface{}, 0, len(list)*7)
	for i, item := range list {
		placeholders[i] = "(?, ?, ?, ?, ?, ?, ?)"
		args = append(args, item.Symbol, item.BarDateTime,
			item.OpenPrice, item.HighPrice, item.LowPrice, item.ClosePrice, item.Volume)
	}
	query := `INSERT INTO fx_bar_load (symbol, bar_date_time, open_price, high_price, low_price, close_price, volume) VALUES ` +
		strings.Join(placeholders, ",")
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *MySQLBarDataRepository) BulkLoadSma(ctx context.Context, list []fxmodel.BarLoadSma) error {
	if len(list) == 0 {
		return nil
	}
	placeholders := make([]string, len(list))
	args := make([]interface{}, 0, len(list)*5)
	for i, item := range list {
		placeholders[i] = "(?, ?, ?, ?, ?)"
		smaCross := 0
		if item.SmaCross {
			smaCross = 1
		}
		args = append(args, item.Symbol, item.BarDateTime, item.SmaRange, item.SmaPrice, smaCross)
	}
	query := `INSERT INTO fx_bar_load_sma (symbol, bar_date_time, sma_range, sma_price, sma_cross) VALUES ` +
		strings.Join(placeholders, ",")
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *MySQLBarDataRepository) BulkLoadRsi(ctx context.Context, list []fxmodel.BarLoadRsi) error {
	if len(list) == 0 {
		return nil
	}
	placeholders := make([]string, len(list))
	args := make([]interface{}, 0, len(list)*5)
	for i, item := range list {
		placeholders[i] = "(?, ?, ?, ?, ?)"
		args = append(args, item.Symbol, item.BarDateTime, item.RsiRange, item.RsiValue, item.RsiMa)
	}
	query := `INSERT INTO fx_bar_load_rsi (symbol, bar_date_time, rsi_range, rsi_value, rsi_ma) VALUES ` +
		strings.Join(placeholders, ",")
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *MySQLBarDataRepository) DeleteLatestLoad(ctx context.Context, symbol string) error {
	// MySQL は DELETE の WHERE 句で同一テーブルをサブクエリに直接使えないためサブクエリを一段ネストする
	queries := []struct{ table, q string }{
		{"fx_bar_load", `DELETE FROM fx_bar_load WHERE symbol = ? AND bar_date_time = (SELECT max_dt FROM (SELECT MAX(bar_date_time) AS max_dt FROM fx_bar_load WHERE symbol = ?) AS t)`},
		{"fx_bar_load_sma", `DELETE FROM fx_bar_load_sma WHERE symbol = ? AND bar_date_time = (SELECT max_dt FROM (SELECT MAX(bar_date_time) AS max_dt FROM fx_bar_load_sma WHERE symbol = ?) AS t)`},
		{"fx_bar_load_rsi", `DELETE FROM fx_bar_load_rsi WHERE symbol = ? AND bar_date_time = (SELECT max_dt FROM (SELECT MAX(bar_date_time) AS max_dt FROM fx_bar_load_rsi WHERE symbol = ?) AS t)`},
	}
	for _, q := range queries {
		if _, err := r.db.ExecContext(ctx, q.q, symbol, symbol); err != nil {
			return err
		}
	}
	return nil
}

func (r *MySQLBarDataRepository) GetLatestLoadBarDateTime(ctx context.Context, symbol string) (string, error) {
	var result string
	err := r.db.GetContext(ctx, &result,
		`SELECT COALESCE(DATE_FORMAT(MAX(bar_date_time), '%Y-%m-%d %H:%i:%s'), '') FROM fx_bar_load WHERE symbol = ?`,
		symbol)
	return result, err
}

func (r *MySQLBarDataRepository) ImportCheck(ctx context.Context, barType fxmodel.BarType, symbol string) (fxmodel.BarCsvImportCheck, error) {
	query := fmt.Sprintf(`
		SELECT
			COUNT(M.bar_date_time) AS existsCount,
			COALESCE(SUM(
				CASE WHEN M.bar_date_time IS NOT NULL
				          AND (M.open_price  != L.open_price
				               OR M.high_price != L.high_price
				               OR M.low_price  != L.low_price
				               OR M.close_price != L.close_price)
				     THEN 1 ELSE 0 END
			), 0) AS diffCount
		FROM fx_bar_load L
		LEFT JOIN %s M ON M.symbol = L.symbol AND M.bar_date_time = L.bar_date_time
		WHERE L.symbol = ?`, barType.TableName())
	var rec fxBarCsvImportCheckRecord
	if err := r.db.GetContext(ctx, &rec, query, symbol); err != nil {
		return fxmodel.BarCsvImportCheck{}, err
	}
	return fxmodel.BarCsvImportCheck{
		Symbol:      symbol,
		ExistsCount: rec.ExistsCount,
		DiffCount:   rec.DiffCount,
	}, nil
}

func (r *MySQLBarDataRepository) InsertFromLoad(ctx context.Context, symbol string, barType fxmodel.BarType) (int, error) {
	suffix := barType.Suffix()
	query := fmt.Sprintf(`
		INSERT INTO fx_bar_%s
		SELECT
			L.symbol,
			L.bar_date_time,
			L.open_price,
			L.high_price,
			L.low_price,
			L.close_price,
			L.volume,
			(L.high_price  - L.open_price),
			(L.low_price   - L.open_price),
			(L.close_price - L.open_price),
			(L.high_price  - L.low_price)
		FROM fx_bar_load L
		WHERE L.symbol = ?
		  AND NOT EXISTS (
			SELECT 1 FROM fx_bar_%s M
			WHERE M.symbol = L.symbol AND M.bar_date_time = L.bar_date_time
		  )`, suffix, suffix)
	result, err := r.db.ExecContext(ctx, query, symbol)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	return int(rows), err
}

func (r *MySQLBarDataRepository) InsertFromLoadSma(ctx context.Context, symbol string, barType fxmodel.BarType) (int, error) {
	suffix := barType.Suffix()
	query := fmt.Sprintf(`
		INSERT INTO fx_bar_%s_sma (symbol, bar_date_time, sma_range, sma_price, sma_cross)
		SELECT L.symbol, L.bar_date_time, L.sma_range, L.sma_price, L.sma_cross
		FROM fx_bar_load_sma L
		WHERE L.symbol = ?
		  AND NOT EXISTS (
			SELECT 1 FROM fx_bar_%s_sma M
			WHERE M.symbol = L.symbol AND M.bar_date_time = L.bar_date_time AND M.sma_range = L.sma_range
		  )`, suffix, suffix)
	result, err := r.db.ExecContext(ctx, query, symbol)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	return int(rows), err
}

func (r *MySQLBarDataRepository) InsertFromLoadRsi(ctx context.Context, symbol string, barType fxmodel.BarType) (int, error) {
	suffix := barType.Suffix()
	query := fmt.Sprintf(`
		INSERT INTO fx_bar_%s_rsi (symbol, bar_date_time, rsi_range, rsi_value, rsi_ma)
		SELECT L.symbol, L.bar_date_time, L.rsi_range, L.rsi_value, L.rsi_ma
		FROM fx_bar_load_rsi L
		WHERE L.symbol = ?
		  AND NOT EXISTS (
			SELECT 1 FROM fx_bar_%s_rsi M
			WHERE M.symbol = L.symbol AND M.bar_date_time = L.bar_date_time AND M.rsi_range = L.rsi_range
		  )`, suffix, suffix)
	result, err := r.db.ExecContext(ctx, query, symbol)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	return int(rows), err
}

func (r *MySQLBarDataRepository) GetDiffBarData(ctx context.Context, symbol string, barType fxmodel.BarType) ([]fxmodel.BarLoadData, error) {
	query := fmt.Sprintf(`
		SELECT
			L.symbol,
			L.bar_date_time AS barDateTime,
			L.open_price    AS openPrice,
			L.high_price    AS highPrice,
			L.low_price     AS lowPrice,
			L.close_price   AS closePrice,
			L.volume
		FROM fx_bar_load L
		INNER JOIN fx_bar_%s M ON M.symbol = L.symbol AND M.bar_date_time = L.bar_date_time
		WHERE L.symbol = ?
		  AND (M.open_price  != L.open_price
		       OR M.high_price != L.high_price
		       OR M.low_price  != L.low_price
		       OR M.close_price != L.close_price)`, barType.Suffix())
	var recs []fxBarLoadDataRecord
	if err := r.db.SelectContext(ctx, &recs, query, symbol); err != nil {
		return nil, err
	}
	result := make([]fxmodel.BarLoadData, len(recs))
	for i, rec := range recs {
		result[i] = fxmodel.BarLoadData{
			Symbol:      rec.Symbol,
			BarDateTime: rec.BarDateTime,
			OpenPrice:   rec.OpenPrice,
			HighPrice:   rec.HighPrice,
			LowPrice:    rec.LowPrice,
			ClosePrice:  rec.ClosePrice,
			Volume:      rec.Volume,
		}
	}
	return result, nil
}

func (r *MySQLBarDataRepository) GetDiffBarSma(ctx context.Context, symbol string, barType fxmodel.BarType) ([]fxmodel.BarLoadSma, error) {
	query := fmt.Sprintf(`
		SELECT
			L.symbol,
			L.bar_date_time AS barDateTime,
			L.sma_range     AS smaRange,
			L.sma_price     AS smaPrice,
			(L.sma_cross+0) AS smaCross
		FROM fx_bar_load_sma L
		INNER JOIN fx_bar_%s_sma M
			ON M.symbol = L.symbol AND M.bar_date_time = L.bar_date_time AND M.sma_range = L.sma_range
		WHERE L.symbol = ?
		  AND (M.sma_price != L.sma_price OR (M.sma_cross+0) != (L.sma_cross+0))`, barType.Suffix())
	var recs []fxBarLoadSmaRecord
	if err := r.db.SelectContext(ctx, &recs, query, symbol); err != nil {
		return nil, err
	}
	result := make([]fxmodel.BarLoadSma, len(recs))
	for i, rec := range recs {
		result[i] = fxmodel.BarLoadSma{
			Symbol:      rec.Symbol,
			BarDateTime: rec.BarDateTime,
			SmaRange:    rec.SmaRange,
			SmaPrice:    rec.SmaPrice,
			SmaCross:    rec.SmaCross != 0,
		}
	}
	return result, nil
}

func (r *MySQLBarDataRepository) GetDiffBarRsi(ctx context.Context, symbol string, barType fxmodel.BarType) ([]fxmodel.BarLoadRsi, error) {
	query := fmt.Sprintf(`
		SELECT
			L.symbol,
			L.bar_date_time AS barDateTime,
			L.rsi_range     AS rsiRange,
			L.rsi_value     AS rsiValue,
			L.rsi_ma        AS rsiMa
		FROM fx_bar_load_rsi L
		INNER JOIN fx_bar_%s_rsi M
			ON M.symbol = L.symbol AND M.bar_date_time = L.bar_date_time AND M.rsi_range = L.rsi_range
		WHERE L.symbol = ?
		  AND (M.rsi_value != L.rsi_value OR M.rsi_ma != L.rsi_ma)`, barType.Suffix())
	var recs []fxBarLoadRsiRecord
	if err := r.db.SelectContext(ctx, &recs, query, symbol); err != nil {
		return nil, err
	}
	result := make([]fxmodel.BarLoadRsi, len(recs))
	for i, rec := range recs {
		result[i] = fxmodel.BarLoadRsi{
			Symbol:      rec.Symbol,
			BarDateTime: rec.BarDateTime,
			RsiRange:    rec.RsiRange,
			RsiValue:    rec.RsiValue,
			RsiMa:       rec.RsiMa,
		}
	}
	return result, nil
}

func (r *MySQLBarDataRepository) UpdateBarData(ctx context.Context, symbol string, barType fxmodel.BarType) (int, error) {
	query := fmt.Sprintf(`
		UPDATE fx_bar_%s M
		INNER JOIN fx_bar_load L ON L.symbol = M.symbol AND L.bar_date_time = M.bar_date_time
		SET M.open_price  = L.open_price,
		    M.high_price  = L.high_price,
		    M.low_price   = L.low_price,
		    M.close_price = L.close_price,
		    M.volume      = L.volume
		WHERE L.symbol = ?
		  AND (M.open_price  != L.open_price
		       OR M.high_price != L.high_price
		       OR M.low_price  != L.low_price
		       OR M.close_price != L.close_price)`, barType.Suffix())
	result, err := r.db.ExecContext(ctx, query, symbol)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	return int(rows), err
}

func (r *MySQLBarDataRepository) UpdateBarSma(ctx context.Context, symbol string, barType fxmodel.BarType) (int, error) {
	query := fmt.Sprintf(`
		UPDATE fx_bar_%s_sma M
		INNER JOIN fx_bar_load_sma L
			ON L.symbol = M.symbol AND L.bar_date_time = M.bar_date_time AND L.sma_range = M.sma_range
		SET M.sma_price = L.sma_price,
		    M.sma_cross = L.sma_cross
		WHERE L.symbol = ?
		  AND (M.sma_price != L.sma_price OR (M.sma_cross+0) != (L.sma_cross+0))`, barType.Suffix())
	result, err := r.db.ExecContext(ctx, query, symbol)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	return int(rows), err
}

func (r *MySQLBarDataRepository) UpdateBarRsi(ctx context.Context, symbol string, barType fxmodel.BarType) (int, error) {
	query := fmt.Sprintf(`
		UPDATE fx_bar_%s_rsi M
		INNER JOIN fx_bar_load_rsi L
			ON L.symbol = M.symbol AND L.bar_date_time = M.bar_date_time AND L.rsi_range = M.rsi_range
		SET M.rsi_value = L.rsi_value,
		    M.rsi_ma    = L.rsi_ma
		WHERE L.symbol = ?
		  AND (M.rsi_value != L.rsi_value OR M.rsi_ma != L.rsi_ma)`, barType.Suffix())
	result, err := r.db.ExecContext(ctx, query, symbol)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	return int(rows), err
}

// appendDateRange はyyyyMMdd形式の日付範囲条件をクエリに追加する。
func appendDateRange(query string, args []any, barDateFrom, barDateTo string) (string, []any) {
	switch {
	case barDateFrom != "" && barDateTo != "":
		query += ` AND B.bar_date_time BETWEEN STR_TO_DATE(CONCAT(?, '000000'), '%Y%m%d%H%i%s') AND STR_TO_DATE(CONCAT(?, '235959'), '%Y%m%d%H%i%s')`
		args = append(args, barDateFrom, barDateTo)
	case barDateFrom != "":
		query += ` AND B.bar_date_time >= STR_TO_DATE(CONCAT(?, '000000'), '%Y%m%d%H%i%s')`
		args = append(args, barDateFrom)
	case barDateTo != "":
		query += ` AND B.bar_date_time <= STR_TO_DATE(CONCAT(?, '235959'), '%Y%m%d%H%i%s')`
		args = append(args, barDateTo)
	}
	return query, args
}
