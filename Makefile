.PHONY: default all

vpn: build_raspi deploy_vpn
local: build_raspi deploy_local
default: build deploy

build_linux_amd64:
	go env -w GOOS=linux && go env -w GOARCH=amd64 && go build -o bin/shagent

build_raspi:
	go env -w GOOS=linux && go env -w GOARCH=arm64 && go build -o bin/shagent

deploy_vpn:
	plink -pw p dima@172.27.208.8 "sudo systemctl stop shagent.service" && pscp -pw p ./bin/shagent dima@172.27.208.8:/home/dima/go/src/shagent/bin && plink -pw p dima@172.27.208.8 "sudo systemctl start shagent.service"

deploy_local:
	plink -pw p dima@10.42.0.1 "sudo systemctl stop shagent.service" && pscp -pw p ./bin/shagent dima@10.42.0.1:/home/dima/go/src/shagent/bin && plink -pw p dima@10.42.0.1 "sudo systemctl start shagent.service"

codegen:
	go env -w GOOS=linux && go env -w GOARCH=amd64 && go generate ./...

test_all:
	env SH_RUN_ALL_TESTS=1 GOOS=linux GOARCH=amd64 go test -count=1 ./...

test:
	env SH_RUN_ALL_TESTS=0 GOOS=linux GOARCH=amd64 go test -count=1 ./...