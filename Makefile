.PHONY: default all

local: build_raspi deploy_local
remote: build_raspi deploy_remote
default: build deploy

build_linux_amd64:
	go env -w GOOS=linux && go env -w GOARCH=amd64 && go build -o bin/shagent

build_raspi:
	go env -w GOOS=linux && go env -w GOARCH=arm64 && go build -o bin/shagent

deploy_remote:
	plink -P 16177 -pw p dima@pitunnel.com "sudo systemctl stop shagent.service" && pscp -P 16177 -pw p ./bin/shagent dima@pitunnel.com:/home/dima/go/src/shagent/bin && plink -P 16177 -pw p dima@pitunnel.com "sudo systemctl start shagent.service"

deploy_local:
	plink -pw p dima@10.42.0.1 "sudo systemctl stop shagent.service" && pscp -pw p ./bin/shagent dima@10.42.0.1:/home/dima/go/src/shagent/bin && plink -pw p dima@10.42.0.1 "sudo systemctl start shagent.service"

codegen:
	go env -w GOOS=linux && go env -w GOARCH=amd64 && go generate ./...

test_all:
	env SH_RUN_ALL_TESTS=1 go test -count=1 ./...

test:
	env SH_RUN_ALL_TESTS=0 go test -count=1 ./...