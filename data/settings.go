package data

type Settings struct {
	Subdomain  string `json:"subdomain,omiempty" bson:"subdomain,omiempty"`
	SSHKey     string `json:"sshKey,omiempty" bson:"sshKey,omiempty"`
	SSHHash    uint32 `json:"sshHash,omiempty" bson:"sshHash,omiempty"`
	ModifiedAt int64  `json:"modifiedAt,omiempty" bson:"modifiedAt,omiempty"`
}
