package model

import (
	"encoding/json"
	"errors"

	"github.com/shopspring/decimal"
)

var (
	Error_AmountNotEnough = errors.New("余额不足")
	Error_VerifyFailed    = errors.New("验证失败")
)

type Product struct {
	ProductID    int64           // 產品ID
	ProductName  string          // 產品說明
	ProductCount int64           // 上架的商品數量
	BaseAmount   decimal.Decimal // 上架初始金額
	Currency     string          // 貨幣
}

func (b *Product) ToPO() *Product_PO {
	return &Product_PO{
		//ProductID:   b.ProductID,
		ProductName:  b.ProductName,
		ProductCount: b.ProductCount,
		BaseAmount:   b.BaseAmount,
		Currency:     b.Currency,
	}
}

type ProductCreateParams struct {
	ProductName  string          `json:"product_name"`  // 商品名稱
	ProductCount int64           `json:"product_count"` // 上架的商品數量
	Currency     string          `json:"currency"`      // 幣種
	BaseAmount   decimal.Decimal `json:"base_amount"`   // 基本價格
}

func (c *ProductCreateParams) ToDomain() (*Product, error) {

	// todo 驗證用戶參數

	return &Product{
		ProductName:  c.ProductName,
		ProductCount: c.ProductCount,
		Currency:     c.Currency,
		BaseAmount:   c.BaseAmount,
	}, nil
}

// 市場價格
type MarketPrice struct {
	ProductName string // 產品說明
	Currency    string // 貨幣
}

func (b *MarketPrice) ToPO() *Product_PO {
	return &Product_PO{
		ProductName: b.ProductName,
		Currency:    b.Currency,
	}
}

// 取得市場價格
type MarketPriceParams struct {
	ProductName string `json:"product_name"` // 商品名稱
	Currency    string `json:"currency"`     // 幣種
}

func (c *MarketPriceParams) ToDomain() (*Product, error) {

	// todo 驗證用戶參數

	return &Product{
		ProductName: c.ProductName,
		Currency:    c.Currency,
	}, nil
}

type MarketPriceRedis struct {
	ProductCount int64           `json:"product_count"` // 上架的商品數量
	Currency     string          `json:"currency"`      // 幣種
	Amount       decimal.Decimal `json:"amount"`        // 基本價格
}

func NewMarketPriceRedis(jsonStr string) (*MarketPriceRedis, error) {

	// byteArray, err := json.Marshal(jsonStr)
	// if err != nil {
	// 	return nil, err
	// }

	var marketPriceRedis MarketPriceRedis
	err := json.Unmarshal([]byte(jsonStr), &marketPriceRedis)
	if err != nil {
		return nil, err
	}
	return &marketPriceRedis, nil
}

func (c *MarketPriceRedis) ToJson() (string, error) {

	byteArray, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	return string(byteArray[:]), nil
}
