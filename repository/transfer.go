package repository

import "github.com/hnimtadd/senditsh/data"

func (repo *repositoryImpl) InsertTransfer(transfer Transfer) error {
	return nil
}

func (repo *repositoryImpl) GetTransfers() ([]data.Transfer, error) {
	return []data.Transfer{}, nil
}

func (repo *repositoryImpl) GetTransfersOfUser(userId string) ([]data.Transfer, error) {
	return []data.Transfer{}, nil
}
