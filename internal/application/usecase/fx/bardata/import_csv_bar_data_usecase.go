package bardata

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	fxcommand "sandbox-api-gin/internal/application/command/fx"
	fxmodel "sandbox-api-gin/internal/domain/model/fx"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

const csvRsiRange = 14

type BarDataImportResult struct {
	Symbol          string `json:"symbol,omitempty"`
	BarDateTime     string `json:"barDateTime,omitempty"`
	FileName        string `json:"fileName"`
	FileSize        int64  `json:"fileSize"`
	ResultStatus    string `json:"resultStatus"`
	ReadCount       int    `json:"readCount"`
	ExistsCount     int    `json:"existsCount"`
	InsertCount     int    `json:"insertCount"`
	DifferenceCount int    `json:"differenceCount"`
	Message         string `json:"message,omitempty"`
}

type ImportCsvBarDataUseCase struct {
	repo            fxrepository.BarDataRepository
	storageBucket   string
	storageFX       string
	bulkLoadSize    int
	importCheckSkip bool
}

func NewImportCsvBarDataUseCase(
	repo fxrepository.BarDataRepository,
	storageBucket, storageFX string,
	bulkLoadSize int,
	importCheckSkip bool,
) *ImportCsvBarDataUseCase {
	return &ImportCsvBarDataUseCase{
		repo:            repo,
		storageBucket:   storageBucket,
		storageFX:       storageFX,
		bulkLoadSize:    bulkLoadSize,
		importCheckSkip: importCheckSkip,
	}
}

func (uc *ImportCsvBarDataUseCase) Execute(ctx context.Context, cmd fxcommand.ImportCsvBarDataCommand) (*BarDataImportResult, error) {
	// 1. ファイル名チェック
	expectedPattern := cmd.Symbol + "_" + cmd.BarType.Keyword()
	if !strings.Contains(cmd.OriginalFileName, expectedPattern) {
		return &BarDataImportResult{
			FileName:     cmd.OriginalFileName,
			FileSize:     cmd.FileSize,
			ResultStatus: "ERROR",
			ReadCount:    0,
			Message:      "file not exists.",
		}, nil
	}

	// 2. ファイル保存（保存後に再読み込みするためパスを返す）
	savedPath, err := uc.saveFile(cmd.FileReader, cmd.OriginalFileName, cmd.UserSub)
	if err != nil {
		return nil, err
	}

	// 3. ロードテーブル初期化
	if err := uc.repo.DeleteLoad(ctx, cmd.Symbol); err != nil {
		return nil, err
	}
	if err := uc.repo.DeleteLoadSma(ctx, cmd.Symbol); err != nil {
		return nil, err
	}
	if err := uc.repo.DeleteLoadRsi(ctx, cmd.Symbol); err != nil {
		return nil, err
	}

	// 4. CSVを読み込みロードテーブルへバルクインサート
	readCount, err := uc.loadCsv(ctx, savedPath, cmd.Symbol, cmd.BarType)
	if err != nil {
		return nil, err
	}

	// skipLatest=true の場合、最新1件を削除
	if cmd.SkipLatest && readCount > 0 {
		if err := uc.repo.DeleteLatestLoad(ctx, cmd.Symbol); err != nil {
			return nil, err
		}
		readCount--
	}

	// ロードテーブルの最新足日時を取得
	latestBarDateTime, err := uc.repo.GetLatestLoadBarDateTime(ctx, cmd.Symbol)
	if err != nil {
		return nil, err
	}

	// 5. 既存データとの整合性チェック
	checkResult, err := uc.repo.ImportCheck(ctx, cmd.BarType, cmd.Symbol)
	if err != nil {
		return nil, err
	}

	if checkResult.ExistsCount == 0 && !uc.importCheckSkip {
		slog.Error("インポートチェックエラー: 既存レコードが0件", "symbol", cmd.Symbol)
		return &BarDataImportResult{
			Symbol:       cmd.Symbol,
			BarDateTime:  latestBarDateTime,
			FileName:     cmd.OriginalFileName,
			FileSize:     cmd.FileSize,
			ResultStatus: "ERROR",
			ReadCount:    readCount,
			Message:      "import check error.",
		}, nil
	}

	if checkResult.DiffCount > 0 {
		slog.Warn("差分データ検出", "symbol", cmd.Symbol, "diffCount", checkResult.DiffCount)
	}

	// 6. ロードテーブルから本テーブルへインサート
	insertCount, err := uc.repo.InsertFromLoad(ctx, cmd.Symbol, cmd.BarType)
	if err != nil {
		return nil, err
	}
	insertSmaCount, err := uc.repo.InsertFromLoadSma(ctx, cmd.Symbol, cmd.BarType)
	if err != nil {
		return nil, err
	}
	insertRsiCount, err := uc.repo.InsertFromLoadRsi(ctx, cmd.Symbol, cmd.BarType)
	if err != nil {
		return nil, err
	}
	slog.Info("インサート完了", "symbol", cmd.Symbol, "bar", insertCount, "sma", insertSmaCount, "rsi", insertRsiCount)

	// 7. 差分更新
	differenceCount, err := uc.processDiffUpdate(ctx, cmd.Symbol, cmd.BarType)
	if err != nil {
		return nil, err
	}

	resultStatus := "SKIP"
	if insertCount > 0 {
		resultStatus = "OK"
	}

	return &BarDataImportResult{
		Symbol:          cmd.Symbol,
		BarDateTime:     latestBarDateTime,
		FileName:        cmd.OriginalFileName,
		FileSize:        cmd.FileSize,
		ResultStatus:    resultStatus,
		ReadCount:       readCount,
		ExistsCount:     checkResult.ExistsCount,
		InsertCount:     insertCount,
		DifferenceCount: differenceCount,
	}, nil
}

func (uc *ImportCsvBarDataUseCase) saveFile(r io.Reader, filename, userSub string) (string, error) {
	dir := filepath.Join(uc.storageBucket, uc.storageFX, "BarDataService", userSub)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("ディレクトリ作成エラー: %w", err)
	}
	path := filepath.Join(dir, filename)
	f, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("ファイル作成エラー: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("ファイルクローズエラー", "error", err)
		}
	}()
	if _, err := io.Copy(f, r); err != nil {
		return "", fmt.Errorf("ファイル書き込みエラー: %w", err)
	}
	return path, nil
}

func (uc *ImportCsvBarDataUseCase) loadCsv(ctx context.Context, path, symbol string, barType fxmodel.BarType) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("ファイルオープンエラー: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("ファイルクローズエラー", "error", err)
		}
	}()

	reader := csv.NewReader(f)

	// ヘッダー行をスキップ
	if _, err := reader.Read(); err != nil {
		return 0, fmt.Errorf("ヘッダー読み込みエラー: %w", err)
	}

	barBuffer := make([]fxmodel.BarLoadData, 0, uc.bulkLoadSize)
	smaBuffer := make([]fxmodel.BarLoadSma, 0, uc.bulkLoadSize*3)
	rsiBuffer := make([]fxmodel.BarLoadRsi, 0, uc.bulkLoadSize)
	count := 0

	flush := func() error {
		if err := uc.repo.BulkLoad(ctx, barBuffer); err != nil {
			return err
		}
		if err := uc.repo.BulkLoadSma(ctx, smaBuffer); err != nil {
			return err
		}
		if err := uc.repo.BulkLoadRsi(ctx, rsiBuffer); err != nil {
			return err
		}
		barBuffer = barBuffer[:0]
		smaBuffer = smaBuffer[:0]
		rsiBuffer = rsiBuffer[:0]
		return nil
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, fmt.Errorf("CSV読み込みエラー: %w", err)
		}

		// col: barDateTime(0) openPrice(1) highPrice(2) lowPrice(3) closePrice(4) volume(5) sma200(6) sma75(7) sma20(8) rsi(9) rsiMa(10)
		barDateTime, err := barType.ParseBarDateTime(row[0])
		if err != nil {
			return 0, fmt.Errorf("barDateTime解析エラー '%s': %w", row[0], err)
		}
		openPrice, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			return 0, fmt.Errorf("openPrice解析エラー: %w", err)
		}
		highPrice, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return 0, fmt.Errorf("highPrice解析エラー: %w", err)
		}
		lowPrice, err := strconv.ParseFloat(row[3], 64)
		if err != nil {
			return 0, fmt.Errorf("lowPrice解析エラー: %w", err)
		}
		closePrice, err := strconv.ParseFloat(row[4], 64)
		if err != nil {
			return 0, fmt.Errorf("closePrice解析エラー: %w", err)
		}

		volume := 0
		if row[5] != "" {
			if v, err := strconv.Atoi(row[5]); err == nil {
				volume = v
			}
		}

		sma200 := parseOptionalFloat64(row[6])
		sma75 := parseOptionalFloat64(row[7])
		sma20 := parseOptionalFloat64(row[8])
		rsiVal := parseOptionalFloat64(row[9])
		rsiMa := parseOptionalFloat64(row[10])

		barBuffer = append(barBuffer, fxmodel.BarLoadData{
			Symbol:      symbol,
			BarDateTime: barDateTime,
			OpenPrice:   openPrice,
			HighPrice:   highPrice,
			LowPrice:    lowPrice,
			ClosePrice:  closePrice,
			Volume:      volume,
		})
		smaBuffer = append(smaBuffer,
			fxmodel.NewBarLoadSma(symbol, barDateTime, 200, sma200, highPrice, lowPrice),
			fxmodel.NewBarLoadSma(symbol, barDateTime, 75, sma75, highPrice, lowPrice),
			fxmodel.NewBarLoadSma(symbol, barDateTime, 20, sma20, highPrice, lowPrice),
		)
		rsiBuffer = append(rsiBuffer, fxmodel.BarLoadRsi{
			Symbol:      symbol,
			BarDateTime: barDateTime,
			RsiRange:    csvRsiRange,
			RsiValue:    rsiVal,
			RsiMa:       rsiMa,
		})
		count++

		if len(barBuffer) >= uc.bulkLoadSize {
			if err := flush(); err != nil {
				return 0, err
			}
		}
	}

	if len(barBuffer) > 0 {
		if err := flush(); err != nil {
			return 0, err
		}
	}

	return count, nil
}

func (uc *ImportCsvBarDataUseCase) processDiffUpdate(ctx context.Context, symbol string, barType fxmodel.BarType) (int, error) {
	diffData, err := uc.repo.GetDiffBarData(ctx, symbol, barType)
	if err != nil {
		return 0, err
	}
	if len(diffData) > 0 {
		for _, d := range diffData {
			slog.Warn("BarData差分", "symbol", d.Symbol, "barDateTime", d.BarDateTime,
				"open", d.OpenPrice, "close", d.ClosePrice)
		}
		if _, err := uc.repo.UpdateBarData(ctx, symbol, barType); err != nil {
			return 0, err
		}
	}

	diffSma, err := uc.repo.GetDiffBarSma(ctx, symbol, barType)
	if err != nil {
		return 0, err
	}
	if len(diffSma) > 0 {
		slog.Warn("BarSma差分", "symbol", symbol, "diffCount", len(diffSma))
		if _, err := uc.repo.UpdateBarSma(ctx, symbol, barType); err != nil {
			return 0, err
		}
	}

	diffRsi, err := uc.repo.GetDiffBarRsi(ctx, symbol, barType)
	if err != nil {
		return 0, err
	}
	if len(diffRsi) > 0 {
		slog.Warn("BarRsi差分", "symbol", symbol, "diffCount", len(diffRsi))
		if _, err := uc.repo.UpdateBarRsi(ctx, symbol, barType); err != nil {
			return 0, err
		}
	}

	return len(diffData), nil
}

func parseOptionalFloat64(s string) *float64 {
	if s == "" {
		return nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return &v
}
