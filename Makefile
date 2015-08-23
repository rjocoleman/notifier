all: build

build:
	go build -o bin/notifier

install:
	go get .

vendor:
	godep save -r ./...
