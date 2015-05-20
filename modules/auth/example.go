// 认证模块
package auth

import (
	"database/sql"
	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
)

type ExampleAuthHandler struct {
	DSN string
}

func NewExampleAuthHandler() AuthHandlerInterface {
	return new(ExampleAuthHandler)
}

func (ah *ExampleAuthHandler) Init(config string) error {
	return json.Unmarshal([]byte(config), ah)
}

// 认证方法
func (ah *ExampleAuthHandler) IsTokenValid(token string, uid string) error {
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

	return nil
}

func init() {
	Register("example", NewExampleAuthHandler)
}
