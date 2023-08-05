package repository

import "github.com/hnimtadd/senditsh/data"

func (repo *repositoryImpl) InsertSession(session Session) error {
	return nil
}

func (repo *repositoryImpl) GetSessions() ([]data.Session, error) {
	return []data.Session{}, nil
}
