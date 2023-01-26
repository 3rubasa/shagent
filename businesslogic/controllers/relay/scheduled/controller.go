package scheduledrelay

import (
	"context"
	"errors"
	"fmt"

	"github.com/3rubasa/shagent/businesslogic/interfaces"
	"github.com/procyon-projects/chrono"
)

type ScheduledRelay struct {
	driver    interfaces.RelayDriver
	scheduler chrono.TaskScheduler
	tasks     []chrono.ScheduledTask
	ontimes   []string
	offtimes  []string
}

func New(driver interfaces.RelayDriver, ontimes, offtimes []string) *ScheduledRelay {
	r := &ScheduledRelay{
		driver:   driver,
		ontimes:  ontimes,
		offtimes: offtimes,
	}

	return r
}

func (r *ScheduledRelay) Start() error {
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

func (r *ScheduledRelay) Stop() error {
	for _, t := range r.tasks {
		t.Cancel()
	}

	sdc := r.scheduler.Shutdown()
	<-sdc

	r.driver.Stop()

	return nil
}

func (r *ScheduledRelay) Get() (int, error) {
	s, err := r.driver.GetState()
	if err != nil {
		fmt.Println("Failed to get state of the room light: ", err)
		return 0, err
	}

	switch s {
	case "on":
		return 1, nil
	case "off":
		return 0, nil
	default:
		fmt.Println("Unexpected state: ", s)
		return 0, errors.New("unexpected state")
	}
}

func (r *ScheduledRelay) TurnOn() error {
	return r.driver.TurnOn()
}

func (r *ScheduledRelay) TurnOff() error {
	return r.driver.TurnOff()
}

func (r *ScheduledRelay) turnOnHandler(ctx context.Context) {
	err := r.driver.TurnOn()
	if err != nil {
		fmt.Println("Error while turning the room lights on: ", err)
	}
}

func (r *ScheduledRelay) turnOffHandler(ctx context.Context) {
	err := r.driver.TurnOff()
	if err != nil {
		fmt.Println("Error while turning the room lights off: ", err)
	}
}
