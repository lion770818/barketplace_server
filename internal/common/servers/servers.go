package servers

import (
	"marketplace_server/internal/common/logs"
	"runtime"
)

type ServerInterface interface {
	GetVersion() string
	GetSystemInfo() string
	AsyncStart()
	Stop()
}

var _ ServerInterface = &Servers{}

type Servers struct {
	Servers []ServerInterface
}

func (s *Servers) GetVersion() string {

	for _, server := range s.Servers {
		ver := server.GetVersion()
		logs.Debugf("version:%s", ver)
	}

	return ""
}

func (s *Servers) GetSystemInfo() string {
	logs.Debugf("cpu count:%v", runtime.NumCPU())
	logs.Debugf("goroutine count:%v", runtime.NumGoroutine())
	return ""
}

func (s *Servers) AsyncStart() {
	for _, server := range s.Servers {
		server.AsyncStart()
	}
}

func (s *Servers) Stop() {
	for _, server := range s.Servers {
		server.Stop()
	}
}

func NewServers() *Servers {
	return &Servers{}
}

func (s *Servers) AddServer(server ServerInterface) {
	s.Servers = append(s.Servers, server)
}
