package tests

import (
	"GoTuiFrontend/models"
	"GoTuiFrontend/operations"
	"os"
	"path/filepath"
	"testing"
)

func setupTestFile(t *testing.T) {
	t.Helper()

	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test-api-data.json")

	oldFileName := operations.FileName
	operations.FileName = testFile

	t.Cleanup(func() {
		operations.FileName = oldFileName
	})
}

func TestWriteAndReadFile(t *testing.T) {
	setupTestFile(t)

	storage := models.Storage{
		Collections: []models.Collection{
			{
				Name: "Test Collection",
			},
		},
	}

	err := operations.WriteFile(storage)
	if err != nil {
		t.Fatalf("expected WriteFile to succeed, got error: %v", err)
	}

	result, err := operations.ReadFile()
	if err != nil {
		t.Fatalf("expected ReadFile to succeed, got error: %v", err)
	}

	if len(result.Collections) != 1 {
		t.Fatalf("expected 1 collection, got %d", len(result.Collections))
	}

	if result.Collections[0].Name != "Test Collection" {
		t.Errorf("expected collection name to be Test Collection, got %q", result.Collections[0].Name)
	}
}

func TestReadFileCreatesFileIfMissing(t *testing.T) {
	setupTestFile(t)

	_, err := os.Stat(operations.FileName)
	if !os.IsNotExist(err) {
		t.Fatalf("expected test file to not exist before ReadFile")
	}

	storage, err := operations.ReadFile()
	if err != nil {
		t.Fatalf("expected ReadFile to create file without error, got: %v", err)
	}

	if len(storage.Collections) != 0 {
		t.Errorf("expected empty collections, got %d", len(storage.Collections))
	}

	_, err = os.Stat(operations.FileName)
	if err != nil {
		t.Fatalf("expected file to be created, got error: %v", err)
	}
}

func TestAddCollection(t *testing.T) {
	setupTestFile(t)

	storage := models.Storage{
		Collections: []models.Collection{},
	}

	err := operations.AddCollection(storage, storage.Collections, "Backend APIs")
	if err != nil {
		t.Fatalf("expected AddCollection to succeed, got error: %v", err)
	}

	result, err := operations.ReadFile()
	if err != nil {
		t.Fatalf("expected ReadFile to succeed, got error: %v", err)
	}

	if len(result.Collections) != 1 {
		t.Fatalf("expected 1 collection, got %d", len(result.Collections))
	}

	if result.Collections[0].Name != "Backend APIs" {
		t.Errorf("expected collection name Backend APIs, got %q", result.Collections[0].Name)
	}
}

func TestAddCollectionEmptyNameReturnsError(t *testing.T) {
	setupTestFile(t)

	storage := models.Storage{
		Collections: []models.Collection{},
	}

	err := operations.AddCollection(storage, storage.Collections, "")
	if err == nil {
		t.Fatal("expected error for empty collection name, got nil")
	}
}

func TestAddApi(t *testing.T) {
	setupTestFile(t)

	storage := models.Storage{
		Collections: []models.Collection{
			{
				Name:     "Test Collection",
				Requests: []models.Api{},
			},
		},
	}

	err := operations.AddApi(storage, 0, storage.Collections[0].Requests, "GET https://api.example.com/users")
	if err != nil {
		t.Fatalf("expected AddApi to succeed, got error: %v", err)
	}

	result, err := operations.ReadFile()
	if err != nil {
		t.Fatalf("expected ReadFile to succeed, got error: %v", err)
	}

	requests := result.Collections[0].Requests

	if len(requests) != 1 {
		t.Fatalf("expected 1 API request, got %d", len(requests))
	}

	if requests[0].Method != "GET" {
		t.Errorf("expected method GET, got %q", requests[0].Method)
	}

	if requests[0].Url != "https://api.example.com/users" {
		t.Errorf("expected URL https://api.example.com/users, got %q", requests[0].Url)
	}
}

func TestAddApiInvalidFormatReturnsError(t *testing.T) {
	setupTestFile(t)

	storage := models.Storage{
		Collections: []models.Collection{
			{
				Name:     "Test Collection",
				Requests: []models.Api{},
			},
		},
	}

	err := operations.AddApi(storage, 0, storage.Collections[0].Requests, "GET")
	if err == nil {
		t.Fatal("expected error for invalid API format, got nil")
	}
}

func TestDeleteCollection(t *testing.T) {
	setupTestFile(t)

	storage := models.Storage{
		Collections: []models.Collection{
			{
				Name: "Collection One",
			},
			{
				Name: "Collection Two",
			},
		},
	}

	remainingCollections, err := operations.DeleteCollection(storage.Collections[0], storage)
	if err != nil {
		t.Fatalf("expected DeleteCollection to succeed, got error: %v", err)
	}

	if len(remainingCollections) != 1 {
		t.Fatalf("expected 1 remaining collection, got %d", len(remainingCollections))
	}

	if remainingCollections[0].Name != "Collection Two" {
		t.Errorf("expected remaining collection to be Collection Two, got %q", remainingCollections[0].Name)
	}
}

func TestAddLocalVariable(t *testing.T) {
	setupTestFile(t)

	storage := models.Storage{
		Collections: []models.Collection{
			{
				Name:           "Auth Collection",
				LocalVariables: []models.LocalVariable{},
			},
		},
	}

	variables := []models.LocalVariable{
		{
			Key:   "token",
			Value: "abc123",
		},
	}

	err := operations.AddLocalVariable(storage, 0, variables)
	if err != nil {
		t.Fatalf("expected AddLocalVariable to succeed, got error: %v", err)
	}

	result, err := operations.ReadFile()
	if err != nil {
		t.Fatalf("expected ReadFile to succeed, got error: %v", err)
	}

	if len(result.Collections[0].LocalVariables) != 1 {
		t.Fatalf("expected 1 local variable, got %d", len(result.Collections[0].LocalVariables))
	}

	if result.Collections[0].LocalVariables[0].Key != "token" {
		t.Errorf("expected variable key token, got %q", result.Collections[0].LocalVariables[0].Key)
	}

	if result.Collections[0].LocalVariables[0].Value != "abc123" {
		t.Errorf("expected variable value abc123, got %q", result.Collections[0].LocalVariables[0].Value)
	}
}

func TestHandleJson(t *testing.T) {
	response := models.ApiResponse{
		Body: `{"token":"abc123","id":1,"active":true}`,
	}

	result, err := operations.HandleJson(response)
	if err != nil {
		t.Fatalf("expected HandleJson to succeed, got error: %v", err)
	}

	if len(result) != 3 {
		t.Fatalf("expected 3 response values, got %d", len(result))
	}

	expectedKeys := []string{"active", "id", "token"}

	for i, expectedKey := range expectedKeys {
		if result[i].Key != expectedKey {
			t.Errorf("expected key %q at index %d, got %q", expectedKey, i, result[i].Key)
		}
	}
}

func TestHandleJsonInvalidJsonReturnsError(t *testing.T) {
	response := models.ApiResponse{
		Body: `{invalid json}`,
	}

	_, err := operations.HandleJson(response)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestHandleJsonNestedArrayValues(t *testing.T) {
	response := models.ApiResponse{
		Body: `{"users":[{"id":1,"name":"Ada"},{"id":2,"name":"Lin"}],"meta":{"count":2}}`,
	}

	result, err := operations.HandleJson(response)
	if err != nil {
		t.Fatalf("expected HandleJson to succeed, got error: %v", err)
	}

	values := map[string]string{}
	for _, item := range result {
		values[item.Key] = item.Value
	}

	expected := map[string]string{
		"meta.count":    "2",
		"users[0].id":   "1",
		"users[0].name": "Ada",
		"users[1].id":   "2",
		"users[1].name": "Lin",
	}

	for key, expectedValue := range expected {
		if values[key] != expectedValue {
			t.Errorf("expected %s=%q, got %q", key, expectedValue, values[key])
		}
	}
}
