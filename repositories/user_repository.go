package repositories

import (
	"classroom_itats_api/entities"
	"classroom_itats_api/entities/input"
	"classroom_itats_api/entities/jwt_claim"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

type UserRepository interface {
	GetAllUser(ctx context.Context) ([]entities.User, error)
	GetUserByUid(ctx context.Context, uid int) (*entities.User, error)
	Login(ctx context.Context, userLogin *input.UserLogin) (*jwt_claim.Claim, error)
	StoreLoginInfo(ctx context.Context, name string, fbt string) error
	Logout(ctx context.Context, name string) error
	CheckNPMIsExist(ctx context.Context, name string) (entities.StudentCheck, error)
	CheckNPMsIsExist(ctx context.Context, name []interface{}) ([]entities.StudentCheck, error)
	GetDataMhs(ctx context.Context) ([]entities.StudentGetAll, error)
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) GetAllUser(ctx context.Context) ([]entities.User, error) {
	users := []entities.User{}

	err := u.db.WithContext(ctx).Limit(10).Find(&users).Error

	return users, err
}

func (u *userRepository) GetUserByUid(ctx context.Context, uid int) (*entities.User, error) {
	user := entities.User{}

	err := u.db.WithContext(ctx).First(&user, "uid = ?", uid).Error

	return &user, err
}

func (u *userRepository) Login(ctx context.Context, userLogin *input.UserLogin) (*jwt_claim.Claim, error) {
	user := entities.User{}

	err := u.db.WithContext(ctx).Where("name = ?", userLogin.Name).Find(&user).Error

	if err != nil || user.UID == 0 {
		return &jwt_claim.Claim{}, errors.New("akun tidak ditemukan")
	}

	if user.Pass == userLogin.Pass {
		res := map[string]interface{}{}

		err := u.db.WithContext(ctx).
			Table("users").
			Select("users.name as name, role.name as role").
			Joins("JOIN users_roles ON users.uid = users_roles.uid").
			Joins("JOIN role ON role.rid = users_roles.rid").Take(&res, "users.uid = ?", user.UID).Error

		if err != nil {
			return &jwt_claim.Claim{}, err
		}

		now := time.Now()

		year := now.AddDate(1, 0, 0)

		claim := jwt_claim.Claim{
			Name: res["name"].(string),
			Role: res["role"].(string),
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(year),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				ID:        uuid.NewString(),
			},
		}

		return &claim, nil
	}

	return &jwt_claim.Claim{}, errors.New("password yang anda masukkan salah")
}

func (u *userRepository) StoreLoginInfo(ctx context.Context, name string, fbt string) error {
	err := u.db.WithContext(ctx).Table("users").Where("name = ?", name).Update("mobile_token", fbt).Error

	return err
}

func (u *userRepository) Logout(ctx context.Context, name string) error {
	err := u.db.WithContext(ctx).Table("users").Where("name = ?", name).Update("mobile_token", nil).Error

	return err
}

func (u *userRepository) CheckNPMIsExist(ctx context.Context, name string) (entities.StudentCheck, error) {
	user := entities.StudentCheck{}

	err := u.db.WithContext(ctx).Select("mhsid, mhsnama").Table("mhs").Where("mhsid = ?", name).Find(&user).Error

	if err != nil {
		return entities.StudentCheck{}, err
	}

	if user.UserID == "" {
		return entities.StudentCheck{}, errors.New("user not found")
	}

	return user, err
}

func (u *userRepository) CheckNPMsIsExist(ctx context.Context, name []interface{}) ([]entities.StudentCheck, error) {
	user := []entities.StudentCheck{}

	err := u.db.WithContext(ctx).Select("mhsid, mhsnama").Table("mhs").Where("mhsid in ?", name).Find(&user).Error

	if err != nil {
		return nil, err
	}

	if len(user) == 0 {
		return nil, errors.New("user not found")
	}

	return user, err
}

func (u *userRepository) GetDataMhs(ctx context.Context) ([]entities.StudentGetAll, error) {
	user := []entities.StudentGetAll{}

	err := u.db.WithContext(ctx).Table("mhs").Select("mhsid", "mhsnama", "mobile", "email", "tgllahir::text", "sexid").Where("aktif not in ?", []string{"2", "3", "4", "N"}).Find(&user).Error

	if err != nil {
		return nil, err
	}

	if len(user) == 0 {
		return nil, errors.New("user not found")
	}

	return user, err
}
