package fxdto

type FileImportResult struct {
	FileName     string `json:"fileName"`
	FileSize     int64  `json:"fileSize"`
	ReadCount    int    `json:"readCount"`
	ResultStatus string `json:"resultStatus"`
}
