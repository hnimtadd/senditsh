package utils

import (
	"encoding/base64"
	"errors"
	"strings"

	"golang.org/x/crypto/ssh"
)

func ParsePublicKey(pubicKey string) (sshHash string, sshKey string, err error) {
	parts := strings.Split(pubicKey, " ")
	if len(parts) < 2 {
		err = errors.New("Public key not valid")
		return
	}
	switch sshHash = parts[0]; sshHash {
	case "ssh-rsa":
		return sshHash, parts[1], nil
	default:
		return "", "", errors.New("sshHash not supported")
	}
}

func GetFingerprint(sshKey string) (string, error) {
	keyByte, err := base64.StdEncoding.DecodeString(sshKey)
	if err != nil {
		return "", err
	}
	key, err := ssh.ParsePublicKey(keyByte)
	if err != nil {
		return "", err
	}
	fingerprint := ssh.FingerprintSHA256(key)
	return fingerprint, nil
}
