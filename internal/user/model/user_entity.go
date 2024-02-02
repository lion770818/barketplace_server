package model

import (
	"errors"

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
)

type User struct {
	ID       int64 // UserID
	Username string
	Password string
	Currency string
	Amount   decimal.Decimal
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
		UserID:   u.ID,
		Username: u.Username,
		Token:    token,
	}
}

func (u *User) ToUserInfo() *S2C_UserInfo {
	return &S2C_UserInfo{
		UserID:   u.ID,
		Username: u.Username,
		Amount:   u.Amount.String(),
		Currency: u.Currency,
	}
}

func (u *User) ToPO() *UserPO {

	return &UserPO{
		ID:       u.ID,
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

func (r *Rate) Exchange(amount decimal.Decimal) decimal.Decimal {
	return amount.Mul(r.rate)
}
