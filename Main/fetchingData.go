package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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

func FetchData(SelectedApi Api, m model) ApiResponse {
	processedApi := processRequest(SelectedApi, m.SelectedCollection.LocalVariables)

	headers := processedApi.Headers
	api := buildURL(processedApi)

	url := strings.TrimSpace(api)
	url = strings.Trim(url, `"`)

	method := processedApi.Method

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

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return ApiResponse{StatusCode: 0, Status: "Failed to read Response : " + err.Error()}
	}

	m.apiResponse = ApiResponse{
		StatusCode:     resp.StatusCode,
		Status:         resp.Status,
		Body:           string(bodyBytes),
		Headers:        resp.Header,
		RequestHeaders: SelectedApi.Headers,
		ContentType:    resp.Header.Get("Content-Type"),
		ContentLength:  resp.ContentLength,
	}

	return m.apiResponse

}

func PostAPiFunc(m model) ApiResponse {
	SelectedApi := processRequest(m.SelectedApi, m.LocalVariables)

	headers := m.SelectedApi.Headers

	data := parseData(SelectedApi)

	Url := buildURL(SelectedApi)
	bodyReader := strings.NewReader(data)

	url := strings.TrimSpace(Url)
	url = strings.Trim(url, `"`)

	method := SelectedApi.Method

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return ApiResponse{StatusCode: 0, Status: err.Error()}
	}

	newHeader := Header{
		Key:   "Content-Type",
		Value: "application/json",
	}

	headers = append(headers, newHeader)

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
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return ApiResponse{StatusCode: 0, Status: "Failed to read Response : " + err.Error()}
	}

	m.apiResponse = ApiResponse{
		StatusCode:     resp.StatusCode,
		Status:         resp.Status,
		Body:           string(bodyBytes),
		Headers:        resp.Header,
		RequestHeaders: m.SelectedApi.Headers,
		ContentType:    resp.Header.Get("Content-Type"),
		ContentLength:  resp.ContentLength,
	}

	return m.apiResponse
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

type apiResponseMsg struct {
	response ApiResponse
}

func fetchApiCommand(api Api, m model) tea.Cmd {
	return func() tea.Msg {
		response := FetchData(api, m)
		return apiResponseMsg{response: response}
	}
}

func postApiCommand(m model) tea.Cmd {
	return func() tea.Msg {
		response := PostAPiFunc(m)
		return apiResponseMsg{response: response}
	}
}
func buildURL(api Api) string {
	if len(api.QueryParams) == 0 {
		return api.Url
	}

	var params []string
	for _, param := range api.QueryParams {
		params = append(params, url.QueryEscape(param.Key)+"="+url.QueryEscape(param.Value))
	}

	return api.Url + "?" + strings.Join(params, "&")
}

func processRequest(api Api, variables []LocalVariable) Api {
	processed := api

	processed.Url = replaceVariables(api.Url, variables)

	for i := range processed.Headers {
		processed.Headers[i].Value = replaceVariables(api.Headers[i].Value, variables)
	}
	for i := range processed.QueryParams {
		processed.QueryParams[i].Value = replaceVariables(api.QueryParams[i].Value, variables)
	}
	for i := range processed.BodyField {
		processed.BodyField[i].Value = replaceVariables(api.BodyField[i].Value, variables)
	}

	return processed
}

func replaceVariables(text string, variables []LocalVariable) string {
	result := text
	for _, variable := range variables {
		placeholder := "{{" + variable.Key + "}}"
		result = strings.ReplaceAll(result, placeholder, variable.Value)
	}
	return result
}
