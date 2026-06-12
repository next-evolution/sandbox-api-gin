package fx

import "sandbox-api-gin/internal/domain/apperror"

type BarType string

const (
	BarTypeM15 BarType = "15M"
	BarTypeH1  BarType = "1H"
	BarTypeH4  BarType = "4H"
	BarTypeD1  BarType = "1D"
)

func (b BarType) Suffix() string {
	switch b {
	case BarTypeM15:
		return "15m"
	case BarTypeH1:
		return "1h"
	case BarTypeH4:
		return "4h"
	case BarTypeD1:
		return "1d"
	default:
		return ""
	}
}

func BarTypeOf(code string) (BarType, error) {
	switch code {
	case "15M":
		return BarTypeM15, nil
	case "1H":
		return BarTypeH1, nil
	case "4H":
		return BarTypeH4, nil
	case "1D":
		return BarTypeD1, nil
	default:
		return "", apperror.NewDomainValidationError("Unknown BarType: " + code)
	}
}
