package main

import (
	"fmt"
	"strings"
)

func main() {
	/*
		airbnb-cli start consumer --workers=10 --queue=<your-queue-server>

	consumer start
	consumer stop

	producer start
	producer stop

		airbnb-cli start producer --data tasks.json --queue=<your-queue-server>
	*/

	var cmd string

	switch strings.ToLower(cmd) {
	case "consumer":



	case "producer":


	default:
		fmt.Println("invalid command")
	}
}

