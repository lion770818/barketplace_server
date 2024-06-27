package application_layer

import (
	Infrastructure_bill "marketplace_server/internal/bill/Infrastructure_layer"
	"marketplace_server/internal/bill/model"
)

type TransactionAppInterface interface {
	CreateTransaction(transaction *model.Transaction) error
}

var _ TransactionAppInterface = &TransactionApp{}

type TransactionApp struct {
	TransactionRepo Infrastructure_bill.TransactionRepo
}

func NewTransactionApp(transactionRepo Infrastructure_bill.TransactionRepo) *TransactionApp {
	return &TransactionApp{
		TransactionRepo: transactionRepo,
	}
}

func (a *TransactionApp) CreateTransaction(transaction *model.Transaction) error {
	return a.TransactionRepo.Save(transaction)
}
