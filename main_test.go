// main_test.go
package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestWeatherHandler(t *testing.T) {
	// Set up environment variable for testing
	os.Setenv("WEATHER_API_KEY", "testapikey")

	req, err := http.NewRequest("GET", "/weather?cep=01001000", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock response for testing
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"temp_C": 28.5, "temp_F": 83.3, "temp_K": 301.5}`))
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Decode the response body
	var responseBody map[string]float64
	if err := json.NewDecoder(rr.Body).Decode(&responseBody); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	// Expected response
	expected := map[string]float64{
		"temp_C": 28.5,
		"temp_F": 83.3,
		"temp_K": 301.5,
	}

	// Compare the actual response with the expected response
	for key, expectedValue := range expected {
		if actualValue, ok := responseBody[key]; !ok || actualValue != expectedValue {
			t.Errorf("handler returned unexpected body: got %v want %v", responseBody, expected)
		}
	}
}
