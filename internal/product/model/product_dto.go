package model

import "github.com/shopspring/decimal"

// dto (data transfer object) 数据传输对象
// [Demain 層]

// C2S_ProductCreate 新增商品
type C2S_ProductCreate struct {
	ProductName string          `json:"product_name"`
	Currency    string          `json:"currency"`
	BaseAmount  decimal.Decimal `json:"base_amount"`
}

func (c *C2S_ProductCreate) ToDomain() (*ProductCreateParams, error) {

	// 驗證用戶參數
	if err := c.Verify(); err != nil {
		return nil, err
	}

	// 將用戶參數轉換為領域對象
	return &ProductCreateParams{
		ProductName: c.ProductName,
		Currency:    c.Currency,
		BaseAmount:  c.BaseAmount,
	}, nil
}

// 驗證商品
func (c *C2S_ProductCreate) Verify() error {
	if len(c.ProductName) == 0 || len(c.Currency) == 0 {
		return Error_VerifyFailed
	}
	// 判斷金額是否 <= 0
	if !c.BaseAmount.GreaterThan(decimal.Zero) {
		return Error_VerifyFailed
	}

	return nil
}

// S2C_ProductCreate 新增商品回應
type S2C_ProductCreate struct {
	ProductID   int64           `json:"product_id"`
	ProductName string          `json:"product_name"`
	Currency    string          `json:"currency"`
	BaseAmount  decimal.Decimal `json:"base_amount"`
}

// C2S_PurchaseProduct 新增商品
type C2S_PurchaseProduct struct {
	ProductName string          `json:"product_name"`
	Currency    string          `json:"currency"`
	BaseAmount  decimal.Decimal `json:"base_amount"`
}

func (c *C2S_PurchaseProduct) ToDomain() (*ProductCreateParams, error) {

	// 驗證用戶參數
	if err := c.Verify(); err != nil {
		return nil, err
	}

	// 將用戶參數轉換為領域對象
	return &ProductCreateParams{
		ProductName: c.ProductName,
		Currency:    c.Currency,
		BaseAmount:  c.BaseAmount,
	}, nil
}

// 驗證商品
func (c *C2S_PurchaseProduct) Verify() error {
	if len(c.ProductName) == 0 || len(c.Currency) == 0 {
		return Error_VerifyFailed
	}
	// 判斷金額是否 <= 0
	if !c.BaseAmount.GreaterThan(decimal.Zero) {
		return Error_VerifyFailed
	}

	return nil
}

// C2S_MarketPrice 取得市場價格
type C2S_MarketPrice struct {
	ProductName string `json:"product_name"` // （可選)
	Currency    string `json:"currency"`     // （可選)
}

func (c *C2S_MarketPrice) ToDomain() (*MarketPriceParams, error) {

	// 驗證參數
	if err := c.Verify(); err != nil {
		return nil, err
	}

	// 將用戶參數轉換為領域對象
	return &MarketPriceParams{
		ProductName: c.ProductName,
		Currency:    c.Currency,
	}, nil
}

// 驗證
func (c *C2S_MarketPrice) Verify() error {

	// 參數可選 不驗證
	return nil
}

// 取得市場價格 的回應
type S2C_MarketPrice struct {
	ProductID   int64           `json:"product_id"`   // 商品ID
	ProductName string          `json:"product_name"` // 商品名稱
	Currency    string          `json:"currency"`     // 幣種
	BaseAmount  decimal.Decimal `json:"base_amount"`  // 基本上市價格
	NowAmount   decimal.Decimal `json:"now_amount"`   // 目前價格
}
