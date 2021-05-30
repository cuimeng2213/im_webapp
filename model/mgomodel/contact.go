package mgomodel

import (
	"time"
)

const (
	CONCAT_CATE_USER     = 0x01
	CONCAT_CATE_COMUNITY = 0x02
)

//好友和群在一个表里面 是否要区分开
type Contact struct {
	Id_ string `form:"id" json:"id" bson:"_id,omitempty"`
	//群id/userid  ObjectId
	Ownerid string `form:"ownerid" json:"ownerid" bson:"ownerid"`
	//对端ObjectId
	Dstobj   string    `form:"dstobj" json:"dstobj" bson:"dstobj"`
	Cate     int       `form:"cate" json:"cate" bson:"cate"`
	Memo     string    `form:"memo" json:"memo" bson:"memo,omitempty"`
	Createat time.Time `form:"createat" json:"createat" bson:"createat"`
}
