build:
	go build -o api ./cmd/app/main.go && go build -o cli ./cmd/client/main.go

build-api:
	go build -o api ./cmd/app/main.go

build-client:
	go build -o cli ./cmd/client/main.go

run-api:
	go run ./cmd/app/main.go

run-client:
	go run ./cmd/client/main.go