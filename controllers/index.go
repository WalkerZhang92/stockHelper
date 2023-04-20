package controllers

import (
	"encoding/json"
	"net/http"
	"quickstart/models"

	"github.com/astaxie/beego"
)

type IndexController struct {
	beego.Controller
}

func (c *IndexController) Get() {
	secid := c.Ctx.Input.Query("secid")
	var secidStr string
	if secid == "" {
		secidStr = "1.600001"
	} else {
		firstStr := string(secid[0])

		switch firstStr {
		case "6":
			secidStr = "1." + string(secid)
		default:
			secidStr = "0." + string(secid)
		}
	}

	url := "http://push2his.eastmoney.com/api/qt/stock/kline/get?fields1=f1,f2,f3,f4,f5,f6&fields2=f51,f52,f53,f54,f55,f56,f57&klt=101&fqt=1&secid=" + secidStr + "&beg=19980000&end=20500000&_=1680767375018"
	resp, err := http.Get(url)
	if err != nil {
		c.Data["json"] = &models.Response{Code: 401, Message: err.Error()}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()
	// 解码 JSON 响应数据
	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		c.Data["json"] = &models.Response{Code: 401, Message: err.Error()}
		c.ServeJSON()
		return
	}

	// 将数据作为 JSON 响应发送回客户端
	result := data["data"]
	c.Data["Code"] = result.(map[string]interface{})["code"]
	c.Data["Name"] = result.(map[string]interface{})["name"]
	c.Data["Klines"] = result.(map[string]interface{})["klines"]
	c.TplName = "index.html"
	err = c.Render()
	if err != nil {
		println(err)
	}
}
