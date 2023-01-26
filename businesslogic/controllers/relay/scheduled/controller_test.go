package scheduledrelay

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/3rubasa/shagent/businesslogic/interfaces/mockrelaydriver"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	if os.Getenv("SH_RUN_ALL_TESTS") != "1" {
		t.Skip("Long test, skipping due to SH_RUN_ALL_TEST != 1 ...")
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockrelay := mockrelaydriver.NewMockRelayDriver(mockCtrl)

	now := time.Now()
	on1 := now.Add(2 * time.Second)
	off1 := now.Add(3 * time.Second)
	on2 := now.Add(4 * time.Second)
	off2 := now.Add(5 * time.Second)

	on1Str := fmt.Sprintf("%d %d %d * * *", on1.Second(), on1.Minute(), on1.Hour())
	on2Str := fmt.Sprintf("%d %d %d * * *", on2.Second(), on2.Minute(), on2.Hour())
	off1Str := fmt.Sprintf("%d %d %d * * *", off1.Second(), off1.Minute(), off1.Hour())
	off2Str := fmt.Sprintf("%d %d %d * * *", off2.Second(), off2.Minute(), off2.Hour())

	t.Log(on1Str)
	t.Log(on2Str)
	t.Log(off1Str)
	t.Log(off2Str)

	var onTimes, offTimes []string
	onTimes = append(onTimes, on1Str)
	onTimes = append(onTimes, on2Str)
	offTimes = append(offTimes, off1Str)
	offTimes = append(offTimes, off2Str)

	mockrelay.EXPECT().Start().Return(nil).Times(1)
	mockrelay.EXPECT().TurnOn().Return(nil).Times(3) // 3 times because Start() calls TurnOn unconditionally
	mockrelay.EXPECT().TurnOff().Return(nil).Times(2)
	mockrelay.EXPECT().Stop().Times(1)

	r := New(mockrelay, onTimes, offTimes)
	err := r.Start()
	assert.NoError(t, err)

	time.Sleep(10 * time.Second)

	r.Stop()
}
