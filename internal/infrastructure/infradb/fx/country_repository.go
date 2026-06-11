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

type MySQLCountryRepository struct {
	db *sqlx.DB
}

func NewMySQLCountryRepository(db *sqlx.DB) fxrepository.CountryRepository {
	return &MySQLCountryRepository{db: db}
}

type fxCountryRecord struct {
	Code         string    `db:"code"`
	Name         string    `db:"name"`
	CurrencyCode string    `db:"currencyCode"`
	NameEn       string    `db:"nameEn"`
	NameShort    string    `db:"nameShort"`
	SortOrder    int16     `db:"sortOrder"`
	Deleted      uint8     `db:"deleted"`
	CreatedAt    time.Time `db:"createdAt"`
	CreatedBy    string    `db:"createdBy"`
	UpdatedAt    time.Time `db:"updatedAt"`
	UpdatedBy    string    `db:"updatedBy"`
}

func (r *MySQLCountryRepository) GetList(ctx context.Context) ([]model.KeyValue, error) {
	query := `
		SELECT code AS ` + "`key`" + `, name_short AS ` + "`value`" + `
		FROM fx_country
		WHERE (deleted+0) = 0
		ORDER BY sort_order`
	var recs []keyValueRecord
	if err := r.db.SelectContext(ctx, &recs, query); err != nil {
		return nil, err
	}
	return toKeyValues(recs), nil
}

func (r *MySQLCountryRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM fx_country WHERE (deleted+0) = 0`
	var count int
	err := r.db.GetContext(ctx, &count, query)
	return count, err
}

func (r *MySQLCountryRepository) Search(ctx context.Context, page, size int) ([]fxmodel.Country, error) {
	query := `
		SELECT
			code          AS code,
			name          AS name,
			currency_code AS currencyCode,
			name_en       AS nameEn,
			name_short    AS nameShort,
			sort_order    AS sortOrder,
			(deleted+0)   AS deleted,
			created_at    AS createdAt,
			created_by    AS createdBy,
			updated_at    AS updatedAt,
			updated_by    AS updatedBy
		FROM fx_country
		WHERE (deleted+0) = 0
		ORDER BY sort_order
		LIMIT ? OFFSET ?`
	offset := (page - 1) * size
	var recs []fxCountryRecord
	if err := r.db.SelectContext(ctx, &recs, query, size, offset); err != nil {
		return nil, err
	}
	result := make([]fxmodel.Country, len(recs))
	for i, rec := range recs {
		result[i] = toCountryDomain(rec)
	}
	return result, nil
}

func (r *MySQLCountryRepository) Get(ctx context.Context, code string) (*fxmodel.Country, error) {
	query := `
		SELECT
			code          AS code,
			name          AS name,
			currency_code AS currencyCode,
			name_en       AS nameEn,
			name_short    AS nameShort,
			sort_order    AS sortOrder,
			(deleted+0)   AS deleted,
			created_at    AS createdAt,
			created_by    AS createdBy,
			updated_at    AS updatedAt,
			updated_by    AS updatedBy
		FROM fx_country WHERE code = ?`
	var rec fxCountryRecord
	if err := r.db.GetContext(ctx, &rec, query, code); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	c := toCountryDomain(rec)
	return &c, nil
}

func (r *MySQLCountryRepository) Exists(ctx context.Context, code string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM fx_country WHERE code = ?)`
	var exists bool
	err := r.db.GetContext(ctx, &exists, query, code)
	return exists, err
}

func (r *MySQLCountryRepository) Add(ctx context.Context, country fxmodel.Country) error {
	query := `
		INSERT INTO fx_country (
			code, name, currency_code, name_en, name_short,
			sort_order, deleted, created_at, created_by, updated_at, updated_by
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query,
		country.Code, country.Name, country.CurrencyCode, country.NameEn, country.NameShort,
		country.SortOrder, country.Deleted, country.CreatedAt, country.CreatedBy,
		country.UpdatedAt, country.UpdatedBy)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return apperror.NewInsertError(country.Code)
	}
	return nil
}

func (r *MySQLCountryRepository) Update(ctx context.Context, country fxmodel.Country) error {
	query := `
		UPDATE fx_country SET
			name          = ?,
			currency_code = ?,
			name_en       = ?,
			name_short    = ?,
			sort_order    = ?,
			updated_at    = ?,
			updated_by    = ?
		WHERE code = ?`
	result, err := r.db.ExecContext(ctx, query,
		country.Name, country.CurrencyCode, country.NameEn, country.NameShort,
		country.SortOrder, country.UpdatedAt, country.UpdatedBy, country.Code)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return apperror.NewUpdateError(country.Code)
	}
	return nil
}

func (r *MySQLCountryRepository) UpdateCode(ctx context.Context, country fxmodel.Country, baseCode string) error {
	query := `
		UPDATE fx_country SET
			code          = ?,
			name          = ?,
			currency_code = ?,
			name_en       = ?,
			name_short    = ?,
			sort_order    = ?,
			updated_at    = ?,
			updated_by    = ?
		WHERE code = ?`
	result, err := r.db.ExecContext(ctx, query,
		country.Code, country.Name, country.CurrencyCode, country.NameEn, country.NameShort,
		country.SortOrder, country.UpdatedAt, country.UpdatedBy, baseCode)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return apperror.NewUpdateError(baseCode)
	}
	return nil
}

// RefreshCache は GetList が DB を直接参照するため no-op。
func (r *MySQLCountryRepository) RefreshCache(_ context.Context) error {
	return nil
}

func toCountryDomain(rec fxCountryRecord) fxmodel.Country {
	return fxmodel.Country{
		Code:         rec.Code,
		Name:         rec.Name,
		CurrencyCode: rec.CurrencyCode,
		NameEn:       rec.NameEn,
		NameShort:    rec.NameShort,
		SortOrder:    rec.SortOrder,
		Deleted:      rec.Deleted != 0,
		CreatedAt:    rec.CreatedAt,
		CreatedBy:    rec.CreatedBy,
		UpdatedAt:    rec.UpdatedAt,
		UpdatedBy:    rec.UpdatedBy,
	}
}
