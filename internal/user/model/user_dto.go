package model

import (
	"encoding/json"
	"marketplace_server/internal/common/logs"

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

type TransferMode int

const (
	Purchase TransferMode = iota // 0:買
	Sell                         // 1:賣
	Cancel                       // 2:取消
)

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
	TransferMode  int             `json:"transaction_mode"` // 交易模式 0:買 1:賣
	TransferType  int             `json:"transaction_type"` // 交易種類 0:限價 1:市價
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

// C2S_TransactionProduct 買商品 賣商品
type C2S_TransactionProduct struct {
	TransferMode int             `json:"transaction_mode"` // 交易模式 0:買 1:賣
	TransferType int             `json:"transaction_type"` // 交易種類 0:限價 1:市價
	ProductName  string          `json:"product_name"`     // 商品名稱
	UserID       int64           `json:"user_id"`          // 發起交易人
	Currency     string          `json:"currency"`         // 幣種
	Amount       decimal.Decimal `json:"amount"`           // 購買價格 LimitPrice 時會參考
	OperateCount int64           `json:"operate_count"`    // 操作數量 ( 買 / 賣)
}

func (c *C2S_TransactionProduct) ToDomain() (*ProductTransactionParams, error) {

	// 驗證用戶參數
	if err := c.Verify(); err != nil {
		return nil, err
	}

	// 將用戶參數轉換為領域對象
	return &ProductTransactionParams{
		TransferMode: c.TransferMode,
		TransferType: c.TransferType,
		ProductName:  c.ProductName,
		UserID:       c.UserID,
		Currency:     c.Currency,
		Amount:       c.Amount,
		OperateCount: c.OperateCount,
	}, nil
}

// 驗證商品
func (c *C2S_TransactionProduct) Verify() error {
	if len(c.ProductName) == 0 || len(c.Currency) == 0 {
		return Error_VerifyFailed
	}
	// 判斷金額是否 <= 0
	if !c.Amount.GreaterThan(decimal.Zero) {
		return Error_VerifyFailed
	}
	// 判斷購買數量 <= 0
	if c.OperateCount <= 0 {
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

// 購買/販賣 單
type ProductTransactionParams struct {
	TransferMode  int             `json:"transaction_mode"` // 交易模式 0:買 1:賣
	TransferType  int             `json:"transaction_type"` // 交易種類 0:限價 1:市價
	TransactionID string          `json:"transaction_id"`   // 交易單號
	ProductName   string          `json:"product_name"`     // 購買的商品名稱
	UserID        int64           `json:"user_id"`          // 購買人
	Currency      string          `json:"currency"`         // 幣種
	Amount        decimal.Decimal `json:"amount"`           // 用戶想購買價格 LimitPrice 時會參考
	OperateCount  int64           `json:"operate_count"`    // 操作數量 (買 / 賣) (但先固定一次買賣一張, 多張的很複雜)
	TimeStamp     int64           `json:"timestamp"`        // 時間搓
}

func (c *ProductTransactionParams) GetPrice(marketPrice decimal.Decimal) (price decimal.Decimal) {

	// 根據現價或市價 取得此用戶想要的價格
	switch TransferType(c.TransferType) {
	case LimitPrice: // 限價

		// 取得買方的現價 價格
		price = c.Amount
	case MarketPrice: // 市價

		// 使用市場價格當買方價格
		price = marketPrice

	default:
		logs.Warnf("錯誤的 transferType data:%+v", c)
	}

	return
}

// func (c *ProductTransactionParams) ToDomain() (*modelProduct.Product, error) {

// 	// todo 驗證用戶參數

// 	return &modelProduct.Product{
// 		ProductName: c.ProductName,
// 		Currency:    c.Currency,
// 		Amount:  c.Amount,
// 	}, nil
// }

// 販賣單 買賣合併了
// type ProductSellParams struct {
// 	TransferType  int             `json:"transaction_type"` // 交易種類 0:限價 1:市價
// 	ProductName   string          `json:"product_name"`     // 購買的商品名稱
// 	UserID        int64           `json:"user_id"`          // 購買人
// 	Currency      string          `json:"currency"`         // 幣種
// 	Amount        decimal.Decimal `json:"amount"`           // 購買價格 LimitPrice 時會參考
// 	PurchaseCount int             `json:"purchase_count"`   // 購買數量
// 	TimeStamp     int64           `json:"timestamp"`        // 時間搓
// }

// C2S_SellProduct 賣商品
type C2S_SellProduct struct {
	TransferMode  int             `json:"transaction_mode"` // 交易模式 0:買 1:賣
	TransferType  int             `json:"transaction_type"` // 交易種類 0:限價 1:市價
	ProductName   string          `json:"product_name"`     // 商品名稱
	UserID        int64           `json:"user_id"`          // 發起交易人
	Currency      string          `json:"currency"`         // 幣種
	Amount        decimal.Decimal `json:"amount"`           // 購買價格 LimitPrice 時會參考
	PurchaseCount int64           `json:"purchase_count"`   // 購買數量
}

func (c *C2S_SellProduct) ToDomain() (*ProductTransactionParams, error) {

	// 驗證用戶參數
	if err := c.Verify(); err != nil {
		return nil, err
	}

	// 將用戶參數轉換為領域對象
	return &ProductTransactionParams{
		TransferMode: c.TransferMode,
		TransferType: c.TransferType,
		ProductName:  c.ProductName,
		UserID:       c.UserID,
		Currency:     c.Currency,
		Amount:       c.Amount,
		OperateCount: c.PurchaseCount,
	}, nil
}

// 驗證商品
func (c *C2S_SellProduct) Verify() error {
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

// C2S_TransactionProduct 買商品 賣商品
type C2S_CancelProduct struct {
	TransferMode  int    `json:"transaction_mode"` // 交易模式 0:買 1:賣
	TransactionID string `json:"transaction_id"`   // 交易清單
	UserID        int64  `json:"user_id"`          // 發起交易人

}

func (c *C2S_CancelProduct) ToDomain() (*ProductCancelParams, error) {

	// 驗證用戶參數
	if err := c.Verify(); err != nil {
		return nil, err
	}

	// 將用戶參數轉換為領域對象
	return &ProductCancelParams{
		TransactionID: c.TransactionID,
		UserID:        c.UserID,
	}, nil
}

// 驗證商品
func (c *C2S_CancelProduct) Verify() error {
	if len(c.TransactionID) == 0 {
		return Error_VerifyFailed
	}

	// 判斷購買數量 <= 0
	if c.UserID <= 0 {
		return Error_VerifyFailed
	}

	return nil
}

// 取消 購買/販賣 單
type ProductCancelParams struct {
	TransactionID string `json:"transaction_id"` // 交易清單
	UserID        int64  `json:"user_id"`        // 購買人
}

// 產生 購買/販賣 單 物件
func NewProductCancelParams(Data interface{}) (*ProductCancelParams, error) {

	// 解析封包
	byteArray, err := json.Marshal(Data)
	if err != nil {
		return nil, err
	}
	// 解析封包
	var productCancelParams ProductCancelParams
	err = json.Unmarshal(byteArray, &productCancelParams)
	if err != nil {
		return nil, err
	}

	return &productCancelParams, nil
}
