FROM golang:1.19
WORKDIR /go/src/mozz
COPY . .
RUN go build -ldflags "-s -w -extldflags=-static" -mod=vendor -o /go/bin/muzz ./cmd/muzz/main.go
ENTRYPOINT ["/go/bin/muzz"]