package main

import (
	"context"
	"fmt"

	"github.com/3rubasa/shagent/pkg/grpcapi"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TurnOnCamLight(c grpcapi.StateProviderClient) {
	res, err := c.TurnOnCamLight(context.Background(), &emptypb.Empty{})
	if err != nil {
		fmt.Println("An error occurred: ", err)
		return
	}

	if res.Error {
		fmt.Println("Failed to turn the cam light on: ", res.ErrorMessage)
		return
	}

	fmt.Println("SUCCESS")
}

func TurnOffCamLight(c grpcapi.StateProviderClient) {
	res, err := c.TurnOffCamLight(context.Background(), &emptypb.Empty{})
	if err != nil {
		fmt.Println("An error occurred: ", err)
		return
	}

	if res.Error {
		fmt.Println("Failed to turn the cam light off: ", res.ErrorMessage)
		return
	}

	fmt.Println("SUCCESS")
}
