package gmail_test

import (
	"encoding/base64"
	"fmt"
	"testing"

	api "google.golang.org/api/gmail/v1"

	"github.com/khisakuni/kommunicake/gmail"
)

type getMessageResult struct {
	body *gmail.DecodedMessage
	err  error
}

var err = fmt.Errorf("uh oh")
var messageBody = "message-foobar"
var textMessagePart = &api.MessagePart{
	MimeType: "text/plain",
	Body: &api.MessagePartBody{
		Data: encodedMessageBody,
	},
}
var htmlMessagePart = &api.MessagePart{
	MimeType: "text/html",
	Body: &api.MessagePartBody{
		Data: encodedMessageBody,
	},
}
var encodedMessageBody = base64.URLEncoding.EncodeToString([]byte(messageBody))
var bodyMessage = &api.Message{
	Payload: textMessagePart,
}
var htmlPartMessage = &api.Message{
	Payload: &api.MessagePart{
		Body: &api.MessagePartBody{},
		Parts: []*api.MessagePart{
			&api.MessagePart{
				MimeType: "text/html",
				Body: &api.MessagePartBody{
					Data: encodedMessageBody,
				},
			},
		},
	},
}
var textPartMessage = &api.Message{
	Payload: &api.MessagePart{
		Body: &api.MessagePartBody{},
		Parts: []*api.MessagePart{
			&api.MessagePart{
				MimeType: "text/plain",
				Body: &api.MessagePartBody{
					Data: encodedMessageBody,
				},
			},
		},
	},
}
var textAndHTMLPartMessage = &api.Message{
	Payload: &api.MessagePart{
		Body:  &api.MessagePartBody{},
		Parts: []*api.MessagePart{htmlMessagePart, textMessagePart},
	},
}
var getMessageTests = []struct {
	input    gmail.MessageRequest
	expected getMessageResult
}{
	{
		input:    func() (*api.Message, error) { return nil, err },
		expected: getMessageResult{body: &gmail.DecodedMessage{MimeType: "", Body: ""}, err: err},
	},
	{
		input:    func() (*api.Message, error) { return bodyMessage, nil },
		expected: getMessageResult{body: &gmail.DecodedMessage{Body: messageBody, MimeType: "text/plain"}, err: nil},
	},
	{
		input:    func() (*api.Message, error) { return htmlPartMessage, nil },
		expected: getMessageResult{body: &gmail.DecodedMessage{Body: messageBody, MimeType: "text/html"}, err: nil},
	},
	{
		input:    func() (*api.Message, error) { return textPartMessage, nil },
		expected: getMessageResult{body: &gmail.DecodedMessage{Body: messageBody, MimeType: "text/plain"}, err: nil},
	},
	{
		input:    func() (*api.Message, error) { return textAndHTMLPartMessage, nil },
		expected: getMessageResult{body: &gmail.DecodedMessage{Body: messageBody, MimeType: "text/html"}, err: nil},
	},
}

func TestGetMessageBody(t *testing.T) {
	client := gmail.NewMessageClient()

	for _, getMessageTest := range getMessageTests {
		messageBody, err := client.GetMessageBody(getMessageTest.input)

		if err != getMessageTest.expected.err {
			t.Errorf("Expected err to be %v\n got %v\n", getMessageTest.expected.err, err)
		}

		if messageBody.Body != getMessageTest.expected.body.Body {
			t.Errorf("Expected Body to be %s, got %s\n", getMessageTest.expected.body.Body, messageBody.Body)
		}

		if messageBody.MimeType != getMessageTest.expected.body.MimeType {
			t.Errorf("Expected MimeType to be %s, got %s\n", getMessageTest.expected.body.MimeType, messageBody.MimeType)
		}
	}
}
