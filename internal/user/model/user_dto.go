package model

import (
	"github.com/shopspring/decimal"
)

// dto (data transfer object) 数据传输对象
// [Demain 層]

// C2S_Login Web 登入請求
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

// S2C_Login Web 登入回應
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

type TransferType int

const (
	LimitPrice  TransferType = iota // 0:限價單
	MarketPrice                     // 1:市價單
)

const (
	TransactionExchange    = "transaction_exchange"        // 通知交换机
	BindKeyPurchaseProduct = "notify_purchase_product_key" // 通用邮件绑定key
)

// 交易下單
type C2S_Transfer struct {
	TransferType  int             `json:"transaction_mode"` // 交易種類 0:限價 1:市價
	ProductID     int64           `json:"product_id"`       // 購買的商品id
	ToUserID      int64           `json:"to_user_id"`       // 購買人
	Currency      string          `json:"currency"`         // 幣種
	Amount        decimal.Decimal `json:"amount"`           // 購買價格 LimitPrice 時會參考
	PurchaseCount int             `json:"purchase_count"`   // 購買數量
}

// 驗證用戶
func (c *C2S_Transfer) Verify() error {

	if c.ToUserID <= 0 || c.Currency == "" || c.Amount.IsNegative() {
		return Error_VerifyFailed
	}

	return nil
}

// C2S_PurchaseProduct 新增商品
type C2S_PurchaseProduct struct {
	TransferType  int             `json:"transaction_type"` // 交易種類 0:限價 1:市價
	ProductName   string          `json:"product_name"`     // 購買的商品名稱
	UserID        int64           `json:"user_id"`          // 購買人
	Currency      string          `json:"currency"`         // 幣種
	Amount        decimal.Decimal `json:"amount"`           // 購買價格 LimitPrice 時會參考
	PurchaseCount int             `json:"purchase_count"`   // 購買數量
}

func (c *C2S_PurchaseProduct) ToDomain() (*ProductPurchaseParams, error) {

	// 驗證用戶參數
	if err := c.Verify(); err != nil {
		return nil, err
	}

	// 將用戶參數轉換為領域對象
	return &ProductPurchaseParams{
		TransferType:  c.TransferType,
		ProductName:   c.ProductName,
		UserID:        c.UserID,
		Currency:      c.Currency,
		Amount:        c.Amount,
		PurchaseCount: c.PurchaseCount,
	}, nil
}

// 驗證商品
func (c *C2S_PurchaseProduct) Verify() error {
	if len(c.ProductName) == 0 || len(c.Currency) == 0 {
		return Error_VerifyFailed
	}
	// 判斷金額是否 <= 0
	if !c.Amount.GreaterThan(decimal.Zero) {
		return Error_VerifyFailed
	}
	// 判斷購買數量 <= 0
	if c.PurchaseCount <= 0 {
		return Error_VerifyFailed
	}

	// 如果交易種類不是 市價 或 現價
	switch TransferType(c.TransferType) {
	case LimitPrice:
	case MarketPrice:
	default:
		return Error_VerifyFailed
	}

	return nil
}

type ProductPurchaseParams struct {
	TransferType  int             `json:"transaction_type"` // 交易種類 0:限價 1:市價
	ProductName   string          `json:"product_name"`     // 購買的商品名稱
	UserID        int64           `json:"user_id"`          // 購買人
	Currency      string          `json:"currency"`         // 幣種
	Amount        decimal.Decimal `json:"amount"`           // 購買價格 LimitPrice 時會參考
	PurchaseCount int             `json:"purchase_count"`   // 購買數量
	TimeStamp     int64           `json:"timestamp"`        // 時間搓
}

// func (c *ProductPurchaseParams) ToDomain() (*modelProduct.Product, error) {

// 	// todo 驗證用戶參數

// 	return &modelProduct.Product{
// 		ProductName: c.ProductName,
// 		Currency:    c.Currency,
// 		Amount:  c.Amount,
// 	}, nil
// }
