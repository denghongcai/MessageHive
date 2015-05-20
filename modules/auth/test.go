// 认证模块
package auth

import (
	"encoding/json"
	"errors"
)

type TestAuthHandler struct {
}

func NewTestAuthHandler() AuthHandlerInterface {
	return new(TestAuthHandler)
}

func (th *TestAuthHandler) Init(config string) error {
	return json.Unmarshal([]byte(config), th)
}

// 认证方法
func (th *TestAuthHandler) IsTokenValid(token string, uid string) error {
	if uid == "00000001" {
		return nil
	}
	if token != "hehe" {
		return errors.New("you lose!")
	}

	return nil
}

func init() {
	Register("test", NewTestAuthHandler)
}
