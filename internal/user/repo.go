package user

import (
	"errors"
	"fmt"
	"marketplace_server/internal/common/logs"
	"marketplace_server/internal/user/model"

	"github.com/jinzhu/gorm"
	redis "github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
)

// [Infrastructure層]
type UserRepo interface {
	GetUserInfo(userID int64) (*model.User, error)
	GetUserByLoginParams(*model.LoginParams) (*model.User, error)
	GetUserByRegisterParams(*model.RegisterParams) (*model.User, error)
	Save(*model.User) (*model.User, error)
	UpdateAmount(user *model.User, changeAmount decimal.Decimal) (*model.User, error)
}

var (
	ErrUserUsernameOrPassword = errors.New("用户名或者密码错误")
	ErrUserNotFound           = errors.New("用户不存在")
	ErrUserParamsInvalid      = errors.New("用户参数无效")
)

var _ UserRepo = &MysqlUserRepo{}

type MysqlUserRepo struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewMysqlUserRepo(db *gorm.DB, redisClient *redis.Client) *MysqlUserRepo {
	return &MysqlUserRepo{db: db, redisClient: redisClient}
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
		logs.Warnf("err:%v", err)
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

// 取得用戶資訊
func (r *MysqlUserRepo) GetUserInfo(userID int64) (*model.User, error) {
	var userPO model.UserPO
	var db = r.db

	if err := db.Where("user_id = ?", userID).First(&userPO).Error; err != nil {
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

func (r *MysqlUserRepo) UpdateAmount(user *model.User, changeAmount decimal.Decimal) (*model.User, error) {
	var userPO = user.ToPO()

	finishAmount := userPO.Amount.Add(changeAmount)
	logs.Debugf("finishAmount:%s", finishAmount.String(), userPO.Amount.String())
	if finishAmount.IsNegative() {
		return nil, fmt.Errorf("IsNegative amount:%s", finishAmount.String())
	}
	if err := r.db.Save(&userPO).Error; err != nil {
		return nil, err
	}

	// userPO.Amount = userPO.Amount.Add(changeAmount)
	// logs.Debugf("amount:%s", userPO.Amount.String())
	// err := r.db.Model(&userPO).Where("user_id = ?", userPO.UserID).Update("amount", userPO.Amount).Error
	// if err != nil {
	// 	logs.Debugf("err:%v", err)
	// }

	return userPO.ToDomain()
}
