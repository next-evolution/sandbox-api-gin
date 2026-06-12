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

type MySQLEconomicIndicatorRepository struct {
	db *sqlx.DB
}

func NewMySQLEconomicIndicatorRepository(db *sqlx.DB) fxrepository.EconomicIndicatorRepository {
	return &MySQLEconomicIndicatorRepository{db: db}
}

type fxEconomicIndicatorRecord struct {
	ID               int64     `db:"id"`
	CountryCode      string    `db:"countryCode"`
	Name             string    `db:"name"`
	Importance       string    `db:"importance"`
	Description      string    `db:"description"`
	UnitOfValue      string    `db:"unitOfValue"`
	CountryName      string    `db:"countryName"`
	CountryNameShort string    `db:"countryNameShort"`
	Deleted          uint8     `db:"deleted"`
	CreatedAt        time.Time `db:"createdAt"`
	CreatedBy        string    `db:"createdBy"`
	UpdatedAt        time.Time `db:"updatedAt"`
	UpdatedBy        string    `db:"updatedBy"`
}

func (r *MySQLEconomicIndicatorRepository) GetList(ctx context.Context, countryCode string) ([]model.KeyValue, error) {
	var (
		query string
		args  []interface{}
	)
	if countryCode == "ALL" {
		query = `
			SELECT CAST(id AS CHAR) AS ` + "`key`" + `, name AS ` + "`value`" + `
			FROM fx_economic_indicator
			WHERE (deleted+0) = 0
			ORDER BY name`
	} else {
		query = `
			SELECT CAST(id AS CHAR) AS ` + "`key`" + `, name AS ` + "`value`" + `
			FROM fx_economic_indicator
			WHERE (deleted+0) = 0 AND country_code = ?
			ORDER BY name`
		args = []interface{}{countryCode}
	}
	var recs []keyValueRecord
	if err := r.db.SelectContext(ctx, &recs, query, args...); err != nil {
		return nil, err
	}
	return toKeyValues(recs), nil
}

func (r *MySQLEconomicIndicatorRepository) Count(ctx context.Context, countryCode, importance, name string) (int, error) {
	query := `SELECT COUNT(id) FROM fx_economic_indicator WHERE (deleted+0) = 0`
	args := make([]interface{}, 0)
	if importance != "" {
		query += ` AND importance = ?`
		args = append(args, importance)
	}
	if countryCode != "" {
		query += ` AND country_code = ?`
		args = append(args, countryCode)
	}
	if name != "" {
		query += ` AND name LIKE ?`
		args = append(args, "%"+name+"%")
	}
	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	return count, err
}

func (r *MySQLEconomicIndicatorRepository) Search(ctx context.Context, page, size int, countryCode, importance, name string) ([]fxmodel.EconomicIndicator, error) {
	query := `
		SELECT
			t.id              AS id,
			t.country_code    AS countryCode,
			t.name            AS name,
			t.importance      AS importance,
			COALESCE(t.description, '') AS description,
			t.unit_of_value   AS unitOfValue,
			c.name            AS countryName,
			c.name_short      AS countryNameShort
		FROM fx_economic_indicator t
		INNER JOIN fx_country c ON c.code = t.country_code
		WHERE t.deleted = 0`
	args := make([]interface{}, 0)
	if importance != "" {
		query += ` AND t.importance = ?`
		args = append(args, importance)
	}
	if countryCode != "" {
		query += ` AND t.country_code = ?`
		args = append(args, countryCode)
	}
	if name != "" {
		query += ` AND t.name LIKE ?`
		args = append(args, "%"+name+"%")
	}
	query += ` ORDER BY c.sort_order, t.id LIMIT ? OFFSET ?`
	args = append(args, size, (page-1)*size)

	var recs []fxEconomicIndicatorRecord
	if err := r.db.SelectContext(ctx, &recs, query, args...); err != nil {
		return nil, err
	}
	result := make([]fxmodel.EconomicIndicator, len(recs))
	for i, rec := range recs {
		result[i] = toEconomicIndicatorDomain(rec)
	}
	return result, nil
}

func (r *MySQLEconomicIndicatorRepository) Get(ctx context.Context, id int64) (*fxmodel.EconomicIndicator, error) {
	query := `
		SELECT
			t.id              AS id,
			t.country_code    AS countryCode,
			t.name            AS name,
			t.importance      AS importance,
			COALESCE(t.description, '') AS description,
			t.unit_of_value   AS unitOfValue,
			(t.deleted+0)     AS deleted,
			t.created_at      AS createdAt,
			t.created_by      AS createdBy,
			t.updated_at      AS updatedAt,
			t.updated_by      AS updatedBy,
			c.name            AS countryName,
			c.name_short      AS countryNameShort
		FROM fx_economic_indicator t
		INNER JOIN fx_country c ON c.code = t.country_code
		WHERE t.id = ?`
	var rec fxEconomicIndicatorRecord
	if err := r.db.GetContext(ctx, &rec, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	m := toEconomicIndicatorDomain(rec)
	return &m, nil
}

func (r *MySQLEconomicIndicatorRepository) Exists(ctx context.Context, countryCode, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM fx_economic_indicator WHERE country_code = ? AND name = ?)`
	var exists bool
	err := r.db.GetContext(ctx, &exists, query, countryCode, name)
	return exists, err
}

func (r *MySQLEconomicIndicatorRepository) Add(ctx context.Context, indicator fxmodel.EconomicIndicator) error {
	query := `
		INSERT INTO fx_economic_indicator (
			country_code, name, importance, description, unit_of_value,
			deleted, created_at, created_by, updated_at, updated_by
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query,
		indicator.CountryCode, indicator.Name, indicator.Importance,
		indicator.Description, indicator.UnitOfValue, indicator.Deleted,
		indicator.CreatedAt, indicator.CreatedBy, indicator.UpdatedAt, indicator.UpdatedBy,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return apperror.NewInsertError(indicator.Name)
	}
	return nil
}

func (r *MySQLEconomicIndicatorRepository) Update(ctx context.Context, indicator fxmodel.EconomicIndicator, countryCode string) error {
	query := `
		UPDATE fx_economic_indicator SET
			name          = ?,
			importance    = ?,
			country_code  = ?,
			description   = ?,
			unit_of_value = ?,
			updated_at    = ?,
			updated_by    = ?
		WHERE id = ? AND country_code = ?`
	result, err := r.db.ExecContext(ctx, query,
		indicator.Name, indicator.Importance, indicator.CountryCode,
		indicator.Description, indicator.UnitOfValue,
		indicator.UpdatedAt, indicator.UpdatedBy,
		indicator.ID, countryCode,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return apperror.NewUpdateError(indicator.Name)
	}
	return nil
}

func (r *MySQLEconomicIndicatorRepository) GetEconomicIndicatorList(ctx context.Context, countryCode string) ([]fxmodel.EconomicIndicator, error) {
	query := `
		SELECT
			t.id            AS id,
			t.country_code  AS countryCode,
			t.name          AS name,
			t.importance    AS importance,
			COALESCE(t.description, '') AS description,
			t.unit_of_value AS unitOfValue,
			c.name          AS countryName,
			c.name_short    AS countryNameShort
		FROM fx_economic_indicator t
		INNER JOIN fx_country c ON c.code = t.country_code
		WHERE (t.deleted+0) = 0 AND t.country_code = ?`
	var recs []fxEconomicIndicatorRecord
	if err := r.db.SelectContext(ctx, &recs, query, countryCode); err != nil {
		return nil, err
	}
	result := make([]fxmodel.EconomicIndicator, len(recs))
	for i, rec := range recs {
		result[i] = toEconomicIndicatorDomain(rec)
	}
	return result, nil
}

// RefreshCache はDBを直接参照するため no-op。
func (r *MySQLEconomicIndicatorRepository) RefreshCache(_ context.Context, _ string) error {
	return nil
}

func toEconomicIndicatorDomain(rec fxEconomicIndicatorRecord) fxmodel.EconomicIndicator {
	return fxmodel.EconomicIndicator{
		ID:               rec.ID,
		CountryCode:      rec.CountryCode,
		Name:             rec.Name,
		Importance:       rec.Importance,
		Description:      rec.Description,
		UnitOfValue:      rec.UnitOfValue,
		CountryName:      rec.CountryName,
		CountryNameShort: rec.CountryNameShort,
		Deleted:          rec.Deleted != 0,
		CreatedAt:        rec.CreatedAt,
		CreatedBy:        rec.CreatedBy,
		UpdatedAt:        rec.UpdatedAt,
		UpdatedBy:        rec.UpdatedBy,
	}
}
