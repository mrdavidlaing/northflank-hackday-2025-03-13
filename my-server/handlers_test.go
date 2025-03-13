package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInfoHandler(t *testing.T) {
	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/info", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(infoHandler)

	// Call the handler with our request and response recorder
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the content type
	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, expectedContentType)
	}

	// Check the response body
	var response InfoResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	// Verify the version
	if response.Version != Version {
		t.Errorf("handler returned unexpected version: got %v want %v",
			response.Version, Version)
	}
}

func TestInfoHandlerMethodNotAllowed(t *testing.T) {
	// Test with methods that should not be allowed
	methods := []string{"POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		// Create a request with the current method
		req, err := http.NewRequest(method, "/info", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a ResponseRecorder to record the response
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(infoHandler)

		// Call the handler with our request and response recorder
		handler.ServeHTTP(rr, req)

		// Check the status code is 405 Method Not Allowed
		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("%s method: handler returned wrong status code: got %v want %v",
				method, status, http.StatusMethodNotAllowed)
		}
	}
} 