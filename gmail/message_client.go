package gmail

import (
	"encoding/base64"

	api "google.golang.org/api/gmail/v1"
)

type MessageClient struct {
}

type MessageRequest func() (*api.Message, error)
type DecodedMessage struct {
	Body     string
	MimeType string
}

func NewMessageClient() *MessageClient {
	return &MessageClient{}
}

func (mc *MessageClient) GetMessageBody(reqFn MessageRequest) (*DecodedMessage, error) {
	message, err := reqFn()
	if err != nil {
		return &DecodedMessage{Body: "", MimeType: ""}, err
	}

	var messageBody string
	var mimeType string
	if len(message.Payload.Body.Data) > 0 {
		decoded, err := decodeMessagePart(message.Payload)
		if err != nil {
			return &DecodedMessage{Body: "", MimeType: ""}, err
		}
		messageBody = decoded
		mimeType = message.Payload.MimeType
	} else {
		for _, part := range message.Payload.Parts {
			decoded, err := decodeMessagePart(part)
			if err != nil {
				return &DecodedMessage{Body: "", MimeType: ""}, err
			}
			messageBody = decoded
			mimeType = part.MimeType
			if mimeType == "text/html" {
				break
			}
		}
	}

	return &DecodedMessage{Body: messageBody, MimeType: mimeType}, nil
}

func decodeMessagePart(part *api.MessagePart) (string, error) {
	decoded, err := decode(part.Body.Data)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

func decode(messageBody string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(messageBody)
}
