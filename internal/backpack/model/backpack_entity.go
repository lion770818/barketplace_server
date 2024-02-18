package model

import (
	"time"
)

type Backpack struct {
	BackpackID   int64     // 背包ID
	UserID       int64     // 持有人
	ProductName  string    // 產品名稱
	ProductCount int64     // 產品數量
	CreatedAt    time.Time // 創建時間
	UodateAt     time.Time // 更新時間
}

func (b *Backpack) ToPO() *Backpack_PO {
	return &Backpack_PO{
		BackpackID:   b.BackpackID,
		UserID:       b.UserID,
		ProductName:  b.ProductName,
		ProductCount: b.ProductCount,
		CreatedAt:    b.CreatedAt,
		UodateAt:     b.UodateAt,
	}
}
