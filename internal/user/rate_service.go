package user

import (
	"errors"
	"marketplace_server/internal/user/model"

	"github.com/shopspring/decimal"
)

var (
	ErrorRateNotFound = errors.New("匯率不存在")
)

const (
	USD = "USD"
	CNY = "CNY"
)

type RateService interface {
	GetRate(fromCurrency string, toCurrency string) (*model.Rate, error)
}

var _ RateService = &RateServiceImpl{}

type RateServiceImpl struct {
}

func NewRateService() *RateServiceImpl {
	return &RateServiceImpl{}
}

func (r *RateServiceImpl) GetRate(fromCurrency string, toCurrency string) (*model.Rate, error) {
	// 匯率獲取 API 可以参考: https://learn.microsoft.com/zh-cn/partner/develop/get-foreign-exchange-rates

	// 這裡 MOCK 數據替代
	// 1 USD = 6.5 CNY
	if fromCurrency == toCurrency {
		return model.NewRate(decimal.NewFromFloat(1))
	} else if fromCurrency == USD && toCurrency == CNY {
		return model.NewRate(decimal.NewFromFloat(6.5))
	} else if fromCurrency == CNY && toCurrency == USD {
		return model.NewRate(decimal.NewFromFloat(0.15))
	}

	return nil, ErrorRateNotFound
}
