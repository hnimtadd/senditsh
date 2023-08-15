package api

import (
	"fmt"
	"time"

	"github.com/hnimtadd/senditsh/data"
	"github.com/hnimtadd/senditsh/utils"
)

func (h *ApiHandlerImpl) CreateUser(user *User) error {
	usr, _ := h.repo.GetUserByUserName(user.Username)
	if usr != nil {
		return fmt.Errorf("Error while creating user with username %v. Err: Username existed", user.Username)
	}
	newUser := ToUserData(user)
	newUser.CreatedAt = time.Now().Unix()
	newUser.LastLoginAt = time.Now().Unix()
	if err := h.repo.CreateUser(newUser); err != nil {
		return err
	}
	return nil
}

func (h *ApiHandlerImpl) RegisterUserSubscription(userId string, subscript data.Subscription) error {
	return nil
}

func (h *ApiHandlerImpl) GetSettingOfUser(userName string) (*Settings, error) {
	setting, err := h.repo.GetSettingOfUser(userName)
	if err != nil {
		return nil, err
	}
	res := FromSettingData(setting)
	return res, nil
}

func (h *ApiHandlerImpl) RegisterUserDomainSetting(userName string, domain string) error {
	settings, err := h.GetSettingOfUser(userName)
	if err != nil {
		return err
	}

	// not need to update
	if domain == settings.Subdomain {
		return nil
	}
	if settings.Subdomain != "" {
		logger.Info("msg", "user overrided subdomain")
	}

	ok, err := h.ValidateDomain(domain)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("This domain already exists\n")
	}
	if err := h.repo.InsertUserDomain(userName, domain); err != nil {
		return err
	}
	return nil
}

func (h *ApiHandlerImpl) RegisterUserSSHKeySetting(userName string, publicKey string) error {
	setting, err := h.GetSettingOfUser(userName)
	logger.Info("msg", "setup new ssh key")
	if err != nil {
		return err
	}
	if setting.FullKey == publicKey {
		return nil
	}
	logger.Info("msg", "setup new ssh key")

	sshHash, sshKey, err := utils.ParsePublicKey(publicKey)
	if err != nil {
		return err
	}
	if setting.SSHHash == sshHash && setting.SSHKey == sshKey {
		return nil
	}

	if setting.SSHHash != "" && setting.SSHKey != "" {
		logger.Info("msg", "user overrided sshKey")
	}

	if err := h.repo.InsertUserSSHKey(userName, publicKey, sshKey, sshHash); err != nil {
		return err
	}
	return nil
}

func (h *ApiHandlerImpl) GetUsers() ([]User, error) {
	usrs, err := h.repo.GetUsers()
	if err != nil {
		return nil, err
	}

	users := []User{}
	for _, usr := range usrs {
		user := FromUserData(&usr)
		users = append(users, *user)
	}
	return users, nil
}

func (h *ApiHandlerImpl) GetUserByUserName(userName string) (*User, error) {
	usr, err := h.repo.GetUserByUserName(userName)
	if err != nil {
		return nil, err
	}
	user := FromUserData(usr)
	return user, nil
}

func (h *ApiHandlerImpl) UpdateUserInformation(userName string, fullName string, email string, location string) error {
	if err := h.repo.UpdateUserInformation(userName, fullName, email, location); err != nil {
		return err
	}
	return nil
}

func (h *ApiHandlerImpl) GetUserSharingLink(userName string) (*Transfer, error) {
	transferd, err := h.repo.GetLastTransfer(userName)
	if err != nil {
		return nil, err
	}
	transfer := FromTransferData(transferd)
	if transfer.Status == "Not done" {
		return transfer, nil
	}
	return nil, nil
}
func (h *ApiHandlerImpl) GetUserByDomain(domain string) (*User, error) {
	user, err := h.repo.GetUserByDomain(domain)
	if err != nil {
		return nil, err
	}
	usr := FromUserData(user)
	return usr, nil
}
