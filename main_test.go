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

	location, ok := w.HeaderMap["Location"]

	if !ok {
		t.Fatal("expected header 'Location' not present")
	}

	if !contains(location, targetUrl) {
		t.Fatal("header Location does not contain expected value ")
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
