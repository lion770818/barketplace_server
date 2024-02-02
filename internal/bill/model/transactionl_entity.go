package model

import (
	"github.com/shopspring/decimal"
)

type Transaction struct {
	TransactionID string          // 交易單號
	FromUserID    int64           // 付款人
	ToUserID      int64           // 收款人
	Amount        decimal.Decimal // 金額
	Currency      string          // 貨幣
}

func (b *Transaction) ToPO() *Transaction_PO {
	return &Transaction_PO{
		TransactionID: b.TransactionID,
		FromUserID:    b.FromUserID,
		ToUserID:      b.ToUserID,
		Amount:        b.Amount,
		Currency:      b.Currency,
	}
}
