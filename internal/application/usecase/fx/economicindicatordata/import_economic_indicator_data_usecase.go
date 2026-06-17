package economicindicatordata

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"sandbox-api-gin/internal/application/dto"
	fxmodel "sandbox-api-gin/internal/domain/model/fx"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type ImportEconomicIndicatorDataUseCase struct {
	dataRepo             fxrepository.EconomicIndicatorDataRepository
	indicatorRepo        fxrepository.EconomicIndicatorRepository
	countryRepo          fxrepository.CountryRepository
	storageBucket        string
	storageFX            string
	indicatorExcludeList []string
}

func NewImportEconomicIndicatorDataUseCase(
	dataRepo fxrepository.EconomicIndicatorDataRepository,
	indicatorRepo fxrepository.EconomicIndicatorRepository,
	countryRepo fxrepository.CountryRepository,
	storageBucket, storageFX string,
	indicatorExcludeList []string,
) *ImportEconomicIndicatorDataUseCase {
	return &ImportEconomicIndicatorDataUseCase{
		dataRepo:             dataRepo,
		indicatorRepo:        indicatorRepo,
		countryRepo:          countryRepo,
		storageBucket:        storageBucket,
		storageFX:            storageFX,
		indicatorExcludeList: indicatorExcludeList,
	}
}

type FileEntry struct {
	FileName string
	Reader   io.Reader
	FileSize int64
}

func (uc *ImportEconomicIndicatorDataUseCase) Execute(ctx context.Context, files []FileEntry, userSub string) ([]dto.FileImportResult, error) {
	countryMap, err := uc.buildCountryMap(ctx)
	if err != nil {
		return nil, err
	}
	indicatorMap, err := uc.buildIndicatorMap(ctx, countryMap)
	if err != nil {
		return nil, err
	}

	results := make([]dto.FileImportResult, 0, len(files))
	for _, entry := range files {
		result, err := uc.processFile(ctx, entry, countryMap, indicatorMap, userSub)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (uc *ImportEconomicIndicatorDataUseCase) processFile(
	ctx context.Context,
	entry FileEntry,
	countryMap map[string]string,
	indicatorMap map[string]map[string]fxmodel.EconomicIndicator,
	userSub string,
) (dto.FileImportResult, error) {
	savedPath, err := uc.saveFile(entry.FileName, entry.Reader, userSub)
	if err != nil {
		return dto.FileImportResult{}, fmt.Errorf("ファイル保存に失敗しました: %s: %w", entry.FileName, err)
	}

	dataList, err := uc.parseFile(savedPath, entry.FileName, countryMap, indicatorMap)
	if err != nil {
		return dto.FileImportResult{}, err
	}

	if err := uc.dataRepo.DeleteLoad(ctx); err != nil {
		return dto.FileImportResult{}, err
	}

	for _, data := range dataList {
		if err := uc.dataRepo.InsertLoad(ctx, data); err != nil {
			return dto.FileImportResult{}, err
		}
	}

	diffList, err := uc.dataRepo.LoadDiff(ctx)
	if err != nil {
		return dto.FileImportResult{}, err
	}

	insertFromLoad := 0
	diffCount := len(diffList)
	if diffCount == 0 {
		if err := uc.dataRepo.InsertFromLoad(ctx); err != nil {
			return dto.FileImportResult{}, err
		}
		insertFromLoad = len(dataList)
	} else {
		slog.Warn("---------- load diff ----------")
		for _, d := range diffList {
			slog.Warn("load diff", "code", d.Code, "countryCode", d.CountryCode, "publication", d.Publication)
		}
	}

	slog.Info("import complete",
		"file", entry.FileName,
		"insert", len(dataList),
		"insertFromLoad", insertFromLoad,
		"diffCount", diffCount,
	)

	return dto.FileImportResult{
		FileName:     entry.FileName,
		FileSize:     entry.FileSize,
		ReadCount:    len(dataList),
		ResultStatus: "OK",
	}, nil
}

func (uc *ImportEconomicIndicatorDataUseCase) saveFile(fileName string, r io.Reader, userSub string) (string, error) {
	uploadDir := filepath.Join(uc.storageBucket, uc.storageFX, "EconomicIndicatorDataService", userSub)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", err
	}
	savedPath := filepath.Join(uploadDir, fileName)
	f, err := os.Create(savedPath)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("ファイルクローズエラー", "error", err)
		}
	}()
	if _, err := io.Copy(f, r); err != nil {
		return "", err
	}
	return savedPath, nil
}

var (
	ptnDate   = regexp.MustCompile(`^[0-9]{1,2}/[0-9]{1,2}\(`)
	ptnTime   = regexp.MustCompile(`^[0-9]{2}:[0-9]{2}`)
	ptnPeriod = regexp.MustCompile(`^[0-9]{1,2}-[0-9]{1,2}月期`)
	ptnMonth  = regexp.MustCompile(`[0-9]{1,2}月`)
	dtfYMD    = "2006-01-02"
	dtfPub    = "2006-01-02 15:04"
)

func (uc *ImportEconomicIndicatorDataUseCase) parseFile(
	path, fileName string,
	countryMap map[string]string,
	indicatorMap map[string]map[string]fxmodel.EconomicIndicator,
) ([]fxmodel.EconomicIndicatorData, error) {
	importance := getImportance(fileName)
	year := getYear(fileName)

	lines, err := readLines(path)
	if err != nil {
		return nil, fmt.Errorf("ファイル読み込みエラー: %s: %w", fileName, err)
	}

	baseDate := ""
	errorCount := 0
	var result []fxmodel.EconomicIndicatorData

	for _, line := range lines {
		if ptnDate.MatchString(line) {
			baseDate = toDate(year, line)
		} else {
			data, err := uc.parseDataLine(baseDate, importance, line, countryMap, indicatorMap)
			if err != nil {
				slog.Error("parseDataLine error", "line", line, "error", err)
				errorCount++
			} else if data != nil {
				result = append(result, *data)
			}
		}
	}

	if errorCount > 0 {
		return nil, fmt.Errorf("parseFile errorCount=%d", errorCount)
	}
	return result, nil
}

func (uc *ImportEconomicIndicatorDataUseCase) parseDataLine(
	baseDate, importance, line string,
	countryMap map[string]string,
	indicatorMap map[string]map[string]fxmodel.EconomicIndicator,
) (*fxmodel.EconomicIndicatorData, error) {
	if uc.isSkip(line) {
		return nil, nil
	}

	if !ptnTime.MatchString(line) {
		switch {
		case strings.Contains(line, "日本"):
			line = "12:00\t" + line
		case strings.Contains(line, "中国"):
			line = "10:00\t" + line
		case strings.Contains(line, "インド"):
			line = "21:00\t" + line
		}
	}
	if !ptnTime.MatchString(line) {
		return nil, nil
	}

	elem := strings.Split(line, "\t")
	if len(elem) < 7 {
		return nil, nil
	}

	if _, ok := countryMap[elem[1]]; !ok {
		return nil, fmt.Errorf("country not found: %s", elem[1])
	}

	subTitle := getSubTitle(elem[2])
	name := normalizeIndicatorName(strings.ReplaceAll(elem[2], subTitle, ""))

	countryIndicators, ok := indicatorMap[elem[1]]
	if !ok {
		return nil, fmt.Errorf("economic-indicator not found: %s", name)
	}
	indicator, ok := countryIndicators[name]
	if !ok {
		return nil, fmt.Errorf("economic-indicator not found: %s", name)
	}

	if importance != indicator.Importance {
		slog.Warn("diff importance", "file", importance, "db", indicator.Importance, "baseDate", baseDate, "line", line)
	}

	unitOfValue := extractUnitOfValue(elem[6])
	resultVal := elem[6]
	if strings.TrimSpace(resultVal) == "" {
		resultVal = "-"
	}

	pub, err := toPublication(baseDate, elem[0])
	if err != nil {
		return nil, err
	}

	return &fxmodel.EconomicIndicatorData{
		Code:          indicator.Code,
		CountryCode:   indicator.CountryCode,
		Publication:   pub,
		SubTitle:      subTitle,
		PreviousValue: removeUnitOfValue(elem[4], unitOfValue),
		ForecastValue: removeUnitOfValue(elem[5], unitOfValue),
		ResultValue:   removeUnitOfValue(resultVal, unitOfValue),
	}, nil
}

func (uc *ImportEconomicIndicatorDataUseCase) isSkip(line string) bool {
	for _, kw := range uc.indicatorExcludeList {
		if strings.Contains(line, kw) {
			return true
		}
	}
	return false
}

func (uc *ImportEconomicIndicatorDataUseCase) buildCountryMap(ctx context.Context) (map[string]string, error) {
	countries, err := uc.countryRepo.CountryAll(ctx)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string, len(countries))
	for _, c := range countries {
		m[c.Name] = c.Code
	}
	return m, nil
}

func (uc *ImportEconomicIndicatorDataUseCase) buildIndicatorMap(ctx context.Context, countryMap map[string]string) (map[string]map[string]fxmodel.EconomicIndicator, error) {
	result := make(map[string]map[string]fxmodel.EconomicIndicator)
	for countryName, countryCode := range countryMap {
		indicators, err := uc.indicatorRepo.GetEconomicIndicatorList(ctx, countryCode)
		if err != nil {
			return nil, err
		}
		byName := make(map[string]fxmodel.EconomicIndicator, len(indicators))
		for _, ind := range indicators {
			byName[ind.Name] = ind
		}
		result[countryName] = byName
	}
	return result, nil
}

func getImportance(fileName string) string {
	parts := strings.Split(fileName, "_")
	if len(parts) < 3 {
		return ""
	}
	return strings.Split(parts[2], ".")[0]
}

func getYear(fileName string) int {
	parts := strings.Split(fileName, "_")
	if len(parts) == 0 {
		return 0
	}
	y, _ := strconv.Atoi(parts[0])
	return y
}

func toDate(year int, dateStr string) string {
	parts := strings.Split(strings.Split(dateStr, "(")[0], "/")
	if len(parts) < 2 {
		return ""
	}
	m, _ := strconv.Atoi(parts[0])
	d, _ := strconv.Atoi(parts[1])
	return time.Date(year, time.Month(m), d, 0, 0, 0, 0, time.UTC).Format(dtfYMD)
}

func toPublication(baseDate, timeStr string) (time.Time, error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) < 2 {
		return time.Time{}, fmt.Errorf("invalid time: %s", timeStr)
	}
	hour, _ := strconv.Atoi(parts[0])
	minute, _ := strconv.Atoi(parts[1])
	if hour > 23 {
		t, err := time.Parse(dtfPub, fmt.Sprintf("%s %02d:%02d", baseDate, hour-24, minute))
		if err != nil {
			return time.Time{}, err
		}
		return t.AddDate(0, 0, 1), nil
	}
	return time.Parse(dtfPub, fmt.Sprintf("%s %02d:%02d", baseDate, hour, minute))
}

func normalizeIndicatorName(name string) string {
	name = strings.ReplaceAll(name, "、", "")
	name = strings.ReplaceAll(name, "､", "")
	name = strings.ReplaceAll(name, "・", "")
	name = strings.ReplaceAll(name, "　", "")
	name = strings.ReplaceAll(name, " ", "")
	return strings.Map(func(r rune) rune {
		if r >= '！' && r <= '～' {
			return r - 0xFEE0
		}
		if r == '　' {
			return ' '
		}
		if unicode.Is(unicode.Katakana, r) {
			if hira := r - 0x60; hira >= 'ぁ' && hira <= 'ん' {
				return hira
			}
		}
		return r
	}, name)
}

func extractUnitOfValue(value string) string {
	s := strings.ReplaceAll(value, "-", "")
	s = strings.ReplaceAll(s, "+", "")
	s = strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return -1
		}
		return r
	}, s)
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "(", "")
	s = strings.ReplaceAll(s, ")", "")
	return s
}

func removeUnitOfValue(value, unit string) string {
	if unit == "" {
		return value
	}
	if value == "" {
		return ""
	}
	s := strings.ReplaceAll(value, unit, "")
	s = strings.ReplaceAll(s, "％", "")
	s = strings.ReplaceAll(s, "億円", "")
	s = strings.ReplaceAll(s, "億元", "")
	return s
}

func getSubTitle(indicatorName string) string {
	var elem string
	if ptnPeriod.MatchString(indicatorName) {
		elem = ptnPeriod.ReplaceAllString(indicatorName, "$0\t")
	} else {
		elem = ptnMonth.ReplaceAllString(indicatorName, "$0\t")
	}
	parts := strings.SplitN(elem, "\t", 2)
	if len(parts) == 1 {
		return ""
	}
	return parts[0]
}

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("ファイルクローズエラー", "error", err)
		}
	}()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "(") && len(lines) > 0 {
			lines[len(lines)-1] += line
		} else {
			lines = append(lines, line)
		}
	}
	return lines, scanner.Err()
}
