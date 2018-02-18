package main

import (
	"fmt"

	_ "github.com/joho/godotenv/autoload"
	"github.com/khisakuni/kommunicake/jobs/queue"
)

func main() {
	fmt.Printf("WORKER MAIN\n")
	q, err := queue.NewQueue("default")
	handleError(err)
	defer q.CleanUp()
	fmt.Println("INITIALIZED WORKER")

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
		fmt.Println("RECEIVED MESSAGE!!!")
		//		bodyMap := args.(map[string]interface{})
		//		fmt.Println(bodyMap)
		//		historyID := uint64(bodyMap["HistoryId"].(float64))
		//		userID := int(bodyMap["UserId"].(float64))
		//
		//		gmail_workers.ProcessGmailHistoryID(historyID, userID)

	})

	err = q.Consume()
	if err != nil {
		handleError(err)
	}

	select {}
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("UH OH %s\n", err.Error())
		panic(err)
	}
}
