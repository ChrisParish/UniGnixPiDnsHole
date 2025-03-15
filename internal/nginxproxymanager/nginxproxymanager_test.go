package nginxproxymanager

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/tokens" {
			t.Errorf("Expected path /api/tokens, got %s", r.URL.Path)
		}
		response := AuthResponse{Token: "test"}

		responseJson, err := json.Marshal(response)
		if err != nil {
			t.Errorf("Error marshalling response: %s", err)
		}
		w.WriteHeader(200)
		w.Write([]byte(responseJson))

	}))
	defer server.Close()
	token, err := auth("user", "pass", server.URL)
	if err != nil {
		t.Errorf("Error authenticating: %s", err)
	}
	if token != "test" {
		t.Errorf("Expected token test, got %s", token)
	}
}
