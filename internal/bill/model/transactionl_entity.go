package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction_Status int8

const (
	Transaction_Status_Wait   Transaction_Status = iota // 0:未完成
	Transaction_Status_Finish                           // 1:已完成
	Transaction_Status_Cancel                           // 2:取消
	Transaction_Status_Error                            // 3:錯誤
)

// 交易清單
type Transaction struct {
	ID                int64           // 流水編號
	TransferMode      int             // 交易模式 0:買 1:賣
	TransferType      int             // 交易種類 0:限價 1:市價
	TransactionID     string          // 交易單號
	FromUserID        int64           // 發起人的用戶ID
	ToUserID          int64           // 交易對象的用戶ID
	ProductName       string          // 產品名稱
	ProductCount      int64           // 產品數量
	ProductNeedAmount decimal.Decimal // 商品需要的預扣金額
	Amount            decimal.Decimal // 成交實際金額
	Currency          string          // 貨幣
	CreatedAt         time.Time       // 創建時間
	UodateAt          time.Time       // 更新時間
	Status            int8            // 交易狀態 0:未完成 1:已完成 2:取消 3:錯誤
}

func (b *Transaction) ToPO() *Transaction_PO {
	return &Transaction_PO{
		ID:                b.ID,
		TransferMode:      b.TransferMode,
		TransferType:      b.TransferType,
		TransactionID:     b.TransactionID,
		FromUserID:        b.FromUserID,
		ToUserID:          b.ToUserID,
		ProductName:       b.ProductName,
		ProductCount:      b.ProductCount,
		ProductNeedAmount: b.ProductNeedAmount,
		Amount:            b.Amount,
		Currency:          b.Currency,
		CreatedAt:         b.CreatedAt,
		UodateAt:          b.UodateAt,
		Status:            b.Status,
	}
}
