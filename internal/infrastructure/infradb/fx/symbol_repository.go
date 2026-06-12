package infradbfx

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"sandbox-api-gin/internal/domain/apperror"
	"sandbox-api-gin/internal/domain/model"
	fxmodel "sandbox-api-gin/internal/domain/model/fx"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type keyValueRecord struct {
	Key   string `db:"key"`
	Value string `db:"value"`
}

type fxSymbolRecord struct {
	Symbol           string    `db:"symbol"`
	SymbolType       string    `db:"symbolType"`
	Name             string    `db:"name"`
	ValidScale       int16     `db:"validScale"`
	TargetVolatility float64   `db:"targetVolatility"`
	SortOrder        int       `db:"sortOrder"`
	Deleted          uint8     `db:"deleted"`
	CreatedAt        time.Time `db:"createdAt"`
	CreatedBy        string    `db:"createdBy"`
	UpdatedAt        time.Time `db:"updatedAt"`
	UpdatedBy        string    `db:"updatedBy"`
}

type MySQLSymbolRepository struct {
	db *sqlx.DB
}

func NewMySQLSymbolRepository(db *sqlx.DB) fxrepository.SymbolRepository {
	return &MySQLSymbolRepository{db: db}
}

func (r *MySQLSymbolRepository) GetList(ctx context.Context, symbolType string) ([]model.KeyValue, error) {
	query := `
		SELECT symbol AS ` + "`key`" + `, name AS ` + "`value`" + `
		FROM fx_symbol
		WHERE (deleted+0) = 0 AND symbol_type = ?
		ORDER BY sort_order`
	var recs []keyValueRecord
	if err := r.db.SelectContext(ctx, &recs, query, symbolType); err != nil {
		return nil, err
	}
	return toKeyValues(recs), nil
}

func (r *MySQLSymbolRepository) GetTradingSymbols(ctx context.Context) ([]string, error) {
	query := `SELECT symbol FROM fx_symbol WHERE symbol_type = ? AND (deleted+0) = 0`
	var symbols []string
	if err := r.db.SelectContext(ctx, &symbols, query, "Trade"); err != nil {
		return nil, err
	}
	return symbols, nil
}

func (r *MySQLSymbolRepository) Count(ctx context.Context, symbolType string) (int, error) {
	query := `SELECT COUNT(symbol) FROM fx_symbol WHERE (deleted+0) = 0`
	args := make([]any, 0)
	if symbolType != "" {
		query += ` AND symbol_type = ?`
		args = append(args, symbolType)
	}
	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	return count, err
}

func (r *MySQLSymbolRepository) Search(ctx context.Context, symbolType string, page, size int) ([]fxmodel.Symbol, error) {
	query := `
		SELECT
			symbol         AS symbol,
			symbol_type    AS symbolType,
			name           AS name,
			valid_scale    AS validScale,
			target_volatility AS targetVolatility,
			sort_order     AS sortOrder,
			(deleted+0)    AS deleted,
			created_at     AS createdAt,
			created_by     AS createdBy,
			updated_at     AS updatedAt,
			updated_by     AS updatedBy
		FROM fx_symbol
		WHERE (deleted+0) = 0`
	args := make([]any, 0)
	if symbolType != "" {
		query += ` AND symbol_type = ?`
		args = append(args, symbolType)
	}
	offset := (page - 1) * size
	query += ` ORDER BY sort_order LIMIT ? OFFSET ?`
	args = append(args, size, offset)

	var recs []fxSymbolRecord
	if err := r.db.SelectContext(ctx, &recs, query, args...); err != nil {
		return nil, err
	}
	result := make([]fxmodel.Symbol, len(recs))
	for i, rec := range recs {
		result[i] = toSymbolDomain(rec)
	}
	return result, nil
}

func (r *MySQLSymbolRepository) Get(ctx context.Context, symbol string) (*fxmodel.Symbol, error) {
	query := `
		SELECT
			symbol         AS symbol,
			symbol_type    AS symbolType,
			name           AS name,
			valid_scale    AS validScale,
			target_volatility AS targetVolatility,
			sort_order     AS sortOrder,
			(deleted+0)    AS deleted,
			created_at     AS createdAt,
			created_by     AS createdBy,
			updated_at     AS updatedAt,
			updated_by     AS updatedBy
		FROM fx_symbol
		WHERE symbol = ?`
	var rec fxSymbolRecord
	if err := r.db.GetContext(ctx, &rec, query, symbol); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	s := toSymbolDomain(rec)
	return &s, nil
}

func (r *MySQLSymbolRepository) Exists(ctx context.Context, symbol string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM fx_symbol WHERE symbol = ?)`
	var exists bool
	err := r.db.GetContext(ctx, &exists, query, symbol)
	return exists, err
}

func (r *MySQLSymbolRepository) Add(ctx context.Context, symbol fxmodel.Symbol) error {
	query := `
		INSERT INTO fx_symbol (
			symbol, symbol_type, name, valid_scale, target_volatility,
			sort_order, deleted, created_at, created_by, updated_at, updated_by
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query,
		symbol.Symbol, symbol.SymbolType, symbol.Name, symbol.ValidScale, symbol.TargetVolatility,
		symbol.SortOrder, symbol.Deleted, symbol.CreatedAt, symbol.CreatedBy,
		symbol.UpdatedAt, symbol.UpdatedBy)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return apperror.NewInsertError(symbol.Symbol)
	}
	return nil
}

func (r *MySQLSymbolRepository) Update(ctx context.Context, symbol fxmodel.Symbol) error {
	query := `
		UPDATE fx_symbol SET
			symbol_type       = ?,
			name              = ?,
			valid_scale       = ?,
			target_volatility = ?,
			sort_order        = ?,
			updated_at        = ?,
			updated_by        = ?
		WHERE symbol = ?`
	result, err := r.db.ExecContext(ctx, query,
		symbol.SymbolType, symbol.Name, symbol.ValidScale, symbol.TargetVolatility,
		symbol.SortOrder, symbol.UpdatedAt, symbol.UpdatedBy, symbol.Symbol)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return apperror.NewUpdateError(symbol.Symbol)
	}
	return nil
}

func (r *MySQLSymbolRepository) UpdateSymbol(ctx context.Context, symbol fxmodel.Symbol, baseSymbol string) error {
	query := `
		UPDATE fx_symbol SET
			symbol            = ?,
			symbol_type       = ?,
			name              = ?,
			valid_scale       = ?,
			target_volatility = ?,
			sort_order        = ?,
			updated_at        = ?,
			updated_by        = ?
		WHERE symbol = ?`
	result, err := r.db.ExecContext(ctx, query,
		symbol.Symbol, symbol.SymbolType, symbol.Name, symbol.ValidScale, symbol.TargetVolatility,
		symbol.SortOrder, symbol.UpdatedAt, symbol.UpdatedBy, baseSymbol)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return apperror.NewUpdateError(baseSymbol)
	}
	return nil
}

// RefreshCache は GetList が DB を直接参照するため no-op。
func (r *MySQLSymbolRepository) RefreshCache(_ context.Context, _ string) error {
	return nil
}

func toKeyValues(recs []keyValueRecord) []model.KeyValue {
	result := make([]model.KeyValue, len(recs))
	for i, r := range recs {
		result[i] = model.KeyValue{Key: r.Key, Value: r.Value}
	}
	return result
}

func toSymbolDomain(rec fxSymbolRecord) fxmodel.Symbol {
	return fxmodel.Symbol{
		Symbol:           rec.Symbol,
		SymbolType:       rec.SymbolType,
		Name:             rec.Name,
		ValidScale:       rec.ValidScale,
		TargetVolatility: rec.TargetVolatility,
		SortOrder:        rec.SortOrder,
		Deleted:          rec.Deleted != 0,
		CreatedAt:        rec.CreatedAt,
		CreatedBy:        rec.CreatedBy,
		UpdatedAt:        rec.UpdatedAt,
		UpdatedBy:        rec.UpdatedBy,
	}
}
