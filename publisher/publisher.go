package main

import (
	"encoding/json"

	"github.com/khisakuni/kommunicake/jobs/queue"
)

func main() {
	q, err := queue.NewQueue("default")
	defer q.CleanUp()

	handleError(err)

	args := struct {
		Name    string
		Message string
	}{Name: "SayHi", Message: "Hey there!"}
	j, _ := json.Marshal(args)
	err = q.Publish(j)
	handleError(err)
}

func handleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
