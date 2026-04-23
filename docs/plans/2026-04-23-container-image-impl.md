# Container Image & GHCR Publish Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a minimal scratch-based container image for the hello HTTP server and publish it to GHCR via GitHub Actions on push to main (latest) and on semver tags.

**Architecture:** Multi-stage Dockerfile compiles a static Go binary in golang:1.26-bookworm then copies it into a scratch image. A single GHA workflow fires on push to main and on v* tags, using docker/metadata-action to derive the correct image tags.

**Tech Stack:** Docker multi-stage builds, GitHub Actions, docker/metadata-action@v5, docker/build-push-action@v6, GHCR

---

### Task 1: Dockerfile

**Files:**
- Create: `Dockerfile`

**Step 1: Write the Dockerfile**

```dockerfile
FROM golang:1.26-bookworm AS builder
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /hello .

FROM scratch
COPY --from=builder /hello /hello
ENTRYPOINT ["/hello"]
```

**Step 2: Verify it builds**

```bash
docker build -t hello:local .
```

Expected: build completes, final image ~10MB or less.

**Step 3: Verify the image runs**

```bash
docker run --rm -p 8080:8080 hello:local &
curl -s http://localhost:8080
```

Expected: response with the greeting message. Kill the container after.

**Step 4: Commit**

```bash
git add Dockerfile
git commit -m "feat: add multi-stage Dockerfile"
```

---

### Task 2: GHA Workflow

**Files:**
- Create: `.github/workflows/docker.yml`

**Step 1: Create the workflow directory**

```bash
mkdir -p .github/workflows
```

**Step 2: Write the workflow**

```yaml
name: Docker

on:
  push:
    branches:
      - main
    tags:
      - 'v*'

jobs:
  build-and-push:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write

    steps:
      - uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Log in to GHCR
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ghcr.io/pboyd/hello
          tags: |
            type=raw,value=latest,enable={{is_default_branch}}
            type=semver,pattern={{version}}
            type=semver,pattern={{major}}.{{minor}}
            type=semver,pattern={{major}}

      - name: Build and push
        uses: docker/build-push-action@v6
        with:
          context: .
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max
```

**Step 3: Validate the YAML is well-formed**

```bash
python3 -c "import yaml, sys; yaml.safe_load(open('.github/workflows/docker.yml'))" && echo OK
```

Expected: `OK`

**Step 4: Commit**

```bash
git add .github/workflows/docker.yml
git commit -m "feat: add GHA workflow to publish image to GHCR"
```

---

## Verification

After pushing to GitHub:

1. Check the Actions tab — the `Docker` workflow should trigger and complete green.
2. Check `https://github.com/pboyd/hello/pkgs/container/hello` — the `latest` tag should appear.
3. Pull and run the published image:

```bash
docker run --rm -p 8080:8080 ghcr.io/pboyd/hello:latest
curl http://localhost:8080
```

To test tagged releases, push a tag:

```bash
git tag v0.1.0
git push origin v0.1.0
```

Then check GHCR for `0.1.0`, `0.1`, and `0` tags.
