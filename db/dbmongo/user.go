package dbmongo

import (
	"context"
	"errors"
	"fmt"
	"im_webapp/model/mgomodel"
	"im_webapp/util"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
}

func (u *UserService) Register(mobile, password, nickName, avatar, sex string) (mgomodel.User, error) {
	user := mgomodel.User{}
	ctx := context.Background()

	ChatUserCol.FindOne(ctx, bson.D{{Key: "mobile", Value: mobile}}).Decode(&user)
	fmt.Printf("have error %v %s %s\n", err, user.Nickname, user.Id_)
	if mobile == user.Mobile && user.Id_ != "" {
		return user, fmt.Errorf("%s monile is regintered", mobile)
	}
	fmt.Println(">>>>>> do register user")
	//注册新用户
	user.Mobile = mobile
	user.Nickname = nickName
	user.Avatar = avatar
	user.Sex = sex
	user.Createat = time.Now()
	rand.Seed(time.Now().Unix())
	user.Salt = fmt.Sprintf("%06d", rand.Int31n(10000))
	user.Passwd = util.MakePasswd(password, user.Salt)
	user.Token = fmt.Sprintf("%08d", rand.Int31())
	ChatUserCol.InsertOne(ctx, &user)
	return user, nil
}

func (u *UserService) Login(mobile, passwd string) (mgomodel.User, error) {
	user := mgomodel.User{}
	ctx := context.Background()
	err := ChatUserCol.FindOne(ctx, bson.D{{Key: "mobile", Value: mobile}}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return user, err
	}
	fmt.Printf("##### %d\n", user.Createat.Unix())
	//校验密码
	if !util.ValidatePasswd(passwd, user.Salt, user.Passwd) {
		return user, errors.New("password not correct")
	}
	return user, nil
}
