package data

import (
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transfer struct {
	Id          primitive.ObjectID `json:"transferId,omiempty" bson:"_id,omiempty"`
	Link        string             `json:"link,omiempty" bson:"link,omiempty"`
	Filename    string             `json:"fileName,omiempty" bson:"fileName,omiempty"`
	UserName    string             `json:"userName,omiempty" bson:"userName,omiempty"`
	ToEmail     string             `json:"toEmail,omiempty" bson:"toEmail,omiempty"`
	Message     string             `json:"message,omiempty" bson:"message,omiempty"`
	IsVerified  bool               `json:"isVerified,omiempty" bson:"isVerified,omiempty"`
	Status      string             `json:"status,omiempty" bson:"status,omiempty"`
	Initiator   string             `json:"initiator,omiempty" bson:"initiator,omiempty"`
	InitiatorIP string             `json:"initiatorIP,omiempty" bson:"initiatorIP,omiempty"`
	CreatedAt   int64              `json:"createdAt,omiempty" bson:"createdAt,omiempty"`
}

func (t Transfer) String() string {
	builder := strings.Builder{}
	builder.WriteString("\n--------Transfer--------\n")
	builder.WriteString(fmt.Sprintf("Id: %v\n", t.Id.String()))
	builder.WriteString(fmt.Sprintf("Filename: %v\n", t.Filename))
	builder.WriteString(fmt.Sprintf("From: %v\n", t.UserName))
	builder.WriteString(fmt.Sprintf("ToEmail: %v\n", t.ToEmail))
	builder.WriteString(fmt.Sprintf("Link: %v\n", t.Link))
	builder.WriteString(fmt.Sprintf("Message: %v\n", t.Message))
	builder.WriteString(fmt.Sprintf("Status: %v\n", t.Status))
	builder.WriteString(fmt.Sprintf("IsVerified: %v", t.IsVerified))
	builder.WriteString(fmt.Sprintf("Inititor: %v\n", t.Initiator))
	builder.WriteString(fmt.Sprintf("InititorIP: %v\n", t.InitiatorIP))
	builder.WriteString(fmt.Sprintf("CreateAt: %v\n", time.UnixMicro(t.CreatedAt).String()))
	builder.WriteString("------------------------\n")
	return builder.String()
}
