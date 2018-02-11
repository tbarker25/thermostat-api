package thermostat

import "fmt"

// FanMode represents the states that a thermostat fan can be in
type FanMode uint8

// possible fan states
const (
	FanOff FanMode = iota
	FanAuto
)

func (m FanMode) String() string {
	switch m {
	case FanOff:
		return "off"

	case FanAuto:
		return "auto"

	default:
		panic("unknown FanMode")
	}
}

// UnmarshalText implements the encoding.TextUnmarshaller interface
func (m *FanMode) UnmarshalText(text []byte) error {
	switch string(text) {
	case "off":
		*m = FanOff
		return nil

	case "auto":
		*m = FanOff
		return nil

	default:
		return fmt.Errorf("invalid fanMode '%s'", text)
	}
}

// MarshalText implements the encoding.TextMarshaller interface
func (m FanMode) MarshalText() ([]byte, error) {
	return []byte(m.String()), nil
}
