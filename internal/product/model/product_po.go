package model

import "github.com/shopspring/decimal"

type Product_PO struct {
	ProductID    int64           `gorm:"primary_key;auto_increment;comment:'產品ID 主鍵'" json:"product_id"`
	ProductName  string          `gorm:"unique;not null; comment:'產品名稱'" json:"product_name"`
	ProductCount int64           `gorm:"type:bigint(20);comment:'產品數量'" json:"product_count"`
	BaseAmount   decimal.Decimal `gorm:"type:decimal(20,2); comment:'上架初始金額'" json:"base_amount"`
	Currency     string          `gorm:"size:32;not null; comment:'幣種'" json:"currency"`
}

func (Product_PO) TableName() string {
	return "product"
}

// 持久層轉網域層
func (p *Product_PO) ToDomain() *Product {

	return &Product{
		ProductID:    p.ProductID,
		ProductCount: p.ProductCount,
		ProductName:  p.ProductName,
		BaseAmount:   p.BaseAmount,
		Currency:     p.Currency,
	}

}
