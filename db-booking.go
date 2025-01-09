package main

import (
	"context"
	"fmt"
	"strings"
)

/*
CREATE TABLE `airbnb-booking` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `hotel_name` varchar(256) NOT NULL DEFAULT '' COMMENT '酒店名称',
  `star` int(11) NOT NULL COMMENT '评分',
  `price` int(11) NOT NULL DEFAULT '' COMMENT '价格',
  `taxes_price` int(11) NOT NULL COMMENT '税前价格',
  `checkin` timestamp(3) NOT NULL COMMENT '入住时间',
  `checkout` timestamp(3) NOT NULL COMMENT '退房时间',
  `guests` tinyint(4) NOT NULL COMMENT '顾客数',
  PRIMARY KEY (`id`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务队列表';
*/

const (
	_addBookingSQL = "INSERT INTO `task_queue`(`hotel_name`, `star`, `price`, `taxes_price`, `checkin`, `checkout`, `guests`) VALUES%s;"
)

func addBooking(ctx context.Context, bookings []*Booking) (affected int64, err error) {
	if len(bookings) <= 0 {
		return

	}

	marks := []string{}
	params := []interface{}{}
	for _, booking := range bookings {
		marks = append(marks, "(?, ?, ?, ?, ?, ?, ?)")
		params = append(params, booking.Hotel, booking.Star, booking.Price, booking.TaxesPrices, booking.CheckIn, booking.CheckOut, booking.Guests)
	}
	sqlStr := fmt.Sprintf(_addBookingSQL, strings.Join(marks, ", "))

	res, err := getDB().ExecContext(ctx, sqlStr, params...)
	if err != nil {
		return
	}

	affected, err = res.RowsAffected()
	if err != nil {
		return
	}
	return
}
