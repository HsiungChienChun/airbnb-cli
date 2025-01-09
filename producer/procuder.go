package producer

func NewProducer(dataPath string, queueID string) (err error) {
	tasks, err := ParseTasks(dataPath)
	if err != nil {
		return
	}

	var sendBatchNum = 10
	for i := 0; i < len(tasks); i += sendBatchNum {
		err = sendTask(dataPath, queueID, tasks[i:i+sendBatchNum])
		if err != nil {
			// TODO return success offset
			return
		}
	}

	return
}

func sendTask(dataID string, queueID string, tasks []*TaskItem) (err error) {
	// TODO send task to mysql
	return
}
