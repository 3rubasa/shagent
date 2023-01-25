package power

import "github.com/3rubasa/shagent/businesslogic/interfaces"

type Controller struct {
	driver interfaces.PowerDriver
}

func New(driver interfaces.PowerDriver) *Controller {
	return &Controller{
		driver: driver,
	}
}

func (c *Controller) Start() error {
	return c.driver.Initialize()
}

func (c *Controller) Stop() error {
	return nil
}

func (c *Controller) Get() (int, error) {
	return c.driver.Get()
}
