.PHONY: default all

default: build deploy

build: 
	env GOOS=linux GOARCH=arm64 go build

deploy:
	scp ./shagent dima@172.20.10.2:/home/dima/go/src/shagent