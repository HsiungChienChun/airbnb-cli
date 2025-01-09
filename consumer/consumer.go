package consumer

import (
	"github.com/go-sql-driver/mysql"
	"sync"
)

var consumerStop bool
func NewConsumer(currency int, queueID string) (err error) {
	if currency <= 0 {
		currency = 1
	}

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(currency)
	for i := 0; i < currency; i++ {
		go func() {
			defer waitGroup.Done()

			for !consumerStop {
				// todo 抢占任务

				// 执行任务
					// 抓取页面信息
					// 解析分页信息
					// 解析详情页连接
					// 解析详情
					// 数据存储

				// 关闭任务
			}
		}()
	}

	waitGroup.Wait()

	return
}
