package utils

import (
	"fmt"
	"net/http"
)

func CreateHTTPResponse(statusCode int, body string) string {
	return fmt.Sprintf("HTTP/1.1 %d %s\r\nContent-Length: %d\r\nContent-Type: text/plain\r\n\r\n%s",
		statusCode, http.StatusText(statusCode), len(body), body)
}
