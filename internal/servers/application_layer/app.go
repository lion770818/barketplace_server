package application_layer

import (
	application_product "marketplace_server/internal/product/application_layer"
	Infrastructure_server "marketplace_server/internal/servers/Infrastructure_layer"
	"marketplace_server/internal/user/application_layer"
	application_user "marketplace_server/internal/user/application_layer"
)

// [Application 層]
type Apps struct {
	UserApp    application_user.UserAppInterface       // 用戶應用層
	ProductAPP application_product.ProductAppInterface // 產品應用層
}

func NewApps(repos *Infrastructure_server.RepositoriesManager) *Apps {

	//  取得產品APP層
	productAPP := application_product.NewProductApp(repos.ProductRepo)

	// 綁定應用層物件, 並回傳
	return &Apps{
		UserApp:    application_layer.NewUserApp(repos.UserRepo, repos.AuthRepo, repos.TransactionRepo, productAPP),
		ProductAPP: productAPP,
	}
}
