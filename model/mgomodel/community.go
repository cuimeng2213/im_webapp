package mgomodel

import (
	"time"
)

type Community struct {
	Id_ string `json:"id" bson:"_id,omitempty"`
	//名称
	Name string `json:"name" bson:"name"`
	//群主ID
	Ownerid string `json:"ownerid" bson:"ownerid"`
	//群logo
	Icon string `json:"icon" bson:"icon"`
	//cate 需要注意如果设置了omitempty, 设置0值时mongo驱动会忽略0值 踩坑路过^_^
	Cate int `json:"cate" bson:"cate"`
	//描述
	Memo string `json:"memo" bson:"memo,omitempty"`
	// 创建时间
	Createat time.Time `json:"createat" bson:"createat"`
	// 群聊ID值
	Uuid string `json:"uuid" bson:"uuid"`
}
