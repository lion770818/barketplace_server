package model

import (
	"github.com/shopspring/decimal"
)

type Product struct {
	ProductID   int             // 產品ID
	ProductName string          // 產品說明
	BaseAmount  decimal.Decimal // 上架初始金額
	Currency    string          // 貨幣
}

func (b *Product) ToPO() *Product_PO {
	return &Product_PO{
		ProductID:   b.ProductID,
		ProductName: b.ProductName,
		BaseAmount:  b.BaseAmount,
		Currency:    b.Currency,
	}
}
