package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction_PO struct {
	ID            int64           `gorm:"primary_key;auto_increment;comment:'流水號 主鍵'" json:"id"`
	TransactionID string          `gorm:"unique;not null; uniqueIndex; comment:'交易訂單'" json:"transaction_id"`
	ToUserID      int64           `gorm:"column:to_user_id; comment:'來源用戶ID'" `
	FromUserID    int64           `gorm:"column:from_user_id; comment:'目的用戶ID'" `
	ProductName   string          `gorm:"unique;not null; comment:'產品名稱'" json:"product_name"`
	ProductCount  int64           `gorm:"type:bigint(20);comment:'產品數量'" json:"product_count"`
	Amount        decimal.Decimal `gorm:"type:decimal(20,2); comment:'金額'" json:"amount"`
	Currency      string          `gorm:"size:32;not null; comment:'幣種'" json:"currency"`
	CreatedAt     time.Time       `gorm:"autoCreateTime;comment:'創建時間'" json:"created_at"`
	UodateAt      time.Time       `gorm:"autoUpdateTime;comment:'更新時間'" json:"update_at"`
	Status        int8            `gorm:"type:tinyint(1);default:0;comment:'交易狀態 0:未完成 1:已完成'" json:"status"`
}

func (Transaction_PO) TableName() string {
	return "transaction"
}
