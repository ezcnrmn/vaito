FROM golang:1.26
WORKDIR /usr/src/api

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o /usr/local/bin/api ./cmd/api/main.go

EXPOSE 4000

CMD ["api"]
