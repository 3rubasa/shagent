.PHONY: default all

local: build deploy_local
default: build deploy

build: 
	go build -o bin/shagent

deploy:
	pscp -P 16177 -pw p ./bin/shagent dima@pitunnel.com:/home/dima/go/src/shagent/bin

deploy_local:
	pscp -pw p ./bin/shagent dima@10.42.0.1:/home/dima/go/src/shagent/bin