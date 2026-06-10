package response

import "sandbox-api-gin/internal/application/dto"

// ReturnCode はJavaのReturnCode enumに対応する
type ReturnCode int

const (
	ReturnCodeOk    ReturnCode = 0
	ReturnCodeWarn  ReturnCode = 1
	ReturnCodeError ReturnCode = 2
	ReturnCodeFatal ReturnCode = 2147483647
)

type ApiResponse struct {
	ReturnCode ReturnCode `json:"returnCode"`
	Message    string     `json:"message,omitempty"`
}

type LoginResponse struct {
	ApiResponse
	User *dto.UserDto `json:"user"`
}

type UserResponse struct {
	ApiResponse
	User *dto.UserDto `json:"user"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Error   string `json:"error"`
	Message string `json:"message"`
}
