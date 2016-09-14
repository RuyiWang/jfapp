package routers

import (
	"jfapp/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/jfapp", &controllers.MainController{}, "get,post:Index")
	beego.Router("/jfapp/bankList", &controllers.AppServerController{}, "get,post:BankList")
	beego.Router("/jfapp/historyList", &controllers.AppServerController{}, "get,post:HistoryList")
	beego.Router("/jfapp/basInfo", &controllers.AppServerController{}, "get,post:BasInfo")
	beego.Router("/jfapp/monthBasInfo", &controllers.AppServerController{}, "get,post:MonthBasInfo")
}
