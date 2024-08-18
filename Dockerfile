FROM golang:1.22-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /main ./cmd/api/
FROM alpine:3
COPY --from=builder main /bin/main
ENTRYPOINT ["/bin/main"]
