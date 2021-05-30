package dbmongo

import (
	"context"
	"errors"
	"fmt"
	"im_webapp/model/mgomodel"
	"math/rand"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ContactService struct {
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// HandleAddFriend
func (c *ContactService) HandleAddFriend(userid, dstMobile string) error {
	user := mgomodel.User{}
	// 根据userid 获取当前对象
	col := MongoClient.Database("chat").Collection("user")
	userObjId, _ := primitive.ObjectIDFromHex(userid)
	col.FindOne(context.Background(), bson.D{{Key: "_id", Value: userObjId}}).Decode(&user)
	if user.Id_ == "" {
		return fmt.Errorf("no have %s user", userid)
	}
	if dstMobile == user.Mobile {
		//不能添加自己
		return errors.New("cant't add yourself")
	}
	dstUser := mgomodel.User{}
	col.FindOne(context.Background(), bson.D{{Key: "mobile", Value: dstMobile}}).Decode(&dstUser)
	if dstUser.Id_ == "" {
		return fmt.Errorf("no have firend mobile:%s", dstMobile)
	}

	contactCol := MongoClient.Database("chat").Collection("contact")
	//判断是否已经添加过好友了
	tmpContact := mgomodel.Contact{}
	contactCol.FindOne(context.Background(), bson.D{{Key: "ownerid", Value: userid}, {Key: "dstobj", Value: dstUser.Id_}}).Decode(&tmpContact)
	fmt.Printf(">>>>>>find contact user=%s dst=%s tmpid=%s\n", userid, dstUser.Id_, tmpContact.Id_)
	if tmpContact.Id_ != "" {
		return fmt.Errorf("%s is alreay added", dstUser.Nickname)
	}

	// 开启事务
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	sessionClient, err := MongoClient.StartSession()
	if err != nil {
		return errors.New("satrt session error")
	}

	sessionClient.StartTransaction()
	_, e1 := contactCol.InsertOne(ctx, &mgomodel.Contact{
		Ownerid:  userid,
		Dstobj:   dstUser.Id_,
		Cate:     mgomodel.CONCAT_CATE_USER,
		Createat: time.Now(),
	})
	_, e2 := contactCol.InsertOne(ctx, &mgomodel.Contact{
		Ownerid:  dstUser.Id_,
		Dstobj:   userid,
		Cate:     mgomodel.CONCAT_CATE_USER,
		Createat: time.Now(),
	})
	if e1 != nil || e2 != nil {
		sessionClient.AbortTransaction(ctx)
	}
	if e1 == nil && e2 == nil {
		sessionClient.CommitTransaction(ctx)
	}

	return nil
}

// 加载好友列表
func (c *ContactService) HandleLoadFriend(userid string) []mgomodel.User {
	//col := MongoClient.Database("chat").Collection("contact")
	filter := bson.D{{Key: "ownerid", Value: userid}, {Key: "cate", Value: mgomodel.CONCAT_CATE_USER}}
	ctx := context.Background()
	cusor, err := ChatContact.Find(ctx, filter)
	if err != nil {
		return []mgomodel.User{}
	}
	friendIds := []primitive.ObjectID{}
	friends := []mgomodel.User{}
	for cusor.Next(ctx) {
		tmpContact := mgomodel.Contact{}
		cusor.Decode(&tmpContact)
		tmpid, _ := primitive.ObjectIDFromHex(tmpContact.Dstobj)
		friendIds = append(friendIds, tmpid)
	}
	fmt.Printf("find friends %v \n", friendIds)
	// 根据dstobj查找
	//userCol := MongoClient.Database("chat").Collection("user")
	//相当于 arr:=make([]interface{}, len(friendIds))
	arr := make(bson.A, len(friendIds))
	for idx, v := range friendIds {
		arr[idx] = v
	}
	filter = bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: arr}}}}
	friendCursor, err := ChatUserCol.Find(ctx, filter)
	if err != nil {
		return []mgomodel.User{}
	}
	for friendCursor.Next(ctx) {
		tmpFriend := mgomodel.User{}
		friendCursor.Decode(&tmpFriend)
		friends = append(friends, tmpFriend)
	}
	return friends
}

func (c *ContactService) HandleCreateCommunity(name, cate, memo, icon, ownerid string) error {
	tmpCommunity := mgomodel.Community{}
	// 查看是否已经存在该群
	ctx := context.Background()
	cateInt, _ := strconv.ParseInt(cate, 10, 64)
	fmt.Printf("#################: cate=%d cates=%s\n", cateInt, cate)
	filter := bson.D{{Key: "ownerid", Value: ownerid}, {Key: "name", Value: name}, {Key: "cate", Value: int(cateInt)}}
	ChatCommunity.FindOne(ctx, filter).Decode(&tmpCommunity)

	fmt.Printf("get community=%v\n", tmpCommunity)
	// 判断是否和传入的参数一致
	if name == tmpCommunity.Name && ownerid == tmpCommunity.Ownerid {
		return errors.New("alreay create this community")
	}
	// 添加群聊
	objId := primitive.NewObjectID().Hex()
	tmpCommunity.Id_ = objId
	tmpCommunity.Name = name
	tmpCommunity.Ownerid = ownerid
	tmpCommunity.Cate = int(cateInt)
	tmpCommunity.Createat = time.Now()
	tmpCommunity.Uuid = fmt.Sprintf("%08d", rand.Int31())[:8]
	tmpCommunity.Memo = memo
	_, err := ChatCommunity.InsertOne(ctx, &tmpCommunity)
	fmt.Printf("err=%v insert id=%v\n", err, tmpCommunity.Id_)
	//创建群聊时应该把自己加入群聊哈
	cont := mgomodel.Contact{
		Ownerid:  ownerid,
		Cate:     mgomodel.CONCAT_CATE_COMUNITY,
		Createat: time.Now(),
		Memo:     memo,
		Dstobj:   tmpCommunity.Id_,
	}
	_, err = ChatContact.InsertOne(ctx, &cont)
	if err != nil {
		return err
	}
	return nil
}
func (c *ContactService) HandleLoadCommunity(userid string) []mgomodel.Community {
	comms := []mgomodel.Community{}
	ctx := context.Background()
	fmt.Printf(">>>>>>>>>>>>>>>>>>>>>: handleLoadCommunity %s \n", userid)
	contCour, err := ChatContact.Find(ctx, bson.D{{Key: "ownerid", Value: userid},
		{Key: "cate", Value: mgomodel.CONCAT_CATE_COMUNITY}})
	if err != nil {
		return comms
	}
	commIds := make(bson.A, 0)
	for contCour.Next(ctx) {
		tmpComm := mgomodel.Contact{}
		contCour.Decode(&tmpComm)
		//objId, _ := primitive.ObjectIDFromHex(tmpComm.Dstobj)
		commIds = append(commIds, tmpComm.Dstobj)
	}
	fmt.Printf("find communitys %v \n", commIds)
	commsCoursour, err := ChatCommunity.Find(ctx, bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: commIds}}}})
	if err != nil {
		return comms
	}
	for commsCoursour.Next(ctx) {
		tmpComm := mgomodel.Community{}
		commsCoursour.Decode(&tmpComm)
		comms = append(comms, tmpComm)
	}

	return comms
}

// dstobj 此处是群聊的uuid值
func (C *ContactService) HandleJoinCommunity(userid, dstobj string) error {
	//检查dstobj是否存在
	tmpComm := mgomodel.Community{}
	filter := bson.D{{Key: "uuid", Value: dstobj}}
	ChatCommunity.FindOne(context.Background(), filter).Decode(&tmpComm)
	if tmpComm.Id_ == "" {
		return errors.New("no have this " + dstobj + " group")
	}
	// 创建contact连接
	fmt.Printf("find community %s \n", tmpComm.Uuid)
	cont := mgomodel.Contact{}
	cont.Ownerid = userid
	cont.Dstobj = tmpComm.Id_
	cont.Cate = mgomodel.CONCAT_CATE_COMUNITY
	cont.Createat = time.Now()
	cont.Memo = tmpComm.Memo
	_, err := ChatContact.InsertOne(context.Background(), &cont)
	if err != nil {
		return err
	}
	return nil
}

func (c *ContactService) HandleSearchCommunities(userid string) []string {
	result := []string{}
	filter := bson.D{{Key: "ownerid", Value: userid}, {Key: "cate", Value: mgomodel.CONCAT_CATE_COMUNITY}}
	cousor, _ := ChatContact.Find(context.Background(), filter)
	for cousor.Next(context.Background()) {
		tmp := mgomodel.Contact{}
		cousor.Decode(&tmp)
		result = append(result, tmp.Dstobj)
	}
	return result
}
