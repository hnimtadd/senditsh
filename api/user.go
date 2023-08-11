package api

import (
	"fmt"
	"time"

	"github.com/hnimtadd/senditsh/data"
	"github.com/hnimtadd/senditsh/utils"
)

func (h *ApiHandlerImpl) CreateUser(user *data.User) error {
	usr, _ := h.repo.GetUserByUserName(user.Username)
	if usr != nil {
		return fmt.Errorf("Error while creating user with username %v. Err: Username existed", user.Username)
	}
	user.CreatedAt = time.Now().Unix()
	user.LastLoginAt = time.Now().Unix()
	if err := h.repo.CreateUser(user); err != nil {
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
	res := Settings{}
	res.FromSettingData(setting)
	return &res, nil
}

func (h *ApiHandlerImpl) RegisterUserDomainSetting(userName string, domain string) error {
	settings, err := h.GetSettingOfUser(userName)
	if err != nil {
		return err
	}
	if settings.Subdomain != "" {
		logger.Info("msg", "user overrided subdomain")
	}
	if err := h.repo.InsertUserDomain(userName, domain); err != nil {
		return err
	}
	return nil
}

func (h *ApiHandlerImpl) RegisterUserSSHKeySetting(userName string, publicKey string) error {
	sshHash, sshKey, err := utils.ParsePublicKey(publicKey)
	if err != nil {
		return err
	}
	setting, err := h.GetSettingOfUser(userName)
	if err != nil {
		return err
	}
	if setting.SSHHash != "" && setting.SSHKey != "" {
		logger.Info("msg", "user overrided sshKey")
	}

	if err := h.repo.InsertUserSSHKey(userName, sshKey, sshHash); err != nil {
		return err
	}
	return nil
}

func (h *ApiHandlerImpl) GetUsers() ([]data.User, error) {
	users, err := h.repo.GetUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (h *ApiHandlerImpl) GetUserByUserName(userName string) (*data.User, error) {
	user, err := h.repo.GetUserByUserName(userName)
	if err != nil {
		return nil, err
	}
	return user, nil
}
