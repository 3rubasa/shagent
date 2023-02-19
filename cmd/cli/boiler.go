package main

import (
	"context"
	"fmt"

	"github.com/3rubasa/shagent/pkg/grpcapi"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TurnOnBoiler(c grpcapi.StateProviderClient) {
	res, err := c.TurnOnBoiler(context.Background(), &emptypb.Empty{})
	if err != nil {
		fmt.Println("An error occurred: ", err)
		return
	}

	if res.Error {
		fmt.Println("Failed to turn the boiler on: ", res.ErrorMessage)
		return
	}

	fmt.Println("SUCCESS")
}

func TurnOffBoiler(c grpcapi.StateProviderClient) {
	res, err := c.TurnOffBoiler(context.Background(), &emptypb.Empty{})
	if err != nil {
		fmt.Println("An error occurred: ", err)
		return
	}

	if res.Error {
		fmt.Println("Failed to turn the boiler off: ", res.ErrorMessage)
		return
	}

	fmt.Println("SUCCESS")
}
