package adminrequest

import "sandbox-api-gin/internal/api/dto/request"

type UserSearchRequest struct {
	request.ApiSearchRequest
	EmailAddress string `json:"emailAddress"`
	Approved     *bool  `json:"approved"`
}
