package ctrl

import (
	"fmt"
	"im_webapp/service"
	"im_webapp/util"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

var ChatServer service.ChatService

// HandleChat ws://localhost:9090/chat?id=xxx&token=xxxx
func HandleChat(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	token := r.URL.Query().Get("token")
	fmt.Printf(">>>>>: HandleChat id=%s token=%s \n", id, token)
	//创建websocket
	err := ChatServer.Upgrade(w, r, id, token)
	fmt.Printf(">>>>> error=%v \n", err)

}

func HandleAttachUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(4096)
	//filetype=.mp3
	fileType := r.FormValue("filetype")
	fileName := ""
	fileUrl := ""
	f, fHeader, err := r.FormFile("file")
	if err != nil {
		fmt.Printf("HandleUpload file error %s\n", err.Error())
		util.RespFail(w, err.Error())
		return
	}
	//如果文件名字有后缀则不需要添加了
	if strings.LastIndex(fHeader.Filename, ".") != -1 {
		pos := strings.LastIndex(fHeader.Filename, ".")
		fileName = fmt.Sprintf("%s%d.%s", fHeader.Filename[0:pos], rand.Int31(), fHeader.Filename[pos+1:])
	} else {
		//拼接fileurl
		fileName = fmt.Sprintf("%s%d%s", fHeader.Filename, rand.Int31(), fileType)
	}

	fileUrl = fmt.Sprintf("http://%s/mnt/%s", r.Host, fileName)
	//保存资源到本地文件夹
	dstFile, err := os.Create("./mnt/" + fileName)
	if err != nil {
		fmt.Println("open file error ", err)
		util.RespFail(w, err.Error())
		return
	}
	_, err = io.Copy(dstFile, f)
	if err != nil {
		fmt.Println("copy file error ", err)
		util.RespFail(w, err.Error())
		return
	}
	dstFile.Close()
	f.Close()
	util.RespOk(w, fileUrl, "upload success")
}
