package producer

import (
	"encoding/json"
	"os"
)

type TaskHeader struct {
	UserAgent string `json:"user-agent"`
}

type TaskItem struct {
	Name    string     `json:"name"`
	URL     string     `json:"url"`
	Headers TaskHeader `json:"headers"`
}

type TaskData struct {
	Tasks []*TaskItem `json:"tasks"`
}

func ParseTasks(dataPath string) (tasks []*TaskItem, err error) {
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
