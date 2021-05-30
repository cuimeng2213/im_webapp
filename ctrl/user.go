package ctrl

import (
	"im_webapp/db/dbmongo"
	"im_webapp/util"
	"net/http"
)

var UserServerMgo *dbmongo.UserService

// curl -v -X POST localhost:9090/user/register -d"mobile=13011259131&password=123456&nickname=xiaofang&avatar=ava.png&sex=M"
func HandleUserRegister(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//获取表单信息
	mobile := r.PostForm.Get("mobile")
	passwd := r.PostForm.Get("password")
	nickName := r.PostForm.Get("nickname")
	sex := r.PostForm.Get("sex")
	avatar := r.PostForm.Get("avatar")

	_, err := UserServerMgo.Register(mobile, passwd, nickName, avatar, sex)

	if err != nil {
		//返回失败信息
		util.RespFail(w, err.Error())
	} else {
		util.RespOk(w, "", "注册成功")
	}
}

func HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	mobile := r.PostForm.Get("mobile")
	passwd := r.PostForm.Get("password")

	user, err := UserServerMgo.Login(mobile, passwd)
	if err != nil {
		util.RespFail(w, err.Error())
		return
	}
	util.RespOk(w, user, "登录成功")
}
