package servers

import (
	//"marketplace_server/internal/product"
	"marketplace_server/internal/user"
)

// [Application 層]
type Apps struct {
	UserApp user.UserAppInterface // 用戶應用層
	//ProductAPP product.ProductAppInterface // 產品應用層
}

func NewApps(repos *RepositoriesManager) *Apps {
	// 綁定應用層物件, 並回傳
	return &Apps{
		UserApp: user.NewUserApp(repos.UserRepo, repos.AuthRepo, repos.BillRepo),
		//ProductAPP: product.NewProductApp(repos.ProductRepo),
	}
}
