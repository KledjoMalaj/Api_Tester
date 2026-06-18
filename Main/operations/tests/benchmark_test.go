package tests

import (
	"GoTuiFrontend/models"
	"GoTuiFrontend/operations"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkFetchDataGET(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"success"}`))
	}))
	defer server.Close()

	variables := []models.LocalVariable{
		{
			Key:   "token",
			Value: "abc123",
		},
		{
			Key:   "search",
			Value: "test",
		},
	}

	api := models.Api{
		Method: "GET",
		Url:    server.URL + "/users",
		Headers: []models.Header{
			{
				Key:   "Authorization",
				Value: "Bearer {{token}}",
			},
		},
		QueryParams: []models.QueryParam{
			{
				Key:   "search",
				Value: "{{search}}",
			},
		},
	}

	m := models.TuiModel{
		LocalVariables: variables,
		SelectedCollection: models.Collection{
			LocalVariables: variables,
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		response := operations.FetchData(api, m)

		if response.StatusCode != http.StatusOK {
			b.Fatalf("expected status 200, got %d", response.StatusCode)
		}
	}
}
