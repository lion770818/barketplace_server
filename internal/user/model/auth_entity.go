package model

import (
	"encoding/json"

	"github.com/shopspring/decimal"
)

// type AuthKey struct {
// 	UserID string `json:"user_id"`
// }

type AuthInfo struct {
	UserID int64           `json:"user_id"`
	Amount decimal.Decimal `json:"amount"`
}

func (s *AuthInfo) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *AuthInfo) UnmarshalBinary(b []byte) error {
	return json.Unmarshal(b, s)
}
