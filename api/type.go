package api

import "io"

type User struct {
	Id                   string   `json:"userId,omiemtpy" bson:"_id,omiempty"`
	Email                string   `json:"email,omiempty" bson:"email,omiempty"`
	FullName             string   `json:"fullName,omiempty" bson:"fullName,omiempty"`
	Username             string   `json:"userName,omiempty" bson:"userName,omiempty"`
	Location             string   `json:"location,omiempty" bson:"location,omiempty"`
	CreateAt             int64    `json:"createdAt,omiempty" bson:"createdAt,omiempty"`
	LastLoginAt          int64    `json:"lastLoginAt,omiempty" bson:"lastLoginAt,omiempty"`
	SubscriptionDuration int64    `json:"duration,omiempty" bson:"duration,omiempty"`
	SubscriptionLevel    Level    `json:"level,omiempty" bson:"level,omiempty"`
	Settings             Settings `json:"settings,omiempty" bson:"settings,omiempty"`
}

type Settings struct {
	Subdomain  string `json:"subdomain,omiempty" bson:"subdomain,omiempty"`
	SSHKey     string `json:"sshKey,omiempty" bson:"sshKey,omiempty"`
	SSHHash    uint32 `json:"sshHash,omiempty" bson:"sshHash,omiempty"`
	ModifiedAt int64  `json:"modifiedAt,omiempty" bson:"modifiedAt,omiempty"`
}
type Level int64

const (
	Free Level = 0
	Vip  Level = 1
)

func (l Level) String() string {
	switch l {
	case Free:
		return "free"
	case Vip:
		return "vip"
	}
	return "unknown"
}

type Transfer struct {
	Filename    string `json:"fileName,omiempty" bson:"fileName,omiempty"`
	From        string `json:"from,omiempty" bson:"from,omiempty"`
	ToEmail     string `json:"toEmail,omiempty" bson:"toEmail,omiempty"`
	Message     string `json:"message,omiempty" bson:"message,omiempty"`
	IsVerified  bool   `json:"isVerified,omiempty" bson:"isVerified,omiempty"`
	Initiator   string `json:"initiator,omiempty" bson:"initiator,omiempty"`
	InitiatorIP string `json:"initiatorIP,omiempty" bson:"initiatorIP,omiempty"`
	CreatedAt   int64  `json:"createdAt,omiempty" bson:"createdAt,omiempty"`
}

type File struct {
	FileName  string
	Extension string
	Mime      string
	Reader    io.Reader
}
