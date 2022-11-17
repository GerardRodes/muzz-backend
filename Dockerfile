FROM golang:alpine as builder
WORKDIR /go/src/muzz
COPY . .
RUN CGO_ENABLED=0 go build -ldflags '-s -w -extldflags=-static' -tags timetzdata -o /go/bin/muzz ./cmd/muzz/main.go

FROM scratch
COPY --from=builder /go/bin/muzz /muzz
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/muzz"]