package controllers

import "github.com/astaxie/beego"

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	c.TplName = "login.html"
	err := c.Render()
	if err != nil {
		println(err)
	}
}

func (c *LoginController) Login() (result int32) {
	return 111
}
