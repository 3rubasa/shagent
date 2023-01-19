.PHONY: default all

local: build deploy_local
remote: build deploy_remote
default: build deploy

build:
	go env -w GOOS=linux && go env -w GOARCH=arm64 && go build -o bin/shagent

deploy_remote:
	plink -P 16177 -pw p dima@pitunnel.com "sudo systemctl stop shagent.service" && pscp -P 16177 -pw p ./bin/shagent dima@pitunnel.com:/home/dima/go/src/shagent/bin && plink -P 16177 -pw p dima@pitunnel.com "sudo systemctl start shagent.service"

deploy_local:
	plink -pw p dima@10.42.0.1 "sudo systemctl stop shagent.service" && pscp -pw p ./bin/shagent dima@10.42.0.1:/home/dima/go/src/shagent/bin && plink -pw p dima@10.42.0.1 "sudo systemctl start shagent.service"
