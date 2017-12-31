package middleware

import "bytes"

type contextKey string

func (c contextKey) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("app context key")
	buffer.WriteString(string(c))
	return buffer.String()
}
