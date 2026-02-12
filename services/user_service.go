package services

import (
	"classroom_itats_api/entities"
	"classroom_itats_api/entities/input"
	"classroom_itats_api/entities/jwt_claim"
	"classroom_itats_api/repositories"
	"context"
)

type userService struct {
	userRepository repositories.UserRepository
}

type UserService interface {
	GetAllUser(ctx context.Context) ([]entities.User, error)
	GetUserByUid(ctx context.Context, uid int) (*entities.User, error)
	Login(ctx context.Context, userLogin *input.UserLogin) (*jwt_claim.Claim, error)
	StoreLoginInfo(ctx context.Context, name string, fbt string) error
	Logout(ctx context.Context, name string) error
	CheckNPMIsExist(ctx context.Context, name string) (entities.StudentCheck, error)
	CheckNPMsIsExist(ctx context.Context, name []interface{}) ([]entities.StudentCheck, error)
	GetDataMhs(ctx context.Context) ([]entities.StudentGetAll, error)
}

func NewUserService(userRepository repositories.UserRepository) *userService {
	return &userService{userRepository: userRepository}
}

func (u *userService) GetAllUser(ctx context.Context) ([]entities.User, error) {
	return u.userRepository.GetAllUser(ctx)
}

func (u *userService) GetUserByUid(ctx context.Context, uid int) (*entities.User, error) {
	return u.userRepository.GetUserByUid(ctx, uid)
}

func (u *userService) Login(ctx context.Context, userLogin *input.UserLogin) (*jwt_claim.Claim, error) {
	return u.userRepository.Login(ctx, userLogin)
}

func (u *userService) StoreLoginInfo(ctx context.Context, name string, fbt string) error {
	return u.userRepository.StoreLoginInfo(ctx, name, fbt)
}

func (u *userService) Logout(ctx context.Context, name string) error {
	return u.userRepository.Logout(ctx, name)
}

func (u *userService) CheckNPMIsExist(ctx context.Context, name string) (entities.StudentCheck, error) {
	return u.userRepository.CheckNPMIsExist(ctx, name)
}

func (u *userService) CheckNPMsIsExist(ctx context.Context, name []interface{}) ([]entities.StudentCheck, error) {
	return u.userRepository.CheckNPMsIsExist(ctx, name)
}

func (u *userService) GetDataMhs(ctx context.Context) ([]entities.StudentGetAll, error) {
	return u.userRepository.GetDataMhs(ctx)
}
