package fxrepository

import (
	"context"

	fxmodel "sandbox-api-gin/internal/domain/model/fx"
)

type BarDataRepository interface {
	StatusList(ctx context.Context, symbolType string, barType fxmodel.BarType) ([]fxmodel.BarDataStatus, error)
	SearchCount(ctx context.Context, symbol string, barType fxmodel.BarType, barDateFrom, barDateTo string) (int, error)
	Search(ctx context.Context, symbol string, barType fxmodel.BarType, barDateFrom, barDateTo string, sortAsc bool, page, size int) ([]fxmodel.BarData, error)

	// ロードテーブル初期化
	DeleteLoad(ctx context.Context, symbol string) error
	DeleteLoadSma(ctx context.Context, symbol string) error
	DeleteLoadRsi(ctx context.Context, symbol string) error

	// ロードテーブルへのバルクロード
	BulkLoad(ctx context.Context, list []fxmodel.BarLoadData) error
	BulkLoadSma(ctx context.Context, list []fxmodel.BarLoadSma) error
	BulkLoadRsi(ctx context.Context, list []fxmodel.BarLoadRsi) error

	// ロードテーブルの最新1件削除
	DeleteLatestLoad(ctx context.Context, symbol string) error

	// ロードテーブルの最新足日時取得
	GetLatestLoadBarDateTime(ctx context.Context, symbol string) (string, error)

	// 既存データとの整合性チェック
	ImportCheck(ctx context.Context, barType fxmodel.BarType, symbol string) (fxmodel.BarCsvImportCheck, error)

	// ロードテーブルから本テーブルへインサート
	InsertFromLoad(ctx context.Context, symbol string, barType fxmodel.BarType) (int, error)
	InsertFromLoadSma(ctx context.Context, symbol string, barType fxmodel.BarType) (int, error)
	InsertFromLoadRsi(ctx context.Context, symbol string, barType fxmodel.BarType) (int, error)

	// 差分データ取得
	GetDiffBarData(ctx context.Context, symbol string, barType fxmodel.BarType) ([]fxmodel.BarLoadData, error)
	GetDiffBarSma(ctx context.Context, symbol string, barType fxmodel.BarType) ([]fxmodel.BarLoadSma, error)
	GetDiffBarRsi(ctx context.Context, symbol string, barType fxmodel.BarType) ([]fxmodel.BarLoadRsi, error)

	// 差分データ更新
	UpdateBarData(ctx context.Context, symbol string, barType fxmodel.BarType) (int, error)
	UpdateBarSma(ctx context.Context, symbol string, barType fxmodel.BarType) (int, error)
	UpdateBarRsi(ctx context.Context, symbol string, barType fxmodel.BarType) (int, error)
}
