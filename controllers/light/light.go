package light

import (
	"context"
	"fmt"
	"log"

	"github.com/3rubasa/shagent/controllers"
	"github.com/procyon-projects/chrono"
	"github.com/stianeikeland/go-rpio"
)

const chan1Pin = 26

type lightController struct {
	scheduler           chrono.TaskScheduler
	lightOnMorningTask  chrono.ScheduledTask
	lightOffMorningTask chrono.ScheduledTask
	lightOnEveningTask  chrono.ScheduledTask
	lightOffEveningTask chrono.ScheduledTask
	pin                 rpio.Pin
}

var controllerSingleton *lightController

func New() controllers.LightController {
	if controllerSingleton == nil {
		controllerSingleton = &lightController{}
	}

	return controllerSingleton
}

func (p *lightController) Initialize() error {
	err := rpio.Open()
	if err != nil {
		fmt.Printf("Failed to open rpio: %s", err.Error())
		return err
	}

	p.pin = rpio.Pin(chan1Pin)
	p.pin.Output()
	// TODO: Init light properly

	p.scheduler = chrono.NewDefaultTaskScheduler()

	return nil
}

func (p *lightController) Start() error {
	var err error

	p.lightOnMorningTask, err = p.scheduler.ScheduleWithCron(func(ctx context.Context) {
		log.Print("Turning lights on - morning")
		p.pin.Low()
	}, "0 45 06 * * *")

	if err != nil {
		log.Println("Error scheduling a morning lights on task: ", err)
		return err
	}

	p.lightOffMorningTask, err = p.scheduler.ScheduleWithCron(func(ctx context.Context) {
		log.Print("Turning lights off - morning")
		p.pin.High()
	}, "0 15 08 * * *")

	if err != nil {
		log.Println("Error scheduling a lights on task: ", err)
		return err
	}

	p.lightOnEveningTask, err = p.scheduler.ScheduleWithCron(func(ctx context.Context) {
		log.Print("Turning lights on - evening")
		p.pin.Low()
	}, "0 10 17 * * *")

	if err != nil {
		log.Println("Error scheduling a evening lights on task: ", err)
		return err
	}

	p.lightOffEveningTask, err = p.scheduler.ScheduleWithCron(func(ctx context.Context) {
		log.Print("Turning lights off - evening")
		p.pin.High()
	}, "0 12 01 * * *")

	if err != nil {
		log.Println("Error scheduling a evening lights off task: ", err)
		return err
	}

	return nil
}

func (p *lightController) Stop() {
	rpio.Close()
	p.lightOffEveningTask.Cancel()
	p.lightOffMorningTask.Cancel()
	p.lightOnEveningTask.Cancel()
	p.lightOnMorningTask.Cancel()
	p.scheduler.Shutdown()
}
