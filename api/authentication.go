package api

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/gliderlabs/ssh"
	"github.com/hnimtadd/senditsh/data"
)

func fetchKeys(ctx context.Context, username string) ([]ssh.PublicKey, error) {
	fi, err := os.Open("authorized_keys")
	if err != nil {
		return nil, err
	}

	keys := []ssh.PublicKey{}
	authorizedKeysBytes, err := io.ReadAll(fi)
	if err != nil {
		return nil, err
	}

	for len(authorizedKeysBytes) > 0 {
		pubKey, _, _, rest, err := ssh.ParseAuthorizedKey(authorizedKeysBytes)
		if err != nil {
			return nil, err
		}
		keys = append(keys, pubKey)
		authorizedKeysBytes = rest
	}
	return keys, nil
}

func (api *ApiHandlerImpl) GetUserWithPubKey(pub string) (*data.User, error) {
	user, err := api.repo.GetUserByPublicKey(pub)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (api *ApiHandlerImpl) AuthenticationPublicKeyFromClient() ssh.PublicKeyHandler {
	// Get authorized publicKey from db, check if that publicKey is accepted
	return func(ctx ssh.Context, key ssh.PublicKey) bool {
		userName := ctx.User()
		pubKey := base64.StdEncoding.EncodeToString(key.Marshal())
		user, err := api.GetUserWithPubKey(pubKey)
		if err != nil {
			logger.Info("msg", fmt.Sprintf("failed while authenticate for user %v", userName), "err", err)
		}

		if user != nil {
			logger.Info("msg", fmt.Sprintf("authenticated for user: %v", user.Username))
			ctx.SetValue("user", user)
			return true
		}

		logger.Info("msg", "anonymous session open")
		return true
	}
}
