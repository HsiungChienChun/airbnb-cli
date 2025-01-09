package main

import (
	"context"
	"flag"
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"
)

func ConsumerStart(stopNotify func() bool) {
	var currency int
	flag.IntVar(&currency, "workers", 1, "消费并发度，默认1")
	queueIDPtr := flag.String("queue", "", "任务消费队列名")
	flag.Parse()

	if queueIDPtr == nil {
		fmt.Println("invalid command: 任务消费队列名")
		helpInfo()
		return
	}
	queueID := *queueIDPtr

	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(currency)
	for i := 1; i <= currency; i++ {
		go func(workerID int) {
			defer func() {
				waitGroup.Done()
			}()

			for {
				if stopNotify() {
					fmt.Printf("服务关闭，消费worker(%d)退出\n", workerID)
					return
				}

				tmpWait := &sync.WaitGroup{}
				tmpWait.Add(1)
				go consumerWorkerDo(tmpWait, stopNotify, queueID)
				tmpWait.Wait()

				// 防止连续出错打垮数据库
				time.Sleep(20 * time.Millisecond)
			}
		}(i)
	}

	waitGroup.Wait()
	fmt.Println("消费者退出！")
	return
}

// 单独routine真正执行任务，防止出错影响导致消费者worker退出
/* 这个routine退出：
case1:	服务关闭
case2:	执行出错
*/
func consumerWorkerDo(waitGroup *sync.WaitGroup, stopNotify func() bool, queueID string) (err error) {
	defer func() {
		if recoverErr := recover(); recoverErr != nil {
			trace := make([]byte, 1<<16)
			n := runtime.Stack(trace, true)
			err = fmt.Errorf("panic: '%v'\n, Stack Trace:\n %s", recoverErr, string(trace[:int(math.Min(float64(n), float64(7000)))]))
		}
		if err != nil {
			fmt.Printf("消费执行出错，错误信息：%+v\n", err)
		}
		waitGroup.Done()
	}()

	var taskIDs []int64
	var nextLastID int64
	for { // 分批次循环 抢占任务 - 读取任务 - 执行任务 - 关闭任务
		if stopNotify() {
			fmt.Println("服务关闭，消费退出")
			return
		}

		// 抢占任务
		ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
		nextLastID, taskIDs, err = waitTasks(ctx, queueID, TaskStateWait, nextLastID, 10)
		if err != nil {
			return
		}
		// 当前没有任务，等待后重试
		if len(taskIDs) <= 0 {
			time.Sleep(20 * time.Millisecond)
			continue
		}

		var affected int64
		var taskCtx context.Context
		var task *TaskItem
		var taskFail bool
		for _, taskID := range taskIDs { // 分批次内部依次尝试执行
			// 抢占任务: 等待 => 锁定
			taskCtx, _ = context.WithTimeout(context.Background(), 100*time.Millisecond)
			affected, err = updateTask(taskCtx, taskID, TaskStateWait, TaskStateLock)
			if err != nil {
				return
			}
			if affected == 0 {
				continue
			}

			// 取任务详情
			taskCtx, _ = context.WithTimeout(context.Background(), 100*time.Millisecond)
			taskFail, task, err = taskDetail(ctx, taskID)
			if err != nil {
				return
			}
			if taskFail { // 脏数据，直接跳过
				taskCtx, _ = context.WithTimeout(context.Background(), 100*time.Millisecond)
				affected, err = updateTask(taskCtx, taskID, TaskStateLock, TaskStateFail)
				continue
			}

			// 解析任务并执行: 抓取 => 解析子页面连接 => 解析详情 => 数据存储
			err = bookingTaskProcess(task)
			if err != nil {
				return
			}

			// 关闭任务：锁定 => 完成
			taskCtx, _ = context.WithTimeout(context.Background(), 100*time.Millisecond)
			affected, err = updateTask(taskCtx, taskID, TaskStateLock, TaskStateSucc)
		}
	}
}
