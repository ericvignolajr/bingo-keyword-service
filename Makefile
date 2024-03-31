run:
	go run cmd/keyword/main.go

watch:
	air

test: 
	go test -v ./...

build:
	go build -o cmd/keyword/bin/ cmd/keyword/main.go

clean:
	rm -rf cmd/keyword/bin

clean-tmp:
	rm -rf tmp

clean-all: clean clean-tmp

.PHONY: run test clean