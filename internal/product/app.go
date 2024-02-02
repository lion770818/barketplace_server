package product

import "marketplace_server/internal/product/model"

// [Application 層]
type ProductAppInterface interface {
	CreateProduct(product *model.Product) error
}

var _ ProductAppInterface = &ProductApp{}

type ProductApp struct {
	ProductRepo ProductRepo
}

func NewProductApp(productRepo ProductRepo) *ProductApp {
	return &ProductApp{
		ProductRepo: productRepo,
	}
}

func (a *ProductApp) CreateProduct(product *model.Product) error {
	return a.ProductRepo.Save(product)
}
