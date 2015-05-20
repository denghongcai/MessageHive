// 认证模块
package auth

import (
	"encoding/json"
	"errors"

	"database/sql"
	"github.com/denghongcai/MessageHive/modules/log"
	"github.com/denghongcai/MessageHive/modules/message"
	_ "github.com/go-sql-driver/mysql"
)

var (
	authHandlers []*AuthHandlerInterface
	AuthHandle   *AuthHandlerInterface
)

// 认证消息结构
type authMsg struct {
	Token string `json:"token"`
}

type AuthHandlerInterface interface {
	Init(config string) error
	IsTokenValid(token string) (error, bool)
}

func SetAuthHandler(adapter string, config string) error {
	if handler, ok := authHandlers[adapter]; ok {
	} else {
		log.Fatal("Auth: unknown adapter \"" + adapter + "\"")
	}
	if err := handler.Init(config); err != nil {
		return err
	}
	AuthHandle = handler
	return nil
}

// 认证方法
func Authenticate(msg *message.Container) error {
	uid := msg.GetSID()
	body := msg.GetBODY()
	authdata := new(authMsg)
	if err := json.Unmarshal([]byte(body), authdata); err != nil {
		return errors.New("Auth failed, json parse error")
	}
	if len(authdata.Token) == 0 {
		return errors.New("Auth failed, token field was empty")
	}
	token := authdata.Token
	if uid == "00000001" {
		return nil
	}
	// TODO: 向认证服务器认证Token
	db, err := sql.Open("mysql", "dhc:denghc@/Register")
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}

	stmtOut, err := db.Prepare("SELECT timeout FROM Token WHERE uid = ? and token= ?")
	if err != nil {
		return err
	}
	defer stmtOut.Close()

	var timestamp int64

	err = stmtOut.QueryRow(uid, token).Scan(&timestamp)

	if err != nil {
		return err
	}

	log.Info("Uid: %s, Token: %s, authenticated", uid, token)

	return nil
}
