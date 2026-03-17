package models

import (
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

type Storage struct {
	Collections []Collection `json:"collections"`
}
type Collection struct {
	Name           string          `json:"name"`
	Requests       []Api           `json:"requests"`
	LocalVariables []LocalVariable `json:"localVariables"`
}

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type BodyField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type QueryParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type LocalVariable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Response struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Api struct {
	Method      string       `json:"method"`
	Url         string       `json:"url"`
	Headers     []Header     `json:"headers"`
	BodyField   []BodyField  `json:"bodyFields"`
	QueryParams []QueryParam `json:"queryParams"`
	Responses   []Response   `json:"responses"`
}

type ApiResponse struct {
	StatusCode     int
	Status         string
	Body           string
	Headers        http.Header
	RequestHeaders []Header
	ContentType    string
	ContentLength  int64
}

type ErrorMsg struct {
	Message string
}

type ApiResponseMsg struct {
	Response ApiResponse
}

func ShowErrorCommand(message string) tea.Cmd {
	return func() tea.Msg {
		return ErrorMsg{Message: message}
	}
}

type TuiModel struct {
	SelectedApi        Api
	SelectedCollection Collection
	LocalVariables     []LocalVariable
	ApiResponse        ApiResponse
}
