package model

import "errors"

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u UserRequest) Validate() error {

	if u.Username == "" {
		return errors.New("username 不可为空")
	}

	if u.Password == "" {
		return errors.New("password 不可为空")
	}

	return nil
}
