package auth

import (
	"encoding/json"
	"errors"

	"github.com/denghongcai/MessageHive/message"
	"github.com/op/go-logging"
)

type authMsg struct {
	Token string `json:"token"`
}

var log = logging.MustGetLogger("main")

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
	log.Info("Uid: %s, Token: %s, authenticated", uid, token)
	//TODO

	return nil
}
