package web

import (
	"context"
	"fmt"
	"marketplace_server/config"
	"marketplace_server/internal/common/logs"
	"marketplace_server/internal/common/servers"
	application_server "marketplace_server/internal/servers/application_layer"

	//application_server "marketplace_server/internal/servers/application_layer"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	Version = "市場交易-1.0.0"
)

type WebServer struct {
	httpServer *http.Server
	Engin      *gin.Engine
	Apps       *application_server.Apps
}

func (s *WebServer) GetVersion() string {
	return Version
}

func (s *WebServer) GetSystemInfo() string {
	s.GetSystemInfo()
	return ""
}

func (s *WebServer) AsyncStart() {
	logs.Debugf("[服务启动] [web] 服务地址: %s", s.httpServer.Addr)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logs.Fatalf("[服务启动] [web] 服务异常: %+v", zap.Error(err))
		}
	}()
}

func (s *WebServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	logs.Debugf("[服务关闭] [web] 关闭服务")
	if err := s.httpServer.Shutdown(ctx); err != nil {
		logs.Fatalf("[服务关闭] [web] 关闭服务异常: %s", zap.Error(err))
	}
}

func NewWebServer(cfg *config.Config, apps *application_server.Apps) servers.ServerInterface {

	logs.Debugf("創建 web server mode:%s poet:%s", cfg.Web.Mode, cfg.Web.Port)

	// 設定gin
	gin.SetMode(cfg.Web.Mode)
	e := gin.Default()
	e.Use(cors.Default())

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Web.Port),
		Handler: e,
	}

	server := &WebServer{
		httpServer: httpServer,
		Engin:      e,
		Apps:       apps,
	}

	// 注册路由
	WithRouter(server)

	return server
}
