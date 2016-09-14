package main

import (
	_ "jfapp/routers"

	"github.com/astaxie/beego"
	// "github.com/astaxie/beego/logs"
)

func main() {
	beego.Run()
	constant.Testt()
}
