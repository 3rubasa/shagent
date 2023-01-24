package roomlight

import (
	"context"
	"fmt"

	"github.com/3rubasa/shagent/businesslogic/interfaces"
	"github.com/procyon-projects/chrono"
)

type RoomLight struct {
	driver    interfaces.RelayDriver
	scheduler chrono.TaskScheduler
	tasks     []chrono.ScheduledTask
	ontimes   []string
	offtimes  []string
}

func New(driver interfaces.RelayDriver, ontimes, offtimes []string) *RoomLight {
	r := &RoomLight{
		driver:   driver,
		ontimes:  ontimes,
		offtimes: offtimes,
	}

	return r
}

func (r *RoomLight) Start() error {
	// We always start with the light on
	// TODO: make this smart
	r.driver.TurnOn()

	r.scheduler = chrono.NewDefaultTaskScheduler()

	for _, time := range r.ontimes {
		task, err := r.scheduler.ScheduleWithCron(r.turnOnHandler, time)
		if err != nil {
			fmt.Printf("Failed to schedule room light turn on task for time %s, error = %s", time, err.Error())
			return err
		}

		r.tasks = append(r.tasks, task)
	}

	for _, time := range r.offtimes {
		task, err := r.scheduler.ScheduleWithCron(r.turnOffHandler, time)
		if err != nil {
			fmt.Printf("Failed to schedule room light turn off task for time %s, error = %s", time, err.Error())
			return err
		}

		r.tasks = append(r.tasks, task)
	}

	return r.driver.Start()
}

func (r *RoomLight) Stop() error {
	for _, t := range r.tasks {
		t.Cancel()
	}

	sdc := r.scheduler.Shutdown()
	<-sdc

	r.driver.Stop()

	return nil
}

func (r *RoomLight) turnOnHandler(ctx context.Context) {
	err := r.driver.TurnOn()
	if err != nil {
		fmt.Println("Error while turning the room lights on: ", err)
	}
}

func (r *RoomLight) turnOffHandler(ctx context.Context) {
	err := r.driver.TurnOff()
	if err != nil {
		fmt.Println("Error while turning the room lights off: ", err)
	}
}
