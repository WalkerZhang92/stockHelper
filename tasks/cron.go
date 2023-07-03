package tasks

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

func StartCron() {
	c := cron.New()

	// 添加每天盯盘任务
	_, err := c.AddFunc("0/10 9-15 * * 1-5", func() {
		go WatchStock()
	})
	if err != nil {
		fmt.Println("Failed to add cron job:", err)
		return
	}

	_, err = c.AddFunc("0 16 * * 1-5", func() {
		go SaveSectorFLow()
	})
	if err != nil {
		fmt.Println("Failed to add cron job:", err)
		return
	}

	// 启动定时任务调度器
	c.Start()
	fmt.Println("Cron job started.")
}
