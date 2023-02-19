package grpcapi

import (
	context "context"

	"github.com/3rubasa/shagent/pkg/businesslogic"
	empty "github.com/golang/protobuf/ptypes/empty"
)

type impl struct {
	UnimplementedStateProviderServer
	mc businesslogic.MainController
}

func NewImpl(mc businesslogic.MainController) *impl {
	return &impl{
		mc: mc,
	}
}

func (im *impl) GetKitchenTemp(ctx context.Context, e *empty.Empty) (*KitchenTempMessage, error) {
	t, err := im.mc.GetKitchenTemp()
	if err != nil {
		// Invalid temperature value indicate the fact that temperature is not available
		t = 998.0
	}

	return &KitchenTempMessage{T: float32(t)}, nil
}

func (im *impl) GetPowerState(ctx context.Context, e *empty.Empty) (*PowerStateMessage, error) {
	res := &PowerStateMessage{}

	s, err := im.mc.GetPowerState()
	if err != nil {
		res.Error = true
		res.ErrorMessage = err.Error()
	} else {
		res.S = int64(s)
	}

	return res, nil
}

func (im *impl) GetBoilerState(ctx context.Context, e *empty.Empty) (*BoilerStateMessage, error) {
	res := &BoilerStateMessage{}

	s, err := im.mc.GetBoilerState()
	if err != nil {
		res.Error = true
		res.ErrorMessage = err.Error()
	} else {
		res.S = int64(s)
	}

	return res, nil
}

func (im *impl) TurnOnBoiler(ctx context.Context, e *empty.Empty) (*BoilerOpResultMessage, error) {
	res := &BoilerOpResultMessage{}

	err := im.mc.TurnOnBoiler()
	if err != nil {
		res.Error = true
		res.ErrorMessage = err.Error()
	}

	return res, nil
}

func (im *impl) TurnOffBoiler(ctx context.Context, e *empty.Empty) (*BoilerOpResultMessage, error) {
	res := &BoilerOpResultMessage{}

	err := im.mc.TurnOffBoiler()
	if err != nil {
		res.Error = true
		res.ErrorMessage = err.Error()
	}

	return res, nil
}

func (im *impl) GetRoomLightState(ctx context.Context, e *empty.Empty) (*RoomLightStateMessage, error) {
	res := &RoomLightStateMessage{}

	s, err := im.mc.GetRoomLightState()
	if err != nil {
		// Invalid temperature value indicate the fact that temperature is not available
		res.Error = true
		res.ErrorMessage = err.Error()
	} else {
		res.S = int64(s)
	}

	return res, nil
}

func (im *impl) TurnOnRoomLight(ctx context.Context, e *empty.Empty) (*RoomLightOpResultMessage, error) {
	res := &RoomLightOpResultMessage{}

	err := im.mc.TurnOnRoomLight()
	if err != nil {
		res.Error = true
		res.ErrorMessage = err.Error()
	}

	return res, nil
}

func (im *impl) TurnOffRoomLight(ctx context.Context, e *empty.Empty) (*RoomLightOpResultMessage, error) {
	res := &RoomLightOpResultMessage{}

	err := im.mc.TurnOffRoomLight()
	if err != nil {
		res.Error = true
		res.ErrorMessage = err.Error()
	}

	return res, nil
}

func (im *impl) GetCamLightState(ctx context.Context, e *empty.Empty) (*CamLightStateMessage, error) {
	res := &CamLightStateMessage{}

	s, err := im.mc.GetCamLightState()
	if err != nil {
		// Invalid temperature value indicate the fact that temperature is not available
		res.Error = true
		res.ErrorMessage = err.Error()
	} else {
		res.S = int64(s)
	}

	return res, nil
}

func (im *impl) TurnOnCamLight(ctx context.Context, e *empty.Empty) (*CamLightOpResultMessage, error) {
	res := &CamLightOpResultMessage{}

	err := im.mc.TurnOnCamLight()
	if err != nil {
		res.Error = true
		res.ErrorMessage = err.Error()
	}

	return res, nil
}

func (im *impl) TurnOffCamLight(ctx context.Context, e *empty.Empty) (*CamLightOpResultMessage, error) {
	res := &CamLightOpResultMessage{}

	err := im.mc.TurnOffCamLight()
	if err != nil {
		res.Error = true
		res.ErrorMessage = err.Error()
	}

	return res, nil
}

func (im *impl) GetCellBalance(ctx context.Context, e *empty.Empty) (*CellBalanceMessage, error) {
	res := &CellBalanceMessage{}

	b, err := im.mc.GetCellAccBalance()
	if err != nil {
		res.Error = true
		res.ErrorMessage = err.Error()
	} else {
		res.B = float32(b)
	}

	return res, nil
}

func (im *impl) GetCellInetBalance(ctx context.Context, e *empty.Empty) (*CellInetBalanceMessage, error) {
	res := &CellInetBalanceMessage{}

	b, err := im.mc.GetCellInetBalance()
	if err != nil {
		res.Error = true
		res.ErrorMessage = err.Error()
	} else {
		res.B = float32(b)
	}

	return res, nil
}

func (im *impl) GetCellTariff(ctx context.Context, e *empty.Empty) (*CellTariffMessage, error) {
	res := &CellTariffMessage{}

	t, err := im.mc.GetCellTariff()
	if err != nil {
		res.Error = true
		res.ErrorMessage = err.Error()
	} else {
		res.T = t
	}

	return res, nil
}
func (im *impl) GetCellPhoneNumber(ctx context.Context, e *empty.Empty) (*CellPhoneNumberMessage, error) {
	res := &CellPhoneNumberMessage{}

	p, err := im.mc.GetCellPhoneNumber()
	if err != nil {
		res.Error = true
		res.ErrorMessage = err.Error()
	} else {
		res.P = p
	}

	return res, nil
}
