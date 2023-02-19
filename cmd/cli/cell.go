package main

import (
	"context"
	"fmt"

	"github.com/3rubasa/shagent/pkg/grpcapi"
	"google.golang.org/protobuf/types/known/emptypb"
)

func PrintCellInfo(c grpcapi.StateProviderClient) {
	PrintCellInetBalance(c)
	PrintCellTariff(c)
	PrintCellPhoneNumber(c)
	PrintCellBalance(c)
}

func PrintCellBalance(c grpcapi.StateProviderClient) {
	b, err := c.GetCellBalance(context.Background(), &emptypb.Empty{})
	if err != nil {
		fmt.Println("ERROR: Could not get Cell Balance: ", err)
	} else {
		fmt.Print("Cell Balance: ............ ")

		if b.Error {
			fmt.Println(b.ErrorMessage)
		} else {
			fmt.Println(b.B, " UAH")
		}
	}
}

func PrintCellInetBalance(c grpcapi.StateProviderClient) {
	inet, err := c.GetCellInetBalance(context.Background(), &emptypb.Empty{})
	if err != nil {
		fmt.Println("ERROR: Could not get Cell Internet Balance: ", err)
	} else {
		fmt.Print("Cell Inet Balance: ......... ")

		if inet.Error {
			fmt.Println(inet.ErrorMessage)
		} else {
			fmt.Println(inet.B, " GB")
		}
	}
}

func PrintCellTariff(c grpcapi.StateProviderClient) {
	tariff, err := c.GetCellTariff(context.Background(), &emptypb.Empty{})
	if err != nil {
		fmt.Println("ERROR: Could not get Cell Tariff: ", err)
	} else {
		fmt.Print("Cell Tariff: ........... ")

		if tariff.Error {
			fmt.Println(tariff.ErrorMessage)
		} else {
			fmt.Println(tariff.T)
		}
	}
}

func PrintCellPhoneNumber(c grpcapi.StateProviderClient) {
	phoneNum, err := c.GetCellPhoneNumber(context.Background(), &emptypb.Empty{})
	if err != nil {
		fmt.Println("ERROR: Could not get Cell Phone Num: ", err)
	} else {
		fmt.Print("Cell Phone Num: ........... ")

		if phoneNum.Error {
			fmt.Println(phoneNum.ErrorMessage)
		} else {
			fmt.Println(phoneNum.P)
		}
	}
}
