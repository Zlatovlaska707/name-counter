package pprof

import (
	"io"
	"net/http"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	addr := ":6060"
	server := NewServer(addr)

	if server.addr != addr {
		t.Errorf("NewServer() addr = %v, want %v", server.addr, addr)
	}

	if server.srv != nil {
		t.Error("NewServer() should not create http server initially")
	}
}

func TestServer_Start(t *testing.T) {
	tests := []struct {
		name      string
		addr      string
		expectErr bool
	}{
		{
			name:      "valid address",
			addr:      ":6061",
			expectErr: false,
		},
		{
			name:      "empty address",
			addr:      "",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewServer(tt.addr)
			err := server.Start()

			if tt.expectErr && err == nil {
				t.Error("Start() expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Start() unexpected error: %v", err)
			}

			// Cleanup
			if server.srv != nil {
				time.Sleep(100 * time.Millisecond)
				defer safeStop(t, server)
			}
		})
	}
}

func TestServer_Start_ServerRunning(t *testing.T) {
	addr := ":6062"
	server := NewServer(addr)

	err := server.Start()
	if err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost" + addr + "/debug/pprof/")
	if err != nil {
		t.Errorf("Failed to connect to pprof server: %v", err)
	}
	if resp != nil {
		safeCloseBody(t, resp.Body)
	}

	err = server.Stop()
	if err != nil {
		t.Errorf("Stop() failed: %v", err)
	}
}

func TestServer_Stop(t *testing.T) {
	server := NewServer(":6063")

	err := server.Stop()
	if err != nil {
		t.Errorf("Stop() on non-started server should not error, got: %v", err)
	}

	err = server.Start()
	if err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	err = server.Stop()
	if err != nil {
		t.Errorf("Stop() failed: %v", err)
	}
}

func TestServer_Start_InvalidAddress(t *testing.T) {
	server := NewServer("invalid:address")
	err := server.Start()

	_ = err

	if server.srv != nil {
		defer safeStop(t, server)
	}
}

func safeStop(t *testing.T, server *Server) {
	if server == nil {
		return
	}
	if err := server.Stop(); err != nil {
		t.Logf("Warning: failed to stop server: %v", err)
	}
}

func safeCloseBody(t *testing.T, body io.ReadCloser) {
	if body == nil {
		return
	}
	if err := body.Close(); err != nil {
		t.Logf("Warning: failed to close response body: %v", err)
	}
}
