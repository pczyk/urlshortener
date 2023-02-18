package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleRegisterRequestHappyPath(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/google", strings.NewReader("https://google.com"))
	w := httptest.NewRecorder()

	applicationState := ApplicationState{registrations: make(map[string]string)}

	applicationState.handleRegisterRequest(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("POST /google expected HTTP 201, but got %d", w.Code)
	}

	_, ok := applicationState.registrations["google"]

	if !ok {
		t.Fatalf("registrations are expected to contain key 'google'")
	}
}

func TestHandleRegisterRequestBadRedirectPath(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/;", nil)
	w := httptest.NewRecorder()

	applicationState := ApplicationState{registrations: make(map[string]string)}

	applicationState.handleRegisterRequest(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("POST /google expected HTTP 400, but got %d", w.Code)
	}

	_, ok := applicationState.registrations[";"]

	if ok {
		t.Fatalf("registrations are NOT expected to contain key ';'")
	}
}

func TestHandleRedirectRequestHappyPath(t *testing.T) {
	redirectPath := "google"
	targetUrl := "https://google.com"

	r := httptest.NewRequest(http.MethodGet, "/"+redirectPath, nil)
	w := httptest.NewRecorder()

	applicationState := ApplicationState{registrations: map[string]string{redirectPath: targetUrl}}

	applicationState.handleRedirectRequest(w, r)

	if w.Code != http.StatusMovedPermanently {
		t.Fatalf("POST %s expected HTTP 301, but got %d", redirectPath, w.Code)
	}

	location, ok := w.Result().Header["Location"]

	if !ok {
		t.Fatal("expected header 'Location' not present")
	}

	if !contains(location, targetUrl) {
		t.Fatal("header Location does not contain expected value ")
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
