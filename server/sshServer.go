package server

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gliderlabs/ssh"
	"github.com/hnimtadd/senditsh/api"
	"github.com/hnimtadd/senditsh/config"
	"github.com/hnimtadd/senditsh/data"
)

type SSHServer interface {
	Listen() error
}

type SSHServerImpl struct {
	config *config.SSHConfig
	api    *api.ApiHandlerImpl
}

func NewSSHServerImpl(api *api.ApiHandlerImpl, config *config.SSHConfig) (SSHServer, error) {
	server := &SSHServerImpl{
		config: config,
		api:    api,
	}
	go server.initHandler()
	return server, nil
}

func (server *SSHServerImpl) Listen() error {
	logger.Info("listen and serve on port", ":"+server.config.Port)
	if err := ssh.ListenAndServe(
		":"+server.config.Port,
		nil,
		ssh.HostKeyFile("sendit"),
		ssh.PublicKeyAuth(server.api.AuthenticationPublicKeyFromClient()),
		ssh.NoPty(),
	); err != nil {
		return err
	}
	return nil
}

func (server *SSHServerImpl) initHandler() {
	ssh.Handle(func(s ssh.Session) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
		defer cancel()
		logger.Info("msg", fmt.Sprintf("New connection to server: %v, remote ip: %v, command: %v", s.Context().SessionID()[:5], s.RemoteAddr(), s.Command()))
		opt, err := parseUserOptions(s.Command())
		ctx = context.WithValue(ctx, "opt", opt)
		logger.Info("msg", "parsedUserOption", "opt", opt)
		if err != nil {
			logger.Error(err)
			s.Exit(1)
			return
		}

		// id := s.Context().SessionID()
		id := "sample"
		ctx = context.WithValue(ctx, "link", id)
		if err := server.api.InitTunnel(id); err != nil {
			log.Println(err)
			s.Exit(1)
			return
		}
		logger.Info("msg", fmt.Sprintf("http://%v:%v/api/v1/transfer/%v", server.config.Host, 3000, id))

		// Wait for user click shared link, if not, drop this session
		tunnel, err := server.api.WaitToGetTunnel(ctx, id)
		if err != nil {
			logger.Error("msg", err)
			s.Exit(1)
			return
		}

		// Get file from ssh connection
		if opt.fileName != "" {
			ctx = context.WithValue(ctx, "fileName", opt.fileName)
		}
		file, err := server.api.GetFileInfo(ctx, s)
		if err != nil {
			logger.Error("msg", err)
			s.Exit(1)
			return
		}
		// logger.Info("file", file)
		ctx = context.WithValue(ctx, "file", file)

		if err := server.api.CopyToTunnel(tunnel, file.Reader); err != nil {
			logger.Error("msg", err)
			s.Exit(1)
			return
		}

		transfer, err := server.createTransfer(ctx, s)
		if err != nil {
			log.Fatal(err)
		}

		defer server.api.FinalizeAndCleanUpAfterTransfer(transfer, tunnel)

		if err := server.api.CreateTransfer(transfer); err != nil {
			logger.Error("msg", err)
			s.Exit(1)
			return
		}
		return
	})
}

// type PeekSession struct {
// 	err error
// 	r   io.Reader
// }

// func peekSession(s ssh.Session) chan PeekSession {
// 	var (
// 		peekch = make(chan PeekSession)
// 		pr     = bufio.NewReader(s)
// 	)
// 	go func(r *bufio.Reader) {
// 		b, err := r.Peek(1)
// 		if err != nil {
// 			peekch <- PeekSession{err: err}
// 			return
// 		}
// 		if len(b) == 0 {
// 			peekch <- PeekSession{err: fmt.Errorf("empty bytes")}
// 			return
// 		} else {
// 			peekch <- PeekSession{r: pr}
// 			return
// 		}
// 	}(pr)
// 	return peekch
// }

func (server *SSHServerImpl) createTransfer(ctx context.Context, s ssh.Session) (*data.Transfer, error) {
	file := &data.File{}
	if fileParam := ctx.Value("file"); fileParam != nil {
		file = fileParam.(*data.File)
	}

	opt := &UserOptions{}
	if optParam := ctx.Value("opt"); optParam != nil {
		opt = optParam.(*UserOptions)
	}

	from := ""
	if fromParam := ctx.Value("from"); fromParam != nil {
		from = fromParam.(string)
	}

	link := ""
	if linkParam := ctx.Value("link"); linkParam != nil {
		link = linkParam.(string)
	}

	transfer := data.Transfer{
		Filename:    file.FileName,
		From:        from,
		Link:        link,
		ToEmail:     opt.toEmail,
		Message:     opt.msg,
		Initiator:   s.User(),
		InitiatorIP: s.RemoteAddr().String(),
		CreatedAt:   time.Now().UnixNano(),
	}
	return &transfer, nil
}

type UserOptions struct {
	fileName string
	toEmail  string
	msg      string
}

func (uOpt UserOptions) String() string {
	str := strings.Builder{}
	str.WriteString("[ ")
	if uOpt.fileName != "" {
		str.WriteString(fmt.Sprintf("fileName: %v ", uOpt.fileName))
	}
	if uOpt.msg != "" {
		str.WriteString(fmt.Sprintf("msg: %v ", uOpt.msg))
	}

	if uOpt.toEmail != "" {
		str.WriteString(fmt.Sprintf("toEmail: %v ", uOpt.toEmail))
	}
	str.WriteString("]")
	return str.String()
}

const (
	fileNameDefault string = "sendit"
	msg             string = "msg"
	fileName        string = "filename"
	toEmail         string = "to-email"
	msgMinLen       int    = 2
	msgMaxLen       int    = 50
	fileNameMinLen  int    = 3
	fileNameMaxLen  int    = 30
)

func parseUserOptions(commands []string) (*UserOptions, error) {
	var usrOpts = &UserOptions{}
	for _, command := range commands {
		parts := strings.Split(command, "=")
		if len(parts) > 2 {
			return nil, fmt.Errorf("option must specified in key=value format. Ex: msg=\"hell\"")
		}
		key := (parts[0])
		switch key {
		case msg:
			if len(parts) != 2 {
				return nil, fmt.Errorf("msg must not null\n")
			}
			msg := parts[1]
			if len(msg) > msgMaxLen || len(msg) < msgMinLen {
				return nil, fmt.Errorf("msg len must in range %v-%v characters\n", msgMinLen, msgMaxLen)
			}
			usrOpts.msg = msg

		case fileName:
			if len(parts) != 2 {
				return nil, fmt.Errorf("file name must not null\n")
			}
			fName := parts[1]
			if len(fName) > fileNameMaxLen || len(fName) < fileNameMinLen {
				return nil, fmt.Errorf("file name len must in range %v-%v characters\n", msgMinLen, msgMaxLen)
			}
			usrOpts.fileName = fName

		case toEmail:
			if len(parts) != 2 {
				return nil, fmt.Errorf("email must not null\n")
			}
			// TODO: parse valid email
			email := parts[1]
			usrOpts.toEmail = email
		}
	}
	return usrOpts, nil
}
