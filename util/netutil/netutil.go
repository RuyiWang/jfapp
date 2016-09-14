package netutil

import (
	"log"
	"net/http"
	"io/ioutil"
	"time"
)

import(
	"github.com/astaxie/beego/httplib"	
)


//请求接口
func CallInterface(url, data string) string {
	var resp *http.Response
	var err error
	
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Runtime error caught: %v", r)
		} else {
			
		}
	}()
	
	req := httplib.Post(url)
	req.SetTimeout(2000 * time.Millisecond, 2000 * time.Millisecond)
	req.Param("data", data)
	resp, err = req.Response()
	if err != nil {
		return ""
	} else {
		if resp.StatusCode == 200 {
			res,_ := ioutil.ReadAll(resp.Body)
			return string(res)
		} else {
			return ""
		}	
	}
	
}