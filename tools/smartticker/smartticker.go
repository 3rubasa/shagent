package smartticker

import "time"

const longDuration = 1000 * time.Hour

type SmartTicker struct {
	firstPeriod time.Duration
	mainPeriod  time.Duration
	ticker      *time.Ticker
	timer       *time.Timer
	C           chan time.Time
	done        chan bool
}

func NewSmartTicker(firstPeriod, mainPeriod time.Duration) *SmartTicker {
	st := &SmartTicker{
		firstPeriod: firstPeriod,
		mainPeriod:  mainPeriod,
		done:        make(chan bool),
		C:           make(chan time.Time),
	}

	st.timer = time.NewTimer(firstPeriod)

	st.ticker = time.NewTicker(longDuration)
	st.ticker.Stop()

	return st
}

func (st *SmartTicker) MainLoop() {
	var t time.Time
	for {
		select {
		case t = <-st.timer.C:
			st.C <- t
			st.ticker.Reset(st.mainPeriod)
		case t = <-st.ticker.C:
			st.C <- t
		case <-st.done:
			st.timer.Stop()
			st.ticker.Stop()
			return
		}
	}
}

func (st *SmartTicker) Stop() {
	st.done <- true
}
