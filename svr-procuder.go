package main

import (
	"context"
	"flag"
	"fmt"
	"time"
)

// 启动生产逻辑
func ProducerStart(stopNotify func() bool) {
	dataPathPtr := flag.String("data", "", "任务文件访问路径path")
	queueIDPtr := flag.String("queue", "", "任务投递队列名")
	var sendBatchNum int
	flag.IntVar(&sendBatchNum, "batch", 1, "任务投递队列名")
	flag.Parse()

	if dataPathPtr == nil {
		fmt.Println("invalid command: 缺失任务文件")
		helpInfo()
		return
	}
	if queueIDPtr == nil {
		fmt.Println("invalid command: 缺失投递队列名")
		helpInfo()
		return
	}
	dataPath := *dataPathPtr
	queueID := *queueIDPtr

	tasks, err := ParseTaskFile(dataPath)
	if err != nil {
		fmt.Printf("invalid command: 任务文件读取或者解析失败\n	错误信息：%+v\n", err)
		return
	}

	fmt.Println("生产者启动，任务开始处理！")
	for i := 0; i < len(tasks); i += sendBatchNum {
		if stopNotify() {
			fmt.Println("服务关闭，生产者退出！")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		_, err = AddTask(ctx, queueID, tasks[i:i+sendBatchNum])
		cancel()
		if err != nil {
			fmt.Printf("任务执行失败，已成功执行: %d\n	错误信息：%+v\n", i, err)
			return
		}
	}

	fmt.Println("任务处理完成，生产者退出！")
	return
}
