package api_test

import (
	"encoding/json"
	"fetcher/internal/api"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPI_Health_Unit(t *testing.T) {
	rr := httptest.NewRecorder()
	r := api.New()

	req := httptest.NewRequest(http.MethodGet, api.HealthEndpoint, nil)
	r.ServeHTTP(rr, req)

	var res api.JSONResponse
	dec := json.NewDecoder(rr.Body)
	err := dec.Decode(&res)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	if rr.Code != http.StatusOK {
		t.Fatalf("expected to get statusOK; got %v", rr.Code)
	}

	if res.Message != "OK" {
		t.Fatalf("expected to get message OK; got %s", res.Message)
	}
}
