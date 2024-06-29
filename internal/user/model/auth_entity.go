package model

import (
	"encoding/json"

	"github.com/shopspring/decimal"
)

// type AuthKey struct {
// 	UserID string `json:"user_id"`
// }

type AuthInfo struct {
	UserID   int64           `json:"user_id"`
	Currency string          `json:"currency"`
	Amount   decimal.Decimal `json:"amount"` // 用戶帳戶的餘額, 買時 會先預扣, 等成交時再回滾調整, 避免餘額不足情況
}

func (s *AuthInfo) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *AuthInfo) UnmarshalBinary(b []byte) error {
	return json.Unmarshal(b, s)
}
