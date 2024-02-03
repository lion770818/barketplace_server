package application_layer

import (
	"errors"
	"marketplace_server/internal/product/Infrastructure_layer"
	"marketplace_server/internal/product/model"
)

var (
	Error_ProductAlreadyExists = errors.New("商品已存在")
	Error_VerifyFailed         = errors.New("验证失败")
)

// [Application 層]
type ProductAppInterface interface {
	CreateProduct(product *model.ProductCreateParams) error
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
