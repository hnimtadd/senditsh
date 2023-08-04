package data

type Session struct {
	Id       string `json:"sessionId,omiempty" db:"sessionId,omiempty"`
	ClientId string `json:"clientId,omiempty" db:"clientId,omiempty"`
}
