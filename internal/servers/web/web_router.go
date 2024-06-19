package web

import (
	interface_product "marketplace_server/internal/product/interface_layer"
	"marketplace_server/internal/user"
)

func WithRouter(s *WebServer) {
	// 新建 handler
	userHandler := user.NewUserHandler(s.Apps.UserApp, s.Apps.ProductAPP)
	authMiddleware := user.NewAuthMiddleware(s.Apps.UserApp)
	productHandler := interface_product.NewProducHandler(s.Apps.ProductAPP)

	// 路由
	auth := s.Engin.Group("/auth")
	auth.POST("/login", userHandler.Login)       // 用戶登入 token ttl=expireTime(2hour)
	auth.POST("/register", userHandler.Register) // 用戶註冊

	// api
	api := s.Engin.Group("/v1")

	// 中间件 檢查 header 是否有夾帶 Authorization=token
	api.Use(authMiddleware.Auth)

	// 路由
	api.GET("/user_info", userHandler.UserInfo)                      // 取得用戶資料
	api.GET("/get_market_price", productHandler.GetMarketPrice)      // 取得市場價格
	api.POST("/create_product", productHandler.CreateProduct)        // 商品上架
	api.POST("/transaction_product", userHandler.TransactionProduct) // 買商品 / 賣商品
	api.POST("/cancel_product", userHandler.CancelProduct)           // 取消交易
}
