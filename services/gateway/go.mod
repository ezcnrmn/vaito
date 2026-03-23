module github.com/ezcnrmn/vaito/services/gateway

go 1.25.0

require (
	github.com/ezcnrmn/vaito/gen/go/listing v0.0.0
	github.com/ezcnrmn/vaito/gen/go/user v0.0.0
	github.com/go-playground/validator/v10 v10.30.1
	github.com/julienschmidt/httprouter v1.3.0
	golang.org/x/crypto v0.46.0
	google.golang.org/grpc v1.79.3
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/gabriel-vasile/mimetype v1.4.12 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
)

replace github.com/ezcnrmn/vaito/gen/go/user => ../../gen/go/user/

replace github.com/ezcnrmn/vaito/gen/go/listing => ../../gen/go/listing/
