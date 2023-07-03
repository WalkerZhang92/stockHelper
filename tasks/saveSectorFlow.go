package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"net/http"
	"quickstart/models"
	"strconv"
	"time"
)

func SaveSectorFLow()  {
	url := "https://push2.eastmoney.com/api/qt/clist/get?fid=f62&po=1&pz=500&pn=1&np=1&fltt=2&invt=2&fs=m:90+t:2&fields=f12,f14,f2,f3,f62,f184,f66,f69,f72,f75,f78,f81,f84,f87,f204,f205,f124,f1,f13"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	result := data["data"]
	if err != nil {
		fmt.Println(err)
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

		newId, _ := o.Insert(sectorFlow)
		fmt.Println(newId)

	}

}
