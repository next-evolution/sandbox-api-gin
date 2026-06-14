package user

import (
	"context"

	admincommand "sandbox-api-gin/internal/application/command/admin"
	"sandbox-api-gin/internal/application/dto"
	"sandbox-api-gin/internal/domain/repository"
)

type SearchUsersUseCase struct {
	userRepo repository.UserRepository
}

func NewSearchUsersUseCase(userRepo repository.UserRepository) *SearchUsersUseCase {
	return &SearchUsersUseCase{userRepo: userRepo}
}

type SearchUsersResult struct {
	TotalCount int
	TotalPage  int
	List       []*dto.UserDto
}

func (uc *SearchUsersUseCase) Execute(ctx context.Context, cmd *admincommand.SearchUsersCommand) (*SearchUsersResult, error) {
	count, err := uc.userRepo.SearchCount(ctx, cmd.EmailAddress, cmd.Approved)
	if err != nil {
		return nil, err
	}

	var list []*dto.UserDto
	if count > 0 {
		users, err := uc.userRepo.Search(ctx, cmd.EmailAddress, cmd.Approved, cmd.Page, cmd.Size)
		if err != nil {
			return nil, err
		}
		list = make([]*dto.UserDto, len(users))
		for i, u := range users {
			list[i] = dto.UserDtoFrom(u)
		}
	}

	totalPage := 0
	if count > 0 {
		totalPage = (count + cmd.Size - 1) / cmd.Size
	}

	return &SearchUsersResult{
		TotalCount: count,
		TotalPage:  totalPage,
		List:       list,
	}, nil
}
