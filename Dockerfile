FROM golang:1.26-bookworm AS builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /hello .

FROM scratch
COPY --from=builder /hello /hello
EXPOSE 8080
USER 65534
ENTRYPOINT ["/hello"]
