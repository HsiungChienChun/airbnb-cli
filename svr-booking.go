package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// 解析任务并执行保存
func bookingTaskProcess(task *TaskItem) (err error) {
	listPageBody, err := getListPage(task)
	if err != nil {
		return
	}
	cursorList := parseSubPageCursor(listPageBody)

	// 依次遍历各个分页，解析数据
	for _, item := range cursorList {
		pageBody := ""

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
	return
}

// 解析分页列表，后去列表项ID（or链接）
func parseDetailIDs(pageBody string) (detailIDs []string) {
	return
}

// 访问第一页，拿到分页游标列表，后续按分页遍历
func parseSubPageCursor(pageBody string) (cursors []string) {
	return
}

// 获取首页列表
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

// /s/Europe/homes?refinement_paths%5B%5D=%2Fhomes&query=Europe&place_id=ChIJhdqtz4aI7UYRefD8s-aZ73I&flexible_trip_lengths%5B%5D=one_week&monthly_start_date=2025-02-01&monthly_length=3&monthly_end_date=2025-05-01&search_mode=regular_search&price_filter_input_type=0&channel=EXPLORE&search_type=user_map_move&price_filter_num_nights=5&ne_lat=61.86486211129339&ne_lng=37.42984862589401&sw_lat=23.59018458925222&sw_lng=-9.683937857471903&zoom=3.3359550737736523&zoom_level=3&search_by_map=true&federated_search_session_id=ec667451-7475-4772-bb10-a897c895251a&pagination_search=true&cursor=eyJzZWN0aW9uX29mZnNldCI6MCwiaXRlbXNfb2Zmc2V0IjozNiwidmVyc2lvbiI6MX0%3D
