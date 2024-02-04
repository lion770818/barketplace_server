package web

import (
	interface_product "marketplace_server/internal/product/interface_layer"
	"marketplace_server/internal/user"
)

func WithRouter(s *WebServer) {
	// 新建 handler
	userHandler := user.NewUserHandler(s.Apps.UserApp)
	authMiddleware := user.NewAuthMiddleware(s.Apps.UserApp)
	productHandler := interface_product.NewUserHandler(s.Apps.ProductAPP)

	// 路由
	auth := s.Engin.Group("/auth")
	auth.POST("/login", userHandler.Login)
	auth.POST("/register", userHandler.Register)

	// api
	api := s.Engin.Group("/v1")

	// 中间件
	api.Use(authMiddleware.Auth)

	// 路由
	api.GET("/user_info", userHandler.UserInfo)
	api.POST("/transfer", userHandler.Transfer) // 買商品

	api.POST("/create_product", productHandler.CreateProduct)     // 商品上架
	api.POST("/purchase_product", productHandler.PurchaseProduct) // 買商品
	// api.POST("/sell_product", productHandler.SellProduct)	// 賣商品
	api.GET("/get_market_price", productHandler.GetMarketPrice) // 取得市場價格
}
