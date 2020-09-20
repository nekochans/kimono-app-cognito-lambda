.PHONY: build clean deploy test

build:
	GOOS=linux GOARCH=amd64 go build -o bin/message/signup ./cmd/message/signup/main.go
	chmod +x bin/message/signup
	cp cmd/message/signup/signup-template.html bin/message/signup-template.html

clean:
	rm -rf ./bin

deploy: clean build
	npm run deploy

remove:
	npm run remove

test:
	go test -v ./...

format:
	gofmt -l -s -w .
