package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pboyd/hello/internal/greeting"
)

func TestHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	resp := rec.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	body := rec.Body.String()
	want := greeting.Message()
	if body != want {
		t.Errorf("body = %q, want %q", body, want)
	}
}
