package infradbfx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	fxmodel "sandbox-api-gin/internal/domain/model/fx"
	"sandbox-api-gin/internal/domain/model/fx/zigzag"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type MySQLZigZagRepository struct {
	db *sqlx.DB
}

func NewMySQLZigZagRepository(db *sqlx.DB) fxrepository.ZigZagRepository {
	return &MySQLZigZagRepository{db: db}
}

// ------------------------------------------------------------------ //
// DBレコード型                                                          //
// ------------------------------------------------------------------ //

type zigzagStatusRecord struct {
	Symbol               string  `db:"symbol"`
	BarDateTimeMin       string  `db:"barDateTimeMin"`
	BarDateTimeMax       string  `db:"barDateTimeMax"`
	BarCount             int     `db:"barCount"`
	BarDateTimeMinZigZag *string `db:"barDateTimeMinZigZag"`
	BarDateTimeMaxZigZag *string `db:"barDateTimeMaxZigZag"`
	ZigzagCount          int     `db:"zigzagCount"`
	BreakResistanceCount int     `db:"breakResistanceCount"`
	BreakSupportCount    int     `db:"breakSupportCount"`
	Depth                int16   `db:"depth"`
}

type zigzagSearchRecord struct {
	Symbol string `db:"symbol"`
	Depth  int    `db:"depth"`

	CurWaveStart time.Time `db:"curWaveStart"`
	CurWaveEnd   time.Time `db:"curWaveEnd"`
	CurWave      int       `db:"curWave"`
	CurResistance float64  `db:"curResistance"`
	CurSupport    float64  `db:"curSupport"`
	CurSma4h200sS float64  `db:"curSma4h200sS"`
	CurSma4h200sE float64  `db:"curSma4h200sE"`
	CurSma4h75sS  float64  `db:"curSma4h75sS"`
	CurSma4h75sE  float64  `db:"curSma4h75sE"`
	CurSma4h20sS  float64  `db:"curSma4h20sS"`
	CurSma4h20sE  float64  `db:"curSma4h20sE"`
	CurSma1h200sS float64  `db:"curSma1h200sS"`
	CurSma1h200sE float64  `db:"curSma1h200sE"`
	CurSma15m200sS float64 `db:"curSma15m200sS"`
	CurSma15m200sE float64 `db:"curSma15m200sE"`

	T4hWaveStart  *time.Time `db:"t4hWaveStart"`
	T4hWaveEnd    *time.Time `db:"t4hWaveEnd"`
	T4hWave       int        `db:"t4hWave"`
	T4hResistance float64    `db:"t4hResistance"`
	T4hSupport    float64    `db:"t4hSupport"`
	T4hSma4h200sS float64    `db:"t4hSma4h200sS"`
	T4hSma4h200sE float64    `db:"t4hSma4h200sE"`
	T4hSma4h75sS  float64    `db:"t4hSma4h75sS"`
	T4hSma4h75sE  float64    `db:"t4hSma4h75sE"`
	T4hSma4h20sS  float64    `db:"t4hSma4h20sS"`
	T4hSma4h20sE  float64    `db:"t4hSma4h20sE"`
	T4hSma1h200sS float64    `db:"t4hSma1h200sS"`
	T4hSma1h200sE float64    `db:"t4hSma1h200sE"`
	T4hSma15m200sS float64   `db:"t4hSma15m200sS"`
	T4hSma15m200sE float64   `db:"t4hSma15m200sE"`

	PrvWaveStart  *time.Time `db:"prvWaveStart"`
	PrvWaveEnd    *time.Time `db:"prvWaveEnd"`
	PrvWave       int        `db:"prvWave"`
	PrvResistance float64    `db:"prvResistance"`
	PrvSupport    float64    `db:"prvSupport"`

	NxtWaveStart  *time.Time `db:"nxtWaveStart"`
	NxtWaveEnd    *time.Time `db:"nxtWaveEnd"`
	NxtWave       int        `db:"nxtWave"`
	NxtResistance float64    `db:"nxtResistance"`
	NxtSupport    float64    `db:"nxtSupport"`

	Nx2WaveStart  *time.Time `db:"nx2WaveStart"`
	Nx2WaveEnd    *time.Time `db:"nx2WaveEnd"`
	Nx2Wave       int        `db:"nx2Wave"`
	Nx2Resistance float64    `db:"nx2Resistance"`
	Nx2Support    float64    `db:"nx2Support"`

	WaveDxy4h float64 `db:"waveDxy4h"`
	WaveDxy1h float64 `db:"waveDxy1h"`
}

type zigzagRecord struct {
	Symbol string    `db:"symbol"`
	Depth  int       `db:"depth"`
	BarDateTime time.Time `db:"barDateTime"`

	Resistance        float64   `db:"resistance"`
	ResistanceFractal float64   `db:"resistanceFractal"`
	Support           float64   `db:"support"`
	SupportFractal    float64   `db:"supportFractal"`
	PriceHigh         float64   `db:"priceHigh"`
	PriceLow          float64   `db:"priceLow"`
	BackStepHigh      float64   `db:"backStepHigh"`
	BackStepLow       float64   `db:"backStepLow"`
	FractalHigh       float64   `db:"fractalHigh"`
	FractalLow        float64   `db:"fractalLow"`

	ResistanceBarDateTime        time.Time `db:"resistanceBarDateTime"`
	ResistanceFractalBarDateTime time.Time `db:"resistanceFractalBarDateTime"`
	SupportBarDateTime           time.Time `db:"supportBarDateTime"`
	SupportFractalBarDateTime    time.Time `db:"supportFractalBarDateTime"`
	PriceHighBarDateTime         time.Time `db:"priceHighBarDateTime"`
	PriceLowBarDateTime          time.Time `db:"priceLowBarDateTime"`
	BackStepHighBarDateTime      time.Time `db:"backStepHighBarDateTime"`
	BackStepLowBarDateTime       time.Time `db:"backStepLowBarDateTime"`

	Wave            int   `db:"wave"`
	UpTrend         uint8 `db:"upTrend"`
	BreakResistance uint8 `db:"breakResistance"`
	BreakSupport    uint8 `db:"breakSupport"`
	BackStepUp      int   `db:"backStepUp"`
	BackStepDown    int   `db:"backStepDown"`

	BarHighPrice  float64 `db:"barHighPrice"`
	BarLowPrice   float64 `db:"barLowPrice"`
	BarClosePrice float64 `db:"barClosePrice"`
	ExistsZigzag  uint8   `db:"existsZigzag"`
}

type zigzagBarDataRecord struct {
	BarDateTime time.Time `db:"barDateTime"`
	OpenPrice   float64   `db:"openPrice"`
	HighPrice   float64   `db:"highPrice"`
	LowPrice    float64   `db:"lowPrice"`
	ClosePrice  float64   `db:"closePrice"`
	Sma200      float64   `db:"sma200"`
	Sma75       float64   `db:"sma75"`
	Sma20       float64   `db:"sma20"`
}

type zigzagWaveRecord struct {
	WaveStart         time.Time `db:"waveStart"`
	WaveEnd           time.Time `db:"waveEnd"`
	Wave              int       `db:"wave"`
	Resistance        float64   `db:"resistance"`
	Support           float64   `db:"support"`
	PreviousWaveStart time.Time `db:"previousWaveStart"`
	PreviousWave      int       `db:"previousWave"`
	WaveMemo          string    `db:"waveMemo"`
}

// ------------------------------------------------------------------ //
// 検索                                                                 //
// ------------------------------------------------------------------ //

func (r *MySQLZigZagRepository) SearchCount(ctx context.Context, barType fxmodel.BarType, symbol string, depth int,
	barDateTimeMin, barDateTimeMax time.Time,
	wave, previousWave, nextWave, next2Wave, wave4h int) (int, error) {

	suffix := barType.Suffix()
	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM fx_zigzag_wave_%s C
		LEFT JOIN fx_zigzag_wave_%s P
		      ON P.symbol = C.symbol AND P.depth = C.depth
		     AND P.wave_start = C.previous_wave_start
		LEFT JOIN fx_zigzag_wave_%s N
		      ON N.symbol = C.symbol AND N.depth = C.depth
		     AND N.previous_wave_start = C.wave_start
		LEFT JOIN fx_zigzag_wave_%s N2
		      ON N2.symbol = C.symbol AND N2.depth = C.depth
		     AND N2.previous_wave_start = N.wave_start
		LEFT JOIN fx_zigzag_wave_4h T4H
		      ON T4H.symbol = C.symbol AND T4H.depth = C.depth
		     AND T4H.wave_start = (
		         SELECT wave_start FROM fx_zigzag_wave_4h
		         WHERE  symbol = C.symbol AND depth = C.depth
		            AND wave_end <= C.wave_start
		            AND MOD(wave, 2) <> 0
		         ORDER BY wave_start DESC LIMIT 1
		     )
		WHERE C.symbol = ?
		  AND C.depth  = ?
		  AND C.wave_start BETWEEN ? AND ?`,
		suffix, suffix, suffix, suffix)

	args := []any{symbol, depth, barDateTimeMin, barDateTimeMax}
	query, args = appendZigzagWhere(query, args, wave, previousWave, nextWave, next2Wave, wave4h)

	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	return count, err
}

func (r *MySQLZigZagRepository) Search(ctx context.Context, barType fxmodel.BarType, symbol string, depth int,
	barDateTimeMin, barDateTimeMax time.Time,
	wave, previousWave, nextWave, next2Wave, wave4h int,
	page, size int) ([]*zigzag.ZigZagSearchRow, error) {

	suffix := barType.Suffix()
	query := fmt.Sprintf(`
		SELECT
		     C.symbol                              AS symbol
		    ,C.depth                               AS depth
		    ,C.wave_start                          AS curWaveStart
		    ,C.wave_end                            AS curWaveEnd
		    ,C.wave                                AS curWave
		    ,C.resistance                          AS curResistance
		    ,C.support                             AS curSupport
		    ,COALESCE(cur4h200s.sma_price,  0)    AS curSma4h200sS
		    ,COALESCE(cur4h200e.sma_price,  0)    AS curSma4h200sE
		    ,COALESCE(cur4h75s.sma_price,   0)    AS curSma4h75sS
		    ,COALESCE(cur4h75e.sma_price,   0)    AS curSma4h75sE
		    ,COALESCE(cur4h20s.sma_price,   0)    AS curSma4h20sS
		    ,COALESCE(cur4h20e.sma_price,   0)    AS curSma4h20sE
		    ,COALESCE(cur1h200s.sma_price,  0)    AS curSma1h200sS
		    ,COALESCE(cur1h200e.sma_price,  0)    AS curSma1h200sE
		    ,COALESCE(cur15m200s.sma_price, 0)    AS curSma15m200sS
		    ,COALESCE(cur15m200e.sma_price, 0)    AS curSma15m200sE
		    ,T4H.wave_start                        AS t4hWaveStart
		    ,T4H.wave_end                          AS t4hWaveEnd
		    ,COALESCE(T4H.wave,        0)          AS t4hWave
		    ,COALESCE(T4H.resistance,  0)          AS t4hResistance
		    ,COALESCE(T4H.support,     0)          AS t4hSupport
		    ,COALESCE(t4h4h200s.sma_price,  0)    AS t4hSma4h200sS
		    ,COALESCE(t4h4h200e.sma_price,  0)    AS t4hSma4h200sE
		    ,COALESCE(t4h4h75s.sma_price,   0)    AS t4hSma4h75sS
		    ,COALESCE(t4h4h75e.sma_price,   0)    AS t4hSma4h75sE
		    ,COALESCE(t4h4h20s.sma_price,   0)    AS t4hSma4h20sS
		    ,COALESCE(t4h4h20e.sma_price,   0)    AS t4hSma4h20sE
		    ,COALESCE(t4h1h200s.sma_price,  0)    AS t4hSma1h200sS
		    ,COALESCE(t4h1h200e.sma_price,  0)    AS t4hSma1h200sE
		    ,COALESCE(t4h15m200s.sma_price, 0)    AS t4hSma15m200sS
		    ,COALESCE(t4h15m200e.sma_price, 0)    AS t4hSma15m200sE
		    ,P.wave_start                          AS prvWaveStart
		    ,P.wave_end                            AS prvWaveEnd
		    ,COALESCE(P.wave, 0)                   AS prvWave
		    ,COALESCE(P.resistance, 0)             AS prvResistance
		    ,COALESCE(P.support, 0)                AS prvSupport
		    ,N.wave_start                          AS nxtWaveStart
		    ,N.wave_end                            AS nxtWaveEnd
		    ,COALESCE(N.wave, 0)                   AS nxtWave
		    ,COALESCE(N.resistance, 0)             AS nxtResistance
		    ,COALESCE(N.support, 0)                AS nxtSupport
		    ,N2.wave_start                         AS nx2WaveStart
		    ,N2.wave_end                           AS nx2WaveEnd
		    ,COALESCE(N2.wave, 0)                  AS nx2Wave
		    ,COALESCE(N2.resistance, 0)            AS nx2Resistance
		    ,COALESCE(N2.support, 0)               AS nx2Support
		    ,COALESCE(dxy4h.wave, 0)               AS waveDxy4h
		    ,COALESCE(dxy1h.wave, 0)               AS waveDxy1h
		FROM fx_zigzag_wave_%s C
		LEFT JOIN fx_zigzag_wave_%s P
		      ON P.symbol = C.symbol AND P.depth = C.depth
		     AND P.wave_start = C.previous_wave_start
		LEFT JOIN fx_zigzag_wave_%s N
		      ON N.symbol = C.symbol AND N.depth = C.depth
		     AND N.previous_wave_start = C.wave_start
		LEFT JOIN fx_zigzag_wave_%s N2
		      ON N2.symbol = C.symbol AND N2.depth = C.depth
		     AND N2.previous_wave_start = N.wave_start
		LEFT JOIN fx_zigzag_wave_4h T4H
		      ON T4H.symbol = C.symbol AND T4H.depth = C.depth
		     AND T4H.wave_start = (
		         SELECT wave_start FROM fx_zigzag_wave_4h
		         WHERE  symbol = C.symbol AND depth = C.depth
		            AND wave_end <= C.wave_start
		            AND MOD(wave, 2) <> 0
		         ORDER BY wave_start DESC LIMIT 1
		     )
		LEFT JOIN fx_bar_4h_sma  cur4h200s  ON cur4h200s.symbol  = C.symbol AND cur4h200s.bar_date_time  = C.wave_start AND cur4h200s.sma_range  = 200
		LEFT JOIN fx_bar_4h_sma  cur4h200e  ON cur4h200e.symbol  = C.symbol AND cur4h200e.bar_date_time  = C.wave_end   AND cur4h200e.sma_range  = 200
		LEFT JOIN fx_bar_4h_sma  cur4h75s   ON cur4h75s.symbol   = C.symbol AND cur4h75s.bar_date_time   = C.wave_start AND cur4h75s.sma_range   = 75
		LEFT JOIN fx_bar_4h_sma  cur4h75e   ON cur4h75e.symbol   = C.symbol AND cur4h75e.bar_date_time   = C.wave_end   AND cur4h75e.sma_range   = 75
		LEFT JOIN fx_bar_4h_sma  cur4h20s   ON cur4h20s.symbol   = C.symbol AND cur4h20s.bar_date_time   = C.wave_start AND cur4h20s.sma_range   = 20
		LEFT JOIN fx_bar_4h_sma  cur4h20e   ON cur4h20e.symbol   = C.symbol AND cur4h20e.bar_date_time   = C.wave_end   AND cur4h20e.sma_range   = 20
		LEFT JOIN fx_bar_1h_sma  cur1h200s  ON cur1h200s.symbol  = C.symbol AND cur1h200s.bar_date_time  = C.wave_start AND cur1h200s.sma_range  = 200
		LEFT JOIN fx_bar_1h_sma  cur1h200e  ON cur1h200e.symbol  = C.symbol AND cur1h200e.bar_date_time  = C.wave_end   AND cur1h200e.sma_range  = 200
		LEFT JOIN fx_bar_15m_sma cur15m200s ON cur15m200s.symbol = C.symbol AND cur15m200s.bar_date_time = C.wave_start AND cur15m200s.sma_range = 200
		LEFT JOIN fx_bar_15m_sma cur15m200e ON cur15m200e.symbol = C.symbol AND cur15m200e.bar_date_time = C.wave_end   AND cur15m200e.sma_range = 200
		LEFT JOIN fx_bar_4h_sma  t4h4h200s  ON t4h4h200s.symbol  = T4H.symbol AND t4h4h200s.bar_date_time  = T4H.wave_start AND t4h4h200s.sma_range  = 200
		LEFT JOIN fx_bar_4h_sma  t4h4h200e  ON t4h4h200e.symbol  = T4H.symbol AND t4h4h200e.bar_date_time  = T4H.wave_end   AND t4h4h200e.sma_range  = 200
		LEFT JOIN fx_bar_4h_sma  t4h4h75s   ON t4h4h75s.symbol   = T4H.symbol AND t4h4h75s.bar_date_time   = T4H.wave_start AND t4h4h75s.sma_range   = 75
		LEFT JOIN fx_bar_4h_sma  t4h4h75e   ON t4h4h75e.symbol   = T4H.symbol AND t4h4h75e.bar_date_time   = T4H.wave_end   AND t4h4h75e.sma_range   = 75
		LEFT JOIN fx_bar_4h_sma  t4h4h20s   ON t4h4h20s.symbol   = T4H.symbol AND t4h4h20s.bar_date_time   = T4H.wave_start AND t4h4h20s.sma_range   = 20
		LEFT JOIN fx_bar_4h_sma  t4h4h20e   ON t4h4h20e.symbol   = T4H.symbol AND t4h4h20e.bar_date_time   = T4H.wave_end   AND t4h4h20e.sma_range   = 20
		LEFT JOIN fx_bar_1h_sma  t4h1h200s  ON t4h1h200s.symbol  = T4H.symbol AND t4h1h200s.bar_date_time  = T4H.wave_start AND t4h1h200s.sma_range  = 200
		LEFT JOIN fx_bar_1h_sma  t4h1h200e  ON t4h1h200e.symbol  = T4H.symbol AND t4h1h200e.bar_date_time  = T4H.wave_end   AND t4h1h200e.sma_range  = 200
		LEFT JOIN fx_bar_15m_sma t4h15m200s ON t4h15m200s.symbol = T4H.symbol AND t4h15m200s.bar_date_time = T4H.wave_start AND t4h15m200s.sma_range = 200
		LEFT JOIN fx_bar_15m_sma t4h15m200e ON t4h15m200e.symbol = T4H.symbol AND t4h15m200e.bar_date_time = T4H.wave_end   AND t4h15m200e.sma_range = 200
		LEFT JOIN fx_zigzag_wave_4h dxy4h
		      ON dxy4h.symbol = 'DXY'
		     AND dxy4h.depth  = 12
		     AND dxy4h.wave_start = (
		         SELECT wave_start FROM fx_zigzag_wave_4h
		         WHERE  symbol = 'DXY' AND depth = 12
		            AND wave_start <= C.wave_start
		         ORDER BY wave_start DESC LIMIT 1
		     )
		LEFT JOIN fx_zigzag_wave_1h dxy1h
		      ON dxy1h.symbol = 'DXY'
		     AND dxy1h.depth  = 12
		     AND dxy1h.wave_start = (
		         SELECT wave_start FROM fx_zigzag_wave_1h
		         WHERE  symbol = 'DXY' AND depth = 12
		            AND wave_start <= C.wave_start
		         ORDER BY wave_start DESC LIMIT 1
		     )
		WHERE C.symbol = ?
		  AND C.depth  = ?
		  AND C.wave_start BETWEEN ? AND ?`,
		suffix, suffix, suffix, suffix)

	args := []any{symbol, depth, barDateTimeMin, barDateTimeMax}
	query, args = appendZigzagWhere(query, args, wave, previousWave, nextWave, next2Wave, wave4h)
	query += ` ORDER BY C.wave_start DESC`

	offset := (page - 1) * size
	if offset > 0 {
		query += ` LIMIT ? OFFSET ?`
		args = append(args, size, offset)
	} else {
		query += ` LIMIT ?`
		args = append(args, size)
	}

	var recs []zigzagSearchRecord
	if err := r.db.SelectContext(ctx, &recs, query, args...); err != nil {
		return nil, err
	}

	result := make([]*zigzag.ZigZagSearchRow, len(recs))
	for i, rec := range recs {
		result[i] = toSearchDomain(&rec)
	}
	return result, nil
}

// ------------------------------------------------------------------ //
// ステータス                                                            //
// ------------------------------------------------------------------ //

func (r *MySQLZigZagRepository) GetStatusList(ctx context.Context, symbolType string, barType fxmodel.BarType, depth int) ([]*zigzag.ZigZagStatus, error) {
	suffix := barType.Suffix()
	query := fmt.Sprintf(`
		SELECT
		     bar.symbol
		    ,bar.barCount
		    ,bar.barDateTimeMin
		    ,bar.barDateTimeMax
		    ,COALESCE(zig.zigzagCount, 0)           AS zigzagCount
		    ,zig.barDateTimeMinZigZag
		    ,zig.barDateTimeMaxZigZag
		    ,COALESCE(zig.breakResistanceCount, 0)  AS breakResistanceCount
		    ,COALESCE(zig.breakSupportCount, 0)      AS breakSupportCount
		    ,? AS depth
		FROM (
		    SELECT
		         symbol AS symbol
		        ,COUNT(*) AS barCount
		        ,DATE_FORMAT(MIN(bar_date_time), '%%Y-%%m-%%d %%H:%%i') AS barDateTimeMin
		        ,DATE_FORMAT(MAX(bar_date_time), '%%Y-%%m-%%d %%H:%%i') AS barDateTimeMax
		    FROM fx_bar_%s
		    GROUP BY symbol
		) bar
		LEFT JOIN (
		    SELECT
		         symbol           AS symbol
		        ,COUNT(*)         AS zigzagCount
		        ,DATE_FORMAT(MIN(bar_date_time), '%%Y-%%m-%%d %%H:%%i') AS barDateTimeMinZigZag
		        ,DATE_FORMAT(MAX(bar_date_time), '%%Y-%%m-%%d %%H:%%i') AS barDateTimeMaxZigZag
		        ,SUM(IF((break_resistance+0) = 1, 1, 0))               AS breakResistanceCount
		        ,SUM(IF((break_support+0)    = 1, 1, 0))               AS breakSupportCount
		    FROM fx_zigzag_%s
		    WHERE depth = ?
		    GROUP BY symbol
		) zig ON zig.symbol = bar.symbol
		INNER JOIN fx_symbol c ON c.symbol = bar.symbol
		WHERE c.symbol_type = ?
		ORDER BY c.sort_order`, suffix, suffix)

	var recs []zigzagStatusRecord
	if err := r.db.SelectContext(ctx, &recs, query, depth, depth, symbolType); err != nil {
		return nil, err
	}

	result := make([]*zigzag.ZigZagStatus, len(recs))
	for i, rec := range recs {
		result[i] = toStatusDomain(&rec, symbolType, barType)
	}
	return result, nil
}

func (r *MySQLZigZagRepository) GetStatus(ctx context.Context, barType fxmodel.BarType, symbol string, depth int) (*zigzag.ZigZagStatus, error) {
	suffix := barType.Suffix()
	query := fmt.Sprintf(`
		SELECT
		     B.barDateTimeMin   AS barDateTimeMin
		    ,B.barDateTimeMax   AS barDateTimeMax
		    ,B.barCount         AS barCount
		    ,Z.barDateTimeMinZigZag   AS barDateTimeMinZigZag
		    ,Z.barDateTimeMaxZigZag   AS barDateTimeMaxZigZag
		    ,COALESCE(Z.zigzagCount, 0)          AS zigzagCount
		    ,COALESCE(Z.breakResistanceCount, 0) AS breakResistanceCount
		    ,COALESCE(Z.breakSupportCount, 0)    AS breakSupportCount
		    ,? AS depth
		    ,? AS symbol
		FROM (
		    SELECT
		         DATE_FORMAT(MIN(bar_date_time), '%%Y-%%m-%%d %%H:%%i') AS barDateTimeMin
		        ,DATE_FORMAT(MAX(bar_date_time), '%%Y-%%m-%%d %%H:%%i') AS barDateTimeMax
		        ,COUNT(*)                                                AS barCount
		    FROM fx_bar_%s
		    WHERE symbol = ?
		) B
		LEFT JOIN (
		    SELECT
		         DATE_FORMAT(MIN(bar_date_time), '%%Y-%%m-%%d %%H:%%i') AS barDateTimeMinZigZag
		        ,DATE_FORMAT(MAX(bar_date_time), '%%Y-%%m-%%d %%H:%%i') AS barDateTimeMaxZigZag
		        ,COUNT(*)                                                AS zigzagCount
		        ,SUM(IF((break_resistance+0) = 1, 1, 0))               AS breakResistanceCount
		        ,SUM(IF((break_support+0)    = 1, 1, 0))               AS breakSupportCount
		    FROM fx_zigzag_%s
		    WHERE symbol = ? AND depth = ?
		) Z ON TRUE`, suffix, suffix)

	var rec zigzagStatusRecord
	if err := r.db.GetContext(ctx, &rec, query, depth, symbol, symbol, symbol, depth); err != nil {
		return nil, err
	}
	rec.Symbol = symbol
	return toStatusDomain(&rec, "", barType), nil
}

// ------------------------------------------------------------------ //
// バーデータ                                                            //
// ------------------------------------------------------------------ //

func (r *MySQLZigZagRepository) GetBarDataList(ctx context.Context, barType fxmodel.BarType, symbol string, depth int, waveStart time.Time) ([]*zigzag.ZigZagBarDataRow, error) {
	suffix := barType.Suffix()

	var rangeSel string
	if suffix == "4h" {
		rangeSel = `COALESCE(prev2.wave_start, prev.wave_start, cur.wave_start) AS barDateTimeMin
		           ,COALESCE(next3.wave_end, next2.wave_end, next.wave_end, cur.wave_end) AS barDateTimeMax`
	} else {
		rangeSel = `COALESCE(cur.wave_start, cur.wave_start) AS barDateTimeMin
		           ,COALESCE(next.wave_end, cur.wave_end) AS barDateTimeMax`
	}

	var extraJoins string
	if suffix == "4h" {
		extraJoins = `LEFT JOIN fx_zigzag_wave_4h prev2
		      ON prev2.symbol = cur.symbol AND prev2.depth = cur.depth
		     AND prev2.wave_start = COALESCE(prev.previous_wave_start, cur.previous_wave_start)
		LEFT JOIN fx_zigzag_wave_4h next2
		      ON next2.symbol = cur.symbol AND next2.depth = cur.depth
		     AND next2.wave_start = COALESCE(next.wave_end, cur.wave_end)
		LEFT JOIN fx_zigzag_wave_4h next3
		      ON next3.symbol = cur.symbol AND next3.depth = cur.depth
		     AND next3.wave_start = COALESCE(next2.wave_end, next.wave_end, cur.wave_end)`
	}

	query := fmt.Sprintf(`
		SELECT
		     B.bar_date_time               AS barDateTime
		    ,B.open_price                  AS openPrice
		    ,B.high_price                  AS highPrice
		    ,B.low_price                   AS lowPrice
		    ,B.close_price                 AS closePrice
		    ,COALESCE(sma200.sma_price, 0) AS sma200
		    ,COALESCE(sma75.sma_price,  0) AS sma75
		    ,COALESCE(sma20.sma_price,  0) AS sma20
		FROM (
		    SELECT %s
		    FROM (
		        SELECT * FROM fx_zigzag_wave_4h
		        WHERE  symbol = ? AND depth = ?
		           AND wave_start <= ?
		        ORDER BY wave_start DESC LIMIT 1
		    ) cur
		    LEFT JOIN fx_zigzag_wave_4h prev
		          ON prev.symbol = cur.symbol AND prev.depth = cur.depth
		         AND prev.wave_start = cur.previous_wave_start
		    %s
		    LEFT JOIN fx_zigzag_wave_4h next
		          ON next.symbol = cur.symbol AND next.depth = cur.depth
		         AND next.wave_start = cur.wave_end
		) range_
		INNER JOIN fx_bar_%s B
		        ON B.symbol = ?
		       AND B.bar_date_time BETWEEN range_.barDateTimeMin AND range_.barDateTimeMax
		LEFT JOIN fx_bar_%s_sma sma200
		       ON sma200.symbol = B.symbol AND sma200.bar_date_time = B.bar_date_time AND sma200.sma_range = 200
		LEFT JOIN fx_bar_%s_sma sma75
		       ON sma75.symbol  = B.symbol AND sma75.bar_date_time  = B.bar_date_time AND sma75.sma_range  = 75
		LEFT JOIN fx_bar_%s_sma sma20
		       ON sma20.symbol  = B.symbol AND sma20.bar_date_time  = B.bar_date_time AND sma20.sma_range  = 20
		ORDER BY B.bar_date_time`,
		rangeSel, extraJoins, suffix, suffix, suffix, suffix)

	var recs []zigzagBarDataRecord
	if err := r.db.SelectContext(ctx, &recs, query, symbol, depth, waveStart, symbol); err != nil {
		return nil, err
	}

	result := make([]*zigzag.ZigZagBarDataRow, len(recs))
	for i, rec := range recs {
		result[i] = &zigzag.ZigZagBarDataRow{
			BarDateTime: rec.BarDateTime,
			OpenPrice:   rec.OpenPrice,
			HighPrice:   rec.HighPrice,
			LowPrice:    rec.LowPrice,
			ClosePrice:  rec.ClosePrice,
			Sma200:      rec.Sma200,
			Sma75:       rec.Sma75,
			Sma20:       rec.Sma20,
		}
	}
	return result, nil
}

// ------------------------------------------------------------------ //
// 生成                                                                 //
// ------------------------------------------------------------------ //

func (r *MySQLZigZagRepository) TargetBarCount(ctx context.Context, barType fxmodel.BarType, symbol string, barDateTime time.Time) (int, error) {
	query := fmt.Sprintf(`SELECT COUNT(*) FROM fx_bar_%s WHERE symbol = ? AND bar_date_time >= ?`, barType.Suffix())
	var count int
	err := r.db.GetContext(ctx, &count, query, symbol, barDateTime)
	return count, err
}

func (r *MySQLZigZagRepository) PreviousList(ctx context.Context, barType fxmodel.BarType, symbol string, depth int, barDateTime time.Time, limit int) ([]*zigzag.ZigZag, error) {
	suffix := barType.Suffix()
	query := fmt.Sprintf(`
		SELECT
		     Z.symbol                              AS symbol
		    ,Z.depth                               AS depth
		    ,Z.bar_date_time                       AS barDateTime
		    ,Z.resistance                          AS resistance
		    ,Z.resistance_fractal                  AS resistanceFractal
		    ,Z.support                             AS support
		    ,Z.support_fractal                     AS supportFractal
		    ,Z.high                                AS priceHigh
		    ,Z.low                                 AS priceLow
		    ,Z.backstep_high                       AS backStepHigh
		    ,Z.backstep_low                        AS backStepLow
		    ,Z.fractal_high                        AS fractalHigh
		    ,Z.fractal_low                         AS fractalLow
		    ,Z.resistance_bar_date_time            AS resistanceBarDateTime
		    ,Z.resistance_fractal_bar_date_time    AS resistanceFractalBarDateTime
		    ,Z.support_bar_date_time               AS supportBarDateTime
		    ,Z.support_fractal_bar_date_time       AS supportFractalBarDateTime
		    ,Z.high_bar_date_time                  AS priceHighBarDateTime
		    ,Z.low_bar_date_time                   AS priceLowBarDateTime
		    ,Z.backstep_high_bar_date_time         AS backStepHighBarDateTime
		    ,Z.backstep_low_bar_date_time          AS backStepLowBarDateTime
		    ,Z.wave                                AS wave
		    ,(Z.up_trend+0)                        AS upTrend
		    ,(Z.break_resistance+0)                AS breakResistance
		    ,(Z.break_support+0)                   AS breakSupport
		    ,Z.backstep_up                         AS backStepUp
		    ,Z.backstep_down                       AS backStepDown
		    ,B.high_price                          AS barHighPrice
		    ,B.low_price                           AS barLowPrice
		    ,B.close_price                         AS barClosePrice
		    ,1                                     AS existsZigzag
		FROM fx_zigzag_%s Z
		INNER JOIN fx_bar_%s B ON B.symbol = Z.symbol AND B.bar_date_time = Z.bar_date_time
		WHERE Z.symbol = ? AND Z.depth = ? AND Z.bar_date_time < ?
		ORDER BY Z.bar_date_time DESC
		LIMIT ?`, suffix, suffix)

	var recs []zigzagRecord
	if err := r.db.SelectContext(ctx, &recs, query, symbol, depth, barDateTime, limit); err != nil {
		return nil, err
	}

	result := make([]*zigzag.ZigZag, len(recs))
	for i, rec := range recs {
		result[i] = toZigZagDomain(&rec)
	}
	return result, nil
}

func (r *MySQLZigZagRepository) TargetList(ctx context.Context, barType fxmodel.BarType, symbol string, depth int, barDateTime time.Time, limit int) ([]*zigzag.ZigZag, error) {
	suffix := barType.Suffix()
	query := fmt.Sprintf(`
		SELECT
		     B.symbol                              AS symbol
		    ,?                                     AS depth
		    ,B.bar_date_time                       AS barDateTime
		    ,COALESCE(Z.resistance,          0)    AS resistance
		    ,COALESCE(Z.resistance_fractal,  0)    AS resistanceFractal
		    ,COALESCE(Z.support,             0)    AS support
		    ,COALESCE(Z.support_fractal,     0)    AS supportFractal
		    ,COALESCE(Z.high,                0)    AS priceHigh
		    ,COALESCE(Z.low,                 0)    AS priceLow
		    ,COALESCE(Z.backstep_high,       0)    AS backStepHigh
		    ,COALESCE(Z.backstep_low,        0)    AS backStepLow
		    ,COALESCE(Z.fractal_high,        0)    AS fractalHigh
		    ,COALESCE(Z.fractal_low,         0)    AS fractalLow
		    ,COALESCE(Z.resistance_bar_date_time,         B.bar_date_time) AS resistanceBarDateTime
		    ,COALESCE(Z.resistance_fractal_bar_date_time, B.bar_date_time) AS resistanceFractalBarDateTime
		    ,COALESCE(Z.support_bar_date_time,            B.bar_date_time) AS supportBarDateTime
		    ,COALESCE(Z.support_fractal_bar_date_time,    B.bar_date_time) AS supportFractalBarDateTime
		    ,COALESCE(Z.high_bar_date_time,               B.bar_date_time) AS priceHighBarDateTime
		    ,COALESCE(Z.low_bar_date_time,                B.bar_date_time) AS priceLowBarDateTime
		    ,COALESCE(Z.backstep_high_bar_date_time,      B.bar_date_time) AS backStepHighBarDateTime
		    ,COALESCE(Z.backstep_low_bar_date_time,       B.bar_date_time) AS backStepLowBarDateTime
		    ,COALESCE(Z.wave,            0)        AS wave
		    ,COALESCE((Z.up_trend+0),        0)    AS upTrend
		    ,COALESCE((Z.break_resistance+0),0)    AS breakResistance
		    ,COALESCE((Z.break_support+0),   0)    AS breakSupport
		    ,COALESCE(Z.backstep_up,     0)        AS backStepUp
		    ,COALESCE(Z.backstep_down,   0)        AS backStepDown
		    ,B.high_price                          AS barHighPrice
		    ,B.low_price                           AS barLowPrice
		    ,B.close_price                         AS barClosePrice
		    ,IF(Z.symbol IS NOT NULL, 1, 0)        AS existsZigzag
		FROM fx_bar_%s B
		LEFT JOIN fx_zigzag_%s Z
		       ON Z.symbol = B.symbol AND Z.bar_date_time = B.bar_date_time AND Z.depth = ?
		WHERE B.symbol = ? AND B.bar_date_time >= ?
		ORDER BY B.bar_date_time
		LIMIT ?`, suffix, suffix)

	var recs []zigzagRecord
	if err := r.db.SelectContext(ctx, &recs, query, depth, depth, symbol, barDateTime, limit); err != nil {
		return nil, err
	}

	result := make([]*zigzag.ZigZag, len(recs))
	for i, rec := range recs {
		result[i] = toZigZagDomain(&rec)
	}
	return result, nil
}

func (r *MySQLZigZagRepository) Insert(ctx context.Context, barType fxmodel.BarType, z *zigzag.ZigZag) (int64, error) {
	suffix := barType.Suffix()
	query := fmt.Sprintf(`
		INSERT INTO fx_zigzag_%s (
		    symbol, depth, bar_date_time,
		    resistance, resistance_fractal, support, support_fractal,
		    high, low, backstep_high, backstep_low,
		    fractal_high, fractal_low,
		    resistance_bar_date_time, resistance_fractal_bar_date_time,
		    support_bar_date_time, support_fractal_bar_date_time,
		    high_bar_date_time, low_bar_date_time,
		    backstep_high_bar_date_time, backstep_low_bar_date_time,
		    backstep_up, backstep_down, wave, up_trend, break_resistance, break_support
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		suffix)

	result, err := r.db.ExecContext(ctx, query,
		z.Symbol, z.Depth, z.BarDateTime,
		z.Resistance, z.ResistanceFractal, z.Support, z.SupportFractal,
		z.PriceHigh, z.PriceLow, z.BackStepHigh, z.BackStepLow,
		z.FractalHigh, z.FractalLow,
		z.ResistanceBarDateTime, z.ResistanceFractalBarDateTime,
		z.SupportBarDateTime, z.SupportFractalBarDateTime,
		z.PriceHighBarDateTime, z.PriceLowBarDateTime,
		z.BackStepHighBarDateTime, z.BackStepLowBarDateTime,
		z.BackStepUp, z.BackStepDown, z.Wave,
		boolToUint8(z.UpTrend), boolToUint8(z.BreakResistance), boolToUint8(z.BreakSupport),
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *MySQLZigZagRepository) Update(ctx context.Context, barType fxmodel.BarType, z *zigzag.ZigZag) (int64, error) {
	suffix := barType.Suffix()
	query := fmt.Sprintf(`
		UPDATE fx_zigzag_%s SET
		    resistance                       = ?
		   ,resistance_fractal               = ?
		   ,support                          = ?
		   ,support_fractal                  = ?
		   ,high                             = ?
		   ,low                              = ?
		   ,backstep_high                    = ?
		   ,backstep_low                     = ?
		   ,fractal_high                     = ?
		   ,fractal_low                      = ?
		   ,resistance_bar_date_time         = ?
		   ,resistance_fractal_bar_date_time = ?
		   ,support_bar_date_time            = ?
		   ,support_fractal_bar_date_time    = ?
		   ,high_bar_date_time               = ?
		   ,low_bar_date_time                = ?
		   ,backstep_high_bar_date_time      = ?
		   ,backstep_low_bar_date_time       = ?
		   ,backstep_up                      = ?
		   ,backstep_down                    = ?
		   ,wave                             = ?
		   ,up_trend                         = ?
		   ,break_resistance                 = ?
		   ,break_support                    = ?
		WHERE symbol = ? AND depth = ? AND bar_date_time = ?`, suffix)

	result, err := r.db.ExecContext(ctx, query,
		z.Resistance, z.ResistanceFractal, z.Support, z.SupportFractal,
		z.PriceHigh, z.PriceLow, z.BackStepHigh, z.BackStepLow,
		z.FractalHigh, z.FractalLow,
		z.ResistanceBarDateTime, z.ResistanceFractalBarDateTime,
		z.SupportBarDateTime, z.SupportFractalBarDateTime,
		z.PriceHighBarDateTime, z.PriceLowBarDateTime,
		z.BackStepHighBarDateTime, z.BackStepLowBarDateTime,
		z.BackStepUp, z.BackStepDown, z.Wave,
		boolToUint8(z.UpTrend), boolToUint8(z.BreakResistance), boolToUint8(z.BreakSupport),
		z.Symbol, z.Depth, z.BarDateTime,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// ------------------------------------------------------------------ //
// Wave                                                                 //
// ------------------------------------------------------------------ //

func (r *MySQLZigZagRepository) DeleteWave(ctx context.Context, barType fxmodel.BarType, symbol string, depth int, barDateTime time.Time) error {
	query := fmt.Sprintf(`DELETE FROM fx_zigzag_wave_%s WHERE symbol = ? AND depth = ? AND wave_start >= ?`, barType.Suffix())
	_, err := r.db.ExecContext(ctx, query, symbol, depth, barDateTime)
	return err
}

func (r *MySQLZigZagRepository) GetLastWave(ctx context.Context, barType fxmodel.BarType, symbol string, depth int) (*zigzag.ZigZagWave, error) {
	query := fmt.Sprintf(`
		SELECT
		     wave_start          AS waveStart
		    ,wave_end            AS waveEnd
		    ,wave                AS wave
		    ,resistance          AS resistance
		    ,support             AS support
		    ,previous_wave_start AS previousWaveStart
		    ,previous_wave       AS previousWave
		    ,wave_memo           AS waveMemo
		FROM fx_zigzag_wave_%s
		WHERE symbol = ? AND depth = ?
		ORDER BY wave_start DESC
		LIMIT 1`, barType.Suffix())

	var rec zigzagWaveRecord
	err := r.db.GetContext(ctx, &rec, query, symbol, depth)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &zigzag.ZigZagWave{
		WaveStart:         rec.WaveStart,
		WaveEnd:           rec.WaveEnd,
		Wave:              rec.Wave,
		Resistance:        rec.Resistance,
		Support:           rec.Support,
		PreviousWaveStart: rec.PreviousWaveStart,
		PreviousWave:      rec.PreviousWave,
		WaveMemo:          rec.WaveMemo,
	}, nil
}

func (r *MySQLZigZagRepository) InsertWaveBulk(ctx context.Context, barType fxmodel.BarType, symbol string, depth int, waveList []*zigzag.ZigZagWave) error {
	if len(waveList) == 0 {
		return nil
	}
	suffix := barType.Suffix()

	query := fmt.Sprintf(`
		INSERT INTO fx_zigzag_wave_%s (
		    symbol, depth, wave_start, wave_end, wave,
		    resistance, support, previous_wave_start, previous_wave, wave_memo
		) VALUES `, suffix)

	args := make([]any, 0, len(waveList)*10)
	for i, w := range waveList {
		if i > 0 {
			query += ","
		}
		query += "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
		args = append(args, symbol, depth, w.WaveStart, w.WaveEnd, w.Wave,
			w.Resistance, w.Support, w.PreviousWaveStart, w.PreviousWave, w.WaveMemo)
	}
	query += ` ON DUPLICATE KEY UPDATE
		wave_end            = VALUES(wave_end),
		wave                = VALUES(wave),
		resistance          = VALUES(resistance),
		support             = VALUES(support),
		previous_wave_start = VALUES(previous_wave_start),
		previous_wave       = VALUES(previous_wave),
		wave_memo           = VALUES(wave_memo)`

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

// ------------------------------------------------------------------ //
// 変換ヘルパー                                                          //
// ------------------------------------------------------------------ //

func toZigZagDomain(rec *zigzagRecord) *zigzag.ZigZag {
	return &zigzag.ZigZag{
		Symbol:      rec.Symbol,
		Depth:       rec.Depth,
		BarDateTime: rec.BarDateTime,

		Resistance:        rec.Resistance,
		ResistanceFractal: rec.ResistanceFractal,
		Support:           rec.Support,
		SupportFractal:    rec.SupportFractal,
		PriceHigh:         rec.PriceHigh,
		PriceLow:          rec.PriceLow,
		BackStepHigh:      rec.BackStepHigh,
		BackStepLow:       rec.BackStepLow,
		FractalHigh:       rec.FractalHigh,
		FractalLow:        rec.FractalLow,

		ResistanceBarDateTime:        rec.ResistanceBarDateTime,
		ResistanceFractalBarDateTime: rec.ResistanceFractalBarDateTime,
		SupportBarDateTime:           rec.SupportBarDateTime,
		SupportFractalBarDateTime:    rec.SupportFractalBarDateTime,
		PriceHighBarDateTime:         rec.PriceHighBarDateTime,
		PriceLowBarDateTime:          rec.PriceLowBarDateTime,
		BackStepHighBarDateTime:      rec.BackStepHighBarDateTime,
		BackStepLowBarDateTime:       rec.BackStepLowBarDateTime,

		Wave:             rec.Wave,
		UpTrend:          rec.UpTrend != 0,
		BreakResistance:  rec.BreakResistance != 0,
		BreakSupport:     rec.BreakSupport != 0,
		WaveFractal:      rec.Wave,
		BreakResistanceFractal: rec.BreakResistance != 0,
		BreakSupportFractal:    rec.BreakSupport != 0,
		BackStepUp:       rec.BackStepUp,
		BackStepDown:     rec.BackStepDown,

		BarHighPrice:  rec.BarHighPrice,
		BarLowPrice:   rec.BarLowPrice,
		BarClosePrice: rec.BarClosePrice,
		ExistsZigzag:  rec.ExistsZigzag != 0,
	}
}

func toStatusDomain(rec *zigzagStatusRecord, symbolType string, barType fxmodel.BarType) *zigzag.ZigZagStatus {
	minZZ := ""
	if rec.BarDateTimeMinZigZag != nil {
		minZZ = *rec.BarDateTimeMinZigZag
	}
	maxZZ := ""
	if rec.BarDateTimeMaxZigZag != nil {
		maxZZ = *rec.BarDateTimeMaxZigZag
	}
	return &zigzag.ZigZagStatus{
		SymbolType:           symbolType,
		BarType:              string(barType),
		Symbol:               rec.Symbol,
		Depth:                rec.Depth,
		BarDateTimeMin:       rec.BarDateTimeMin,
		BarDateTimeMax:       rec.BarDateTimeMax,
		BarCount:             rec.BarCount,
		BarDateTimeMinZigZag: minZZ,
		BarDateTimeMaxZigZag: maxZZ,
		ZigzagCount:          rec.ZigzagCount,
		BreakResistanceCount: rec.BreakResistanceCount,
		BreakSupportCount:    rec.BreakSupportCount,
	}
}

func toSearchDomain(rec *zigzagSearchRecord) *zigzag.ZigZagSearchRow {
	row := &zigzag.ZigZagSearchRow{
		Symbol:    rec.Symbol,
		Depth:     rec.Depth,
		WaveDxy4h: rec.WaveDxy4h,
		WaveDxy1h: rec.WaveDxy1h,
	}

	row.Current = &zigzag.WaveWithSma{
		WaveStart:  rec.CurWaveStart,
		WaveEnd:    rec.CurWaveEnd,
		Wave:       rec.CurWave,
		Resistance: rec.CurResistance,
		Support:    rec.CurSupport,
		Sma4h200s:  zigzag.SmaPrice{PriceS: rec.CurSma4h200sS, PriceE: rec.CurSma4h200sE},
		Sma4h75s:   zigzag.SmaPrice{PriceS: rec.CurSma4h75sS, PriceE: rec.CurSma4h75sE},
		Sma4h20s:   zigzag.SmaPrice{PriceS: rec.CurSma4h20sS, PriceE: rec.CurSma4h20sE},
		Sma1h200s:  zigzag.SmaPrice{PriceS: rec.CurSma1h200sS, PriceE: rec.CurSma1h200sE},
		Sma15m200s: zigzag.SmaPrice{PriceS: rec.CurSma15m200sS, PriceE: rec.CurSma15m200sE},
	}

	if rec.T4hWaveStart != nil {
		row.Target4h = &zigzag.WaveWithSma{
			WaveStart:  *rec.T4hWaveStart,
			Wave:       rec.T4hWave,
			Resistance: rec.T4hResistance,
			Support:    rec.T4hSupport,
			Sma4h200s:  zigzag.SmaPrice{PriceS: rec.T4hSma4h200sS, PriceE: rec.T4hSma4h200sE},
			Sma4h75s:   zigzag.SmaPrice{PriceS: rec.T4hSma4h75sS, PriceE: rec.T4hSma4h75sE},
			Sma4h20s:   zigzag.SmaPrice{PriceS: rec.T4hSma4h20sS, PriceE: rec.T4hSma4h20sE},
			Sma1h200s:  zigzag.SmaPrice{PriceS: rec.T4hSma1h200sS, PriceE: rec.T4hSma1h200sE},
			Sma15m200s: zigzag.SmaPrice{PriceS: rec.T4hSma15m200sS, PriceE: rec.T4hSma15m200sE},
		}
		if rec.T4hWaveEnd != nil {
			row.Target4h.WaveEnd = *rec.T4hWaveEnd
		}
	}

	if rec.PrvWaveStart != nil {
		row.Previous = &zigzag.WaveInfo{
			WaveStart:  *rec.PrvWaveStart,
			Wave:       rec.PrvWave,
			Resistance: rec.PrvResistance,
			Support:    rec.PrvSupport,
		}
		if rec.PrvWaveEnd != nil {
			row.Previous.WaveEnd = *rec.PrvWaveEnd
		}
	}
	if rec.NxtWaveStart != nil {
		row.Next = &zigzag.WaveInfo{
			WaveStart:  *rec.NxtWaveStart,
			Wave:       rec.NxtWave,
			Resistance: rec.NxtResistance,
			Support:    rec.NxtSupport,
		}
		if rec.NxtWaveEnd != nil {
			row.Next.WaveEnd = *rec.NxtWaveEnd
		}
	}
	if rec.Nx2WaveStart != nil {
		row.Next2 = &zigzag.WaveInfo{
			WaveStart:  *rec.Nx2WaveStart,
			Wave:       rec.Nx2Wave,
			Resistance: rec.Nx2Resistance,
			Support:    rec.Nx2Support,
		}
		if rec.Nx2WaveEnd != nil {
			row.Next2.WaveEnd = *rec.Nx2WaveEnd
		}
	}

	row.FractalWaveList = []zigzag.FractalWaveInfo{}
	return row
}

func appendZigzagWhere(query string, args []any, wave, previousWave, nextWave, next2Wave, wave4h int) (string, []any) {
	if wave != 0 {
		query += " AND C.wave = ?"
		args = append(args, wave)
	}
	if previousWave != 0 {
		query += " AND P.wave = ?"
		args = append(args, previousWave)
	}
	if nextWave != 0 {
		query += " AND N.wave = ?"
		args = append(args, nextWave)
	}
	if next2Wave != 0 {
		query += " AND N2.wave = ?"
		args = append(args, next2Wave)
	}
	if wave4h != 0 {
		query += " AND T4H.wave = ?"
		args = append(args, wave4h)
	}
	return query, args
}

func boolToUint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}
