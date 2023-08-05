package api

import "github.com/gliderlabs/ssh"

func (api *ApiHandlerImpl) AuthenticationPublicKeyFromClient() ssh.PublicKeyHandler {
	// Get authorized publicKey from db, check if that publicKey is accepted
	return func(ctx ssh.Context, key ssh.PublicKey) bool {
		return true
	}
}
