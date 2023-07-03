package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"quickstart/common"
	"quickstart/models"
	"quickstart/services"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"

	"github.com/astaxie/beego"
)

type ApiLoginController struct {
	beego.Controller
}
type loginForm struct {
	Phone    int64
	Password string
}

func (c *ApiLoginController) Login() {
	var loginF loginForm
	var err error
	err = json.Unmarshal(c.Ctx.Input.RequestBody, &loginF)
	if err != nil {
		c.Data["json"] = &models.Response{Code: 400, Message: "请求参数格式错误" + err.Error()}
		c.ServeJSON()
		return
	}
	userRes := models.User{Phone: loginF.Phone}
	o := orm.NewOrm()

	err = o.Read(&userRes, "Phone")
	encodePass := common.MD5(loginF.Password)
	if err == orm.ErrNoRows || userRes.Password != encodePass {
		c.Data["json"] = &models.Response{Code: 401, Message: "用户名或密码错误"}
		c.ServeJSON()
		return
	}
	c.Data["json"] = &models.Response{Code: 200, Message: "登录成功"}
	c.ServeJSON()
}

func (this *ApiLoginController) SetStockDaily() {
	redisService := services.NewRedisService()
	defer redisService.Close()

	// 在需要时使用redisService执行Redis命令
	result, err := redisService.Do("SET", "mykey", "myvalue")
	if err != nil {
		this.Data["json"] = &models.Response{Code: 401, Message: "redis 保存失败"}
		this.ServeJSON()
		return
	}
	this.Data["json"] = &models.Response{Code: 200, Message: "保存成功", Data: result}
	this.ServeJSON()
}

func (this *ApiLoginController) Task() {
	url := "https://data.eastmoney.com/dataapi/xuangu/list?st=CHANGE_RATE&sr=-1&sty=SECUCODE,SECURITY_CODE,SECURITY_NAME_ABBR,NEW_PRICE,CHANGE_RATE,VOLUME_RATIO,HIGH_PRICE,LOW_PRICE,PRE_CLOSE_PRICE,VOLUME,DEAL_AMOUNT,TURNOVERRATE,PE9,BASIC_EPS,ROE_WEIGHT,NETPROFIT_YOY_RATIO,TOI_YOY_RATIO&source=SELECT_SECURITIES&client=WEB&ps=50&p=1"
	resp, err := http.Get(url)
	if err != nil {
		this.Data["json"] = &models.Response{Code: 401, Message: err.Error()}
		this.ServeJSON()
		return
	}
	defer resp.Body.Close()
	// 解码 JSON 响应数据
	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		this.Data["json"] = &models.Response{Code: 401, Message: err.Error()}
		this.ServeJSON()
		return
	}

	// 将数据作为 JSON 响应发送回客户端
	this.Data["json"] = &data
	this.ServeJSON()
}

func (this *ApiLoginController) Recommend() {
	url := "https://data.eastmoney.com/dataapi/xuangu/list?st=CHANGE_RATE&sr=-1&ps=50&p=1&sty=SECUCODE%2CSECURITY_CODE%2CSECURITY_NAME_ABBR%2CNEW_PRICE%2CCHANGE_RATE%2CVOLUME_RATIO%2CHIGH_PRICE%2CLOW_PRICE%2CPRE_CLOSE_PRICE%2CVOLUME%2CDEAL_AMOUNT%2CTURNOVERRATE%2CKDJ_GOLDEN_FORK%2CPOWER_FULGUN&filter=(KDJ_GOLDEN_FORKZ%3D%221%22)(POWER_FULGUN%3D%221%22)&source=SELECT_SECURITIES&client=WEB"
	resp, err := http.Get(url)
	if err != nil {
		this.Data["json"] = &models.Response{Code: 401, Message: err.Error()}
		this.ServeJSON()
		return
	}
	defer resp.Body.Close()
	// 解码 JSON 响应数据
	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		this.Data["json"] = &models.Response{Code: 401, Message: err.Error()}
		this.ServeJSON()
		return
	}
	result := data["result"]
	list := result.(map[string]interface{})["data"].([]interface{})

	now := time.Now()
	dateFormat := "20060102"
	dateString := now.Format(dateFormat)
	o := orm.NewOrm()
	for _, item := range list {
		stock := item.(map[string]interface{})
		stockSelect := &models.StockSelect{}
		stockSelect.User = 1
		stockSelect.Code = stock["SECURITY_CODE"].(string)
		stockSelect.Title = stock["SECURITY_NAME_ABBR"].(string)
		stockSelect.Date = dateString
		o.Insert(stockSelect)
		services.SendDingMsg("股票提醒：" + stockSelect.Code + "  " + stockSelect.Title + "   这只股票不错哦！")

	}

	// 将数据作为 JSON 响应发送回客户端
	this.Data["json"] = &data

	this.ServeJSON()
}

func (this *ApiLoginController) SaveSectorFlow()  {
	url := "https://push2.eastmoney.com/api/qt/clist/get?fid=f62&po=1&pz=500&pn=1&np=1&fltt=2&invt=2&fs=m:90+t:2&fields=f12,f14,f2,f3,f62,f184,f66,f69,f72,f75,f78,f81,f84,f87,f204,f205,f124,f1,f13"
	resp, err := http.Get(url)
	if err != nil {
		this.Data["json"] = &models.Response{Code: 401, Message: err.Error()}
		this.ServeJSON()
		return
	}
	defer resp.Body.Close()
	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	result := data["data"]
	if err != nil {
		this.Data["json"] = &models.Response{Code: 401, Message: err.Error()}
		this.ServeJSON()
		return
	}
	diff := result.(map[string]interface{})["diff"].([]interface{})
	fmt.Println(diff)
	now := time.Now()
	dateFormat := "20060102"
	dateString := now.Format(dateFormat)
	o := orm.NewOrm()

	for _, item := range diff {
		sectorFlow := &models.SectorFlow{}
		sector := item.(map[string]interface{})
		//板块编码
		if _, ok := sector["f12"].(string); !ok {
			sectorFlow.SectorCode = ""
		} else {
			sectorFlow.SectorCode = sector["f12"].(string)
		}

		//板块名称
		if _, ok := sector["f14"].(string); !ok {
			sectorFlow.SectorName = ""
		} else {
			sectorFlow.SectorName = sector["f14"].(string)
		}

		//流入最多的股
		if _, ok := sector["f205"].(string); !ok {
			sectorFlow.BestStock = ""
		} else {
			sectorFlow.BestStock = sector["f205"].(string)
		}
		//流入最多的股名称
		if _, ok := sector["f204"].(string); !ok {
			sectorFlow.BestStockName = ""
		} else {
			sectorFlow.BestStockName = sector["f204"].(string)
		}

		//流入金额
		if _, ok := sector["f62"]; !ok {
			sectorFlow.FlowIn = 0
		} else {
			sectorFlow.FlowIn = sector["f62"].(float64)
		}
		num, _ := strconv.Atoi(dateString)
		sectorFlow.Date = num

		newId, err := o.Insert(sectorFlow)
		fmt.Println(newId)
		if err == nil {
			this.Data["json"] = &models.Response{Code: 200, Message: "插入成功"}
			this.ServeJSON()
		} else {
			fmt.Println(err)
			this.Data["json"] = &models.Response{Code: 500, Message: "插入失败"}
			this.ServeJSON()
		}
	}

	// 将数据作为 JSON 响应发送回客户端
	this.Data["json"] = &data

	this.ServeJSON()
}
