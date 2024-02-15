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
	TWD = "TWD"
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
	} else if fromCurrency == TWD && toCurrency == USD {
		return model.NewRate(decimal.NewFromFloat(0.032)) // 台幣 對 美金
	} else if fromCurrency == USD && toCurrency == TWD {
		return model.NewRate(decimal.NewFromFloat(32)) // 美金 對 台幣
	} else if fromCurrency == TWD && toCurrency == CNY {
		return model.NewRate(decimal.NewFromFloat(0.23)) // 台幣 對 人民幣
	} else if fromCurrency == CNY && toCurrency == TWD {
		return model.NewRate(decimal.NewFromFloat(4.4)) // 人民幣 對 台幣
	}
	return nil, ErrorRateNotFound
}
