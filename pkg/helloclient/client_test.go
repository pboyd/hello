package helloclient_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pboyd/hello/pkg/helloclient"
)

func TestClientHello(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}))
	defer srv.Close()

	client := helloclient.New(srv.URL)
	got, err := client.Hello()
	if err != nil {
		t.Fatalf("Hello() error: %v", err)
	}
	want := "Hello, World!"
	if got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}
