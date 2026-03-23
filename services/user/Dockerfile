FROM golang:1.26

# NB: all paths are relative to project's root
WORKDIR /usr/src/api/gen/go/user
COPY gen/go/user ./
RUN go mod download

WORKDIR /usr/src/api/services/user
COPY services/user/go.mod services/user/go.sum ./
RUN go mod download

COPY services/user/ .
RUN go build -v -o /usr/local/bin/api ./cmd/api/main.go

CMD ["api"]
