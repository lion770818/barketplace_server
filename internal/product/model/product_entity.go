package model

import (
	"errors"

	"github.com/shopspring/decimal"
)

var (
	Error_AmountNotEnough = errors.New("余额不足")
	Error_VerifyFailed    = errors.New("验证失败")
)

type Product struct {
	//ProductID   int             // 產品ID
	ProductName string          // 產品說明
	BaseAmount  decimal.Decimal // 上架初始金額
	Currency    string          // 貨幣
}

func (b *Product) ToPO() *Product_PO {
	return &Product_PO{
		//ProductID:   b.ProductID,
		ProductName: b.ProductName,
		BaseAmount:  b.BaseAmount,
		Currency:    b.Currency,
	}
}

type ProductCreateParams struct {
	//ProductID   int             `json:"product_id"`
	ProductName string          `json:"product_name"`
	Currency    string          `json:"currency"`
	BaseAmount  decimal.Decimal `json:"base_amount"`
}

func (c *ProductCreateParams) ToDomain() (*Product, error) {

	// todo 驗證用戶參數

	return &Product{
		//ProductID:   c.ProductID,
		ProductName: c.ProductName,
		Currency:    c.Currency,
		BaseAmount:  c.BaseAmount,
	}, nil
}
