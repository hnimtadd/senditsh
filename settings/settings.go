package settings

import "time"

const (
	Timeout                = time.Minute * 5
	FileNameDefault string = "sendit"
	MsgOption       string = "msg"
	FileNameOption  string = "filename"
	ToEmailOption   string = "toemail"
	TimeoutOption   string = "timeout"
	MsgMinLen       int    = 2
	MsgMaxLen       int    = 50
	FileNameMinLen  int    = 3
	FileNameMaxLen  int    = 30
)
