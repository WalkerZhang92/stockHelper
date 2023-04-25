package tasks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"quickstart/common"
	"quickstart/models"
	"quickstart/services"
	"time"

	"github.com/astaxie/beego/orm"

	"strconv"

	"github.com/astaxie/beego"
)

func WatchStock() {
	var stocks []*models.WatchStock
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
		Timeout: 600 * time.Second,
	}
	url := "https://27.push2.eastmoney.com/api/qt/ulist/sse?invt=3&pi=0&pz=10&mpi=2000&secids=" + secidAllStr + "&ut=6d2ffaa6a585d612eda28417681d58fb&fields=f2,f3,f12,f14,f30,f31,f32&po=1"
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
