package mgomodel

import (
	"time"
)

const (
	SexWoman   = "W"
	SexMan     = "W"
	SexUnknown = "U"
)

type User struct {
	Id_    string `json:"id" bson:"_id,omitempty"`
	Mobile string `json:"mobile"`
	//用户密码 plainowd+salt
	Passwd   string    `json:"-"`
	Avatar   string    `json:"avatar"`
	Sex      string    `json:"sex"`
	Nickname string    `json:"nickname"`
	Salt     string    `json:"-"`
	Online   int       `json:"online"`
	Token    string    `json:"token"`
	Memo     string    `json:"memo"`
	Createat time.Time `json:"createat"`
}
