package api

import (
	"context"
	"fmt"
	"io"
)

type Tunnel struct {
	ReaderCh chan io.Reader
	WriterCh chan io.Writer
	DoneCh   chan struct{}
	// Indicate that writer already pipe to this tunnel
	Reader io.Reader
	Writer io.Writer
}

func (api *ApiHandlerImpl) InitTunnelWithID(ctx context.Context, id string) (*Tunnel, error) {
	_, ok := api.tunnels[id]
	if ok {
		return nil, fmt.Errorf("Tunnel with id (%v) already clared", id)
	}
	tunnel := &Tunnel{
		ReaderCh: make(chan io.Reader),
		WriterCh: make(chan io.Writer),
		DoneCh:   make(chan struct{}),
		Writer:   nil,
		Reader:   nil,
	}
	api.tunnels[id] = tunnel

	return tunnel, nil
}

func (api *ApiHandlerImpl) GetTunnelWithID(id string) (*Tunnel, error) {
	tunnel, ok := api.tunnels[id]
	if !ok {
		return nil, fmt.Errorf("Tunnel with id (%v) not exists\n", id)
	}
	return tunnel, nil
}

func (tunnel *Tunnel) PipeReader(r io.Reader) error {
	tunnel.Reader = r
	return nil
}

func (tunnel *Tunnel) PipeWriter(w io.Writer) error {
	tunnel.WriterCh <- w
	logger.Info("writer pipe")
	return nil
}

// Wait for reader and writer pipe, and send to ReadyCh

func (api *ApiHandlerImpl) WaitForWriterPipeShake(ctx context.Context, id string) error {
	tunnel, ok := api.tunnels[id]
	if !ok {
		return fmt.Errorf("Tunnel not found")
	}
	for {
		select {
		case <-ctx.Done():
			api.DestroyTunnel(id)
			return fmt.Errorf("Timeout")
		case w := <-tunnel.WriterCh:
			tunnel.Writer = w
			logger.Info("checkpoint")
			return nil
		}
	}
}

func (api *ApiHandlerImpl) WaitForCopyDone(ctx context.Context, id string) error {
	select {
	case <-api.tunnels[id].DoneCh:
		return nil
	case <-ctx.Done():
		api.DestroyTunnel(id)
		return fmt.Errorf("Timeout")
	}
}

// Copy from reader to writer
func (tunnel *Tunnel) CopyInTunnel() error {
	if tunnel.Writer == nil || tunnel.Reader == nil {
		return fmt.Errorf("Reader or writer not pipe")
	}
	_, err := io.Copy(tunnel.Writer, tunnel.Reader)
	if err != nil {
		logger.Error("msg", err)
		return err
	}
	tunnel.DoneCh <- struct{}{}
	return nil
}

func (api *ApiHandlerImpl) DestroyTunnel(id string) error {
	tunnel, ok := api.tunnels[id]
	if !ok {
		return nil
	}
	close(tunnel.WriterCh)
	close(tunnel.ReaderCh)
	close(tunnel.DoneCh)
	delete(api.tunnels, id)
	return nil
}
