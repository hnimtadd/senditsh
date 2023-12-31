package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/gliderlabs/ssh"
	"github.com/hnimtadd/senditsh/api"
	"github.com/hnimtadd/senditsh/config"
	"github.com/hnimtadd/senditsh/data"
	"github.com/hnimtadd/senditsh/settings"
	"github.com/hnimtadd/senditsh/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	ssh.Handle(server.TransferFileSessionHandler())
}

func (server *SSHServerImpl) TransferFileSessionHandler() ssh.Handler {
	return func(s ssh.Session) {
		session, err := server.initSSHSession(s)
		if session.User != nil {
			if err := generateVerifiedUserMessage(session, s); err != nil {
				logger.Error("err", err)
				s.Exit(1)
				return
			}

		}
		if err != nil {
			io.WriteString(s, err.Error())
			return
		}

		timeout := settings.Timeout
		if session.Opt.Expired != 0 {
			timeout = session.Opt.Expired
		}

		var (
			ctx, cancel = context.WithTimeout(context.Background(), timeout)
			id          = s.Context().SessionID()
		)
		defer cancel()

		session.SetContext(ctx)
		session.Link = fmt.Sprintf("http://mysendit.sh/api/v1/transfer/%v", id)

		tunnel, err := server.api.InitTunnelWithID(ctx, id)
		if err != nil {
			logger.Error("err", err)
			io.WriteString(s, err.Error())
			s.Exit(1)
			return
		}
		defer server.api.DestroyTunnel(id)

		if err := generateSetupDoneMessage(session, s); err != nil {
			logger.Error("err", err)
			s.Exit(1)
			return
		}
		transfer, err := server.createTransfer(*session)
		if err != nil {
			log.Fatal(err)
		}

		logger.Info("created Transfer ", transfer)
		if err := server.api.CreateTransfer(transfer); err != nil {
			logger.Error("msg", err)
			s.Exit(1)
			return
		}

		logger.Info("inited tunnel")

		file, err := server.api.GetFileInfo(*session, ctx, s)
		if err != nil {
			generateCustomMessage(session, s, err.Error())
			logger.Error("msg", err)
			s.Exit(1)
			return
		}

		session.File = file
		if err := server.api.WaitForWriterPipeShake(ctx, id); err != nil {
			if err := generateExpiredTransferMessage(session, s); err != nil {
				logger.Error("error", err)
			}

			if err := server.api.UpdateTransferStatus(transfer.Id, "expired"); err != nil {
				logger.Error("error", err)
			}
			s.Exit(1)
			return
		}


		ctx = context.WithValue(ctx, "file", file)

		if err := tunnel.PipeReader(file.Reader); err != nil {
			logger.Error("err", err)
			s.Exit(1)
			return
		}

		if err := tunnel.CopyInTunnel(); err != nil {
			s.Exit(1)
			return
		}
		logger.Info("msg", "copied")


		defer func() {
			if err := generateTransferDoneMessage(session, s); err != nil {
				logger.Error("error", err)
			}
			server.api.FinalizeAndCleanUpAfterTransfer(transfer)
		}()

		logger.Info("Transfer done")
		if server.api.UpdateTransferStatus(transfer.Id, "done"); err != nil {
			logger.Error("err", err.Error())
		}
		return
	}
}

func (sever *SSHServerImpl) initSSHSession(s ssh.Session) (*api.SSHSession, error) {
	opt, err := api.ParseUserOptions(s.Command())
	if err != nil {
		return nil, err
	}
	user := utils.GetContextVariableWithType[*api.User](s.Context(), "user", nil)
	if user == nil {
		logger.Info("msg", fmt.Sprintf("New connection from anonymous user: %v, remote ip: %v, command: %v", s.User(), s.RemoteAddr(), s.Command()))
	} else {
		logger.Info("msg", fmt.Sprintf("New connection from user: %v, remote ip: %v, command: %v", user.Username, s.RemoteAddr(), s.Command()))
	}
	session, err := api.NewSSHSession(s, user, opt)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (server *SSHServerImpl) createTransfer(s api.SSHSession) (*data.Transfer, error) {
	var (
		from = "anonymous"
	)

	if s.User != nil {
		from = s.User.Username
	}

	fileName := settings.FileNameDefault
	if s.Opt.FileName != "" {
		fileName = s.Opt.FileName
	}
	transfer := data.Transfer{
		Id:          primitive.NewObjectID(),
		Filename:    fileName,
		UserName:    from,
		Link:        s.Link,
		Status:      "Not done",
		ToEmail:     s.Opt.ToEmail,
		Message:     s.Opt.Msg,
		Initiator:   s.Session.User(),
		InitiatorIP: s.Session.RemoteAddr().String(),
		CreatedAt:   time.Now().Unix(),
	}
	return &transfer, nil
}
func generateSetupDoneMessage(session *api.SSHSession, w io.Writer) error {
	str := strings.Builder{}
	str.WriteString("Direct download link:\n")
	str.WriteString("\t"+session.Link + "\n")

	user := session.User
	if user != nil && user.Settings.Subdomain != "" {
		str.WriteString("Download page:\n")
		str.WriteString("\t"+session.GetUserDomain() + "\n")
	}

	if _, err := io.WriteString(w, str.String()); err != nil {
		return err
	}
	return nil
}

func generateExpiredTransferMessage(session *api.SSHSession, w io.Writer) error {
	str := strings.Builder{}
	str.WriteString("Link expired!!!")
	str.WriteString("\n")

	if _, err := io.WriteString(w, str.String()); err != nil {
		return err
	}
	return nil
}

func generateTransferDoneMessage(session *api.SSHSession, w io.Writer) error {
	str := strings.Builder{}
	str.WriteString("Transfer done")
	str.WriteString("\n")

	if _, err := io.WriteString(w, str.String()); err != nil {
		return err
	}
	return nil
}
func generateVerifiedUserMessage(session *api.SSHSession, w io.Writer) error {
	str := strings.Builder{}
	str.WriteString("Detected verified user\n")
	str.WriteString(fmt.Sprintf("\tUserName: %s\n", session.User.Username))

	if _, err := io.WriteString(w, str.String()); err != nil {
		return err
	}
	return nil
}
func generateCustomMessage(session *api.SSHSession, w io.Writer, msg string) error {
	if _, err := io.WriteString(w, msg) ;err != nil{
		return err
	}
	return nil

}
