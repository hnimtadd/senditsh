package api

import "go.mongodb.org/mongo-driver/mongo"

func (api *ApiHandlerImpl) ValidateDomain(domain string) (bool, error){
	_, err :=  api.repo.GetUserByDomain(domain)
	if err != nil {
		if err == mongo.ErrNoDocuments{
			return true, nil
		}
		return false, err
	}
	return false ,nil
}
