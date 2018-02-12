package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/valyala/fasthttp"
)

func TestThermostatAPI(t *testing.T) {
	handler := getHandler()

	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	addr := "http://" + ln.Addr().String()

	go func() {
		err = (&fasthttp.Server{
			Handler: handler,
		}).Serve(ln)

		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	}()

	resp, err := http.Get(addr + "/thermostat")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	var body struct {
		Status string
		Data   []map[string]interface{}
	}

	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	want := struct {
		Status string
		Data   []map[string]interface{}
	}{
		Status: "ok",
		Data: []map[string]interface{}{{
			"id":            1,
			"name":          "thermostat-1",
			"currentTemp":   "57.7°F",
			"operatingMode": "off",
			"heatPoint":     "65.0°F",
			"coolPoint":     "80.0°F",
			"fanMode":       "auto",
		}, {
			"id":            2,
			"name":          "thermostat-2",
			"currentTemp":   "68.7°F",
			"operatingMode": "off",
			"heatPoint":     "65.0°F",
			"coolPoint":     "80.0°F",
			"fanMode":       "auto",
		}},
	}

	err = myDeepEquals(&body, &want)
	if err != nil {
		t.Fatalf("readmessage output incorrect: %s", err)
	}

	type Change struct {
		ID    uint32
		Field string
		Value string
	}

	changes := []Change{
		{1, "name", "upstairs bathroom"},
		{1, "operatingMode", "cool"},
		{1, "heatPoint", "65.0°F"},
		{1, "coolPoint", "75.0°F"},
		{1, "fanMode", "auto"},

		{2, "name", "downstairs kitchen"},
		{2, "operatingMode", "heat"},
		{2, "heatPoint", "60.0°F"},
		{2, "coolPoint", "70.0°F"},
		{2, "fanMode", "off"},
	}

	var wg sync.WaitGroup
	wg.Add(len(changes))
	for _, c := range changes {
		go func(c Change) {
			defer wg.Done()
			msg := fmt.Sprintf(`{"%s": "%s"}`, c.Field, c.Value)
			url := fmt.Sprintf(`%s/thermostat/%d`, addr, c.ID)

			req, err := http.NewRequest(
				http.MethodPatch, url, strings.NewReader(msg))
			if err != nil {
				t.Fatalf("readmessage output incorrect: %s", err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("readmessage output incorrect: %s", err)
			}

			var body struct {
				Status       string
				ErrorMessage string
			}

			err = json.NewDecoder(resp.Body).Decode(&body)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if body.Status != "ok" {
				t.Fatalf("got error response: %s", body.ErrorMessage)
			}
		}(c)
	}

	wg.Wait()

	resp, err = http.Get(addr + "/thermostat")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	want = struct {
		Status string
		Data   []map[string]interface{}
	}{
		Status: "ok",
		Data: []map[string]interface{}{{
			"id":            1,
			"name":          "upstairs bathroom",
			"currentTemp":   body.Data[0]["currentTemp"],
			"operatingMode": "cool",
			"heatPoint":     "65.0°F",
			"coolPoint":     "75.0°F",
			"fanMode":       "off",
		}, {
			"id":            2,
			"name":          "downstairs kitchen",
			"currentTemp":   body.Data[1]["currentTemp"],
			"operatingMode": "heat",
			"heatPoint":     "60.0°F",
			"coolPoint":     "70.0°F",
			"fanMode":       "off",
		}},
	}

	err = myDeepEquals(&body, &want)
	if err != nil {
		t.Fatalf("readmessage output incorrect: %s", err)
	}
}
