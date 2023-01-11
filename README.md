# shagent
Smart Home Agent, running on Raspberry Pi

# building
go env -w GOOS=linux
go env -w GOARCH=arm64
make build

# deploying
make deploy
