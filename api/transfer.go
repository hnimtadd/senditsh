package api

import (
	"github.com/hnimtadd/senditsh/data"
)

func (handler *ApiHandlerImpl) CreateTransfer(transfer *data.Transfer) error {
	return nil
}

func (handler *ApiHandlerImpl) GetTransfers() ([]data.Transfer, error) {
	return []data.Transfer{}, nil
}
