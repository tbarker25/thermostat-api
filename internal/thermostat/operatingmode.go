package thermostat

import "fmt"

// OperatingMode represents the states that a thermostat can be operating in
type OperatingMode uint8

// possible operating modes
const (
	ModeOff OperatingMode = iota
	ModeCool
	ModeHeat
)

func (m OperatingMode) String() string {
	switch m {
	case ModeOff:
		return "off"

	case ModeCool:
		return "cool"

	case ModeHeat:
		return "heat"

	default:
		panic("unknown OperatingMode")
	}
}

// UnmarshalText implements the encoding.TextUnmarshaller interface
func (m *OperatingMode) UnmarshalText(text []byte) error {
	switch string(text) {
	case "off":
		*m = ModeOff
		return nil

	case "cool":
		*m = ModeCool
		return nil

	case "heat":
		*m = ModeHeat
		return nil

	default:
		return fmt.Errorf("invalid operatingMode '%s'", text)
	}
}

// MarshalText implements the encoding.TextMarshaller interface
func (m OperatingMode) MarshalText() ([]byte, error) {
	return []byte(m.String()), nil
}
