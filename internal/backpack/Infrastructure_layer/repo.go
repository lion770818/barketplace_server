package Infrastructure_layer

import (
	"marketplace_server/internal/backpack/model"
	"marketplace_server/internal/common/logs"

	"github.com/jinzhu/gorm"
)

// 用戶背包
type BackpackRepo interface {
	Save(backpack *model.Backpack) error
	GetBackpackById(backpackId int64) (*model.Backpack, error)
	GetBackpackByUserId(userId int64, productName string) (*model.Backpack, error)
	FindAll(userId int64) (list []*model.Backpack, err error)
}

type MysqlBackpackRepo struct {
	db *gorm.DB
}

func NewMysqlBackpackRepo(db *gorm.DB) *MysqlBackpackRepo {
	return &MysqlBackpackRepo{db: db}
}

func (r *MysqlBackpackRepo) Save(backpack *model.Backpack) error {
	backpackPO := backpack.ToPO()
	return r.db.Save(backpackPO).Error
}

// 取得背包內持有商品資訊
func (r *MysqlBackpackRepo) GetBackpackById(backpackId int64) (*model.Backpack, error) {
	var backpackPO model.Backpack_PO
	var db = r.db

	if err := db.Where("backpack_id = ?", backpackId).First(&backpackPO).Error; err != nil {
		return nil, err
	}

	return backpackPO.ToDomain()
}

// 取得背包內持有商品資訊
func (r *MysqlBackpackRepo) GetBackpackByUserId(userId int64, productName string) (*model.Backpack, error) {
	var backpackPO model.Backpack_PO
	var db = r.db

	if err := db.Where("user_id = ? AND product_name = ?", userId, productName).First(&backpackPO).Error; err != nil {
		return nil, err
	}

	return backpackPO.ToDomain()
}

func (r *MysqlBackpackRepo) FindAll(userId int64) (list []*model.Backpack, err error) {
	var poList []model.Backpack_PO
	var db = r.db
	if err = db.Where("user_id = ? AND product_name = ?", userId).Find(&poList).Error; err != nil {
		return nil, err
	}

	// 轉成領域物件
	for _, data := range poList {
		domainObj, err := data.ToDomain()
		if err != nil {
			logs.Debugf("err:%v")
			continue
		}
		list = append(list, domainObj)
	}

	return
}
