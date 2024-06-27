package Infrastructure_layer

import (
	"marketplace_server/internal/bill/model"

	"github.com/jinzhu/gorm"
)

type TransactionRepo interface {
	Save(transaction *model.Transaction) error
	GetTransactionInfo(transactionId string) (*model.Transaction, error)
	GetLastInsterId() (int64, error)
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

// 取得交易資訊
func (r *MysqlTransactionRepo) GetTransactionInfo(transactionId string) (*model.Transaction, error) {
	var transactionPO model.Transaction_PO
	var db = r.db

	if err := db.Where("transaction_id = ?", transactionId).First(&transactionPO).Error; err != nil {
		return nil, err
	}

	return transactionPO.ToDomain()
}

func (r *MysqlTransactionRepo) GetLastInsterId() (int64, error) {
	var transactionPO model.Transaction_PO
	var db = r.db

	if err := db.Last(&transactionPO).Error; err != nil {
		return 0, err
	}

	return transactionPO.ID, nil
}
