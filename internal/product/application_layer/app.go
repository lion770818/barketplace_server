package application_layer

import (
	"encoding/json"
	"errors"
	"marketplace_server/internal/common/logs"
	"marketplace_server/internal/product/Infrastructure_layer"
	"marketplace_server/internal/product/model"
)

var (
	Error_VerifyFailed         = errors.New("验证失败")
	Error_ProductAlreadyExists = errors.New("商品已存在")
	Error_RedisFail            = errors.New("取得redis失敗")
)

// [Application 層]
type ProductAppInterface interface {
	CreateProduct(product *model.ProductCreateParams) error                                                   // 建立商品
	GetMarketPrice(marketPrice *model.MarketPriceParams) ([]*model.S2C_MarketPrice, map[string]string, error) // 取得市場價格
}

var _ ProductAppInterface = &ProductApp{}

type ProductApp struct {
	ProductRepo Infrastructure_layer.ProductRepo
}

func NewProductApp(productRepo Infrastructure_layer.ProductRepo) *ProductApp {
	return &ProductApp{
		ProductRepo: productRepo,
	}
}

func (a *ProductApp) CreateProduct(product *model.ProductCreateParams) error {

	// 转换参数
	params, err := product.ToDomain()
	if err != nil {
		return Error_ProductAlreadyExists
	}

	return a.ProductRepo.Save(params)
}

// 取得市場價格
func (a *ProductApp) GetMarketPrice(marketPrice *model.MarketPriceParams) ([]*model.S2C_MarketPrice, map[string]string, error) {

	// 取得商品清單
	productList, err := a.ProductRepo.GetProductList()
	if err != nil {
		return nil, nil, Error_ProductAlreadyExists
	}
	// 取得redis緩存
	dataMap, err := a.ProductRepo.RedisGetMarketPrice(Infrastructure_layer.Redis_MarketPrice)
	if err != nil {
		return nil, nil, Error_RedisFail
	}

	logs.Debugf("productList:%+v", productList)
	logs.Debugf("dataMap:%+v", dataMap)
	if len(dataMap) == 0 {

		// 無緩存 建立一個
		var marketPriceMap = make(map[string]string)
		for _, data := range productList {

			marketPriceRedis := model.MarketPriceRedis{
				ProductCount: data.ProductCount, // 商名數量
				Currency:     data.Currency,     // 幣值
				Amount:       data.BaseAmount,   // 使用初始價格 當 市場價格
			}
			marketPriceRedisStr, err := marketPriceRedis.ToJson()
			if err != nil {
				logs.Warnf("to json fail data:%+v, err:%v", data, err)
				continue
			}
			marketPriceMap[data.ProductName] = marketPriceRedisStr
		}

		// 設定redis緩存
		err = a.ProductRepo.RedisSetMarketPrice(Infrastructure_layer.Redis_MarketPrice, marketPriceMap)
		if err != nil {
			return nil, nil, Error_RedisFail
		}

		// 取得redis緩存
		dataMap, err = a.ProductRepo.RedisGetMarketPrice(Infrastructure_layer.Redis_MarketPrice)
		if err != nil {
			return nil, nil, Error_RedisFail
		}
	}

	// 網域層物件轉換
	var s2cList []*model.S2C_MarketPrice
	for _, data := range productList {

		var marketPriceRedis model.MarketPriceRedis
		err := json.Unmarshal([]byte(dataMap[data.ProductName]), &marketPriceRedis)
		if err != nil {
			logs.Warnf("unmarshal fail key=%s, value=%v err=%v",
				data.ProductName, dataMap[data.ProductName], err)
			continue
		}

		// 組合商品清單(包含目前價格)
		s2c := &model.S2C_MarketPrice{
			ProductID:    int64(data.ProductID),
			ProductName:  data.ProductName,
			ProductCount: marketPriceRedis.ProductCount,
			Currency:     data.Currency,
			BaseAmount:   data.BaseAmount,
			NowAmount:    marketPriceRedis.Amount, // 目前價格
		}
		s2cList = append(s2cList, s2c)
	}

	return s2cList, dataMap, nil
}
