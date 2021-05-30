package ctrl

import (
	"fmt"
	"im_webapp/db/dbmongo"
	"im_webapp/util"
	"net/http"
)

var MgoContactServer dbmongo.ContactService

func HandleAddFriend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.RespFail(w, fmt.Sprintf("not support method:%s", r.Method))
		return
	}
	r.ParseForm()
	//获取userid dstobj
	userid := r.PostForm.Get("userid")
	dstobj := r.PostForm.Get("dstid")
	if userid == "" || dstobj == "" {
		util.RespFail(w, "userid or dstid param invalid")
		return
	}

	fmt.Printf("HandleAddFriend userid=%s dstobj=%s \n", userid, dstobj)
	err := MgoContactServer.HandleAddFriend(userid, dstobj)
	if err != nil {
		util.RespFail(w, err.Error())
		return
	}
	util.RespOk(w, "", "add success")
}

func HandleLoadFriend(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	userid := r.PostForm.Get("userid")
	fmt.Printf(">>>>> HandleLoadFriend %s \n", userid)
	us := MgoContactServer.HandleLoadFriend(userid)
	util.RespOk(w, us, "")
}

/*type Args struct {
	userid int64 `json:"userid,omitempty"`
	dstid  int64 `json:"dstid,omitempty"`
}*/
func HandleCreateCommunity(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	/*
		name: 你好
		cate: 1
		memo: 测试用的
		icon:
		ownerid: 6
	*/
	name := r.Form.Get("name")
	cate := r.Form.Get("cate")
	memo := r.Form.Get("memo")
	icon := r.Form.Get("icon")
	ownerid := r.Form.Get("ownerid")
	err := MgoContactServer.HandleCreateCommunity(name, cate, memo, icon, ownerid)
	if err != nil {
		util.RespFail(w, err.Error())
		return
	}

	util.RespOk(w, "", "Create Community Success")

}
func HandleLoadCommunity(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userid := r.Form.Get("userid")
	fmt.Printf("%s formvalue=%s\n", r.Form.Get("userid"), r.FormValue("userid"))

	communites := MgoContactServer.HandleLoadCommunity(userid)

	util.RespOk(w, communites, "find communites success")
}

func HandleJoinCommunity(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	/*
		userid=6   //当前用户id
		dstid=2    //群聊id
	*/
	userid := r.FormValue("userid")
	dstid := r.FormValue("dstid")
	err := MgoContactServer.HandleJoinCommunity(userid, dstid)
	if err != nil {
		util.RespFail(w, err.Error())
		return
	}
	util.RespOk(w, "", "Join Success")
}
