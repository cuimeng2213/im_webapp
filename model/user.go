package model

import (
	"time"
)

const (
	SexWoman   = "W"
	SexMan     = "W"
	SexUnknown = "U"
)

type User struct {
	Id     int64  `xorm:"pk autoincr bigint(64)" form:"id" json:"id"`
	Mobile string `xorm:"varchar(20)" form:"mobile" json:"mobile"`
	//用户密码 plainowd+salt
	Passwd   string    `xorm:"varchar(40)" form:"passwd" json:"-"`
	Avatar   string    `xorm:"varchar(150)" form:"avatar" json:"avatar"`
	Sex      string    `xorm:"varchar(2)" form:"sex" json:"sex"`
	Nickname string    `xorm:"varchar(20)" form:"nickname" json:"nickname"`
	Salt     string    `xorm:"varchar(10)" form:"salt" json:"-"`
	Online   int       `xorm:"int(10)" form:"online" json:"online"`
	Token    string    `xorm:"varchar(40)" from:"token" json:"token"`
	Memo     string    `xorm:"varchar(140)" form:"memo" json:"memo"`
	Createat time.Time `xorm:"datatime" form:"createat" json:"createat"`
}
