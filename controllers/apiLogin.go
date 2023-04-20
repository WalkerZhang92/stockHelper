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

func (c *ApiLoginController) StockRealTime() {
	var stocks []*models.WatchStock
	pageSize := "10"
	o := orm.NewOrm()
	num, err := o.QueryTable("WatchStock").Filter("is_del", 0).All(&stocks)
	fmt.Printf("Returned Rows Num: %s, %s", num, err)
	redisService := services.NewRedisService()
	defer redisService.Close()

	var secidAllStr string
	var secidStr string
	for index, stock := range stocks {
		firstStr := string(stock.Code[0])

		switch firstStr {
		case "6":
			secidStr = "1." + string(stock.Code)
		default:
			secidStr = "0." + string(stock.Code)
		}
		if num >= 1 {
			if index+1 == int(num) {
				secidAllStr += secidStr
			} else {
				secidAllStr += secidStr + ","
			}

		}

	}
	client := http.Client{
		Timeout: 18000 * time.Second,
	}
	url := "https://27.push2.eastmoney.com/api/qt/ulist/sse?invt=3&pi=0&pz=" + pageSize + "&mpi=2000&secids=" + secidAllStr + "&ut=6d2ffaa6a585d612eda28417681d58fb&fields=f2,f3,f12,f14,f30,f31,f32&po=1"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		beego.Error(err)
		return
	}

	req.Header.Set("Accept", "text/event-stream")

	resp, err := client.Do(req)
	if err != nil {
		beego.Error(err)
		return
	}

	defer resp.Body.Close()

	for {
		event, err := common.ReadEvent(resp.Body)
		if err != nil {
			beego.Error(err)
			return
		}

		// 根据 event 做不同的处理
		beego.Info(event)
		var eventJson map[string]interface{}
		err = json.Unmarshal([]byte(event), &eventJson)
		if err != nil {
			beego.Error("Failed to parse JSON string:", err)
			return
		}
		if eventJson["data"] != nil {
			var data = eventJson["data"].(map[string]interface{})
			var diff = data["diff"].(map[string]interface{})
			for _, eachDiff := range diff {
				code := eachDiff.(map[string]interface{})["f12"] // f12 股票代码
				if currentStockPrice, ok := eachDiff.(map[string]interface{})["f2"]; ok {
					for _, stock := range stocks {
						if code == stock.Code && stock.Type == 1 {
							if currentStockPrice.(float64)/100 < stock.Min {
								// 在需要时使用redisService执行Redis命令
								keyString := stock.Code + "_min"
								result, _ := redisService.Do("GET", keyString)
								if result != nil {
									continue
								} else {
									services.SendDingMsg("股票提醒: " + stock.Code + " " + stock.Name + " 当前股价" + strconv.FormatFloat(currentStockPrice.(float64)/100, 'f', -1, 64) + " 已经低于最低设定值，是否考虑割肉？")
									result, err = redisService.Do("SET", keyString, 1)
									result, _ = redisService.Do("EXPIRE", keyString, 600)
								}

							}
							if currentStockPrice.(float64)/100 > stock.Max {
								keyString := stock.Code + "_max"
								result, _ := redisService.Do("GET", keyString)
								if result != nil {
									continue
								} else {
									services.SendDingMsg("股票提醒: " + stock.Code + " " + stock.Name + " 当前股价" + strconv.FormatFloat(currentStockPrice.(float64)/100, 'f', -1, 64) + " 已经高于最高设定值，是否考虑小赚一波？")
									result, _ = redisService.Do("SET", keyString, 1)
									result, _ = redisService.Do("EXPIRE", keyString, 600)
								}
							}
						}
					}
				}
				if currentPercentage, ok := eachDiff.(map[string]interface{})["f3"]; ok {
					for _, stock := range stocks {
						if code == stock.Code && stock.Type == 2 {
							if currentPercentage.(float64)/100 < stock.Min {
								keyString := stock.Code + "_min"
								result, _ := redisService.Do("GET", keyString)
								if result != nil {
									continue
								} else {
									services.SendDingMsg("股票提醒: " + stock.Code + " " + stock.Name + " 当前跌幅" + strconv.FormatFloat(currentPercentage.(float64)/100, 'f', -1, 64) + " 已经低于最低设定值，是否考虑割肉？")
									result, err = redisService.Do("SET", keyString, 1)
									result, err = redisService.Do("EXPIRE", keyString, 600)
								}
							}
							if currentPercentage.(float64)/100 > stock.Max {
								keyString := stock.Code + "_max"
								result, _ := redisService.Do("GET", keyString)
								if result != nil {
									continue
								} else {
									services.SendDingMsg("股票提醒: " + stock.Code + " " + stock.Name + " 当前涨幅" + strconv.FormatFloat(currentPercentage.(float64)/100, 'f', -1, 64) + " 已经高于最高设定值，是否考虑小赚一波？")
									result, _ = redisService.Do("SET", keyString, 1)
									result, _ = redisService.Do("EXPIRE", keyString, 600)
								}
							}
						}
					}
				}

			}
		}

	}
}
