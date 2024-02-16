package bill

import (
	"marketplace_server/internal/bill/model"

	"github.com/jinzhu/gorm"
)

type TransactionRepo interface {
	Save(transaction *model.Transaction) error
}

type MysqlTransactionRepo struct {
	db *gorm.DB
}

func NewMysqlTransactionRepo(db *gorm.DB) *MysqlTransactionRepo {
	return &MysqlTransactionRepo{db: db}
}

func (r *MysqlTransactionRepo) Save(transaction *model.Transaction) error {
	transactionPO := transaction.ToPO()
	return r.db.Save(transactionPO).Error
}
