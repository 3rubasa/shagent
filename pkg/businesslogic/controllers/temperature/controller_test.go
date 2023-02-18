package temperature

import (
	"errors"
	"testing"
	"time"

	mocktempsenordriver "github.com/3rubasa/shagent/pkg/businesslogic/interfaces/mocktempsensordriver"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	mockTDrv := mocktempsenordriver.NewMockTempSensorDriver(mockCtrl)
	tController := New(mockTDrv, 10*time.Minute, 200*time.Millisecond)

	mockTDrv.EXPECT().Start().Times(1)
	mockTDrv.EXPECT().Get().Return(25.0, nil).MinTimes(1)

	err := tController.Start()
	assert.NoError(t, err)

	time.Sleep(time.Second)

	temp, err := tController.Get()
	assert.NoError(t, err)
	assert.Equal(t, temp, 25.0)
}

func Test_CacheTTLExceeded(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	mockTDrv := mocktempsenordriver.NewMockTempSensorDriver(mockCtrl)
	tController := New(mockTDrv, time.Second, 200*time.Millisecond)

	callCount := 0
	mockTDrv.EXPECT().Start().Times(1)
	mockTDrv.EXPECT().Get().DoAndReturn(func() (float64, error) {
		callCount++

		if callCount == 1 {
			return 25.0, nil
		} else {
			return -253.15, errors.New("dummy-error")
		}
	}).MinTimes(2)

	err := tController.Start()
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)

	_, err = tController.Get()
	assert.Error(t, err)
}

func Test_CacheTTLExceededThenSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	mockTDrv := mocktempsenordriver.NewMockTempSensorDriver(mockCtrl)
	tController := New(mockTDrv, time.Second, 200*time.Millisecond)

	callCount := 0
	mockTDrv.EXPECT().Start().Times(1)
	mockTDrv.EXPECT().Get().DoAndReturn(func() (float64, error) {
		callCount++

		if callCount == 1 || callCount == 7 {
			return 25.0, nil
		} else {
			return -253.15, errors.New("dummy-error")
		}
	}).MinTimes(8)

	err := tController.Start()
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)

	temp, err := tController.Get()
	assert.NoError(t, err)
	assert.Equal(t, temp, 25.0)
}
