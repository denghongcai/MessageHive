// 认证模块
package auth

import (
	"encoding/json"
	"errors"

	"github.com/denghongcai/MessageHive/modules/log"
	"github.com/denghongcai/MessageHive/modules/message"
	_ "github.com/go-sql-driver/mysql"
)

type AuthHandlerInterface interface {
	Init(config string) error
	IsTokenValid(token string, uid string) error
}

type authHandlerType func() AuthHandlerInterface

var (
	authHandlers = make(map[string]authHandlerType)
	AuthHandle   AuthHandlerInterface
)

// 认证消息结构
type authMsg struct {
	Token string `json:"token"`
}

func SetAuthHandler(adapter string, config string) error {
	if handler, ok := authHandlers[adapter]; ok {
		authHandler := handler()
		if err := authHandler.Init(config); err != nil {
			return err
		}
		AuthHandle = authHandler
		return nil
	} else {
		log.Fatal("Auth: unknown adapter \"" + adapter + "\"")
		return nil
	}
}

func Register(name string, authHandler authHandlerType) {
	if authHandler == nil {
		log.Fatal("Auth: register provider is null")
	}
	if _, dup := authHandlers[name]; dup {
		log.Fatal("Auth: register called twice for provider \"" + name + "\"")
	}
	authHandlers[name] = authHandler
}

// 认证方法
func Authenticate(msg *message.Container) error {
	uid := msg.GetSID()
	body := msg.GetBODY()
	authdata := new(authMsg)
	if err := json.Unmarshal([]byte(body), authdata); err != nil {
		return errors.New("Auth: failed, json parse error")
	}
	if len(authdata.Token) == 0 {
		return errors.New("Auth: failed, token field was empty")
	}
	token := authdata.Token
	if err := AuthHandle.IsTokenValid(token, uid); err != nil {
		return err
	} else {
		log.Info("Auth: uid: %s, token: %s, authenticated", uid, token)
		return nil
	}
}
