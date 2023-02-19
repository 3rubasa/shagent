package ltemodulecontroller

import (
	"log"

	"github.com/3rubasa/shagent/pkg/businesslogic/interfaces"
)

type Controller struct {
	driver interfaces.LTEModuleDriver
}

func New(driver interfaces.LTEModuleDriver) *Controller {
	return &Controller{
		driver: driver,
	}
}
func (c *Controller) Start() error {
	return nil
}

func (c *Controller) Stop() error {
	return nil
}

func (c *Controller) GetAccountBalance() (float64, error) {
	b, err := c.driver.SendUSSD(accBalanceUSSD)
	if err != nil {
		log.Println("NOTICE: failed to get a cell account balance: ", err)
		return 0.0, err
	}

	return ParseAccBalance(b)
}

func (c *Controller) GetInetBalance() (float64, error) {
	b, err := c.driver.SendUSSD(inetBalanceUSSD)
	if err != nil {
		log.Println("NOTICE: failed to get a cell inet balance: ", err)
		return 0.0, err
	}

	return ParseInetBalance(b)
}

func (c *Controller) GetTariff() (string, error) {
	t, err := c.driver.SendUSSD(tariffUSSD)
	if err != nil {
		log.Println("NOTICE: failed to get tariff: ", err)
		return "", err
	}

	return ParseTariff(t)
}

func (c *Controller) GetPhoneNumber() (string, error) {
	t, err := c.driver.SendUSSD(tariffUSSD)
	if err != nil {
		log.Println("NOTICE: failed to get phone number: ", err)
		return "", err
	}

	return ParsePhoneNumber(t)
}
