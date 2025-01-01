package utils

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/defskela/httpServer/http"
)

func ReadHTTPRequest(conn net.Conn) (*http.HTTPRequest, error) {
	scanner := bufio.NewScanner(conn)
	request := &http.HTTPRequest{Headers: make(map[string]string)}

	if scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("malformed request line: %s", line)
		}
		request.Method, request.Path, request.Version = parts[0], parts[1], parts[2]
	}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			request.Headers[parts[0]] = parts[1]
		}
	}

	if scanner.Scan() {
		request.Body = scanner.Text()
	}

	return request, nil
}
