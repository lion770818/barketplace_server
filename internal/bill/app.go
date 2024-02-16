package bill

import "marketplace_server/internal/bill/model"

type TransactionAppInterface interface {
	CreateTransaction(transaction *model.Transaction) error
}

var _ TransactionAppInterface = &TransactionApp{}

type TransactionApp struct {
	TransactionRepo TransactionRepo
}

func NewTransactionApp(transactionRepo TransactionRepo) *TransactionApp {
	return &TransactionApp{
		TransactionRepo: transactionRepo,
	}
}

func (a *TransactionApp) CreateTransaction(transaction *model.Transaction) error {
	return a.TransactionRepo.Save(transaction)
}
