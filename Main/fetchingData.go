package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ApiResponse struct {
	StatusCode     int
	Status         string
	Body           string
	Headers        http.Header
	RequestHeaders []Header
	ContentType    string
	ContentLength  int64
}

func FetchData(SelectedApi Api) ApiResponse {
	headers := SelectedApi.Headers

	url := strings.TrimSpace(SelectedApi.Url)
	url = strings.Trim(url, `"`)

	method := SelectedApi.Method

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return ApiResponse{StatusCode: 0, Status: err.Error()}
	}

	for i := 0; i < len(headers); i++ {
		req.Header.Set(headers[i].Key, headers[i].Value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ApiResponse{StatusCode: 0, Status: err.Error()}
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	return ApiResponse{
		StatusCode:     resp.StatusCode,
		Status:         resp.Status,
		Body:           string(bodyBytes),
		Headers:        resp.Header,
		RequestHeaders: SelectedApi.Headers,
		ContentType:    resp.Header.Get("Content-Type"),
		ContentLength:  resp.ContentLength,
	}

}

func PostAPiFunc(m model) ApiResponse {
	SelectedApi := m.SelectedApi
	headers := m.SelectedApi.Headers

	data := parseData(SelectedApi)

	bodyReader := strings.NewReader(data)

	url := strings.TrimSpace(SelectedApi.Url)
	url = strings.Trim(url, `"`)

	method := SelectedApi.Method

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return ApiResponse{StatusCode: 0, Status: err.Error()}
	}

	for i := 0; i < len(headers); i++ {
		req.Header.Set(headers[i].Key, headers[i].Value)
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ApiResponse{StatusCode: 0, Status: err.Error()}
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, _ := io.ReadAll(resp.Body)

	// Return structured response
	return ApiResponse{
		StatusCode:     resp.StatusCode,
		Status:         resp.Status,
		Body:           string(bodyBytes),
		Headers:        resp.Header,
		RequestHeaders: m.SelectedApi.Headers,
		ContentType:    resp.Header.Get("Content-Type"),
		ContentLength:  resp.ContentLength,
	}
}
func parseData(selectedApi Api) string {
	if len(selectedApi.BodyField) == 0 {
		return "{}"
	}

	// Build JSON object
	var b strings.Builder
	b.WriteString("{\n")

	for i, field := range selectedApi.BodyField {
		// Add key-value pair
		b.WriteString(fmt.Sprintf("  \"%s\": \"%s\"", field.Key, field.Value))

		// Add comma if not the last item
		if i < len(selectedApi.BodyField)-1 {
			b.WriteString(",")
		}
		b.WriteString("\n")
	}

	b.WriteString("}")
	return b.String()
}
