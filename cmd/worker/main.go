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

	fmt.Printf("q >>>>>>> %v\n", q)

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
