package main

import (
	"context"
	"fmt"
	"strings"
	"time"
)

/*
CREATE TABLE `task_queue` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `queue` varchar(256) NOT NULL DEFAULT '' COMMENT '队列ID',
  `info` varchar(8096) NOT NULL DEFAULT '' COMMENT '任务信息',
  `state` tinyint NOT NULL COMMENT '任务状态：1待处理，2执行中，3执行成功，4执行失败',
  `locktime` timestamp(3) NOT NULL COMMENT '任务抢占时间，用于延时释放',
  PRIMARY KEY (`id`),
  KEY queue_state(`queue`, `state`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务队列表';
*/

const (
	_addTaskSQL    = "INSERT INTO `task_queue`(`queue`, `info`, `state`) VALUES%s;"
	_waitTasksSQL  = "SELECT `id` FROM `task_queue` WHERE `queue` = ? AND `state` = ? AND `id` > ? LIMIT ?;"
	_lockTaskSQL   = "UPDATE `task_queue` SET `state` = ?, locktime = ? WHERE `id` = ? AND `state` = ?;"
	_taskDetailSQL = "SELECT `info` FROM `task_queue` WHERE `id` = ?;"
)

func taskDetail(ctx context.Context, taskID int64) (fail bool, task *TaskItem, err error) {
	if taskID <= 0 {
		return
	}
	db := getDB()
	res, err := db.QueryContext(ctx, _taskDetailSQL, taskID)
	if err != nil {
		return
	}
	defer res.Close()

	for res.Next() {
		var tmpInfo string
		if err = res.Scan(tmpInfo); err != nil {
			return
		}
		tmpTask := &TaskItem{}
		err = JsonUnmarshal([]byte(tmpInfo), &tmpTask)
		if err != nil {
			fail = true
			return
		}
		task = tmpTask
	}
	err = res.Err()
	return
}

func updateTask(ctx context.Context, taskID, sourceState, targetState int64) (affected int64, err error) {
	res, err := getDB().ExecContext(ctx, _lockTaskSQL, targetState, time.Now().Format("2006-01-02 15:04:05"), taskID, sourceState)
	if err != nil {
		return
	}
	affected, _ = res.RowsAffected()
	return
}

func AddTask(ctx context.Context, queue string, tasks []*TaskItem) (affected int64, err error) {
	if len(tasks) <= 0 {
		return

	}

	marks := []string{}
	params := []interface{}{}
	for _, task := range tasks {
		marks = append(marks, "(?, ?, ?)")
		params = append(params, queue, JsonString(task), TaskStateWait)
	}

	sqlStr := fmt.Sprint(_addTaskSQL, strings.Join(marks, ", "))

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

func waitTasks(ctx context.Context, queueID string, state, lastID, scanSize int64) (nextLastID int64, taskIDs []int64, err error) {
	taskIDs = make([]int64, 0)

	sqlStr := _waitTasksSQL
	params := []interface{}{queueID, state, lastID, scanSize}

	rows, err := getDB().QueryContext(ctx, sqlStr, params...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var taskID int64
		fieldList := []interface{}{&taskID}
		if err = rows.Scan(fieldList...); err != nil {
			return
		}
		taskIDs = append(taskIDs, taskID)
		nextLastID = taskID
	}
	err = rows.Err()
	return
}
