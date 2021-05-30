## 原理
前端通过websocket发送 json数据

### 核心数据结构
```go
type Message struct {
    Id int64 `json:"id,omitempty" form:"id"`
    //谁发送的消息
    Userid int64 `json:"userid,omitempty" form:"userid"`
    Cmd int `json:"cmd,omitempty" form:"cmd"`
    //要发送给谁
    Dstid int64 `json:"dstid,omitempty" form:"dstid"`
    // 数据类型
    Media int `json:"media,omitempty" form:"media"`
    Content string 
}
cosnt (
    CMD_SINGLE_CMD = 10
    CMD_ROOM_MSG = 11
    CMD_HEART = 0
    CMD_ACK = 1
    CMD_ENTRY_ROOM = 2
    CMD_EXIT_ROOM =3
)
cosnt (
    MEDIA_TYPE_TEXT = 1
    MEDIA_TYPE_NEWS = 2
    MEDIA_TYPE_VOICE=3
    MEDIA_TYPE_IMG=4
    MEDIA_TYPE_REDPACKAGE=5
    MEDIA_TYPE_EMOJ=6

)
```
/**
1.MEDIA_TYPE_TEXT
{id:1, userid:2,dstid:3,cmd:10,media:1,content:"hello"}

2. MEDIA_TYPE_VOICE
{id:1,userid:2,dstid:3,cmd:10,media:3,url:"http://xxxx.com/dsturl.mp3", amount:15}
**/

## 实现发送文字 表情包
前端user1拼接好数据对象Message
msg={id:1,userid:2,dstid:3,cmd:10,media:1,context="fuck"}
转化成jsonstr
jsonstr = JSON.stringify(msg)
通过websocket.send(jsonstr) 发送数据
后端server在recvproc中接收数据data
并做相应得逻辑处理dispatch(data) -- 转发给user2
user2通过websocket.onmessage收到消息后解析并展示

```sql
update user set name="lucy" where id=7;
```

实现群聊
分析群ID，找到加了这个群的用户，把消息发出去

方案-、
map<userid><qunid1,qunid2>
又是是锁的频次低
劣势是要轮询全部map

方案2、
map<qunid><uid1，uid2>

### 需要处理的问题
JavaScript
1. 当用户接入的时候初始化groupset
2. 当用户加入群的时候刷新groupset
3. 完成信息分发
