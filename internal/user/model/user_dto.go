package model

import "github.com/shopspring/decimal"

// dto (data transfer object) 数据传输对象
// [Demain 層]

// C2S_Login Web登录请求
type C2S_Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *C2S_Login) ToDomain() (*LoginParams, error) {

	// 驗證用戶參數
	if err := c.Verify(); err != nil {
		return nil, err
	}

	// 將用戶參數轉換為領域對象
	return &LoginParams{
		Username: c.Username,
		Password: c.Password,
	}, nil
}

// 驗證用戶
func (c *C2S_Login) Verify() error {
	if c.Username == "" || c.Password == "" {
		return Error_VerifyFailed
	}

	return nil
}

// S2C_Login Web登录响应
type S2C_Login struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

// 獲得用戶資訊
type S2C_UserInfo struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

// 註冊用戶
type C2S_Register struct {
	Username string          `json:"username"`
	Password string          `json:"password"`
	Currency string          `json:"currency"`
	Amount   decimal.Decimal `json:"amount"`
}

func (c *C2S_Register) ToDomain() (*RegisterParams, error) {

	// 驗證用戶參數
	if err := c.Verify(); err != nil {
		return nil, err
	}

	return &RegisterParams{
		Username: c.Username,
		Password: c.Password,
		Currency: c.Currency,
		Amount:   c.Amount,
	}, nil
}

// 驗證用戶
func (c *C2S_Register) Verify() error {

	if c.Username == "" || c.Password == "" || c.Currency == "" || c.Amount.IsNegative() {
		return Error_VerifyFailed
	}

	return nil
}

type C2S_Transfer struct {
	ToUserID int64           `json:"to_user_id"`
	Amount   decimal.Decimal `json:"amount"`
	Currency string          `json:"currency"`
}

// 驗證用戶
func (c *C2S_Transfer) Verify() error {

	if c.ToUserID <= 0 || c.Currency == "" || c.Amount.IsNegative() {
		return Error_VerifyFailed
	}

	return nil
}
