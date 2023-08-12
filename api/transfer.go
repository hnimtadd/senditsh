package api

import (
	"github.com/hnimtadd/senditsh/data"
)

func (handler *ApiHandlerImpl) CreateTransfer(transfer *data.Transfer) error {
	if err := handler.repo.InsertTransfer(transfer); err != nil {
		return err
	}
	return nil
}

func (handler *ApiHandlerImpl) GetTransfers() ([]*Transfer, error) {
	transfers, err := handler.repo.GetTransfers()
	if err != nil {
		return nil, err
	}
	res := []*Transfer{}
	for _, transfer := range transfers {
		t := FromTransferData(&transfer)
		res = append(res, t)
	}
	return res, nil

}

func (handler *ApiHandlerImpl) GetTransfersOfUser(userId string) ([]*Transfer, error) {
	transfers, err := handler.repo.GetTransfersOfUser(userId)
	if err != nil {
		return nil, err
	}
	res := []*Transfer{}
	for _, transfer := range transfers {
		t := FromTransferData(&transfer)
		res = append(res, t)
	}
	return res, nil
}
