package controllers

import (
	"encoding/json"
	"fmt"

	"quickstart/models"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type WatchStockController struct {
	beego.Controller
}

type WatchStockForm struct {
	Code string
	Name string
	Type string
	Min  string
	Max  string
	Id   int
}

type WatchStockDelete struct {
	Id int
}

func (c *WatchStockController) Get() {
	c.TplName = "watch-config.html"
	o := orm.NewOrm()
	var stocks []*models.WatchStock
	_, err := o.QueryTable("watch_stock").Filter("is_del", 0).All(&stocks)
	if err != nil {
		println(err)
	}

	c.Data["Stocks"] = stocks

	err = c.Render()
	if err != nil {
		println(err)
	}
}

func (c *WatchStockController) Detail() {
	watchId := c.Ctx.Input.Query("id")
	watchIdInt, err := strconv.ParseInt(watchId, 0, 32)
	if err != nil {
		fmt.Println(err)
	}
	o := orm.NewOrm()
	watchStock := models.WatchStock{Id: int(watchIdInt)}
	err = o.Read(&watchStock)
	if err != nil {
		fmt.Println(err)
	}
	var result = make(map[string]interface{})
	result["code"] = 200
	result["data"] = watchStock
	c.Data["json"] = result

	c.ServeJSON()
}

func (c *WatchStockController) Save() {
	var watchStockForm WatchStockForm
	var err error
	err = json.Unmarshal(c.Ctx.Input.RequestBody, &watchStockForm)
	if err != nil {
		c.Data["json"] = &models.Response{Code: 400, Message: "请求参数格式错误" + err.Error()}
		c.ServeJSON()
		return
	}
	o := orm.NewOrm()

	if watchStockForm.Id != 0 {
		var watchStock = models.WatchStock{Id: watchStockForm.Id}
		watchStock.Code = watchStockForm.Code
		watchStock.Name = watchStockForm.Name

		stockType, _ := strconv.ParseInt(watchStockForm.Type, 10, 32)
		min, _ := strconv.ParseFloat(watchStockForm.Min, 64)

		max, _ := strconv.ParseFloat(watchStockForm.Max, 64)
		watchStock.Type = int(stockType)
		watchStock.Min = min
		watchStock.Max = max
		if num, err := o.Update(&watchStock); err == nil {
			fmt.Println(num)
		}
		if err == nil {
			c.Data["json"] = &models.Response{Code: 200, Message: "修改成功"}
			c.ServeJSON()
		} else {
			c.Data["json"] = &models.Response{Code: 500, Message: "修改失败"}
			c.ServeJSON()
		}
	} else {
		var watchStock models.WatchStock
		watchStock.Code = watchStockForm.Code
		watchStock.Name = watchStockForm.Name

		stockType, _ := strconv.ParseInt(watchStockForm.Type, 10, 32)
		min, _ := strconv.ParseFloat(watchStockForm.Min, 64)

		max, _ := strconv.ParseFloat(watchStockForm.Max, 64)
		watchStock.Type = int(stockType)
		watchStock.Min = min
		watchStock.Max = max
		newId, err := o.Insert(&watchStock)
		fmt.Println(newId)
		if err == nil {
			c.Data["json"] = &models.Response{Code: 200, Message: "插入成功"}
			c.ServeJSON()
		} else {
			c.Data["json"] = &models.Response{Code: 500, Message: "插入失败"}
			c.ServeJSON()
		}
	}

}

func (c *WatchStockController) Delete() {
	var watchStockDelete WatchStockDelete
	var err error
	err = json.Unmarshal(c.Ctx.Input.RequestBody, &watchStockDelete)
	if err != nil {
		c.Data["json"] = &models.Response{Code: 400, Message: "请求参数格式错误" + err.Error()}
		c.ServeJSON()
		return
	}
	o := orm.NewOrm()
	var watchStock = models.WatchStock{Id: watchStockDelete.Id}
	watchStock.IsDel = 1
	if num, err := o.Update(&watchStock, "isDel"); err == nil {
		fmt.Println(num)
	}
	if err == nil {
		c.Data["json"] = &models.Response{Code: 200, Message: "删除成功"}
		c.ServeJSON()
	} else {
		c.Data["json"] = &models.Response{Code: 500, Message: "删除失败"}
		c.ServeJSON()
	}
}
