package main

import (
	"context"
	"fmt"

	"github.com/3rubasa/shagent/pkg/grpcapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	conn, err := grpc.Dial("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("ERROR: Failed to connect to the GRPC server: ", err)
		return
	}

	c := grpcapi.NewStateProviderClient(conn)

	s, err := c.GetState(context.Background(), &emptypb.Empty{})
	if err != nil {
		fmt.Println("ERROR: Could not get state: ", err)
		return
	}

	fmt.Println("Kitchen Temperature: .. ", s.KitchenTemp, " deg, C")
	fmt.Println("Room Light: ........... ", s.RoomLightState)
	fmt.Println("Cam Light: ............ ", s.CamLightState)
	fmt.Println("Boiler State: ......... ", s.BoilerState)
	fmt.Println("Power State: .......... ", s.PowerState)
}
