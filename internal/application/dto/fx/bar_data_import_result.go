package fxdto

type BarDataImportResult struct {
	Symbol          string  `json:"symbol"`
	BarDateTime     *string `json:"barDateTime"`
	FileName        string  `json:"fileName,omitempty"`
	FileSize        int64   `json:"fileSize,omitempty"`
	ResultStatus    string  `json:"resultStatus,omitempty"`
	ReadCount       int     `json:"readCount,omitempty"`
	ExistsCount     int     `json:"existsCount"`
	InsertCount     int     `json:"insertCount,omitempty"`
	DifferenceCount int     `json:"differenceCount,omitempty"`
	Message         *string `json:"message"`
}
