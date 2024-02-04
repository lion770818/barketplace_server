package user

import (
	"errors"
	"fmt"
	"marketplace_server/internal/bill"
	bill_model "marketplace_server/internal/bill/model"
	"marketplace_server/internal/common/logs"
	"marketplace_server/internal/product/Infrastructure_layer"
	"marketplace_server/internal/user/model"
	"time"

	"github.com/shopspring/decimal"
)

var (
	Error_UserAlreadyExists = errors.New("用户已存在")
	Error_VerifyFailed      = errors.New("验证失败")
)

// [應用層]
type UserAppInterface interface {
	Login(login *model.LoginParams) (*model.S2C_Login, error)
	GetAuthInfo(token string) (*model.AuthInfo, error)
	GetUserInfo(userID int64) (*model.S2C_UserInfo, error)
	Register(register *model.RegisterParams) (*model.S2C_Login, error)

	Transfer(fromUserID, toUserID int64, amount decimal.Decimal, currencyStr string) error
}

type UserApp struct {
	userRepo        UserRepo
	authRepo        AuthInterface
	transferService TransferService
	rateService     RateService
	billApp         bill.BillAppInterface
}

func NewUserApp(userRepo UserRepo, authRepo AuthInterface, billRepo bill.BillRepo) UserAppInterface {
	return &UserApp{
		userRepo:        userRepo,
		authRepo:        authRepo,
		transferService: NewTransferService(),
		rateService:     NewRateService(),
		billApp:         bill.NewBillApp(billRepo),
	}
}

// Login
func (u *UserApp) Login(login *model.LoginParams) (*model.S2C_Login, error) {
	// 登录
	user, err := u.userRepo.GetUserByLoginParams(login)
	if err != nil {
		return nil, err
	}

	// 生成 token
	authInfo := &model.AuthInfo{
		UserID: user.ID,
	}
	token, err := u.authRepo.Set(authInfo)
	if err != nil {
		return nil, err
	}

	return user.ToLoginResp(token), nil
}

// GetAuthInfo 从 token 中获取用户信息
func (u *UserApp) GetAuthInfo(token string) (*model.AuthInfo, error) {
	return u.authRepo.Get(token)
}

// GetUserInfo 获取用户信息
func (u *UserApp) GetUserInfo(userID int64) (*model.S2C_UserInfo, error) {
	user, err := u.userRepo.GetUserInfo(userID)
	if err != nil {
		return nil, err
	}

	// 領域層物件轉換
	return user.ToUserInfo(), nil
}

// Register 注册 + 自动登录
func (u *UserApp) Register(register *model.RegisterParams) (*model.S2C_Login, error) {
	// 检查是否已经注册
	getUser, err := u.userRepo.GetUserByRegisterParams(register)
	if getUser != nil || err == nil {
		return nil, Error_UserAlreadyExists
	}

	// 转换参数
	params, err := register.ToDomain()
	if err != nil {
		return nil, Error_UserAlreadyExists
	}

	// 注册
	user, err := u.userRepo.Save(params)
	if err != nil {
		return nil, err
	}

	// 生成 token
	authInfo := &model.AuthInfo{
		UserID: user.ID,
	}
	token, err := u.authRepo.Set(authInfo)
	if err != nil {
		return nil, err
	}

	return user.ToLoginResp(token), nil
}

func (u *UserApp) Transfer(fromUserID, toUserID int64, amount decimal.Decimal, toCurrency string) error {
	// 讀取db用戶數據 (來源)
	fromUser, err := u.userRepo.GetUserInfo(fromUserID)
	if err != nil {
		return err
	}

	// 讀取db用戶數據 (目的)
	toUser, err := u.userRepo.GetUserInfo(toUserID)
	if err != nil {
		return err
	}

	logs.Debugf("111:%s", Infrastructure_layer.Redis_MarketPrice)
	//Infrastructure_layer.

	// 讀取匯率
	rate, err := u.rateService.GetRate(fromUser.Currency, toCurrency)
	if err != nil {
		return err
	}

	// 轉帳
	err = u.transferService.Transfer(fromUser, toUser, amount, rate)
	if err != nil {
		return err
	}

	// 保存轉帳號金幣回DB
	u.userRepo.Save(fromUser)
	u.userRepo.Save(toUser)

	// 建立交易單
	bill := &bill_model.Transaction{
		TransactionID: fmt.Sprintf("%d-%d-%s-%d", fromUser.ID, toUser.ID, toCurrency, time.Now().UnixNano()), // 交易單號
		FromUserID:    fromUser.ID,
		ToUserID:      toUser.ID,
		Amount:        amount,
		Currency:      toCurrency,
	}
	err = u.billApp.CreateBill(bill)
	if err != nil {
		return err
	}

	return nil
}
