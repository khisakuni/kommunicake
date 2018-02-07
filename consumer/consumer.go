package main

import (
	"fmt"

	"github.com/khisakuni/kommunicake/jobs/queue"
)

func main() {
	q, err := queue.NewQueue("test")
	defer q.CleanUp()
	handleError(err)

	q.Consume(func(input []byte) {
		fmt.Printf("GOT A MESSAGE >> %s\n", string(input))
	})

	select {}

}

func handleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
