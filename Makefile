run:
	go run cmd/keyword/main.go

test: 
	go test -v ./...

build:
	go build -o cmd/keyword/bin/ cmd/keyword/main.go

clean: 
	rm -rf cmd/keyword/bin

.PHONY: run test clean