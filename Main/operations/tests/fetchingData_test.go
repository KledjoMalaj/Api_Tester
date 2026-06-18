package tests

import (
	"GoTuiFrontend/models"
	"GoTuiFrontend/operations"
	"encoding/json"
	"testing"
)

func TestReplaceVariables(t *testing.T) {
	variables := []models.LocalVariable{
		{
			Key:   "token",
			Value: "tokenValue",
		},
		{
			Key:   "base",
			Value: "http://baseUrl",
		},
	}
	input := "Bearer {{token}} from {{base}}"
	expected := "Bearer tokenValue from http://baseUrl"

	result := operations.ReplaceVariables(input, variables)

	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestFormatJSONValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "string value",
			input:    "hello",
			expected: `"hello"`,
		}, {
			name:     "boolean true",
			input:    "true",
			expected: "true",
		}, {
			name:     "boolean false",
			input:    "false",
			expected: "false",
		}, {
			name:     "null value",
			input:    "null",
			expected: "null",
		}, {
			name:     "array value",
			input:    "[1,2,3]",
			expected: "[1,2,3]",
		}, {
			name:     "object value",
			input:    "{role: admin}",
			expected: "{role: admin}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := operations.FormatJSONValue(tt.input)

			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestBuildURLWithQueryParams(t *testing.T) {
	api := models.Api{
		Url: "http://testingApi",
		QueryParams: []models.QueryParam{
			{
				Key:   "page",
				Value: "1",
			},
			{
				Key:   "search",
				Value: "{{searchValue}}",
			},
		},
	}

	m := models.TuiModel{
		LocalVariables: []models.LocalVariable{
			{
				Key:   "searchValue",
				Value: "apples",
			},
		},
	}

	result := operations.BuildURL(api, m)

	expected := "http://testingApi?page=1&search=apples"

	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestParseDataCreatesValidJSON(t *testing.T) {
	api := models.Api{
		BodyField: []models.BodyField{
			{
				Key:   "name",
				Value: "John",
			},
			{
				Key:   "active",
				Value: "true",
			},
			{
				Key:   "role",
				Value: "admin",
			},
		},
	}

	m := models.TuiModel{}

	result := operations.ParseData(api, m)

	var parsed map[string]interface{}

	err := json.Unmarshal([]byte(result), &parsed)
	if err != nil {
		t.Fatalf("expected valid JSON, got error: %v", err)
	}

	if parsed["name"] != "John" {
		t.Errorf("expected name to be John, got %v", parsed["name"])
	}

	if parsed["active"] != true {
		t.Errorf("expected active to be true, got %v", parsed["active"])
	}

	if parsed["role"] != "admin" {
		t.Errorf("expected role to be admin, got %v", parsed["role"])
	}
}
