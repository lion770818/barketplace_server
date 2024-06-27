package domain_layer

import (
	"marketplace_server/internal/user/model"

	"github.com/shopspring/decimal"
)

type TransferService interface {
	Transfer(fromUser *model.User, toUser *model.User, amount decimal.Decimal, rate *model.Rate) error
}

var _ TransferService = &TransferServiceImpl{}

type TransferServiceImpl struct {
}

func NewTransferService() *TransferServiceImpl {
	return &TransferServiceImpl{}
}

func (t *TransferServiceImpl) Transfer(fromUser *model.User, toUser *model.User, amount decimal.Decimal, rate *model.Rate) error {
	var err error

	// 通过汇率转换金额
	fromAmount := rate.Exchange(amount)

	// 根据用户不同的 vip 等级, 计算手续费
	fee := fromUser.CalcFee(fromAmount)

	fromTotalAmount := fromAmount.Add(fee)

	// 转账
	err = fromUser.Pay(fromTotalAmount)
	if err != nil {
		return err
	}
	err = toUser.Receive(amount)
	if err != nil {
		return err
	}

	return nil
}
