package grpcapi

import (
	context "context"

	"github.com/3rubasa/shagent/pkg/businesslogic"
	empty "github.com/golang/protobuf/ptypes/empty"
)

type impl struct {
	UnimplementedStateProviderServer
	bl *businesslogic.BusinessLogic
}

func NewImpl(bl *businesslogic.BusinessLogic) *impl {
	return &impl{
		bl: bl,
	}
}
func (im *impl) GetState(ctx context.Context, e *empty.Empty) (*StateT, error) {
	s := im.bl.GetState()
	state := &StateT{}

	if s.BoilerStateValid {
		state.BoilerState = int64(s.BoilerState)
	} else {
		state.BoilerState = -1
	}

	if s.RoomLightStateValid {
		state.RoomLightState = int64(s.RoomLightState)
	} else {
		state.RoomLightState = -1
	}

	if s.CamLightStateValid {
		state.CamLightState = int64(s.CamLightState)
	} else {
		state.CamLightState = -1
	}

	if s.KitchenTempValid {
		state.KitchenTemp = float32(s.KitchenTemp)
	} else {
		state.KitchenTemp = -275.0
	}

	if s.PowerValid {
		state.PowerState = int64(s.Power)
	} else {
		state.PowerState = -1
	}

	return state, nil
}
