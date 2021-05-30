package main

import (
	"encoding/json"
	"fmt"
	"im_webapp/ctrl"
	"log"
	"net/http"
	"text/template"
)

func userLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	mobile := r.PostForm.Get("mobile")
	passwd := r.PostForm.Get("passwd")
	loginOk := false
	fmt.Printf("%s %s \n", mobile, passwd)
	if (mobile == "186123") && (passwd == "123456") {
		loginOk = true
	}
	/*if loginOk {
		str := `{"code":0,"data":{"id":1,"token":"test"}}`
		//设置header为JSON ， 默认为text/html, 需要修改为 application/json
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]bye(str))
		// 返回成功json
	} else {
		//失败的json
		str := `{"code":-1,"msg":"密码不正确"}`
		//设置header为JSON ， 默认为text/html, 需要修改为 application/json
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]bye(str))
	}*/
	//代码优化
	/*str := `{"code":0,"data":{"id":1,"token":"test"}}`
	if !loginOk {
		str = `{"code":-1,"msg":"密码不正确"}`
	}*/
	//设置header为JSON ， 默认为text/html, 需要修改为 application/json
	/*w.Header().Set("Content-Type","application/json") //封装
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(str))*/
	if loginOk {
		data := make(map[string]interface{})
		data["id"] = 12
		data["token"] = "asdc213"
		Resp(w, 0, data, "")
	} else {
		Resp(w, -1, nil, "密码错误")
	}

}
func userRegister(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//获取表单信息
	username := r.PostForm.Get("username")
	plainpasswd := r.PostForm.Get("passwd")
	sex := r.PostForm.Get("sex")
	fmt.Printf("%s %s %s\n", username, plainpasswd, sex)
}

type H struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func Resp(w http.ResponseWriter, code int, data interface{}, msg string) {
	//设置header为JSON ， 默认为text/html, 需要修改为 application/json
	w.Header().Set("Content-Type", "application/json") //封装
	w.WriteHeader(http.StatusOK)

	// 定义一个结构体
	h := H{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	//
	ret, err := json.Marshal(h)
	fmt.Printf("Marshal:[%s]\n", string(ret))
	if err != nil {
		fmt.Printf(">>> %v\n", err)
	}

	w.Write(ret)
}

func RegisterView() {
	tpl, err := template.ParseGlob("./view/**/*")
	if err != nil {
		log.Fatal("ParseGlob failed")
		return
	}
	for _, v := range tpl.Templates() {
		tplName := v.Name()
		http.HandleFunc(tplName, func(w http.ResponseWriter, r *http.Request) {
			v.ExecuteTemplate(w, tplName, nil)
		})
	}
}

func main() {
	http.HandleFunc("/user/login", ctrl.HandleUserLogin)
	http.HandleFunc("/user/register", ctrl.HandleUserRegister)

	//
	http.HandleFunc("/contact/addfriend", ctrl.HandleAddFriend)
	http.HandleFunc("/contact/loadfriend", ctrl.HandleLoadFriend)
	//群聊
	http.HandleFunc("/contact/loadcommunity", ctrl.HandleLoadCommunity)
	http.HandleFunc("/contact/createcommunity", ctrl.HandleCreateCommunity)
	http.HandleFunc("/contact/joincommunity", ctrl.HandleJoinCommunity)

	//chat
	http.HandleFunc("/chat", ctrl.HandleChat)

	//file upload
	http.HandleFunc("/attach/upload", ctrl.HandleAttachUpload)

	//1. 提供静态资源文件支持
	//http.Handle("/", http.FileServer(http.Dir(".")))
	//2. 提供指定目录的静态文件支持
	http.Handle("/asset/", http.FileServer(http.Dir(".")))
	http.Handle("/mnt/", http.FileServer(http.Dir(".")))

	//偷懒初始化template
	RegisterView()

	fmt.Println("start server")
	if err := http.ListenAndServe("127.0.0.1:9090", nil); err != nil {
		panic("start failed")
	}
}
