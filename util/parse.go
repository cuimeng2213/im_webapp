package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Bind(obj interface{}, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	fmt.Printf("request content-type %s\n", contentType)
	if contentType == "application/json" {
		bindJson(obj, r)
	}
}

func bindJson(obj interface{}, r *http.Request) {
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, obj)
	if err != nil {
		fmt.Printf("json.Unmarshal error %s\n", err.Error())
		return
	}
}
