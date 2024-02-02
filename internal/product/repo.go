package product

import (
	"marketplace_server/internal/product/model"

	"github.com/jinzhu/gorm"
)

type ProductRepo interface {
	Save(bill *model.Product) error
}

type MysqlProductRepo struct {
	db *gorm.DB
}

func NewMysqlProductRepo(db *gorm.DB) *MysqlProductRepo {
	return &MysqlProductRepo{db: db}
}

func (r *MysqlProductRepo) Save(bill *model.Product) error {
	billPO := bill.ToPO()
	return r.db.Save(billPO).Error
}
