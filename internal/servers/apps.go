package servers

import (
	"marketplace_server/internal/user"
)

type Apps struct {
	UserApp user.UserAppInterface
}

func NewApps(repos *RepositoriesManager) *Apps {
	return &Apps{
		UserApp: user.NewUserApp(repos.UserRepo, repos.AuthRepo, repos.BillRepo),
	}
}
