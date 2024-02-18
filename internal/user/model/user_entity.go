package model

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

// domain 领域对象

var (
	DefaultUserIDValue   = "0"
	DefaultUsernameValue = ""
	DefaultPasswordValue = ""
	DefaultCurrencyValue = "CNY"
	DefaultAmountValue   = decimal.NewFromFloat(0)
	DefaultFeeValue      = decimal.NewFromFloat(0)
)

var (
	Error_AmountNotEnough = errors.New("余额不足")
	Error_VerifyFailed    = errors.New("验证失败")
	Error_GetDAtaFailed   = errors.New("取得資料失敗")
)

type User struct {
	UserID    int64 // 用戶ID
	Username  string
	Password  string
	Currency  string
	Amount    decimal.Decimal
	CreatedAt time.Time
	UpdateAt  time.Time
}

func (u *User) CalcFee(fromAmount decimal.Decimal) decimal.Decimal {
	return fromAmount.Mul(DefaultFeeValue)
}

// 付款
func (u *User) Pay(amount decimal.Decimal) error {
	// 省略参数检查
	if u.Amount.LessThan(amount) {
		return Error_AmountNotEnough
	}
	// 付款運算
	u.Amount = u.Amount.Sub(amount)

	return nil
}

// 收款
func (u *User) Receive(amount decimal.Decimal) error {
	// 省略参数检查

	u.Amount = u.Amount.Add(amount)

	return nil
}

func (u *User) ToLoginResp(token string) *S2C_Login {
	return &S2C_Login{
		UserID:   u.UserID,
		Username: u.Username,
		Token:    token,
	}
}

func (u *User) ToUserInfo() *S2C_UserInfo {
	return &S2C_UserInfo{
		UserID:   u.UserID,
		Username: u.Username,
		Amount:   u.Amount.String(),
		Currency: u.Currency,
	}
}

func (u *User) ToPO() *UserPO {

	return &UserPO{
		UserID:   u.UserID,
		Username: u.Username,
		Password: u.Password,
		Currency: u.Currency,
		Amount:   u.Amount,
	}
}

type LoginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterParams struct {
	Username string          `json:"username"`
	Password string          `json:"password"`
	Currency string          `json:"currency"`
	Amount   decimal.Decimal `json:"amount"`
}

func (c *RegisterParams) ToDomain() (*User, error) {

	// todo 驗證用戶參數

	return &User{

		Username: c.Username,
		Password: c.Password,
		Currency: c.Currency,
		Amount:   c.Amount,
	}, nil
}

type Rate struct {
	rate decimal.Decimal
}

func NewRate(rate decimal.Decimal) (*Rate, error) {
	// 省略参数检查
	return &Rate{
		rate: rate,
	}, nil
}

// 通過匯率轉換金額
func (r *Rate) Exchange(amount decimal.Decimal) decimal.Decimal {
	return amount.Mul(r.rate)
}

// 取得匯率
func (r *Rate) Get() decimal.Decimal {
	return r.rate
}

type Notify_Cmd int

const (
	Notify_Cmd_Unknow   Notify_Cmd = iota // 未定義
	Notify_Cmd_Purchase                   // 買商品
	Notify_Cmd_Sell                       // 賣商品
	Notify_Cmd_Cancel                     // 取消商品
)

// 產品交易通知封包
type ProductTransactionNotify struct {
	Cmd  Notify_Cmd  `json:"cmd"`
	Data interface{} `json:"data"`
}
