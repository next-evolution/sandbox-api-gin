package infradbfx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"sandbox-api-gin/internal/domain/apperror"
	fxmodel "sandbox-api-gin/internal/domain/model/fx"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type MySQLSummerTimeRepository struct {
	db *sqlx.DB
}

func NewMySQLSummerTimeRepository(db *sqlx.DB) fxrepository.SummerTimeRepository {
	return &MySQLSummerTimeRepository{db: db}
}

type fxSummerTimeRecord struct {
	TargetYear int16     `db:"targetYear"`
	ApplyStart time.Time `db:"applyStart"`
	ApplyEnd   time.Time `db:"applyEnd"`
}

func (r *MySQLSummerTimeRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM fx_summer_time`
	var count int
	err := r.db.GetContext(ctx, &count, query)
	return count, err
}

func (r *MySQLSummerTimeRepository) Search(ctx context.Context, page, size int) ([]fxmodel.SummerTime, error) {
	query := `
		SELECT
			target_year AS targetYear,
			apply_start AS applyStart,
			apply_end   AS applyEnd
		FROM fx_summer_time
		ORDER BY target_year
		LIMIT ? OFFSET ?`
	offset := (page - 1) * size
	var recs []fxSummerTimeRecord
	if err := r.db.SelectContext(ctx, &recs, query, size, offset); err != nil {
		return nil, err
	}
	result := make([]fxmodel.SummerTime, len(recs))
	for i, rec := range recs {
		result[i] = toSummerTimeDomain(rec)
	}
	return result, nil
}

func (r *MySQLSummerTimeRepository) Get(ctx context.Context, targetYear int16) (*fxmodel.SummerTime, error) {
	query := `
		SELECT
			target_year AS targetYear,
			apply_start AS applyStart,
			apply_end   AS applyEnd
		FROM fx_summer_time WHERE target_year = ?`
	var rec fxSummerTimeRecord
	if err := r.db.GetContext(ctx, &rec, query, targetYear); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	s := toSummerTimeDomain(rec)
	return &s, nil
}

func (r *MySQLSummerTimeRepository) Exists(ctx context.Context, targetYear int16) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM fx_summer_time WHERE target_year = ?)`
	var exists bool
	err := r.db.GetContext(ctx, &exists, query, targetYear)
	return exists, err
}

func (r *MySQLSummerTimeRepository) Add(ctx context.Context, s fxmodel.SummerTime) error {
	query := `INSERT INTO fx_summer_time (target_year, apply_start, apply_end) VALUES (?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, s.TargetYear, s.ApplyStart, s.ApplyEnd)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return apperror.NewInsertError(fmt.Sprintf("%d", s.TargetYear))
	}
	return nil
}

func (r *MySQLSummerTimeRepository) Update(ctx context.Context, s fxmodel.SummerTime) error {
	query := `UPDATE fx_summer_time SET apply_start = ?, apply_end = ? WHERE target_year = ?`
	result, err := r.db.ExecContext(ctx, query, s.ApplyStart, s.ApplyEnd, s.TargetYear)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return apperror.NewUpdateError(fmt.Sprintf("%d", s.TargetYear))
	}
	return nil
}

func (r *MySQLSummerTimeRepository) UpdateYear(ctx context.Context, s fxmodel.SummerTime, baseYear int16) error {
	query := `UPDATE fx_summer_time SET target_year = ?, apply_start = ?, apply_end = ? WHERE target_year = ?`
	result, err := r.db.ExecContext(ctx, query, s.TargetYear, s.ApplyStart, s.ApplyEnd, baseYear)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return apperror.NewUpdateError(fmt.Sprintf("%d", baseYear))
	}
	return nil
}

func toSummerTimeDomain(rec fxSummerTimeRecord) fxmodel.SummerTime {
	return fxmodel.SummerTime{
		TargetYear: rec.TargetYear,
		ApplyStart: rec.ApplyStart,
		ApplyEnd:   rec.ApplyEnd,
	}
}
