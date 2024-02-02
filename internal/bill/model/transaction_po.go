package model

import "github.com/shopspring/decimal"

type Transaction_PO struct {
	ID            int64           `gorm:"primary_key;auto_increment;comment:'流水號 主鍵'" json:"id"`
	TransactionID string          `gorm:"unique;not null; uniqueIndex; comment:'交易訂單'" json:"transaction_id"`
	ToUserID      int64           `gorm:"column:to_user_id; comment:'來源用戶ID'" `
	FromUserID    int64           `gorm:"column:from_user_id; comment:'目的用戶ID'" `
	Amount        decimal.Decimal `gorm:"type:decimal(20,2); comment:'金額'" json:"amount"`
	Currency      string          `gorm:"size:32;not null; comment:'幣種'" json:"currency"`
}

func (Transaction_PO) TableName() string {
	return "transaction"
}
