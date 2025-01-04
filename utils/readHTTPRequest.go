package utils

import (
	"bufio"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net"
	"strings"

	"github.com/defskela/httpServer/http"
	"github.com/defskela/httpServer/logger"
)

func ReadHTTPRequest(conn net.Conn) (*http.HTTPRequest, error) {

	reader := bufio.NewReader(conn)
	request := &http.HTTPRequest{
		Headers:  make(map[string]string),
		FormData: make(map[string]string),
	}

	// Парсим первую строку запроса, где содержится метод, путь и версия протокола
	line, err := reader.ReadString('\n')

	if err != nil {
		return nil, fmt.Errorf("ошибка чтения строки: %w", err)
	}

	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("malformed request line: %s", line)
	}
	request.Method, request.Path, request.Version = parts[0], parts[1], parts[2]

	// Парсим заголовки запроса
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения строки: %w", err)
		}

		line = strings.TrimSpace(line)
		if line == "" {
			break // Конец заголовков
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			request.Headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	contentType, ok := request.Headers["Content-Type"]
	if ok && strings.HasPrefix(contentType, "multipart/form-data") {
		// Разбор multipart/form-data
		_, params, err := mime.ParseMediaType(contentType)
		if err != nil {
			return nil, fmt.Errorf("error parsing Content-Type: %w", err)
		}
		boundary, ok := params["boundary"]
		if !ok {
			return nil, fmt.Errorf("boundary not found in Content-Type")
		}

		mr := multipart.NewReader(reader, boundary)
		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("error reading part: %w", err)
			}

			data, err := io.ReadAll(part)
			if err != nil {
				return nil, fmt.Errorf("error reading part body: %w", err)
			}

			// Сохраняем данные части
			name := part.FormName()
			if name != "" {
				request.FormData[name] = string(data)
			}

			fileName := part.FileName()
			if fileName != "" {
				request.FormData["fileName"] = fileName
			}
		}
	} else if request.Method == "POST" || request.Method == "PUT" || request.Method == "PATCH" || request.Method == "DELETE" {
		logger.Debug(1)
		body, err := io.ReadAll(reader)
		logger.Debug(2)
		if err != nil {
			return nil, fmt.Errorf("error reading body: %w", err)
		}
		request.Body = string(body)
	}

	return request, nil
}
