FROM golang:1.26

WORKDIR /usr/src/api/services/notification
COPY services/notification/go.mod services/notification/go.sum ./
RUN go mod download

COPY services/notification/ .
RUN go build -v -o /usr/local/bin/api ./cmd/api/main.go

CMD ["api"]
