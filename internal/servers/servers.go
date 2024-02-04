package servers

import "marketplace_server/internal/common/logs"

type ServerInterface interface {
	GetVersion() string
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
