package main

import "encoding/json"

func JsonString(data interface{}) string {
	bs, _ := json.Marshal(data)
	return string(bs)
}

func JsonUnmarshal(data interface{}, result interface{}) (err error) {
	var bs []byte
	if dataStr, ok := data.(string); ok {
		bs = ([]byte)(dataStr)
	} else {
		bs, _ = json.Marshal(data)
	}
	err = json.Unmarshal(bs, result)
	return
}
