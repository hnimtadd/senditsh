package api

import (
	"io"
	"time"

	"github.com/hnimtadd/senditsh/data"
	"github.com/hnimtadd/senditsh/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	SubscriptionLevel    string   `json:"level,omiempty" bson:"level,omiempty"`
	PublicKey            string   `json:"publicKey,omiempty"`
	Settings             Settings `json:"settings,omiempty" bson:"settings,omiempty"`
	Domain               string   `json:"domain,omiempty"`
	Fingerprint          string   `json:"fingerprint,omiempty"`
}

type Settings struct {
	Subdomain  string `json:"subdomain,omiempty" bson:"subdomain,omiempty"`
	SSHKey     string `json:"sshKey,omiempty" bson:"sshKey,omiempty"`
	SSHHash    string `json:"sshHash,omiempty" bson:"sshHash,omiempty"`
	FullKey    string `json:"fullKey,omiempty" bson:"fullKey,omiempty"`
	ModifiedAt string `json:"modifiedAt,omiempty" bson:"modifiedAt,omiempty"`
}

type Transfer struct {
	Id         string `json:"transferId,omiemtpy" bson:"_id,omiempty"`
	Filename   string `json:"fileName,omiempty" bson:"fileName,omiempty"`
	From       string `json:"from,omiempty" bson:"from,omiempty"`
	Link       string `json:"link,omiempty"`
	ToEmail    string `json:"toEmail,omiempty" bson:"toEmail,omiempty"`
	Message    string `json:"message,omiempty" bson:"message,omiempty"`
	Status     string `json:"status,omiempty"`
	IsVerified bool   `json:"isVerified,omiempty" bson:"isVerified,omiempty"`
	CreatedAt  string `json:"createdAt,omiempty" bson:"createdAt,omiempty"`
}

func FromTransferData(td *data.Transfer) *Transfer {
	t := &Transfer{
		Id:         td.Id.Hex(),
		Filename:   td.Filename,
		From:       td.UserName,
		Link:       td.Link,
		ToEmail:    td.ToEmail,
		Message:    td.Message,
		Status:     td.Status,
		IsVerified: td.IsVerified,
		CreatedAt:  time.Unix(td.CreatedAt, 0).Format("15:04:05 01-02-2006"),
	}
	return t
}

func FromSettingData(sd *data.Settings) *Settings {
	s := &Settings{
		SSHHash:    sd.SSHHash,
		SSHKey:     sd.SSHKey,
		FullKey:    sd.FullKey,
		Subdomain:  sd.Subdomain,
		ModifiedAt: time.Unix(sd.ModifiedAt, 0).Format("15:04:05 01-02-2006"),
	}
	return s
}

func ToSettingData(s *Settings) *data.Settings {
	sd := &data.Settings{
		SSHKey:     s.SSHKey,
		SSHHash:    s.SSHHash,
		Subdomain:  s.Subdomain,
		FullKey:    s.FullKey,
		ModifiedAt: time.Now().Unix(),
	}
	return sd
}

func FromUserData(ud *data.User) *User {
	fingerprint := ""
	if ud.Settings.FullKey != "" {
		var err error
		fingerprint, err = utils.GetFingerprint(ud.Settings.SSHKey)
		if err != nil {
			panic(err)
		}
	}
	u := &User{
		Id:                ud.Id.String(),
		Username:          ud.Username,
		Email:             ud.Email,
		LastLoginAt:       time.Unix(ud.LastLoginAt, 0).Format("15:04:05 01-02-2006"),
		Location:          ud.Location,
		FullName:          ud.FullName,
		CreateAt:          time.Unix(ud.CreatedAt, 0).Format("15:04:05 01-02-2006"),
		Settings:          *FromSettingData(&ud.Settings),
		Domain:            ud.Settings.Subdomain,
		PublicKey:         ud.Settings.FullKey,
		Fingerprint:       fingerprint,
		SubscriptionLevel: ud.Subscription.Level.String(),
	}
	return u
}

func ToSubscriptionData(SubscriptionLevel, SubscriptionDuration string) *data.Subscription {
	sd := &data.Subscription{}
	return sd
}

func ToUserData(u *User) *data.User {
	ud := &data.User{
		Id:           primitive.NewObjectID(),
		Email:        u.Email,
		FullName:     u.FullName,
		Location:     u.Location,
		LastLoginAt:  time.Now().Unix(),
		Username:     u.Username,
		CreatedAt:    time.Now().Unix(),
		Settings:     *ToSettingData(&u.Settings),
		Subscription: *ToSubscriptionData(u.SubscriptionLevel, u.SubscriptionDuration),
	}
	return ud
}

type File struct {
	FileName  string
	Extension string
	Mime      string
	Reader    io.Reader
}
