package ltemodulecontroller

import "github.com/3rubasa/shagent/businesslogic/interfaces"

type Controller struct {
	driver interfaces.LTEModuleDriver
}

func New(driver interfaces.LTEModuleDriver) *Controller {
	return &Controller{
		driver: driver,
	}
}

func (c *Controller) GetAccountBalance() (string, error) {
	return c.driver.SendUSSD(accBalanceUSSD)
}

func (c *Controller) GetInetBalance() (string, error) {
	return c.driver.SendUSSD(inetBalanceUSSD)
}

func (c *Controller) GetTariff() (string, error) {
	return c.driver.SendUSSD(tariffUSSD)
}
