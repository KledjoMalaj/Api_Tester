package operations

import (
	"GoTuiFrontend/models"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFetchDataIntegration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/users" {
			t.Errorf("expected path /users, got %s", r.URL.Path)
		}

		if r.URL.Query().Get("search") != "searchTest" {
			t.Errorf("expected query param search to be searchTest, got %s", r.URL.Query().Get("search"))
		}

		if r.Header.Get("Authorization") != "Bearer tokenTest" {
			t.Errorf("expected Authorization header to be Bearer tokenTest, got %s", r.Header.Get("Authorization"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"success"}`))
	}))
	defer server.Close()

	variables := []models.LocalVariable{
		{
			Key:   "token",
			Value: "tokenTest",
		},
		{
			Key:   "search",
			Value: "searchTest",
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

	response := FetchData(api, m)

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", response.StatusCode)
	}
	if !strings.Contains(response.Body, "success") {
		t.Errorf("expected response body to contain success, got %s", response.Body)
	}
	if response.ContentType != "application/json" {
		t.Errorf("expected content type application/json, got %s", response.ContentType)
	}
}

func TestPostApiFuncIntegration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}
		if r.URL.Path != "/login" {
			t.Errorf("expected path /login, got %s", r.URL.Path)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-type application/json, got %s", r.Header.Get("Content-Type"))
		}
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		body := string(bodyBytes)

		if !strings.Contains(body, `"email": "test@mail.com"`) {
			t.Errorf("expected body to contain email, got %s", body)
		}
		if !strings.Contains(body, `"active": true`) {
			t.Errorf("expected body to contain action true, got %s", body)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`"token": "tokenTest"`))
	}))
	defer server.Close()

	m := models.TuiModel{
		LocalVariables: []models.LocalVariable{
			{
				Key:   "email",
				Value: "test@mail.com",
			},
		},
		SelectedApi: models.Api{
			Method: "POST",
			Url:    server.URL + "/login",
			BodyField: []models.BodyField{
				{
					Key:   "email",
					Value: "{{email}}",
				},
				{
					Key:   "active",
					Value: "true",
				},
			},
		},
	}

	res := PostAPiFunc(m)

	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected status code 201, got %d", res.StatusCode)
	}
	if !strings.Contains(res.Body, "token") {
		t.Errorf("expected response body to contain token, got %s", res.Body)
	}
}
