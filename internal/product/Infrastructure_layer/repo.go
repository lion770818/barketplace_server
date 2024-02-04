package Infrastructure_layer

import (
	"context"
	"errors"
	"marketplace_server/internal/product/model"

	//"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	redis "github.com/redis/go-redis/v9"
)

const (
	Redis_MarketPrice = "product:market_price" // rediskey 商品市場價格
)

// 持久層 產品
type ProductRepo interface {
	Save(product *model.Product) error
	GetProductList() ([]*model.Product, error)

	RedisGetMarketPrice(key string) (data map[string]string, err error) // 取得市場價格
	RedisSetMarketPrice(key string, data map[string]string) (err error) // 設定市場價格
}

type ProductRepoManager struct {
	db          *gorm.DB      // 資料庫
	redisClient *redis.Client // redis
}

func NewProductRepoManager(db *gorm.DB, redisDb *redis.Client) *ProductRepoManager {
	return &ProductRepoManager{db: db, redisClient: redisDb}
}

func (p *ProductRepoManager) Save(product *model.Product) error {
	productPO := product.ToPO()
	return p.db.Save(productPO).Error
}

// 取得商品清單 db
func (p *ProductRepoManager) GetProductList() ([]*model.Product, error) {
	var productPoList []model.Product_PO

	// db 撈取 產品清單
	err := p.db.Debug().Find(&productPoList).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("product not found")
	}

	// 持久層轉換領域層
	var productList []*model.Product
	for _, data := range productPoList {
		productList = append(productList, data.ToDomain())
	}

	return productList, nil
}

// 取得商品價格 redis
func (p *ProductRepoManager) RedisGetMarketPrice(key string) (data map[string]string, err error) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	data, err = p.redisClient.HGetAll(ctx, Redis_MarketPrice).Result()
	if err != nil {
		return
	}

	return
}

// 設定商品價格 redis
func (p *ProductRepoManager) RedisSetMarketPrice(key string, data map[string]string) (err error) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ret, err := p.redisClient.HMSet(ctx, key, data).Result()
	if !ret || err != nil {
		return
	}

	return
}
