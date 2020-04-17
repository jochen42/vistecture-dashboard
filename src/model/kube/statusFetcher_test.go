package kube

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckHealth_AllHealthy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{\"services\": []}"))

	}))
	defer server.Close()

	healthStatusOfService, reason, _ := checkHealth(AppDeploymentInfo{}, server.URL, "/")
	if !healthStatusOfService {
		t.Errorf("healthStatusOfService should be true %v", reason)
	}
}

func TestCheckHealth_UnhealthyService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{\"services\": [{\"name\": \"dummy\", \"alive\": false, \"details\": \"dummy\"}]}"))
	}))
	defer server.Close()

	healthStatusOfService, _, _ := checkHealth(AppDeploymentInfo{}, server.URL, "/nonexistingpath")
	if healthStatusOfService {
		t.Errorf("healthStatusOfService should be false")
	}
}

func TestCheckHealth_UserAgentIsSet(t *testing.T) {
	expectedUA := "VistectureDashboard"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ua := r.Header.Get("User-Agent")
		if ua != expectedUA {
			w.WriteHeader(500)
			t.Errorf("expected user-agent to be '%s' but was '%s'", expectedUA, ua)
		} else {
			w.WriteHeader(200)
		}
	}))
	defer server.Close()

	healthStatusOfService, _, _ := checkHealth(AppDeploymentInfo{}, server.URL, "/")
	if healthStatusOfService {
		t.Errorf("user-agent assertion failed")
	}
}
