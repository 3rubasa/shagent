package main

import (
	"context"
	"fmt"

	"github.com/3rubasa/shagent/pkg/grpcapi"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TurnOnRoomLight(c grpcapi.StateProviderClient) {
	res, err := c.TurnOnRoomLight(context.Background(), &emptypb.Empty{})
	if err != nil {
		fmt.Println("An error occurred: ", err)
		return
	}

	if res.Error {
		fmt.Println("Failed to turn the room light on: ", res.ErrorMessage)
		return
	}

	fmt.Println("SUCCESS")
}

func TurnOffRoomLight(c grpcapi.StateProviderClient) {
	res, err := c.TurnOffRoomLight(context.Background(), &emptypb.Empty{})
	if err != nil {
		fmt.Println("An error occurred: ", err)
		return
	}

	if res.Error {
		fmt.Println("Failed to turn the room light off: ", res.ErrorMessage)
		return
	}

	fmt.Println("SUCCESS")
}
