.PHONY: build clean deploy test

build:
	GOOS=linux GOARCH=amd64 go build -o bin/custom-message ./cmd/custom-message/main.go
	chmod +x bin/custom-message
	cp cmd/custom-message/signup-template.html bin/signup-template.html
	cp cmd/custom-message/forgot-password-template.html bin/forgot-password-template.html

clean:
	rm -rf ./bin

deploy: clean build
	npm run deploy

remove:
	npm run remove

test:
	go clean -testcache
	go test -v $$(go list ./... | grep -v /node_modules/)

format:
	gofmt -l -s -w .
