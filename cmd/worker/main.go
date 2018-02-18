package main

import (
	"fmt"

	_ "github.com/joho/godotenv/autoload"
	"github.com/khisakuni/kommunicake/jobs/queue"

	gmail_workers "github.com/khisakuni/kommunicake/jobs/workers/gmail"
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

	/*
		TODO:
			- extract worker code into it's own file
			- clean up controller code
			- do something useful in ProcessGmailHistoryId worker
				- get entire message
	*/

	q.RegisterWorker("ProcessGmailHistoryID", func(args interface{}) {
		fmt.Println("received message!")
		bodyMap := args.(map[string]interface{})
		fmt.Println(bodyMap)
		historyId := uint64(bodyMap["HistoryId"].(float64))
		userId := int(bodyMap["UserId"].(float64))

		gmail_workers.ProcessGmailHistoryID(historyId, userId)

	})

	q.Consume()

	select {}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
