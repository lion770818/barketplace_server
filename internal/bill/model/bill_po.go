package model

import "github.com/shopspring/decimal"

type BillPO struct {
	ID         string          `gorm:"column:id"`
	FromUserID int64           `gorm:"column:from_user_id"`
	ToUserID   int64           `gorm:"column:to_user_id"`
	Amount     decimal.Decimal `gorm:"column:amount"`
	Currency   string          `gorm:"column:currency"`
}

func (BillPO) TableName() string {
	return "bill"
}
