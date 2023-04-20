package main

import (
	_ "quickstart/routers"
	"quickstart/tasks"

	"github.com/astaxie/beego/orm"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
)

func init() {

	//数据库连接
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:root@/stockmarket?charset=utf8")

	//初始化定时任务
	tasks.StartCron()
}
func main() {
	beego.BConfig.WebConfig.AutoRender = false //关闭自动渲染

	beego.BConfig.CopyRequestBody = true

	beego.SetStaticPath("/images", "img")
	beego.SetStaticPath("/css", "css")
	beego.SetStaticPath("/js", "js")
	beego.Run()
}
