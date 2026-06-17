package infradbfx

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"sandbox-api-gin/internal/domain/apperror"
	fxmodel "sandbox-api-gin/internal/domain/model/fx"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type MySQLEconomicIndicatorDataRepository struct {
	db *sqlx.DB
}

func NewMySQLEconomicIndicatorDataRepository(db *sqlx.DB) fxrepository.EconomicIndicatorDataRepository {
	return &MySQLEconomicIndicatorDataRepository{db: db}
}

type fxEconomicIndicatorDataRecord struct {
	Code             string    `db:"code"`
	CountryCode      string    `db:"countryCode"`
	Name             string    `db:"name"`
	Importance       string    `db:"importance"`
	Description      string    `db:"description"`
	Publication      time.Time `db:"publication"`
	PublicationDate  string    `db:"publicationDate"`
	PublicationTime  string    `db:"publicationTime"`
	DayOfWeek        int16     `db:"dayOfWeek"`
	SubTitle         string    `db:"subTitle"`
	ResultValue      string    `db:"resultValue"`
	ForecastValue    string    `db:"forecastValue"`
	PreviousValue    string    `db:"previousValue"`
	UnitOfValue      string    `db:"unitOfValue"`
	Memo             string    `db:"memo"`
	CountryName      string    `db:"countryName"`
	CountryNameShort string    `db:"countryNameShort"`
}

const selectEconomicIndicatorDataColumns = `
	T.code                                  AS code,
	T.country_code                          AS countryCode,
	E.importance                            AS importance,
	C.name                                  AS countryName,
	C.name_short                            AS countryNameShort,
	E.name                                  AS name,
	COALESCE(E.description, '')             AS description,
	T.publication                           AS publication,
	DATE_FORMAT(T.publication, '%Y-%m-%d')  AS publicationDate,
	DATE_FORMAT(T.publication, '%H:%i')     AS publicationTime,
	DAYOFWEEK(T.publication) - 1           AS dayOfWeek,
	COALESCE(T.sub_title, '')               AS subTitle,
	T.result_value                          AS resultValue,
	COALESCE(T.forecast_value, '')          AS forecastValue,
	COALESCE(T.previous_value, '')          AS previousValue,
	COALESCE(E.unit_of_value, '')           AS unitOfValue,
	COALESCE(T.memo, '')                    AS memo`

func (r *MySQLEconomicIndicatorDataRepository) buildWhere(code, importance, countryCode, publicationBaseDate string) (string, []interface{}) {
	where := ``
	args := make([]interface{}, 0)
	if code != "" {
		where += ` AND T.code = ?`
		args = append(args, code)
	}
	if importance != "" {
		where += ` AND E.importance = ?`
		args = append(args, importance)
	}
	if countryCode != "" {
		where += ` AND T.country_code = ?`
		args = append(args, countryCode)
	}
	if publicationBaseDate != "" {
		where += ` AND T.publication >= STR_TO_DATE(?, '%Y-%m-%d')`
		args = append(args, publicationBaseDate)
	}
	return where, args
}

func (r *MySQLEconomicIndicatorDataRepository) Count(ctx context.Context, code, importance, countryCode, publicationBaseDate string) (int, error) {
	where, args := r.buildWhere(code, importance, countryCode, publicationBaseDate)
	query := `
		SELECT COUNT(*)
		FROM fx_economic_indicator_data T
		INNER JOIN fx_economic_indicator E ON E.code = T.code AND E.country_code = T.country_code
		WHERE 1 = 1` + where
	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	return count, err
}

func (r *MySQLEconomicIndicatorDataRepository) Search(ctx context.Context, code, importance, countryCode, publicationBaseDate string, page, size int, sortAsc bool) ([]fxmodel.EconomicIndicatorData, error) {
	where, args := r.buildWhere(code, importance, countryCode, publicationBaseDate)
	order := `ORDER BY T.publication DESC`
	if sortAsc {
		order = `ORDER BY T.publication`
	}
	query := `
		SELECT ` + selectEconomicIndicatorDataColumns + `
		FROM fx_economic_indicator_data T
		INNER JOIN fx_economic_indicator E ON E.code = T.code AND E.country_code = T.country_code
		INNER JOIN fx_country C ON C.code = T.country_code
		WHERE 1 = 1` + where + ` ` + order + ` LIMIT ? OFFSET ?`
	args = append(args, size, (page-1)*size)

	var recs []fxEconomicIndicatorDataRecord
	if err := r.db.SelectContext(ctx, &recs, query, args...); err != nil {
		return nil, err
	}
	return toEconomicIndicatorDataList(recs), nil
}

func (r *MySQLEconomicIndicatorDataRepository) Get(ctx context.Context, code, countryCode string, publication time.Time) (*fxmodel.EconomicIndicatorData, error) {
	query := `
		SELECT ` + selectEconomicIndicatorDataColumns + `
		FROM fx_economic_indicator_data T
		INNER JOIN fx_economic_indicator E ON E.code = T.code AND E.country_code = T.country_code
		INNER JOIN fx_country C ON C.code = T.country_code
		WHERE T.code = ? AND T.country_code = ? AND T.publication = ?`
	var rec fxEconomicIndicatorDataRecord
	if err := r.db.GetContext(ctx, &rec, query, code, countryCode, publication); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	m := toEconomicIndicatorDataDomain(rec)
	return &m, nil
}

func (r *MySQLEconomicIndicatorDataRepository) Exists(ctx context.Context, code, countryCode string, publication time.Time) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM fx_economic_indicator_data WHERE code = ? AND country_code = ? AND publication = ?)`
	var exists bool
	err := r.db.GetContext(ctx, &exists, query, code, countryCode, publication)
	return exists, err
}

func (r *MySQLEconomicIndicatorDataRepository) Add(ctx context.Context, data fxmodel.EconomicIndicatorData) error {
	query := `
		INSERT INTO fx_economic_indicator_data
		(code, country_code, publication, sub_title, result_value, forecast_value, previous_value, memo)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query,
		data.Code, data.CountryCode, data.Publication, data.SubTitle,
		data.ResultValue, data.ForecastValue, data.PreviousValue, data.Memo,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return apperror.NewInsertError(data.Publication.Format("2006-01-02 15:04:05"))
	}
	return nil
}

func (r *MySQLEconomicIndicatorDataRepository) Update(ctx context.Context, data fxmodel.EconomicIndicatorData, publication time.Time) error {
	query := `
		UPDATE fx_economic_indicator_data SET
			publication    = ?,
			sub_title      = ?,
			result_value   = ?,
			forecast_value = ?,
			previous_value = ?,
			memo           = ?
		WHERE code = ? AND country_code = ? AND publication = ?`
	result, err := r.db.ExecContext(ctx, query,
		data.Publication, data.SubTitle, data.ResultValue,
		data.ForecastValue, data.PreviousValue, data.Memo,
		data.Code, data.CountryCode, publication,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return apperror.NewUpdateError(publication.Format("2006-01-02 15:04:05"))
	}
	return nil
}

func (r *MySQLEconomicIndicatorDataRepository) UpdateCode(ctx context.Context, data fxmodel.EconomicIndicatorData, code, countryCode string, publication time.Time) error {
	query := `
		UPDATE fx_economic_indicator_data SET
			code           = ?,
			country_code   = ?,
			publication    = ?,
			sub_title      = ?,
			result_value   = ?,
			forecast_value = ?,
			previous_value = ?,
			memo           = ?
		WHERE code = ? AND country_code = ? AND publication = ?`
	result, err := r.db.ExecContext(ctx, query,
		data.Code, data.CountryCode, data.Publication, data.SubTitle, data.ResultValue,
		data.ForecastValue, data.PreviousValue, data.Memo,
		code, countryCode, publication,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return apperror.NewUpdateError(publication.Format("2006-01-02 15:04:05"))
	}
	return nil
}

func (r *MySQLEconomicIndicatorDataRepository) DeleteLoad(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM fx_economic_indicator_data_load`)
	return err
}

func (r *MySQLEconomicIndicatorDataRepository) InsertLoad(ctx context.Context, data fxmodel.EconomicIndicatorData) error {
	query := `
		INSERT INTO fx_economic_indicator_data_load
		(code, country_code, publication, sub_title, result_value, forecast_value, previous_value, memo)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		data.Code, data.CountryCode, data.Publication, data.SubTitle,
		data.ResultValue, data.ForecastValue, data.PreviousValue, data.Memo,
	)
	return err
}

func (r *MySQLEconomicIndicatorDataRepository) LoadDiff(ctx context.Context) ([]fxmodel.EconomicIndicatorData, error) {
	query := `
		SELECT
			T.code            AS code,
			T.country_code    AS countryCode,
			E.importance      AS importance,
			C.name            AS countryName,
			C.name_short      AS countryNameShort,
			E.name            AS name,
			COALESCE(E.description, '')   AS description,
			T.publication     AS publication,
			''                AS publicationDate,
			''                AS publicationTime,
			0                 AS dayOfWeek,
			COALESCE(T.sub_title, '')      AS subTitle,
			T.result_value    AS resultValue,
			COALESCE(T.forecast_value, '') AS forecastValue,
			COALESCE(T.previous_value, '') AS previousValue,
			COALESCE(E.unit_of_value, '')  AS unitOfValue,
			COALESCE(T.memo, '')           AS memo
		FROM fx_economic_indicator_data T
		INNER JOIN fx_economic_indicator E ON E.code = T.code AND E.country_code = T.country_code
		INNER JOIN fx_country C ON C.code = T.country_code
		INNER JOIN fx_economic_indicator_data_load L
			ON T.code = L.code AND T.country_code = L.country_code AND T.publication = L.publication
		WHERE T.result_value != L.result_value
		   OR T.forecast_value != L.forecast_value
		   OR T.previous_value != L.previous_value`
	var recs []fxEconomicIndicatorDataRecord
	if err := r.db.SelectContext(ctx, &recs, query); err != nil {
		return nil, err
	}
	return toEconomicIndicatorDataList(recs), nil
}

func (r *MySQLEconomicIndicatorDataRepository) InsertFromLoad(ctx context.Context) error {
	query := `
		INSERT INTO fx_economic_indicator_data
		SELECT L.code, L.country_code, L.publication, L.sub_title, L.result_value, L.forecast_value, L.previous_value, L.memo
		FROM fx_economic_indicator_data_load L
		LEFT JOIN fx_economic_indicator_data D
			ON D.code = L.code AND D.country_code = L.country_code AND D.publication = L.publication
		WHERE D.code IS NULL`
	_, err := r.db.ExecContext(ctx, query)
	return err
}

func toEconomicIndicatorDataList(recs []fxEconomicIndicatorDataRecord) []fxmodel.EconomicIndicatorData {
	result := make([]fxmodel.EconomicIndicatorData, len(recs))
	for i, rec := range recs {
		result[i] = toEconomicIndicatorDataDomain(rec)
	}
	return result
}

func toEconomicIndicatorDataDomain(rec fxEconomicIndicatorDataRecord) fxmodel.EconomicIndicatorData {
	return fxmodel.EconomicIndicatorData{
		Code:             rec.Code,
		CountryCode:      rec.CountryCode,
		Name:             rec.Name,
		Importance:       rec.Importance,
		Description:      rec.Description,
		Publication:      rec.Publication,
		PublicationDate:  rec.PublicationDate,
		PublicationTime:  rec.PublicationTime,
		DayOfWeek:        rec.DayOfWeek,
		SubTitle:         rec.SubTitle,
		ResultValue:      rec.ResultValue,
		ForecastValue:    rec.ForecastValue,
		PreviousValue:    rec.PreviousValue,
		UnitOfValue:      rec.UnitOfValue,
		Memo:             rec.Memo,
		CountryName:      rec.CountryName,
		CountryNameShort: rec.CountryNameShort,
	}
}
