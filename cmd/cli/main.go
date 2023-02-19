package main

import (
	"fmt"
	"os"

	"github.com/3rubasa/shagent/pkg/grpcapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("ERROR: Failed to connect to the GRPC server: ", err)
		return
	}

	c := grpcapi.NewStateProviderClient(conn)

	if len(os.Args) == 1 {
		PrintState(c)
	} else {
		switch os.Args[1] {
		case "roomlight_on":
			TurnOnRoomLight(c)
		case "roomlight_off":
			TurnOffRoomLight(c)
		case "camlight_on":
			TurnOnCamLight(c)
		case "camlight_off":
			TurnOffCamLight(c)
		case "cell":
			PrintCellInfo(c)
		case "cell_balance":
			PrintCellBalance(c)
		case "cell_inet":
			PrintCellInetBalance(c)
		case "cell_tariff":
			PrintCellTariff(c)
		case "cell_phone":
			PrintCellPhoneNumber(c)
		case "boiler_on":
			TurnOnBoiler(c)
		case "boiler_off":
			TurnOffBoiler(c)
		default:
			fmt.Println("Invalid argument")
		}
	}
}
