package controllers

import (
	"encoding/json"
	"net/http"
	"quickstart/models"

	"github.com/astaxie/beego"
)

type ApiSearchStockController struct {
	beego.Controller
}

type StockParam struct {
	Input string
}

func (c *ApiSearchStockController) SearchStock() {
	var stockP StockParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &stockP)

	if err != nil {
		c.Data["json"] = &models.Response{Code: 400, Message: "请求参数格式错误" + err.Error()}
		c.ServeJSON()
		return
	}

	url := "https://searchadapter.eastmoney.com/api/suggest/get?type=8&count=6&input=" + stockP.Input
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
	c.Data["json"] = &data
	c.ServeJSON()
}
