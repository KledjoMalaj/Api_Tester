package main

import (
	"io"
	"net/http"
	"strings"
)

type ApiResponse struct {
	StatusCode    int
	Status        string
	Body          string
	Headers       http.Header
	ContentType   string
	ContentLength int64
}

func FetchData(SelectedApi Api) ApiResponse {
	url := strings.TrimSpace(SelectedApi.Url)
	url = strings.Trim(url, `"`)

	method := SelectedApi.Method

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return ApiResponse{StatusCode: 0, Status: err.Error()}
	}

	// Optional: add headers if needed
	req.Header.Set("User-Agent", "Go-Client")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ApiResponse{StatusCode: 0, Status: err.Error()}
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	return ApiResponse{
		StatusCode:    resp.StatusCode,
		Status:        resp.Status,
		Body:          string(bodyBytes),
		Headers:       resp.Header,
		ContentType:   resp.Header.Get("Content-Type"),
		ContentLength: resp.ContentLength,
	}

}

func PostAPiFunc(m model) ApiResponse {
	SelectedApi := m.SelectedApi
	data := m.jsonInput.Value()
	bodyReader := strings.NewReader(data)

	url := strings.TrimSpace(SelectedApi.Url)
	url = strings.Trim(url, `"`)

	method := SelectedApi.Method

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return ApiResponse{StatusCode: 0, Status: err.Error()}
	}

	req.Header.Set("Content-Type", "application/json")
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
		StatusCode:    resp.StatusCode,
		Status:        resp.Status,
		Body:          string(bodyBytes),
		Headers:       resp.Header,
		ContentType:   resp.Header.Get("Content-Type"),
		ContentLength: resp.ContentLength,
	}
}
