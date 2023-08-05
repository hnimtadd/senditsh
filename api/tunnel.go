package api

import (
	"context"
	"fmt"
	"io"
)

type Tunnel struct {
	Writer io.Writer
	DoneCh chan struct{}
}

func (api *ApiHandlerImpl) InitTunnel(id string) error {
	if api.tunnels[id] != nil {
		return fmt.Errorf("Tunnel with id (%v) already clared", id)
	}
	api.tunnels[id] = make(chan Tunnel)
	return nil
}

func (api *ApiHandlerImpl) WaitToGetTunnel(ctx context.Context, id string) (*Tunnel, error) {
	select {
	case tunnel := <-api.tunnels[id]:
		return &tunnel, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("Timeout")
	}
}

func (api *ApiHandlerImpl) CopyToTunnel(tunnel *Tunnel, r io.Reader) error {
	_, err := io.Copy(tunnel.Writer, r)
	if err != nil {
		return err
	}
	return nil
}
