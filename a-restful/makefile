generate-swagger:
	s"$(go env GOPATH)/bin/swag" init -g cmd/server/main.go

run:
	go run cmd/server/main.go

test:
	go test ./... -v