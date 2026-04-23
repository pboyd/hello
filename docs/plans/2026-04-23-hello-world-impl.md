# Hello World Web Service Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a minimal Go HTTP service with an internal greeting package and a public HTTP client package to demonstrate module rename propagation.

**Architecture:** The `main` package runs an HTTP server and delegates message construction to `internal/greeting`. A standalone `pkg/helloclient` package provides a typed HTTP client external consumers can import. No external dependencies — standard library only.

**Tech Stack:** Go (stdlib only) — `net/http`, `testing`, `net/http/httptest`

---

### Task 1: Initialize the Go module

**Files:**
- Create: `go.mod`

**Step 1: Initialize the module**

```bash
go mod init github.com/pboyd/hello
```

Expected output: creates `go.mod` with content:
```
module github.com/pboyd/hello

go 1.24
```

(The Go version will reflect your local toolchain — that's fine.)

**Step 2: Commit**

```bash
git add go.mod
git commit -m "chore: initialize go module"
```

---

### Task 2: Implement `internal/greeting`

**Files:**
- Create: `internal/greeting/greeting.go`
- Create: `internal/greeting/greeting_test.go`

**Step 1: Write the failing test**

Create `internal/greeting/greeting_test.go`:

```go
package greeting_test

import (
	"testing"

	"github.com/pboyd/hello/internal/greeting"
)

func TestMessage(t *testing.T) {
	got := greeting.Message()
	want := "Hello, World!"
	if got != want {
		t.Errorf("Message() = %q, want %q", got, want)
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./internal/greeting/...
```

Expected: build error — `greeting` package does not exist yet.

**Step 3: Write minimal implementation**

Create `internal/greeting/greeting.go`:

```go
package greeting

func Message() string {
	return "Hello, World!"
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./internal/greeting/...
```

Expected: `ok  	github.com/pboyd/hello/internal/greeting`

**Step 5: Commit**

```bash
git add internal/greeting/
git commit -m "feat: add internal greeting package"
```

---

### Task 3: Implement `pkg/helloclient`

**Files:**
- Create: `pkg/helloclient/client.go`
- Create: `pkg/helloclient/client_test.go`

**Step 1: Write the failing test**

Create `pkg/helloclient/client_test.go`:

```go
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
```

**Step 2: Run test to verify it fails**

```bash
go test ./pkg/helloclient/...
```

Expected: build error — `helloclient` package does not exist yet.

**Step 3: Write minimal implementation**

Create `pkg/helloclient/client.go`:

```go
package helloclient

import (
	"io"
	"net/http"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

func (c *Client) Hello() (string, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./pkg/helloclient/...
```

Expected: `ok  	github.com/pboyd/hello/pkg/helloclient`

**Step 5: Commit**

```bash
git add pkg/helloclient/
git commit -m "feat: add helloclient package"
```

---

### Task 4: Implement `main.go`

**Files:**
- Create: `main.go`
- Create: `main_test.go`

**Step 1: Write the failing test**

Create `main_test.go`:

```go
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
```

**Step 2: Run test to verify it fails**

```bash
go test .
```

Expected: build error — `handler` is not defined yet.

**Step 3: Write minimal implementation**

Create `main.go`:

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/pboyd/hello/internal/greeting"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, greeting.Message())
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
```

**Step 4: Run test to verify it passes**

```bash
go test .
```

Expected: `ok  	github.com/pboyd/hello`

**Step 5: Run all tests**

```bash
go test ./...
```

Expected: all three packages pass.

**Step 6: Commit**

```bash
git add main.go main_test.go
git commit -m "feat: add HTTP server"
```

---

### Task 5: Verify the binary builds and runs

**Step 1: Build**

```bash
go build -o hello .
```

Expected: produces a `hello` binary with no errors.

**Step 2: Run and smoke-test**

In one terminal:
```bash
./hello
```

In another:
```bash
curl http://localhost:8080/
```

Expected response: `Hello, World!`

**Step 3: Clean up binary and commit**

```bash
rm hello
git add .gitignore   # if you create one
git commit -m "chore: verify build and smoke test"
```

> Tip: add a `.gitignore` with `hello` (the binary) so it doesn't get committed accidentally.
