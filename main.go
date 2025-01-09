package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("invalid command: miss args")
		helpInfo()
		return
	}

	if strings.ToLower(os.Args[1]) == "help" {
		helpInfo()
		return
	}

	action := os.Args[1]
	workerType := os.Args[2]
	cmdType := strings.ToLower(action + ":" + workerType)

	// 初始化db连接
	getDB()

	switch cmdType {
	case "start:consumer":
		// airbnb-cli start consumer --workers=10 --queue=<your-queue-server>

		ConsumerStart(serviceStopNotify)

	case "start:producer":
		// airbnb-cli start producer --data tasks.json --queue=<your-queue-server>
		ProducerStart(serviceStopNotify)

	default:
		fmt.Println("invalid command")
		return
	}

	notify := make(chan os.Signal, 1)
	signal.Notify(notify, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-notify
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			serviceStop = true
			time.Sleep(1 * time.Second) // 等待业务退出后，进程退出
			return

		case syscall.SIGHUP:
		default:
			return
		}
	}
}

// 控制服务平滑退出
var serviceStop bool

func serviceStopNotify() bool {
	return serviceStop
}

// 输出命令行帮助信息
func helpInfo() {
	fmt.Print(`

	usage: 
		启动consumer: 
			命令：
				airbnb-cli start consumer --workers=10 --queue=<your-queue-server>
			
			参数：
				workers： 当前进程并行消费worker数，最小为1
				queue: 	  消费队列名称，请确保存在

		启动producer:
			命令：
				airbnb-cli start producer --data tasks.json --queue=<your-queue-server>

			参数：
				data：	待发送的任务信息文件，请确保文件可访问，且任务信息格式符合要求
				queue:  任务投递队列名，请确保队列存在

		关闭服务：
			命令：
				kill -9 {进程ID}

`)
}
