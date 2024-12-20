.PHONY: default all

production: build_raspi stop_svc_vpn deploy_agent_vpn deploy_config_vpn start_svc_vpn
production_cli: build_cli_raspi deploy_cli_vpn
stage: build_raspi deploy_local
default: build deploy
stage_cli: build_cli_raspi deploy_cli_local
stage_all: stage stage_cli
production_all: production production_cli

build_linux_amd64:
	go env -w GOOS=linux && go env -w GOARCH=amd64 && go build -o bin/shagent ./cmd/shagent/...

build_raspi:
	go env -w GOOS=linux && go env -w GOARCH=arm64 && go build -o bin/shagent ./cmd/shagent/...

stop_svc_vpn:
	plink -pw p dima@172.27.240.3 "sudo systemctl stop shagent.service"

deploy_agent_vpn:
	pscp -pw p ./bin/shagent dima@172.27.240.3:/opt/shagent/shagent

deploy_config_vpn:
	pscp -pw p ./config/prod_config.json dima@172.27.240.3:/opt/shagent/shagent.json

start_svc_vpn:
	plink -pw p dima@172.27.240.3 "sudo systemctl start shagent.service"

deploy_local:
	plink -pw p dima@10.42.0.1 "sudo systemctl stop shagent.service" && pscp -pw p ./bin/shagent dima@10.42.0.1:/opt/shagent/shagent && pscp -pw p ./config/stage_config.json dima@10.42.0.1:/opt/shagent/shagent.json && plink -pw p dima@10.42.0.1 "sudo systemctl start shagent.service"

codegen:
	go env -w GOOS=linux && go env -w GOARCH=amd64 && go generate ./...

test_all:
	env SH_RUN_ALL_TESTS=1 GOOS=linux GOARCH=amd64 go test -count=1 ./...

test:
	env SH_RUN_ALL_TESTS=0 GOOS=linux GOARCH=amd64 go test -count=1 ./...

proto:
	protoc --go_out=./ --go_opt=paths=source_relative --go-grpc_out=./ --go-grpc_opt=paths=source_relative ./pkg/grpcapi/grpcapi.proto

cli_linux_amd64:
	go env -w GOOS=linux && go env -w GOARCH=amd64 && go build -o bin/shagent_cli ./cmd/cli/...

build_cli_raspi:
	go env -w GOOS=linux && go env -w GOARCH=arm64 && go build -o bin/shagent_cli ./cmd/cli/...

deploy_cli_local:
	pscp -pw p ./bin/shagent_cli dima@10.42.0.1:/opt/shagent/shagent_cli

deploy_cli_vpn:
	pscp -pw p ./bin/shagent_cli dima@172.27.240.3:/opt/shagent/shagent_cli

run_cli_local:
	plink -pw p dima@10.42.0.1 "/opt/shagent/shagent_cli"

win:
	go env -w GOOS=windows && go env -w GOARCH=amd64