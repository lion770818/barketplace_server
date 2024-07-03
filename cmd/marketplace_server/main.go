package main

import (
	"fmt"
	"marketplace_server/config"
	"marketplace_server/internal/common/logs"
	"marketplace_server/internal/common/servers"
	"marketplace_server/internal/common/signals"
	Infrastructure_server "marketplace_server/internal/servers/Infrastructure_layer"
	application_server "marketplace_server/internal/servers/application_layer"
	"marketplace_server/internal/servers/web"
	"os"
)

func main() {

	cfg := &config.Config{}
	env_flag := os.Getenv("env_flag")
	fmt.Println("env_flag=", env_flag)
	if env_flag == "1" {
		fmt.Printf("啟動env")
		cfg = config.NewEnvConfig()
	} else {
		cfg = config.NewYmlConfig("./config.yaml")
	}

	// 初始化配置

	// 初始化日志
	logs.Init(cfg.Log)

	// 获取 servers, 比如 WebServer, RpcServer
	servers := NewServers(cfg)
	// 顯示版本
	servers.GetVersion()

	// 启动 servers
	servers.AsyncStart()

	// 优雅退出
	signals.WaitWith(servers.Stop)
}

// NewServers 通过配置文件初始化 Repo 依赖, 然后初始化 App, 最后组装为 Server
// 比如 UserRepo -> UserApp -> WebServer
func NewServers(cfg *config.Config) servers.ServerInterface {

	// 建立 db連線 和 redis連線
	repos := Infrastructure_server.NewRepositories(cfg)
	repos.Automigrate()
	// 建立 應用層 管理物件
	apps := application_server.NewApps(repos)

	servers := servers.NewServers()
	servers.AddServer(web.NewWebServer(cfg, apps))

	return servers
}
