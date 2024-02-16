package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	TransactionID string          // 交易單號
	FromUserID    int64           // 付款人
	ToUserID      int64           // 收款人
	ProductName   string          // 產品名稱
	ProductCount  int64           // 產品數量
	Amount        decimal.Decimal // 金額
	Currency      string          // 貨幣
	CreatedAt     time.Time       // 創建時間
	UodateAt      time.Time       // 更新時間
	Status        int8            // 交易狀態 0:未完成 1:已完成
}

func (b *Transaction) ToPO() *Transaction_PO {
	return &Transaction_PO{
		TransactionID: b.TransactionID,
		FromUserID:    b.FromUserID,
		ToUserID:      b.ToUserID,
		ProductName:   b.ProductName,
		ProductCount:  b.ProductCount,
		Amount:        b.Amount,
		Currency:      b.Currency,
		CreatedAt:     b.CreatedAt,
		UodateAt:      b.UodateAt,
		Status:        b.Status,
	}
}
