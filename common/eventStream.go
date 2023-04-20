package common

import (
	"bufio"
	"io"
	"strings"
)

func ReadEvent(body io.Reader) (string, error) {
	// 读取一行数据
	reader := bufio.NewReader(body)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	// 去掉行末的换行符
	line = strings.TrimSpace(line)

	// 判断是否是事件行
	if !strings.HasPrefix(line, "data:") {
		return "", nil
	}

	// 获取事件数据
	eventData := strings.TrimSpace(strings.TrimPrefix(line, "data:"))

	return eventData, nil
}
