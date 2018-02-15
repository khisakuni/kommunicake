package main

import (
	"github.com/khisakuni/kommunicake/jobs/queue"
)

func main() {
	q, err := queue.NewQueue("test")
	defer q.CleanUp()
	handleError(err)

	select {}

}

func handleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
