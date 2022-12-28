.PHONY: default all

default: build deploy

build: 
	env GOOS=linux GOARCH=arm64 go build -o bin/shagent

deploy:
	scp -P 47210 ./bin/shagent dima@pitunnel.com:/home/dima/go/src/shagent/bin