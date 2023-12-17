FROM golang:1.21-alpine AS builder

RUN apk add -U --no-cache ca-certificates

WORKDIR /build/
COPY . . 
RUN go mod download
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \ 
    go build -ldflags="-s -w" \
    -o main ./cmd/lists-server/main.go

FROM alpine:3.6 as alpine

RUN apk add -U --no-cache ca-certificates

# build actual image ( with root certs from alpine )

FROM scratch

WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/main /dist/
EXPOSE 8080
CMD ["/dist/main"]