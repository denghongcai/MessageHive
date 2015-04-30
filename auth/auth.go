// 认证模块
package auth

import (
	"encoding/json"
	"errors"

	"database/sql"
	"github.com/denghongcai/MessageHive/message"
	_ "github.com/go-sql-driver/mysql"
	"github.com/op/go-logging"
)

// 认证消息结构
type authMsg struct {
	Token string `json:"token"`
}

var log = logging.MustGetLogger("main")

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
