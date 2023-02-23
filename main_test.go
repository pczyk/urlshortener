package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// POST /{redirectPath}
func TestHandleRegisterRequestHappyPath(t *testing.T) {
	redirectPath := "google"
	r := httptest.NewRequest(http.MethodPost, "/"+redirectPath, strings.NewReader("https://google.com"))
	w := httptest.NewRecorder()

	applicationState := ApplicationState{registrations: make(map[string]string)}
	applicationState.handleRegisterRequest(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("POST /%s expected HTTP 201, but got %d", redirectPath, w.Code)
	}

	_, ok := applicationState.registrations[redirectPath]

	if !ok {
		t.Fatalf("registrations are expected to contain key '%s'", redirectPath)
	}
}

func TestHandleRegisterRequestMalformedRedirectPath(t *testing.T) {
	redirectPath := ";"
	r := httptest.NewRequest(http.MethodPost, "/"+redirectPath, nil)
	w := httptest.NewRecorder()

	applicationState := ApplicationState{registrations: make(map[string]string)}
	applicationState.handleRegisterRequest(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("POST /%s expected HTTP 400, but got %d", redirectPath, w.Code)
	}

	_, ok := applicationState.registrations[redirectPath]

	if ok {
		t.Fatalf("registrations are NOT expected to contain key '%s'", redirectPath)
	}
}

// GET /{redirectPath}
func TestHandleRedirectRequestHappyPath(t *testing.T) {
	redirectPath := "google"
	targetUrl := "https://google.com"

	r := httptest.NewRequest(http.MethodGet, "/"+redirectPath, nil)
	w := httptest.NewRecorder()

	applicationState := ApplicationState{registrations: map[string]string{redirectPath: targetUrl}}
	applicationState.handleRedirectRequest(w, r)

	if w.Code != http.StatusMovedPermanently {
		t.Fatalf("GET /%s expected HTTP 301, but got %d", redirectPath, w.Code)
	}

	location, ok := w.Result().Header["Location"]

	if !ok {
		t.Fatal("expected header 'Location' not present")
	}

	if !contains(location, targetUrl) {
		t.Fatalf("header Location does not contain expected value %s", targetUrl)
	}
}

func TestHandleRedirectRequestUnregisteredPath(t *testing.T) {
	redirectPath := "google"

	r := httptest.NewRequest(http.MethodGet, "/"+redirectPath, nil)
	w := httptest.NewRecorder()

	applicationState := ApplicationState{registrations: map[string]string{}}
	applicationState.handleRedirectRequest(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("GET /%s expected HTTP 404, but got %d", redirectPath, w.Code)
	}
}

// DELETE /{redirectPath}
func TestHandleUnregisterRequestHappyPath(t *testing.T) {
	redirectPath := "google"

	r := httptest.NewRequest(http.MethodDelete, "/"+redirectPath, nil)
	w := httptest.NewRecorder()

	applicationState := ApplicationState{registrations: map[string]string{redirectPath: "https://google.com"}}

	applicationState.handleUnregisterRequest(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("DELETE /%s expected HTTP 200, but got %d", redirectPath, w.Code)
	}

	_, ok := applicationState.registrations[redirectPath]

	if ok {
		t.Fatalf("registrations are NOT expected to contain key '%s'", redirectPath)
	}
}

// GET /health
func TestHandleHealthRequest(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	applicationState := ApplicationState{}

	applicationState.handleHealthRequest(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /health expected HTTP 200, but got %d", w.Code)
	}
}

func TestValidateRedirectPath(t *testing.T) {
	type test struct {
		input string
		want  bool
	}

	tests := []test{
		{"regular", true},
		{"rEgUl4r", true},
		{"x", true},
		{"", false},
		{"/", false},
		{"regular!", false},
		{";", false},
	}

	for _, tc := range tests {
		got := validateRedirectPath(tc.input)
		if got != tc.want {
			t.Fatalf("expected %v, got %v for input %s", tc.want, got, tc.input)
		}
	}
}

func contains(s []string, v string) bool {
	for _, t := range s {
		if t == v {
			return true
		}
	}
	return false
}
