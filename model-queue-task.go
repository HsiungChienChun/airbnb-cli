package main

import (
	"encoding/json"
	"os"
)

const (
	TaskStateWait = 1
	TaskStateLock = 2
	TaskStateSucc = 3
	TaskStateFail = 4
)

type TaskItem struct {
	Name    string            `json:"name"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

type TaskData struct {
	Tasks []*TaskItem `json:"tasks"`
}

func ParseTaskFile(dataPath string) (tasks []*TaskItem, err error) {
	dataBts, err := os.ReadFile(dataPath)
	if err != nil {
		// Todo errors.Wrap
		return
	}

	taskData := &TaskData{
		Tasks: make([]*TaskItem, 0),
	}
	err = json.Unmarshal(dataBts, &taskData)
	if err != nil {
		// Todo errors.Wrap
		return
	}

	tasks = taskData.Tasks
	return
}
