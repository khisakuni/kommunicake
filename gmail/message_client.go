package gmail

import (
	"fmt"

	api "google.golang.org/api/gmail/v1"
)

type MessageClient struct {
}

type MessageRequest func() (*api.Message, error)

func NewMessageClient() *MessageClient {
	return &MessageClient{}
}

func (mc *MessageClient) GetMessageBody(reqFn MessageRequest) (string, error) {
	message, err := reqFn()
	if err != nil {
		return "", err
	}

	fmt.Println(message)

	return "", nil
}
