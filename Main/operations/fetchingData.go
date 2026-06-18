package operations

import (
	"GoTuiFrontend/models"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func FetchData(SelectedApi models.Api, m models.TuiModel) models.ApiResponse {
	processedApi := ProcessRequest(SelectedApi, m.SelectedCollection.LocalVariables)

	headers := processedApi.Headers
	api := BuildURL(processedApi, m)

	url := strings.TrimSpace(api)
	url = strings.Trim(url, `"`)

	method := processedApi.Method

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return models.ApiResponse{StatusCode: 0, Status: err.Error()}
	}

	for i := 0; i < len(headers); i++ {
		req.Header.Set(headers[i].Key, ReplaceVariables(headers[i].Value, m.LocalVariables))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.ApiResponse{StatusCode: 0, Status: err.Error()}
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.ApiResponse{StatusCode: 0, Status: "Failed to read Response : " + err.Error()}
	}

	m.ApiResponse = models.ApiResponse{
		StatusCode:     resp.StatusCode,
		Status:         resp.Status,
		Body:           string(bodyBytes),
		Headers:        resp.Header,
		RequestHeaders: SelectedApi.Headers,
		ContentType:    resp.Header.Get("Content-Type"),
		ContentLength:  resp.ContentLength,
	}

	return m.ApiResponse

}

func PostAPiFunc(m models.TuiModel) models.ApiResponse {
	SelectedApi := ProcessRequest(m.SelectedApi, m.LocalVariables)

	headers := m.SelectedApi.Headers

	data := ParseData(SelectedApi, m)

	Url := BuildURL(SelectedApi, m)
	bodyReader := strings.NewReader(data)

	url := strings.TrimSpace(Url)
	url = strings.Trim(url, `"`)

	method := SelectedApi.Method

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return models.ApiResponse{StatusCode: 0, Status: err.Error()}
	}

	newHeader := models.Header{
		Key:   "Content-Type",
		Value: "application/json",
	}

	headers = append(headers, newHeader)

	for i := 0; i < len(headers); i++ {
		req.Header.Set(headers[i].Key, ReplaceVariables(headers[i].Value, m.LocalVariables))
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.ApiResponse{StatusCode: 0, Status: err.Error()}
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.ApiResponse{StatusCode: 0, Status: "Failed to read Response : " + err.Error()}
	}

	m.ApiResponse = models.ApiResponse{
		StatusCode:     resp.StatusCode,
		Status:         resp.Status,
		Body:           string(bodyBytes),
		Headers:        resp.Header,
		RequestHeaders: m.SelectedApi.Headers,
		ContentType:    resp.Header.Get("Content-Type"),
		ContentLength:  resp.ContentLength,
	}

	return m.ApiResponse
}
func ParseData(selectedApi models.Api, m models.TuiModel) string {
	if len(selectedApi.BodyField) == 0 {
		return "{}"
	}

	var b strings.Builder
	b.WriteString("{\n")

	for i, field := range selectedApi.BodyField {
		key := field.Key

		value := ReplaceVariables(field.Value, m.LocalVariables)

		formattedValue := FormatJSONValue(value)

		b.WriteString(fmt.Sprintf("  \"%s\": %s", key, formattedValue))

		if i < len(selectedApi.BodyField)-1 {
			b.WriteString(",")
		}
		b.WriteString("\n")
	}

	b.WriteString("}")

	return b.String()
}

func FormatJSONValue(value string) string {
	if value == "null" {
		return "null"
	}

	if value == "true" || value == "false" {
		return value
	}

	trimmed := strings.TrimSpace(value)
	if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
		return value
	}

	if strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}") {
		return value
	}

	return fmt.Sprintf("\"%s\"", value)
}

func FetchApiCommand(api models.Api, m models.TuiModel) tea.Cmd {
	return func() tea.Msg {
		response := FetchData(api, m)
		return models.ApiResponseMsg{Response: response}
	}
}

func PostApiCommand(m models.TuiModel) tea.Cmd {
	return func() tea.Msg {
		response := PostAPiFunc(m)
		return models.ApiResponseMsg{Response: response}
	}
}

func BuildURL(api models.Api, m models.TuiModel) string {
	if len(api.QueryParams) == 0 {
		return api.Url
	}

	var params []string
	for _, param := range api.QueryParams {
		params = append(params, url.QueryEscape(param.Key)+"="+url.QueryEscape(ReplaceVariables(param.Value, m.LocalVariables)))
	}

	return api.Url + "?" + strings.Join(params, "&")
}

func ProcessRequest(api models.Api, variables []models.LocalVariable) models.Api {
	processed := api
	processed.Url = ReplaceVariables(api.Url, variables)
	return processed
}

func ReplaceVariables(text string, variables []models.LocalVariable) string {
	result := text
	for _, variable := range variables {
		placeholder := "{{" + variable.Key + "}}"
		result = strings.ReplaceAll(result, placeholder, variable.Value)
	}
	return result
}
