package sim7600

import (
	"fmt"
	"time"

	"github.com/sitec-systems/gmodem"
)

type Sim7600 struct {
	device      string
	readTimeOut time.Duration
}

func New(device string, readTimeOut time.Duration) *Sim7600 {
	return &Sim7600{
		device:      device,
		readTimeOut: readTimeOut,
	}
}

func (s *Sim7600) SendUSSD(ussd string) (string, error) {
	mdm := &gmodem.Modem{
		DevFile:     s.device,
		ReadTimeout: s.readTimeOut,
	}

	err := mdm.Open()
	if err != nil {
		return "", err
	}
	defer mdm.Close()

	at := fmt.Sprintf("AT+CUSD=1,\"%s\",15", ussd)
	r, err := mdm.SendAt(at)
	if err != nil {
		return "", err
	}

	return r, nil
}
