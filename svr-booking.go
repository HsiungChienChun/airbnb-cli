package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// 数据爬取任务执行逻辑
func bookingTaskProcess(task *TaskItem) (err error) {
	listPageBody, err := getListPage(task)
	if err != nil {
		return
	}
	cursorList := parseSubPageCursor(listPageBody)

	// 依次遍历各个分页，解析数据
	for _, cursor := range cursorList {
		pageBody := getListPageByCursor(cursor)

		// 获取列表页内详情链接信息学
		detailList := parseDetailIDs(pageBody)

		// 分页内依次解析详情信息
		bookings := []*Booking{}
		for _, detail := range detailList {
			booking := parseDetailInfo(detail)
			bookings = append(bookings, booking)
		}

		// 按分页存储数据
		ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
		_, err = addBooking(ctx, bookings)
		if err != nil {
			return
		}
	}
	return
}

func parseDetailInfo(detailBody string) (detail *Booking) {
	// TODO 解析详情页中的客房预定信息
	return
}

// 解析分页列表，后去列表项ID（or链接）
func parseDetailIDs(pageBody string) (detailIDs []string) {
	// TODO 解析搜索结果页的列表item中的详情页访问连接ID列表
	return
}

// 访问第一页，拿到分页游标列表，后续按分页遍历
func parseSubPageCursor(pageBody string) (cursors []string) {
	// TODO 解析搜索结果页的分页cursor列表
	return
}

// 根据cursor获取搜索结果列表页数据
func getListPageByCursor(cursor string) (pageInfo string) {
	// TODO 访问搜索结果列表分页，获取所有数据后再解析
	return
}

// 根绝任务信息获取搜索结果首页列表
func getListPage(task *TaskItem) (pageInfo string, err error) {
	req, err := http.NewRequest("GET", task.URL, strings.NewReader(""))
	if err != nil {
		return
	}

	for key, val := range task.Headers {
		req.Header.Add(key, val)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	dataBs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	pageInfo = string(dataBs)
	return
}
