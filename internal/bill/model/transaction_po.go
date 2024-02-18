package model

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

var (
	Error_UserIDIsEmpty        = errors.New("user_id is empty")
	Error_TransactionIDIsEmpty = errors.New("transaction_id is empty")
	Error_ConvertFailed        = errors.New("convert failed")
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
	Status        int8            `gorm:"type:tinyint(1);default:0;comment:'交易狀態 0:未完成 1:已完成 2:取消 3:錯誤'" json:"status"`
}

func (Transaction_PO) TableName() string {
	return "transaction"
}

func (t *Transaction_PO) ToDomain() (*Transaction, error) {

	if t.ToUserID == 0 || t.FromUserID == 0 {
		return nil, Error_UserIDIsEmpty
	}
	if len(t.TransactionID) == 0 {
		return nil, Error_TransactionIDIsEmpty
	}

	user := &Transaction{
		TransactionID: t.TransactionID,
		FromUserID:    t.FromUserID,
		ToUserID:      t.ToUserID,
		ProductName:   t.ProductName,
		ProductCount:  t.ProductCount,
		Amount:        t.Amount,
		Currency:      t.Currency,
		CreatedAt:     t.CreatedAt,
		UodateAt:      t.UodateAt,
		Status:        t.Status,
	}

	return user, nil
}
