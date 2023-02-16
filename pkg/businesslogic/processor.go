package businesslogic

import "fmt"

type Processor struct {
	r *BusinessRules
	c *BusinessLogicComponents
}

const (
	defaultBoilerTempOn  = 8.0
	defaultBoilerTempOff = 100.0
)

func NewProcessor(r *BusinessRules, c *BusinessLogicComponents) *Processor {
	return &Processor{
		r: r,
		c: c,
	}
}

func (p *Processor) Process(s State) {
	// if both temperatures or weather temp or forecast temp are invalid - we just turn the boiler on and return
	if !s.KitchenTempValid && !s.WindowTempValid || !s.WeatherTempValid || !s.ForecastedTempValid {
		p.c.Boiler.TurnOn()
		return
	}

	var minRoomTemp float64

	if s.KitchenTempValid {
		minRoomTemp = s.KitchenTemp
	}

	if s.WindowTempValid {
		if s.WindowTemp < minRoomTemp {
			minRoomTemp = s.WindowTemp
		}
	}

	// Detect current range
	var boilerTempOn = defaultBoilerTempOn
	var boilerTempOff = defaultBoilerTempOff

	tableFound := false

	for rf, rwrr := range p.r.TempControlTable {
		if s.ForecastedTemp > rf.Min && s.ForecastedTemp <= rf.Max {
			for rw, rr := range rwrr {
				if s.WeatherTemp > rw.Min && s.WeatherTemp <= rw.Max {
					boilerTempOn = rr.Min
					boilerTempOff = rr.Max

					tableFound = true
					break
				}
			}
			break
		}
	}

	if !tableFound {
		p.c.Boiler.TurnOn()
		fmt.Printf("Processor: temperature table has not been found for ForecastTemp: %f and WeatherTemp: %f", s.ForecastedTemp, s.WeatherTemp)
		return
	}

	if minRoomTemp > boilerTempOff {
		p.c.Boiler.TurnOff()
	} else if minRoomTemp < boilerTempOn {
		p.c.Boiler.TurnOn()
	}
}
