package api

import (
	"io"
	"time"

	"github.com/hnimtadd/senditsh/data"
)

type User struct {
	Id                   string   `json:"userId,omiemtpy" bson:"_id,omiempty"`
	Email                string   `json:"email,omiempty" bson:"email,omiempty"`
	FullName             string   `json:"fullName,omiempty" bson:"fullName,omiempty"`
	Username             string   `json:"userName,omiempty" bson:"userName,omiempty"`
	Location             string   `json:"location,omiempty" bson:"location,omiempty"`
	CreateAt             string   `json:"createdAt,omiempty" bson:"createdAt,omiempty"`
	LastLoginAt          string   `json:"lastLoginAt,omiempty" bson:"lastLoginAt,omiempty"`
	SubscriptionDuration string   `json:"duration,omiempty" bson:"duration,omiempty"`
	SubscriptionLevel    Level    `json:"level,omiempty" bson:"level,omiempty"`
	Settings             Settings `json:"settings,omiempty" bson:"settings,omiempty"`
}

type Settings struct {
	Subdomain  string `json:"subdomain,omiempty" bson:"subdomain,omiempty"`
	SSHKey     string `json:"sshKey,omiempty" bson:"sshKey,omiempty"`
	SSHHash    string `json:"sshHash,omiempty" bson:"sshHash,omiempty"`
	ModifiedAt string `json:"modifiedAt,omiempty" bson:"modifiedAt,omiempty"`
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
	Filename   string `json:"fileName,omiempty" bson:"fileName,omiempty"`
	From       string `json:"from,omiempty" bson:"from,omiempty"`
	ToEmail    string `json:"toEmail,omiempty" bson:"toEmail,omiempty"`
	Message    string `json:"message,omiempty" bson:"message,omiempty"`
	IsVerified bool   `json:"isVerified,omiempty" bson:"isVerified,omiempty"`
	CreatedAt  string `json:"createdAt,omiempty" bson:"createdAt,omiempty"`
}

func (t *Transfer) FromTransferData(td *data.Transfer) {
	t.Filename = td.Filename
	t.From = td.UserName
	t.ToEmail = td.ToEmail
	t.Message = td.Message
	t.IsVerified = td.IsVerified
	t.CreatedAt = time.Unix(td.CreatedAt, 0).Format("15:04:05 01-02-2006")
}

func (s *Settings) FromSettingData(sd *data.Settings) {
	s.SSHHash = sd.SSHHash
	s.SSHKey = sd.SSHKey
	s.Subdomain = sd.Subdomain
	s.ModifiedAt = time.Unix(sd.ModifiedAt, 0).Format("15:04:05 01-02-2006")
}

type File struct {
	FileName  string
	Extension string
	Mime      string
	Reader    io.Reader
}
