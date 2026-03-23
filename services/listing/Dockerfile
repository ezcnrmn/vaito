FROM golang:1.26

# NB: all paths are relative to project's root
WORKDIR /usr/src/api/gen/go/listing
COPY gen/go/listing ./
RUN go mod download

WORKDIR /usr/src/api/services/listing
COPY services/listing/go.mod services/listing/go.sum ./
RUN go mod download

COPY services/listing/ .
RUN go build -v -o /usr/local/bin/api ./cmd/api/main.go

CMD ["api"]
