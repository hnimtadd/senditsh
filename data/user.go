package data

import (
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id           primitive.ObjectID `json:"userId,omiempty" bson:"_id,omiempty"`
	Email        string             `json:"email,omiempty" bson:"email,omiempty"`
	FullName     string             `json:"fullName,omiempty" bson:"fullName,omiempty"`
	Username     string             `json:"userName,omiempty" bson:"userName,omiempty"`
	Location     string             `json:"location,omiempty" bson:"location,omiempty"`
	CreatedAt    int64              `json:"createdAt,omiempty" bson:"createdAt,omiempty"`
	LastLoginAt  int64              `json:"lastLoginAt,omiempty" bson:"lastLoginAt,omiempty"`
	Subscription Subscription       `json:"subscription,omiempty" bson:"subscription,omiempty"`
	Settings     Settings           `json:"settings,omiempty" bson:"settings,omiempty"`
}

func (u User) String() string {
	str := strings.Builder{}
	str.WriteString(fmt.Sprintf("UserId: %v", u.Id))
	str.WriteString(fmt.Sprintf("Email: %v", u.Email))
	str.WriteString(fmt.Sprintf("FullName: %v", u.FullName))
	str.WriteString(fmt.Sprintf("Username: %v", u.Username))
	str.WriteString(fmt.Sprintf("Location: %v", u.Location))
	str.WriteString(fmt.Sprintf("CreatedAt: %v", time.Unix(u.CreatedAt, 0)))
	str.WriteString(fmt.Sprintf("LastLoginAt: %v", time.Unix(u.LastLoginAt, 0)))
	str.WriteString(fmt.Sprintf("SubScription: %v", u.Subscription))
	str.WriteString(fmt.Sprintf("Settings: %v", u.Settings))
	return str.String()
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

type Subscription struct {
	Level     Level `json:"level,omiempty" bson:"level,omiempty"`
	CreatedAt int64 `json:"createdAt,omiempty" bson:"createdAt,omiempty"`
	ExpiredAt int64 `json:"expired,omiempty" bson:"expired,omiempty"`
}
type Settings struct {
	Subdomain  string `json:"subdomain,omiempty" bson:"subdomain,omiempty"`
	SSHKey     string `json:"sshKey,omiempty" bson:"sshKey,omiempty"`
	SSHHash    string `json:"sshHash,omiempty" bson:"sshHash,omiempty"`
	ModifiedAt int64  `json:"modifiedAt,omiempty" bson:"modifiedAt,omiempty"`
}

func (s Settings) String() string {
	str := strings.Builder{}
	str.WriteString(fmt.Sprintf("SubDomain: %v", s.Subdomain))
	str.WriteString(fmt.Sprintf("SSHKey: %v", s.SSHKey))
	str.WriteString(fmt.Sprintf("SSHHash: %v", s.SSHHash))
	str.WriteString(fmt.Sprintf("ModifiedAt: %v", time.Unix(s.ModifiedAt, 0)))
	return str.String()
}
