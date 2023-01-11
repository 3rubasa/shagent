.PHONY: default all

default: build deploy

build: 
	go build -o bin/shagent

deploy:
	pscp -P 16177 -pw p ./bin/shagent dima@pitunnel.com:/home/dima/go/src/shagent/bin