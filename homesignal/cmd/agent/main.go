package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

var (
	version   = "dev"
	commit    = "unknown"
	buildTime = "unknown"
)

const (
	defaultListenAddr = ":8099"
	defaultConfigDir  = "/config"
	defaultDataDir    = "/data"
	coreAPIBaseURL    = "http://supervisor/core/api/"
)

type DeviceIdentity struct {
	InstallationID string    `json:"installation_id"`
	CreatedAt      time.Time `json:"created_at"`
}

type OptionsState struct {
	Present bool                   `json:"present"`
	Options map[string]interface{} `json:"options,omitempty"`
}

type RuntimeState struct {
	Identity        DeviceIdentity
	Options         OptionsState
	SupervisorToken bool
	CoreAPI         CoreAPIClient
}

type CoreAPIClient struct {
	BaseURL  string
	HasToken bool
}

type readyResponse struct {
	Ready           bool   `json:"ready"`
	Degraded        bool   `json:"degraded"`
	InstallationID  string `json:"installation_id,omitempty"`
	OptionsLoaded   bool   `json:"options_loaded"`
	SupervisorToken bool   `json:"supervisor_token"`
	CoreAPIBaseURL  string `json:"core_api_base_url"`
	Status          string `json:"status"`
	Version         string `json:"version"`
	DegradedReason  string `json:"degraded_reason,omitempty"`
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	state, err := loadRuntimeState(os.Getenv("CONFIG_DIR"), os.Getenv("DATA_DIR"), os.Getenv("SUPERVISOR_TOKEN"))
	if err != nil {
		logger.Error("failed to initialize agent", "error", err)
		os.Exit(1)
	}

	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = defaultListenAddr
	}

	server := &http.Server{
		Addr:              addr,
		Handler:           newRouter(state),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("homesignal agent listening", "addr", addr, "installation_id", state.Identity.InstallationID)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http server failed", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("http server shutdown failed", "error", err)
		os.Exit(1)
	}
}

func loadRuntimeState(configDir, dataDir, supervisorToken string) (RuntimeState, error) {
	if configDir == "" {
		configDir = defaultConfigDir
	}
	if dataDir == "" {
		dataDir = defaultDataDir
	}

	identity, err := ensureIdentity(filepath.Join(configDir, "device.json"))
	if err != nil {
		return RuntimeState{}, err
	}

	options, err := loadOptions(filepath.Join(dataDir, "options.json"))
	if err != nil {
		return RuntimeState{}, err
	}

	hasToken := supervisorToken != ""
	return RuntimeState{
		Identity:        identity,
		Options:         options,
		SupervisorToken: hasToken,
		CoreAPI: CoreAPIClient{
			BaseURL:  coreAPIBaseURL,
			HasToken: hasToken,
		},
	}, nil
}

func ensureIdentity(path string) (DeviceIdentity, error) {
	existing, err := os.ReadFile(path)
	if err == nil {
		var identity DeviceIdentity
		if err := json.Unmarshal(existing, &identity); err != nil {
			return DeviceIdentity{}, fmt.Errorf("read identity: %w", err)
		}
		if identity.InstallationID == "" {
			return DeviceIdentity{}, fmt.Errorf("read identity: installation_id is empty")
		}
		return identity, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return DeviceIdentity{}, fmt.Errorf("read identity: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return DeviceIdentity{}, fmt.Errorf("create identity directory: %w", err)
	}

	identity := DeviceIdentity{
		InstallationID: newInstallationID(),
		CreatedAt:      time.Now().UTC(),
	}

	payload, err := json.MarshalIndent(identity, "", "  ")
	if err != nil {
		return DeviceIdentity{}, fmt.Errorf("encode identity: %w", err)
	}
	payload = append(payload, '\n')

	if err := os.WriteFile(path, payload, 0o600); err != nil {
		return DeviceIdentity{}, fmt.Errorf("write identity: %w", err)
	}

	return identity, nil
}

func newInstallationID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Sprintf("generate installation id: %v", err))
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		hex.EncodeToString(b[0:4]),
		hex.EncodeToString(b[4:6]),
		hex.EncodeToString(b[6:8]),
		hex.EncodeToString(b[8:10]),
		hex.EncodeToString(b[10:16]),
	)
}

func loadOptions(path string) (OptionsState, error) {
	payload, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return OptionsState{Present: false}, nil
	}
	if err != nil {
		return OptionsState{}, fmt.Errorf("read options: %w", err)
	}

	options := map[string]interface{}{}
	if len(payload) > 0 {
		if err := json.Unmarshal(payload, &options); err != nil {
			return OptionsState{}, fmt.Errorf("parse options: %w", err)
		}
	}

	return OptionsState{Present: true, Options: options}, nil
}

func newRouter(state RuntimeState) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthHandler)
	mux.HandleFunc("/readyz", readyHandler(state))
	mux.HandleFunc("/version", versionHandler)
	mux.HandleFunc("/ui", uiHandler(state))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ui", http.StatusFound)
	})
	return mux
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}

func readyHandler(state RuntimeState) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		response := readiness(state)
		writeJSON(w, http.StatusOK, response)
	}
}

func readiness(state RuntimeState) readyResponse {
	degraded := !state.SupervisorToken
	reason := ""
	status := "ready"
	if degraded {
		status = "degraded"
		reason = "SUPERVISOR_TOKEN is not present; Supervisor and Core API calls are disabled"
	}
	return readyResponse{
		Ready:           true,
		Degraded:        degraded,
		InstallationID:  state.Identity.InstallationID,
		OptionsLoaded:   state.Options.Present,
		SupervisorToken: state.SupervisorToken,
		CoreAPIBaseURL:  state.CoreAPI.BaseURL,
		Status:          status,
		Version:         version,
		DegradedReason:  reason,
	}
}

func versionHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"version":    version,
		"commit":     commit,
		"build_time": buildTime,
	})
}

func uiHandler(state RuntimeState) http.HandlerFunc {
	tmpl := template.Must(template.New("ui").Parse(uiHTML))
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, readiness(state)); err != nil {
			http.Error(w, "failed to render status page", http.StatusInternalServerError)
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

const uiHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>HomeSignal</title>
  <style>
    :root { color-scheme: light dark; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; }
    body { margin: 0; padding: 2rem; background: Canvas; color: CanvasText; }
    main { max-width: 720px; margin: 0 auto; }
    h1 { margin: 0 0 0.5rem; font-size: 1.75rem; }
    .status { display: inline-block; margin: 1rem 0; padding: 0.35rem 0.65rem; border-radius: 0.4rem; border: 1px solid ButtonBorder; }
    dl { display: grid; grid-template-columns: minmax(8rem, 14rem) 1fr; gap: 0.75rem 1rem; }
    dt { font-weight: 650; }
    dd { margin: 0; overflow-wrap: anywhere; }
    code { font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; }
  </style>
</head>
<body>
  <main>
    <h1>HomeSignal</h1>
    <p>Home Assistant add-on status and pairing placeholder.</p>
    <div class="status">{{ .Status }}</div>
    <dl>
      <dt>Installation ID</dt><dd><code>{{ .InstallationID }}</code></dd>
      <dt>Version</dt><dd><code>{{ .Version }}</code></dd>
      <dt>Options loaded</dt><dd>{{ .OptionsLoaded }}</dd>
      <dt>Supervisor token</dt><dd>{{ .SupervisorToken }}</dd>
      <dt>Core API</dt><dd><code>{{ .CoreAPIBaseURL }}</code></dd>
      {{ if .DegradedReason }}<dt>Degraded reason</dt><dd>{{ .DegradedReason }}</dd>{{ end }}
    </dl>
  </main>
</body>
</html>
`
