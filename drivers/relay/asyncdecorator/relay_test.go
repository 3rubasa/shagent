package asyncdecorator

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/3rubasa/shagent/drivers/relay/asyncdecorator/mockdeviceapi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_GetState(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	devAPI := mockdeviceapi.NewMockDeviceAPI(mockCtrl)
	s := New(devAPI, 200*time.Millisecond)

	devAPI.EXPECT().GetState().Return("on", nil).Times(1)

	expectedState := "on"

	err := s.Start()
	defer s.Stop()

	assert.NoError(t, err)

	actualState, err := s.GetState()

	assert.NoError(t, err)
	assert.Equal(t, expectedState, actualState)
}

func Test_TurnOn_FirstTime(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	devAPI := mockdeviceapi.NewMockDeviceAPI(mockCtrl)
	s := New(devAPI, 200*time.Millisecond)

	devAPI.EXPECT().GetState().Return("on", nil).AnyTimes()
	devAPI.EXPECT().TurnOn().Return(nil).Times(1)

	err := s.Start()
	defer s.Stop()

	assert.NoError(t, err)

	err = s.TurnOn()
	assert.NoError(t, err)
}

func Test_TurnOn_SecondTime(t *testing.T) {
	if os.Getenv("SH_RUN_ALL_TESTS") != "1" {
		t.Skip("Long test, skipping due to SH_RUN_ALL_TEST != 1 ...")
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	devAPI := mockdeviceapi.NewMockDeviceAPI(mockCtrl)
	s := New(devAPI, 200*time.Millisecond)

	tryNum := 0

	devAPI.EXPECT().GetState().DoAndReturn(func() (string, error) {
		if tryNum == 1 {
			return "off", nil
		} else {
			return "on", nil
		}
	}).AnyTimes()

	devAPI.EXPECT().TurnOn().DoAndReturn(func() error {
		tryNum++
		switch tryNum {
		case 1:
			return errors.New("dummy_error")
		case 2:
			return nil
		}

		return nil
	}).Times(2)

	err := s.Start()
	defer s.Stop()

	assert.NoError(t, err)

	err = s.TurnOn()
	assert.Error(t, err)

	time.Sleep(time.Second)
}

func Test_TurnOff_FirstTime(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	devAPI := mockdeviceapi.NewMockDeviceAPI(mockCtrl)
	s := New(devAPI, 200*time.Millisecond)

	devAPI.EXPECT().GetState().Return("off", nil).AnyTimes()
	devAPI.EXPECT().TurnOff().Return(nil).Times(1)

	err := s.Start()
	defer s.Stop()

	assert.NoError(t, err)

	err = s.TurnOff()
	assert.NoError(t, err)
}

func Test_TurnOff_SecondTime(t *testing.T) {
	if os.Getenv("SH_RUN_ALL_TESTS") != "1" {
		t.Skip("Long test, skipping due to SH_RUN_ALL_TEST != 1 ...")
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	devAPI := mockdeviceapi.NewMockDeviceAPI(mockCtrl)
	s := New(devAPI, 200*time.Millisecond)

	tryNum := 0

	devAPI.EXPECT().GetState().DoAndReturn(func() (string, error) {
		if tryNum == 1 {
			return "on", nil
		} else {
			return "off", nil
		}
	}).AnyTimes()

	devAPI.EXPECT().TurnOff().DoAndReturn(func() error {
		tryNum++
		switch tryNum {
		case 1:
			return errors.New("dummy_error")
		case 2:
			return nil
		}

		return nil
	}).Times(2)

	err := s.Start()
	defer s.Stop()

	assert.NoError(t, err)

	err = s.TurnOff()
	assert.Error(t, err)

	time.Sleep(time.Second)
}
