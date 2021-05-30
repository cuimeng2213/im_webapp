package service

import (
	"context"
	"encoding/json"
	"im_webapp/db/dbmongo"
	"im_webapp/model/mgomodel"
	"net"
	"runtime"
	"sync"

	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/fatih/set.v0"
)

const (
	CMD_SINGLE_CMD = 10
	CMD_ROOM_MSG   = 11
	CMD_HEART      = 0
	CMD_ACK        = 1
	CMD_ENTRY_ROOM = 2
	CMD_EXIT_ROOM  = 3
)
const (
	MEDIA_TYPE_TEXT       = 1
	MEDIA_TYPE_NEWS       = 2
	MEDIA_TYPE_VOICE      = 3
	MEDIA_TYPE_IMG        = 4
	MEDIA_TYPE_REDPACKAGE = 5
	MEDIA_TYPE_EMOJ       = 6
)

type ChatService struct {
}

var udpSendChan chan []byte = make(chan []byte, 1024)

type WsClient struct {
	Conn *websocket.Conn
	// 并行数据转化为串行数据
	DataQueue chan []byte
	UserId    string
	GroupSet  set.Interface //储存所有群id
}

type Message struct {
	Id string `json:"id,omitempty" form:"id"`
	//谁发送的消息
	Userid string `json:"userid,omitempty" form:"userid"`
	Cmd    int    `json:"cmd,omitempty" form:"cmd"`
	//要发送给谁
	Dstid string `json:"dstid,omitempty" form:"dstid"`
	// 数据类型
	Media   int    `json:"media,omitempty" form:"media"`
	Content string `json:"content,omitempty" form:"content"`
	Url     string `json:"url,omitempty" form:"url"`
	Amount  int    `json:"amount,omitempty" form:"amount"`
}

//全局连接对象
var Clients = make(map[string]*WsClient)

var rwlocker sync.RWMutex
var cont dbmongo.ContactService

func CheckToken(userid string, token string) bool {
	//校验token是否合法
	user := mgomodel.User{}
	objId, _ := primitive.ObjectIDFromHex(userid)
	filter := bson.D{{Key: "_id", Value: objId}}
	dbmongo.ChatUserCol.FindOne(context.Background(), filter).Decode(&user)

	return user.Token == token
}
func CheckTokenByMgo(userid, token string) bool {
	user := mgomodel.User{}
	coll := dbmongo.MongoClient.Database("chat").Collection("user")
	//通过_id查询数据时需要将userid转换为ObjectID
	filterId, _ := primitive.ObjectIDFromHex(userid)
	filter := bson.D{{Key: "_id", Value: filterId}}
	err := coll.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		fmt.Printf(">>>>>>>CheckTookenByMgo error %v \n", err)
		return false
	}
	fmt.Printf("token= %s user.Token=%s\n", token, user.Token)
	return token == user.Token
}

func broadcast(d []byte) {
	udpSendChan <- d
}
func sendproc(client *WsClient) {
	defer func() {
		fmt.Println("sendproc exit>>>>>>> ", client.UserId)
	}()
	for {
		select {
		case data := <-client.DataQueue:
			err := client.Conn.WriteMessage(websocket.TextMessage, data)
			//broadcast(data)
			if err != nil {
				fmt.Println(">>>>>>> sendproc error : ", err)
				return
			}
		}
	}
}
func recvproc(client *WsClient) {
	defer func() {
		fmt.Println("recvproc exit>>>>>>> ", client.UserId)
	}()

	for {
		_, data, err := client.Conn.ReadMessage()
		if err != nil {
			fmt.Println(">>>>> recvproc error ", err)
			return
		}
		fmt.Println(">>>> recvproc data: ", string(data))
		//dispatch(data)
		broadcast(data)
	}
}

func sendmsg(userid string, msg []byte) {
	rwlocker.RLock()
	client, ok := Clients[userid]
	rwlocker.RUnlock()
	if ok {
		client.DataQueue <- msg
	}
}
func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	_, fileName, line, _ := runtime.Caller(1)
	fmt.Printf("fileName=[%s] line=%d \n", fileName, line)
	if err != nil {
		fmt.Println(">>>> unmarshal error: ", err)
		return
	}
	switch msg.Cmd {
	case CMD_SINGLE_CMD:
		sendmsg(msg.Dstid, data)
	case CMD_ROOM_MSG:
		rwlocker.RLock()
		//TODO 群聊转发
		//根据dstid找到群内所有用户
		for userId, v := range Clients {
			fmt.Printf("userid=%s msgUserid=%s has=%t\n", userId, msg.Userid, v.GroupSet.Has(msg.Dstid))
			if userId != msg.Userid && v.GroupSet.Has(msg.Dstid) {
				fmt.Printf("sendto group %s\n", msg.Dstid)
				v.DataQueue <- data
			}
		}
		rwlocker.RUnlock()
	case CMD_HEART:

	}
}
func (c *ChatService) Upgrade(w http.ResponseWriter, r *http.Request, userid, token string) error {

	isOK := CheckTokenByMgo(userid, token)

	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isOK
		},
	}).Upgrade(w, r, nil)

	if err != nil {
		fmt.Printf("#### err=%v\n", err)
		return err
	}

	client := &WsClient{
		Conn:      conn,
		DataQueue: make(chan []byte),
		UserId:    userid,
		GroupSet:  set.New(set.ThreadSafe),
	}

	//TODO 获取用户全部群ID
	ids := cont.HandleSearchCommunities(userid)
	for _, v := range ids {
		client.GroupSet.Add(v)
	}
	rwlocker.Lock()
	Clients[userid] = client
	rwlocker.Unlock()

	go sendproc(client)
	go recvproc(client)
	fmt.Printf(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>: handle %s\n", userid)
	return nil
}
func udpSendProc() {
	// 1.创建连接句柄
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 0, 255),
		Port: 3000,
	})
	if err != nil {
		fmt.Printf("DialUDP error: %v\n", err)
		return
	}
	defer conn.Close()
	for {
		data, _ := <-udpSendChan
		_, err = conn.Write(data)
		if err != nil {
			fmt.Printf("udp write data error %v\n", err)
			return
		}
	}
}
func udpRevProc() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		fmt.Printf("listen udp error %v\n", err)
		return
	}
	defer conn.Close()
	buf := [1024]byte{}
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			fmt.Printf("duprecvproc read error %v\n", err)
			return
		}
		dispatch(buf[:n])
	}
}
func init() {
	go udpSendProc()
	go udpRevProc()
}
