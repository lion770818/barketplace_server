package rpc

import (
	pb_user "marketplace_server/internal/servers/rpc/protos/in/user"
	"marketplace_server/internal/user"
)

func WithRouter(s *RpcServer) {
	// 新建 server
	userServer := user.NewUserServer(s.Apps.UserApp)

	// 注册路由
	pb_user.RegisterUserServer(s.srv, userServer)
}
