package infradbfx

import (
	"context"

	"github.com/jmoiron/sqlx"

	"sandbox-api-gin/internal/domain/model"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type MySQLEconomicIndicatorRepository struct {
	db *sqlx.DB
}

func NewMySQLEconomicIndicatorRepository(db *sqlx.DB) fxrepository.EconomicIndicatorRepository {
	return &MySQLEconomicIndicatorRepository{db: db}
}

func (r *MySQLEconomicIndicatorRepository) GetList(ctx context.Context, countryCode string) ([]model.KeyValue, error) {
	// JavaのEconomicIndicatorMapper.getList(): countryCode="ALL"のとき全件取得
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
