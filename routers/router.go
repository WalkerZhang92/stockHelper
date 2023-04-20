package routers

import (
	"quickstart/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/index", &controllers.IndexController{})
	beego.Router("/login", &controllers.LoginController{})

	beego.Router("/watchConfig", &controllers.WatchStockController{})
	beego.Router("/watchConfig/detail", &controllers.WatchStockController{}, "get:Detail")
	beego.Router("/saveWatchStock", &controllers.WatchStockController{}, "post:Save")
	beego.Router("/watchConfig/delete", &controllers.WatchStockController{}, "post:Delete")

	beego.Router("/login/login", &controllers.ApiLoginController{}, "post:Login")
	beego.Router("/login/setStockDaily", &controllers.ApiLoginController{}, "post:SetStockDaily")

	beego.Router("/login/task", &controllers.ApiLoginController{}, "post:Task")
	beego.Router("/login/recommend", &controllers.ApiLoginController{}, "post:Recommend")
	beego.Router("/login/stockRealTime", &controllers.ApiLoginController{}, "post:StockRealTime")

	beego.Router("/searchStock", &controllers.ApiSearchStockController{}, "post:SearchStock")

}
