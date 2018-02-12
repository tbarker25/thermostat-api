package thermostat

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"

	"github.com/tbarker25/thermostat-api/internal/temp"
)

// nextID stores the next ID to assign to a new thermostat. Note that the ID
// field should be unique.
var nextID uint32 = 1

// Thermostat represents the current state of a thermostat
type Thermostat struct {
	id            uint32
	name          string
	operatingMode OperatingMode
	heatPoint     temp.Temp
	coolPoint     temp.Temp
	fanMode       FanMode
	lock          sync.RWMutex
}

// New thermostat with default values
func New() *Thermostat {
	id := nextID
	nextID++
	return &Thermostat{
		id:            id,
		name:          fmt.Sprintf("thermostat-%d", id),
		operatingMode: ModeOff,
		heatPoint:     temp.Fahrenheit(65),
		coolPoint:     temp.Fahrenheit(80),
		fanMode:       FanAuto,
	}
}

// GetID retrieves a thermostat's unique identifier
func (t *Thermostat) GetID() uint32 {
	return t.id
}

// GetName retrieves a thermostat's human-readable name
func (t *Thermostat) GetName() string {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.name
}

// SetName assigns a thermostat's human-readable name
func (t *Thermostat) SetName(s string) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.name = s
}

// GetOperatingMode retrieves a thermostat's operating mode
func (t *Thermostat) GetOperatingMode() OperatingMode {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.operatingMode
}

// SetOperatingMode assigns a thermostat's operating mode
func (t *Thermostat) SetOperatingMode(mode OperatingMode) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.operatingMode = mode
}

// GetFanMode retrieves a thermostat's fan mode
func (t *Thermostat) GetFanMode() FanMode {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.fanMode
}

// SetFanMode assigns a thermostat's fan mode
func (t *Thermostat) SetFanMode(mode FanMode) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.fanMode = mode
}

// GetHeatPoint retrieves the target temperature a thermostat is set to heat to
func (t *Thermostat) GetHeatPoint() temp.Temp {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.heatPoint
}

// SetHeatPoint sets the target temperature a thermostat is set to heat to
func (t *Thermostat) SetHeatPoint(heatPoint temp.Temp) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.heatPoint = heatPoint
}

// GetCoolPoint retrieves the target temperature a thermostat is set to cool to
func (t *Thermostat) GetCoolPoint() temp.Temp {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.coolPoint
}

// SetCoolPoint sets the target temperature a thermostat is set to cool to
func (t *Thermostat) SetCoolPoint(coolPoint temp.Temp) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.coolPoint = coolPoint
}

// GetCurrentTemp retrieves the current temperature of a thermostat
func (t *Thermostat) GetCurrentTemp() temp.Temp {
	// Geometric mean of temperature distribution
	const average = 70

	// Standard Deviation of temperature distribution
	const deviation = 10

	return temp.Fahrenheit(rand.NormFloat64()*deviation + average)
}

// MarshalJSON implements the json.Marshaller interface
func (t *Thermostat) MarshalJSON() ([]byte, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	v := struct {
		ID            uint32        `json:"id"`
		Name          string        `json:"name"`
		CurrentTemp   temp.Temp     `json:"currentTemp"`
		OperatingMode OperatingMode `json:"operatingMode"`
		HeatPoint     temp.Temp     `json:"heatPoint"`
		CoolPoint     temp.Temp     `json:"coolPoint"`
		FanMode       FanMode       `json:"fanMode"`
	}{
		ID:            t.id,
		Name:          t.name,
		CurrentTemp:   t.GetCurrentTemp(),
		OperatingMode: t.operatingMode,
		HeatPoint:     t.heatPoint,
		CoolPoint:     t.coolPoint,
		FanMode:       t.fanMode,
	}

	return json.MarshalIndent(v, "", "  ")
}
