package model

import (
	"errors"
	"time"
)

var (
	Error_UserIDIsEmpty     = errors.New("user_id is empty")
	Error_BackpackIDIsEmpty = errors.New("backpack_id is empty")
	Error_ConvertFailed     = errors.New("convert failed")
)

type Backpack_PO struct {
	BackpackID   int64     `gorm:"primary_key;auto_increment;comment:'流水號 背包ID 主鍵'" json:"backpack_id"`
	UserID       int64     `gorm:"column:user_id; comment:'用戶ID'" `
	ProductName  string    `gorm:"size:256;not null; comment:'產品名稱'" json:"product_name"`
	ProductCount int64     `gorm:"type:bigint(20);comment:'產品數量'" json:"product_count"`
	CreatedAt    time.Time `gorm:"autoCreateTime;comment:'創建時間'" json:"created_at"`
	UodateAt     time.Time `gorm:"autoUpdateTime;comment:'更新時間'" json:"update_at"`
}

func (Backpack_PO) TableName() string {
	return "backpack"
}

func (b *Backpack_PO) ToDomain() (*Backpack, error) {

	if b.UserID == 0 {
		return nil, Error_UserIDIsEmpty
	}
	if b.BackpackID <= 0 {
		return nil, Error_BackpackIDIsEmpty
	}

	backpack := &Backpack{
		BackpackID:   b.BackpackID,
		UserID:       b.UserID,
		ProductName:  b.ProductName,
		ProductCount: b.ProductCount,
		CreatedAt:    b.CreatedAt,
		UodateAt:     b.UodateAt,
	}

	return backpack, nil
}
