package temp

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
)

// Temp represents an absolute temperature value. It is abstracted to allow
// conversions between units.
// Note that this value represents _absolute_ temperature, not temperate
// difference. Arithmetic on temperature values probably don't do what you
// expect.
type Temp float64 // in Kelvin

var tempPattern = regexp.MustCompile(`^\s*(\d+\.?\d*)\s*°?\s*(F|C|K)\s*$`)

// Kelvin makes a new Temp from a value in Kelvin
func Kelvin(x float64) Temp {
	return Temp(x)
}

// Celcius makes a new Temp from a value in degrees Celcius
func Celcius(x float64) Temp {
	return Temp(x + 273.15)
}

// Fahrenheit makes a new Temp from a value in degrees Fahrenheit
func Fahrenheit(x float64) Temp {
	return Temp((x-32)*5/9 + 273)
}

// Kelvin retrieves the value of Temp in Kelvin
func (t Temp) Kelvin() float64 {
	return float64(t)
}

// Celcius retrieves the value of Temp in degrees Celcius
func (t Temp) Celcius() float64 {
	return float64(t) - 273.15
}

// Fahrenheit retrieves the value of Temp in degrees Fahrenheit
func (t Temp) Fahrenheit() float64 {
	return float64(t-273.15)*1.8 + 32
}

// MarshalText implements the encoding.TextMarshaller interface
func (t Temp) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t Temp) String() string {
	return fmt.Sprintf("%0.1f°F", t.Fahrenheit())
}

// UnmarshalText implements the encoding.TextUnmarshaller interface
func (t *Temp) UnmarshalText(s []byte) error {
	ss := tempPattern.FindSubmatch(bytes.ToUpper(s))
	if ss == nil {
		return fmt.Errorf("could not parse temperature string of '%s' (is the unit included?)", s)
	}

	valStr, unit := ss[1], ss[2]
	val, err := strconv.ParseFloat(string(valStr), 64)
	if err != nil {
		return err
	}

	switch string(unit) {
	case "F":
		*t = Fahrenheit(val)
		return nil

	case "C":
		*t = Celcius(val)
		return nil

	case "K":
		*t = Kelvin(val)
		return nil

	default:
		panic("unknown unit in temperature")
	}
}
