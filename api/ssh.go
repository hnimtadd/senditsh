package api

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gliderlabs/ssh"
	"github.com/hnimtadd/senditsh/data"
	"github.com/hnimtadd/senditsh/settings"
)

func (api *ApiHandlerImpl) FinalizeAndCleanUpAfterTransfer(transfer *data.Transfer) error {
	logger.Info("msg", "Deleted tunnel")
	return nil
}

type SSHSession struct {
	Context context.Context
	Session ssh.Session
	Link    string
	User    *User
	Opt     *UserOptions
	Tunnel  *Tunnel
	File    *data.File
}

func NewSSHSession(s ssh.Session, user *User, opt *UserOptions) (*SSHSession, error) {
	session := &SSHSession{
		Session: s,
		Opt:     opt,
		User:    user,
	}
	opt, err := ParseUserOptions(session.Session.Command())
	if err != nil {
		return nil, err
	}

	session.Opt = opt
	return session, nil
}

func (s *SSHSession) GetID() string {
	return s.Session.Context().SessionID()
}

type UserOptions struct {
	FileName string
	ToEmail  string
	Msg      string
	Expired  time.Duration
}

func (uOpt UserOptions) String() string {
	str := strings.Builder{}
	str.WriteString("[ ")
	if uOpt.FileName != "" {
		str.WriteString(fmt.Sprintf("fileName: %v ", uOpt.FileName))
	}
	if uOpt.Msg != "" {
		str.WriteString(fmt.Sprintf("msg: %v ", uOpt.Msg))
	}

	if uOpt.ToEmail != "" {
		str.WriteString(fmt.Sprintf("toEmail: %v ", uOpt.ToEmail))
	}
	str.WriteString("]")
	return str.String()
}

func (s *SSHSession) GetUserDomain() string {
	if s.User != nil {
		return fmt.Sprintf("http://%v.mysendit.sh", s.User.Domain)
	}
	return ""
}

func ParseUserOptions(commands []string) (*UserOptions, error) {
	var usrOpts = &UserOptions{}
	for _, command := range commands {
		parts := strings.Split(command, "=")
		if len(parts) > 2 {
			return nil, fmt.Errorf("option must specified in key=value format. Ex: msg=\"hell\"")
		}
		key := (parts[0])
		switch key {
		case settings.MsgOption:
			if len(parts) != 2 {
				return nil, fmt.Errorf("msg must not null\n")
			}
			msg := parts[1]
			if len(msg) > settings.MsgMaxLen || len(msg) < settings.MsgMinLen {
				return nil, fmt.Errorf("msg len must in range %v-%v characters\n", settings.MsgMinLen, settings.MsgMaxLen)
			}
			usrOpts.Msg = msg

		case settings.FileNameOption:
			if len(parts) != 2 {
				return nil, fmt.Errorf("file name must not null\n")
			}
			fName := parts[1]
			if len(fName) > settings.FileNameMaxLen || len(fName) < settings.FileNameMinLen {
				return nil, fmt.Errorf("file name len must in range %v-%v characters\n", settings.FileNameMinLen, settings.FileNameMaxLen)
			}
			usrOpts.FileName = fName

		case settings.ToEmailOption:
			if len(parts) != 2 {
				return nil, fmt.Errorf("email must not null\n")
			}
			// TODO: parse valid email
			email := parts[1]
			usrOpts.ToEmail = email
		case settings.TimeoutOption:
			if len(parts) != 2 {
				return nil, fmt.Errorf("time must not null\n")
			}
			duration := parts[1]
			unit := duration[len(duration)-1:]
			numstr := duration[0 : len(duration)-1]
			num, err := strconv.Atoi(numstr)
			if err != nil {
				return nil, fmt.Errorf("Time must in formation [num][unit], in which num is int and unit is one of (s,m,h)\n")
			}

			switch unit {
			case "h":
				usrOpts.Expired = time.Hour * time.Duration(num)
			case "m":
				usrOpts.Expired = time.Minute * time.Duration(num)
			case "s":
				usrOpts.Expired = time.Second * time.Duration(num)
			default:
				return nil, fmt.Errorf("Time must in formation [num][unit], in which num is int and unit is one of (s,m,h)\n")
			}
		}
	}
	return usrOpts, nil
}
func (s *SSHSession) SetFile(f *data.File) {
	s.File = f

}

func (s *SSHSession) SetContext(ctx context.Context) {
	s.Context = ctx
}
func (s *SSHSession) GetContext(ctx context.Context) context.Context {
	return s.Context
}

func (s *SSHSession) GetFile() *data.File {
	return s.File
}
