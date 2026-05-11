package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadOptionsMissingFile(t *testing.T) {
	options, err := loadOptions(filepath.Join(t.TempDir(), "options.json"))
	if err != nil {
		t.Fatalf("loadOptions returned error: %v", err)
	}
	if options.Present {
		t.Fatal("expected missing options file to be tolerated")
	}
}

func TestLoadOptionsParsesJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "options.json")
	if err := os.WriteFile(path, []byte(`{"log_level":"debug"}`), 0o600); err != nil {
		t.Fatalf("write options: %v", err)
	}

	options, err := loadOptions(path)
	if err != nil {
		t.Fatalf("loadOptions returned error: %v", err)
	}
	if !options.Present {
		t.Fatal("expected options file to be present")
	}
	if options.Options["log_level"] != "debug" {
		t.Fatalf("expected log_level option, got %#v", options.Options)
	}
}

func TestEnsureIdentityCreatesAndReusesInstallationID(t *testing.T) {
	path := filepath.Join(t.TempDir(), "device.json")

	first, err := ensureIdentity(path)
	if err != nil {
		t.Fatalf("ensureIdentity first run: %v", err)
	}
	if first.InstallationID == "" {
		t.Fatal("expected generated installation_id")
	}

	second, err := ensureIdentity(path)
	if err != nil {
		t.Fatalf("ensureIdentity second run: %v", err)
	}
	if second.InstallationID != first.InstallationID {
		t.Fatalf("expected identity reuse, first=%q second=%q", first.InstallationID, second.InstallationID)
	}
}

func TestHealthEndpoint(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	healthHandler(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
}

func TestVersionEndpoint(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/version", nil)

	versionHandler(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
	var response map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response["version"] == "" || response["commit"] == "" || response["build_time"] == "" {
		t.Fatalf("expected version metadata, got %#v", response)
	}
}

func TestReadyEndpointDegradedWithoutSupervisorToken(t *testing.T) {
	state := RuntimeState{
		Identity: DeviceIdentity{InstallationID: "test-installation"},
		CoreAPI:  CoreAPIClient{BaseURL: coreAPIBaseURL},
	}
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/readyz", nil)

	readyHandler(state)(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
	var response readyResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !response.Ready {
		t.Fatal("expected ready response")
	}
	if !response.Degraded {
		t.Fatal("expected degraded response without supervisor token")
	}
	if response.SupervisorToken {
		t.Fatal("expected supervisor token to be absent")
	}
}
