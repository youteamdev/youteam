package server

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"youteam.dev/youteam/allinone/internal/config"
)

func TestEmbeddedIndexIsServedAtRoot(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	NewHTTPHandler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if !strings.Contains(recorder.Body.String(), "YouTeam is running.") {
		t.Fatalf("body does not contain embedded index content")
	}
}

func TestRunHTTPServerStopsWhenContextIsCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)

	go func() {
		errCh <- RunHTTPServer(ctx, config.Config{Host: "127.0.0.1", Port: "0"})
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("RunHTTPServer() error = %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("RunHTTPServer() did not stop after context cancellation")
	}
}

func TestRunHTTPServerReturnsListenError(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Listen() error = %v", err)
	}
	defer listener.Close()

	_, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		t.Fatalf("SplitHostPort() error = %v", err)
	}

	err = RunHTTPServer(context.Background(), config.Config{Host: "127.0.0.1", Port: port})
	if err == nil {
		t.Fatal("RunHTTPServer() error = nil, want listen error")
	}
}
