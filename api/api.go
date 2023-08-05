package api

import log "github.com/hnimtadd/senditsh/logger"

var logger = log.GetLogger(log.Info, "API")

type ApiHandlerImpl struct {
	tunnels map[string]chan Tunnel
}

func NewAPIHandlerImpl() (*ApiHandlerImpl, error) {
	handler := &ApiHandlerImpl{
		tunnels: map[string]chan Tunnel{},
	}
	logger.Info("msg")
	return handler, nil
}
