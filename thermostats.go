package main

import (
	"sort"

	"github.com/tbarker25/thermostat-api/internal/thermostat"
)

var thermostats []*thermostat.Thermostat

func init() {
	thermostats = append(thermostats, thermostat.New())
	thermostats = append(thermostats, thermostat.New())
	sort.Slice(thermostats, func(i, j int) bool {
		return thermostats[i].GetID() < thermostats[j].GetID()
	})
}

func getThermostatByID(id uint32) *thermostat.Thermostat {
	i := sort.Search(len(thermostats), func(i int) bool {
		return thermostats[i].GetID() >= id
	})

	if i >= len(thermostats) || thermostats[i].GetID() != id {
		return nil
	}

	return thermostats[i]
}
