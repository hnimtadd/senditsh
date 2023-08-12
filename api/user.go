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
