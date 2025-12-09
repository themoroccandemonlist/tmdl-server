FROM golang:1.25-alpine AS builder
WORKDIR /app
RUN go install github.com/a-h/templ/cmd/templ@latest
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN templ generate && \
    CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/web/main.go
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/static ./static
EXPOSE 8080
RUN adduser -D -u 1000 appuser
USER appuser
CMD ["./server"]
