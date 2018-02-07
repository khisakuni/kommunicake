package main

import (
	"fmt"

	"github.com/khisakuni/kommunicake/jobs/queue"
)

func main() {
	q, err := queue.NewQueue("default")
	defer q.CleanUp()
	handleError(err)

	q.RegisterWorker("SayHi", func(args interface{}) {
		bodyMap := args.(map[string]interface{})
		message := bodyMap["Message"].(string)
		fmt.Printf("hi!! %s\n", message)
	})

	q.Consume()

	select {}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
