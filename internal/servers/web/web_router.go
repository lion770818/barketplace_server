package web

import (
	interface_product "marketplace_server/internal/product/interface_layer"
	"marketplace_server/internal/user"
)

func WithRouter(s *WebServer) {
	// 新建 handler
	userHandler := user.NewUserHandler(s.Apps.UserApp)
	//authMiddleware := user.NewAuthMiddleware(s.Apps.UserApp)
	productHandler := interface_product.NewUserHandler(s.Apps.ProductAPP)

	//
	auth := s.Engin.Group("/auth")
	auth.POST("/login", userHandler.Login)
	auth.POST("/register", userHandler.Register)

	// api
	api := s.Engin.Group("/v1")

	// 中间件
	//api.Use(authMiddleware.Auth)

	// 路由
	api.GET("/user_info", userHandler.UserInfo)
	api.POST("/transfer", userHandler.Transfer)

	api.POST("/create_product", productHandler.CreateProduct)
}
