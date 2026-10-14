package client

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/api"
)

func writeTestResponse(t *testing.T, w http.ResponseWriter, status int, body string) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write([]byte(body)); err != nil {
		t.Error(err)
	}
}

func TestListUIMounts(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantErr bool
		wantLen int
	}{
		{name: "empty response", body: `{}`, wantErr: true},
		{name: "missing secret", body: `{"data":{"auth":{}}}`, wantErr: true},
		{name: "null secret", body: `{"data":{"secret":null}}`, wantErr: true},
		{name: "invalid secret", body: `{"data":{"secret":[]}}`, wantErr: true},
		{name: "empty mounts", body: `{"data":{"secret":{}}}`},
		{
			name:    "visible mounts",
			body:    `{"data":{"secret":{"KV2/":{"type":"kv","options":{"version":"2"}}}}}`,
			wantLen: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Method != http.MethodGet || r.URL.Path != "/v1/sys/internal/ui/mounts" {
						t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
					}
					writeTestResponse(t, w, http.StatusOK, tt.body)
				}),
			)
			defer server.Close()
			vault, err := api.NewClient(&api.Config{Address: server.URL})
			if err != nil {
				t.Fatal(err)
			}
			mounts, err := listUIMounts(vault)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error for missing or invalid secret mounts")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if mounts == nil || len(mounts) != tt.wantLen {
				t.Fatalf("expected non-nil map with %d mounts, got %#v", tt.wantLen, mounts)
			}
		})
	}
}

func TestNewClientMountDiscovery(t *testing.T) {
	tests := []struct {
		name      string
		status    int
		uiBody    string
		kv1Mounts string
		kv2Mounts string
		want      map[string]int
		wantUI    int
		wantErr   bool
	}{
		{
			name:   "direct discovery without capabilities check",
			status: http.StatusOK,
			want:   map[string]int{"KV1/": 1, "KV2/": 2},
		},
		{
			name:   "403 falls back to UI mounts",
			status: http.StatusForbidden,
			uiBody: `{"data":{"secret":{"KV1/":{"type":"kv"},"KV2/":{"type":"kv","options":{"version":"2"}}}}}`,
			want:   map[string]int{"KV1/": 1, "KV2/": 2},
			wantUI: 1,
		},
		{
			name:   "no visible KV mounts uses default",
			status: http.StatusForbidden,
			uiBody: `{"data":{"secret":{}}}`,
			want:   map[string]int{"secrets/": 2},
			wantUI: 1,
		},
		{
			name:      "invalid UI response preserves explicit mounts",
			status:    http.StatusForbidden,
			uiBody:    `{"data":{"secret":null}}`,
			kv1Mounts: "custom-kv1",
			kv2Mounts: "custom-kv2",
			want:      map[string]int{"custom-kv1/": 1, "custom-kv2/": 2},
			wantUI:    1,
		},
		{name: "404 does not fall back", status: http.StatusNotFound, wantErr: true},
		{name: "500 does not fall back", status: http.StatusInternalServerError, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("VAULT_KV1_MOUNTS", tt.kv1Mounts)
			t.Setenv("VAULT_KV2_MOUNTS", tt.kv2Mounts)
			t.Setenv("VAULT_MAX_RETRIES", "0")
			mountCalls, uiCalls, capabilityCalls := 0, 0, 0
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					switch r.URL.Path {
					case "/v1/sys/mounts":
						mountCalls++
						body := `{"errors":["mount lookup failed"]}`
						if tt.status == http.StatusOK {
							body = `{"data":{"KV1/":{"type":"kv"},"KV2/":{"type":"kv","options":{"version":"2"}}}}`
						}
						writeTestResponse(t, w, tt.status, body)
					case "/v1/sys/internal/ui/mounts":
						uiCalls++
						writeTestResponse(t, w, http.StatusOK, tt.uiBody)
					case "/v1/sys/capabilities-self":
						capabilityCalls++
						writeTestResponse(
							t,
							w,
							http.StatusForbidden,
							`{"errors":["permission denied"]}`,
						)
					default:
						t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
						writeTestResponse(t, w, http.StatusNotFound, `{}`)
					}
				}),
			)
			defer server.Close()
			client, err := NewClient(&VaultConfig{Addr: server.URL, Token: "test-token"})
			if tt.wantErr {
				var responseErr *api.ResponseError
				if !errors.As(err, &responseErr) || responseErr.StatusCode != tt.status {
					t.Errorf("expected HTTP %d error, got %v", tt.status, err)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if !reflect.DeepEqual(client.KVBackends, tt.want) {
				t.Errorf("expected backends %v, got %v", tt.want, client.KVBackends)
			}
			if mountCalls != 1 || uiCalls != tt.wantUI || capabilityCalls != 0 {
				t.Errorf("request counts: mounts=%d UI=%d capabilities=%d; want 1, %d, 0",
					mountCalls, uiCalls, capabilityCalls, tt.wantUI)
			}
		})
	}
}
