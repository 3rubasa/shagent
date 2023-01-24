package businesslogic

import (
	"fmt"
)

type BusinessLogic struct {
	//boiler    interfaces.RelayDriver
	roomLight RoomLightController
	//camLight  interfaces.RelayDriver
}

func New(roomLight RoomLightController) *BusinessLogic {
	return &BusinessLogic{
		//boiler:    boiler,
		roomLight: roomLight,
		//camLight:  camLight,
	}
}

func (b *BusinessLogic) Start() error {
	var err error

	// err = b.boiler.Start()
	// if err != nil {
	// 	fmt.Println("Failed to start boiler relay: ", err)
	// }

	err = b.roomLight.Start()
	if err != nil {
		fmt.Println("Failed to start room light relay controller: ", err)
	}

	// err = b.camLight.Start()
	// if err != nil {
	// 	fmt.Println("Failed to start cam light relay controller: ", err)
	// }

	return nil
}
