package api

import (
	log "github.com/hnimtadd/senditsh/logger"
	"github.com/hnimtadd/senditsh/repository"
)

var logger = log.GetLogger(log.Info, "API")

type ApiHandlerImpl struct {
	tunnels map[string]*Tunnel
	repo    repository.Repository
}

func NewAPIHandlerImpl(repo repository.Repository) (*ApiHandlerImpl, error) {
	handler := &ApiHandlerImpl{
		tunnels: map[string]*Tunnel{},
		repo:    repo,
	}
	return handler, nil
}
