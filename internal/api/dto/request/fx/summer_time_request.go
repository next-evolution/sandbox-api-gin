package fxrequest

import fxdto "sandbox-api-gin/internal/application/dto/fx"

type SummerTimeRequest struct {
	SummerTime fxdto.SummerTimeDto `json:"summerTime" binding:"required"`
}
