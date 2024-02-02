package model

import "github.com/shopspring/decimal"

type Product_PO struct {
	ID          int64           `gorm:"primary_key;auto_increment;comment:'流水號 主鍵'" json:"id"`
	ProductID   int             `gorm:"column:product_id; uniqueIndex; comment:'產品ID'" `
	ProductName string          `gorm:"unique;not null; comment:'產品名稱'" json:"product_id"`
	BaseAmount  decimal.Decimal `gorm:"type:decimal(20,2); comment:'上架初始金額'" json:"base_amount"`
	Currency    string          `gorm:"size:32;not null; comment:'幣種'" json:"currency"`
}

func (Product_PO) TableName() string {
	return "product"
}
