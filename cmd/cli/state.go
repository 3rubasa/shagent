package main

import (
	"context"
	"fmt"

	"github.com/3rubasa/shagent/pkg/grpcapi"
	"google.golang.org/protobuf/types/known/emptypb"
)

func PrintState(c grpcapi.StateProviderClient) {

	// Kitchen Temperature
	kt, err := c.GetKitchenTemp(context.Background(), &emptypb.Empty{})
	if err != nil {
		fmt.Println("ERROR: Could not get kitchen temperature: ", err)
	}

	fmt.Print("Kitchen Temperature: .. ")
	if kt.T == 998 {
		fmt.Println("N/A", " deg, C")
	} else {
		fmt.Println(kt.T, " deg, C")
	}

	// Power State
	pws, err := c.GetPowerState(context.Background(), &emptypb.Empty{})
	if err != nil {
		fmt.Println("ERROR: Could not get power state: ", err)
	} else {
		fmt.Print("Power: .............. ")

		if pws.Error {
			fmt.Println("N/A, ", pws.ErrorMessage)
		} else {
			switch pws.S {
			case 0:
				fmt.Println("OFF")
			case 1:
				fmt.Println("ON")
			default:
				fmt.Println("ERROR: Unexpected state value: ", pws.S)
			}
		}
	}

	// RoomLight
	rls, err := c.GetRoomLightState(context.Background(), &emptypb.Empty{})
	if err != nil {
		fmt.Println("ERROR: Could not get room light state: ", err)
	} else {
		fmt.Print("Room Light: ........... ")

		if rls.Error {
			fmt.Println("N/A, ", rls.ErrorMessage)
		} else {
			switch rls.S {
			case 0:
				fmt.Println("OFF")
			case 1:
				fmt.Println("ON")
			default:
				fmt.Println("ERROR: Unexpected state value: ", rls.S)
			}
		}
	}

	// CamLight
	cls, err := c.GetCamLightState(context.Background(), &emptypb.Empty{})
	if err != nil {
		fmt.Println("ERROR: Could not get Cam light state: ", err)
	} else {
		fmt.Print("Cam Light: ............ ")

		if cls.Error {
			fmt.Println("N/A, ", cls.ErrorMessage)
		} else {
			switch cls.S {
			case 0:
				fmt.Println("OFF")
			case 1:
				fmt.Println("ON")
			default:
				fmt.Println("ERROR: Unexpected state value: ", cls.S)
			}
		}
	}

	// Boiler State
	bs, err := c.GetBoilerState(context.Background(), &emptypb.Empty{})
	if err != nil {
		fmt.Println("ERROR: Could not get boiler state: ", err)
	} else {
		fmt.Print("Boiler: .............. ")

		if bs.Error {
			fmt.Println("N/A, ", bs.ErrorMessage)
		} else {
			switch bs.S {
			case 0:
				fmt.Println("OFF")
			case 1:
				fmt.Println("ON")
			default:
				fmt.Println("ERROR: Unexpected state value: ", bs.S)
			}
		}
	}

	// Cell
	PrintCellBalance(c)

	// fmt.Println("Pantry Temperature: ... ", s.PantryTemp, " deg, C")
}
