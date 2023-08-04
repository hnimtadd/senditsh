package data

type User struct {
	Id           string       `json:"userId,omiemtpy" bson:"_id,omiempty"`
	Email        string       `json:"email,omiempty" bson:"email,omiempty"`
	FullName     string       `json:"fullName,omiempty" bson:"fullName,omiempty"`
	Username     string       `json:"userName,omiempty" bson:"userName,omiempty"`
	Location     string       `json:"location,omiempty" bson:"location,omiempty"`
	CreateAt     int64        `json:"createdAt,omiempty" bson:"createdAt,omiempty"`
	LastLoginAt  int64        `json:"lastLoginAt,omiempty" bson:"lastLoginAt,omiempty"`
	Subscription Subscription `json:"subscription,omiempty" bson:"subscription,omiempty"`
	Settings     Settings     `json:"settings,omiempty" bson:"settings,omiempty"`
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

// func FindUserByID(hex string) (*User, error) {
// 	id, err := primitive.ObjectIDFromHex(hex)
// 	if err != nil {
// 		return nil, err
// 	}
// }
