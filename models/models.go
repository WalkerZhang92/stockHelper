package models

import (
	"github.com/astaxie/beego/orm"
)

type User struct {
	Id       int
	Name     string
	NickName string
	Phone    int64  `json:"phone"`
	Password string `json:"password"`
	Avatar   string
}

type Stock struct {
	Id        int
	Code      int
	Title     string
	OffMarket int8 `json:"off_market"`
	Date      int
}

type StockSelect struct {
	Id           int
	User         int
	Code         string
	Title        string
	Date         string
	MaxTimes     int `json:"max_times"`
	CurrentTimes int `json:"current_times"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    interface{}
}

type WatchStock struct {
	Id    int
	Code  string
	Name  string
	Type  int
	Min   float64
	Max   float64
	IsDel int8 `json:"is_del"`
}

func (o *WatchStock) TableName() string {
	return "watch_stock"
}

func init() {
	// 需要在init中注册定义的model
	orm.RegisterModel(new(User), new(Stock), new(StockSelect), new(WatchStock))
}
