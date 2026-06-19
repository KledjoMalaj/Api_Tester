package tests

import (
	"GoTuiFrontend/models"
	"GoTuiFrontend/operations"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func setupE2ETestFile(t *testing.T) {
	t.Helper()

	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "e2e-test-data.json")

	oldFileName := operations.FileName
	operations.FileName = testFile

	t.Cleanup(func() {
		operations.FileName = oldFileName
	})
}

func TestE2ERequestAndResponse(t *testing.T) {
	setupE2ETestFile(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET method, got %s", r.Method)
		}

		if r.URL.Path != "/profile" {
			t.Errorf("expected path /profile, got %s", r.URL.Path)
		}

		if r.URL.Query().Get("user") != "user1" {
			t.Errorf("expected query user=user1, got %s", r.URL.Query().Get("user"))
		}

		if r.Header.Get("Authorization") != "Bearer tokenTest" {
			t.Errorf("expected Bearer tokenTest, got %s", r.Header.Get("Authorization"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":1,"name":"User","token":"accessToken"}`))
	}))
	defer server.Close()

	storage := models.Storage{
		Collections: []models.Collection{},
	}

	err := operations.AddCollection(storage, storage.Collections, "E2E Collection")
	if err != nil {
		t.Fatalf("failed to add collection: %v", err)
	}

	storage, err = operations.ReadFile()
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	err = operations.AddLocalVariable(storage, 0, []models.LocalVariable{
		{
			Key:   "token",
			Value: "tokenTest",
		},
		{
			Key:   "username",
			Value: "user1",
		},
	})
	if err != nil {
		t.Fatalf("failed to add local variables: %v", err)
	}

	storage, err = operations.ReadFile()
	if err != nil {
		t.Fatalf("filed to read after adding local variables: %v", err)
	}

	err = operations.AddApi(storage, 0, storage.Collections[0].Requests, "GET "+server.URL+"/profile")
	if err != nil {
		t.Fatalf("failed to add API: %v", err)
	}

	storage, err = operations.ReadFile()
	if err != nil {
		t.Fatalf("failed to read file after add ing api: %v", err)
	}

	headers := []models.Header{
		{
			Key:   "Authorization",
			Value: "Bearer {{token}}",
		},
	}

	err = operations.AddHeader(headers, storage, 0, 0)
	if err != nil {
		t.Fatalf("failed to add headers: %v", err)
	}

	storage, err = operations.ReadFile()
	if err != nil {
		t.Fatalf("failed to read file after adding header: %v", err)
	}

	queryParams := []models.QueryParam{
		{
			Key:   "user",
			Value: "{{username}}",
		},
	}

	err = operations.AddQueryParam(queryParams, storage, 0, 0)
	if err != nil {
		t.Fatalf("failed to add query params: %v", err)
	}

	storage, err = operations.ReadFile()
	if err != nil {
		t.Fatalf("failed to read file after adding query params %v", err)
	}

	selectedCollection := storage.Collections[0]
	selectedApi := selectedCollection.Requests[0]

	m := models.TuiModel{
		SelectedCollection: selectedCollection,
		SelectedApi:        selectedApi,
		LocalVariables:     selectedCollection.LocalVariables,
	}

	response := operations.FetchData(selectedApi, m)

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", response.StatusCode)
	}

	values, err := operations.HandleJson(response)
	if err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	if len(values) != 3 {
		t.Fatalf("expected 3 parsed response values, got %d", len(values))
	}

	foundToken := false

	for _, value := range values {
		if value.Key == "token" && value.Value == "accessToken" {
			foundToken = true
		}
	}

	if !foundToken {
		fmt.Errorf("expected parsed response to contain token=accessToken")
	}
}
