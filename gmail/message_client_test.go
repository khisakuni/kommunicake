package gmail_test

import (
	"fmt"
	"testing"

	api "google.golang.org/api/gmail/v1"

	"github.com/khisakuni/kommunicake/gmail"
)

type getMessageResult struct {
	body string
	err  error
}

var err = fmt.Errorf("uh oh")
var getMessageTests = []struct {
	input    gmail.MessageRequest
	expected getMessageResult
}{
	{
		input:    func() (*api.Message, error) { return nil, err },
		expected: getMessageResult{body: "", err: err},
	},
}

func TestGetMessageBody(t *testing.T) {
	client := gmail.NewMessageClient()

	for _, getMessageTest := range getMessageTests {
		messageBody, err := client.GetMessageBody(getMessageTest.input)

		if err != getMessageTest.expected.err {
			t.Errorf("Expected err to be %v\n got %v\n", getMessageTest.expected.err, err)
		}

		if messageBody != getMessageTest.expected.body {
			t.Errorf("Expected messageBody to be %s, got %s\n", getMessageTest.expected.body, messageBody)
		}
	}
}
