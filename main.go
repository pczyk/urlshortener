package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"sync"
)

type ApplicationState struct {
	sync.Mutex
	registrations map[string]string
}

func main() {
	applicationState := ApplicationState{registrations: make(map[string]string)}

	http.HandleFunc("/", applicationState.handleRequest)
	http.HandleFunc("/list/", applicationState.handleListRequest)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (as *ApplicationState) handleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		as.handleRedirectRequest(w, r)
	case http.MethodPost:
		as.handleRegisterRequest(w, r)
	case http.MethodDelete:
		as.handleUnregisterRequest(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (as *ApplicationState) handleRedirectRequest(w http.ResponseWriter, r *http.Request) {
	redirectPath := r.URL.String()[1:]
	as.Lock()
	targetUrl, ok := as.registrations[redirectPath]
	as.Unlock()

	if ok {
		http.Redirect(w, r, targetUrl, http.StatusMovedPermanently)
		log.Printf("redirecting request %s -> %s\n", redirectPath, targetUrl)
	} else {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("cannot redirect %s: redirect not found\n", redirectPath)
	}
}

func (as *ApplicationState) handleRegisterRequest(w http.ResponseWriter, r *http.Request) {
	redirectPath := r.URL.String()[1:]

	if !validateRedirectPath(redirectPath) {
		log.Printf("error while registering redirect: invalid redirect path %v", redirectPath)
		http.Error(w, "error while registering redirect: invalid redirect path "+redirectPath, http.StatusBadRequest)
		return
	}

	targetUrl, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println("error while registering redirect: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		as.Lock()
		as.registrations[redirectPath] = string(targetUrl)
		as.Unlock()
		log.Printf("registered redirect %s -> %s\n", redirectPath, targetUrl)
		w.WriteHeader(http.StatusCreated)
	}
}

func (as *ApplicationState) handleUnregisterRequest(w http.ResponseWriter, r *http.Request) {
	redirectPath := r.URL.String()[1:]
	as.Lock()
	_, ok := as.registrations[redirectPath]
	as.Unlock()

	if ok {
		delete(as.registrations, redirectPath)
		log.Printf("unregistered redirect %s\n", redirectPath)
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("cannot unregister redirect %s: redirect not found\n", redirectPath)
	}
}

func (as *ApplicationState) handleListRequest(w http.ResponseWriter, r *http.Request) {
	as.Lock()
	json.NewEncoder(w).Encode(as.registrations)
	as.Unlock()
}

func validateRedirectPath(path string) bool {
	validPath := regexp.MustCompile("^[a-zA-Z0-9]+$")
	return validPath.MatchString(path)
}
