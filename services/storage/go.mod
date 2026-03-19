module github.com/ezcnrmn/vaito/services/storage

go 1.25.0

require github.com/lib/pq v1.11.2

require (
	github.com/ezcnrmn/vaito/gen/go/storage v0.0.0
	google.golang.org/grpc v1.79.2
)

require (
	github.com/golang/protobuf v1.5.4 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/ezcnrmn/vaito/gen/go/storage => ../../gen/go/storage/
