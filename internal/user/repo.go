package user

import (
	"errors"
	"marketplace_server/internal/user/model"

	"github.com/jinzhu/gorm"
)

// [Infrastructure層]
type UserRepo interface {
	GetUserInfo(userID int64) (*model.User, error)
	GetUserByLoginParams(*model.LoginParams) (*model.User, error)
	GetUserByRegisterParams(*model.RegisterParams) (*model.User, error)
	Save(*model.User) (*model.User, error)
}

var (
	ErrUserUsernameOrPassword = errors.New("用户名或者密码错误")
	ErrUserNotFound           = errors.New("用户不存在")
	ErrUserParamsInvalid      = errors.New("用户参数无效")
)

var _ UserRepo = &MysqlUserRepo{}

type MysqlUserRepo struct {
	db *gorm.DB
}

func NewMysqlUserRepo(db *gorm.DB) *MysqlUserRepo {
	return &MysqlUserRepo{db: db}
}

func (r *MysqlUserRepo) GetUserByLoginParams(params *model.LoginParams) (*model.User, error) {
	var userPO model.UserPO
	var db = r.db
	var err error

	// 參數檢查
	if len(params.Username) == 0 || len(params.Password) == 0 {
		return nil, ErrUserParamsInvalid
	}

	// 數據庫 查找用戶
	err = db.Where("username = ? AND password = ?", params.Username, params.Password).First(&userPO).Error
	if err != nil {
		return nil, ErrUserUsernameOrPassword
	}

	return userPO.ToDomain()
}

func (r *MysqlUserRepo) GetUserByRegisterParams(params *model.RegisterParams) (*model.User, error) {
	var userPO model.UserPO
	var db = r.db
	var err error

	if len(params.Username) == 0 || len(params.Password) == 0 {
		return nil, ErrUserParamsInvalid
	}

	err = db.Where("username = ?", params.Username).First(&userPO).Error
	if err != nil {
		return nil, ErrUserNotFound
	}

	return userPO.ToDomain()
}

func (r *MysqlUserRepo) GetUserInfo(userID int64) (*model.User, error) {
	var userPO model.UserPO
	var db = r.db

	if err := db.Where("id = ?", userID).First(&userPO).Error; err != nil {
		return nil, ErrUserNotFound
	}

	return userPO.ToDomain()
}

func (r *MysqlUserRepo) Save(user *model.User) (*model.User, error) {
	var userPO = user.ToPO()

	if err := r.db.Save(&userPO).Error; err != nil {
		return nil, err
	}

	return userPO.ToDomain()
}
