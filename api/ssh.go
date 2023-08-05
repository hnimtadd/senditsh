package api

import (
	"github.com/hnimtadd/senditsh/data"
)

func (api *ApiHandlerImpl) FinalizeAndCleanUpAfterTransfer(transfer *data.Transfer, tunnel *Tunnel) error {
	close(tunnel.DoneCh)
	delete(api.tunnels, transfer.Link)
	return nil
}
