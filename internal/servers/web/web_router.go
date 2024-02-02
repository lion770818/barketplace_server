package web

import (
	"marketplace_server/internal/user"
	//"marketplace_server/internal/product"
)

func WithRouter(s *WebServer) {
	// 新建 handler
	userHandler := user.NewUserHandler(s.Apps.UserApp)
	authMiddleware := user.NewAuthMiddleware(s.Apps.UserApp)
	//product.n

	//
	auth := s.Engin.Group("/auth")
	auth.POST("/login", userHandler.Login)
	auth.POST("/register", userHandler.Register)

	// api
	api := s.Engin.Group("/v1")

	// 中间件
	api.Use(authMiddleware.Auth)

	// 路由
	api.GET("/user_info", userHandler.UserInfo)
	api.POST("/transfer", userHandler.Transfer)

	//api.POST("/transfer", userHandler.Transfer)
}
