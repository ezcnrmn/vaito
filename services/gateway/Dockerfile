FROM golang:1.26

# NB: all paths are relative to project's root
WORKDIR /usr/src/api/gen/go/storage
COPY gen/go/storage ./
RUN go mod download

WORKDIR /usr/src/api/services/gateway
COPY services/gateway/go.mod services/gateway/go.sum ./
RUN go mod download

COPY services/gateway/ .
RUN go build -v -o /usr/local/bin/api ./cmd/api/main.go

CMD ["api"]
