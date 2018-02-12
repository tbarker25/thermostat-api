package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"reflect"
	"strconv"
	"strings"

	"github.com/buaazp/fasthttprouter"
	"github.com/tbarker25/thermostat-api/internal/temp"
	"github.com/tbarker25/thermostat-api/internal/thermostat"
	"github.com/valyala/fasthttp"
)

func main() {
	addr := flag.String("address", "0.0.0.0:80", "Address to listen on")
	flag.Parse()

	handler := getHandler()

	ln, err := net.Listen("tcp4", *addr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("listening on %s\n", ln.Addr())

	err = (&fasthttp.Server{
		Handler: handler,
	}).Serve(ln)

	if err != nil {
		log.Fatal(err)
	}
}

func getHandler() fasthttp.RequestHandler {
	router := fasthttprouter.New()
	router.PanicHandler = handlePanic
	router.GET("/thermostat", handleThermostatList)
	router.GET("/thermostat/:id", handleThermostatGet)
	router.PATCH("/thermostat/:id", handleThermostatSet)
	return router.Handler
}

func handlePanic(ctx *fasthttp.RequestCtx, err interface{}) {
	log.Print(err)
	sendError(ctx, fasthttp.StatusInternalServerError, "unexpected error: %+v\n", err)
}

func handleThermostatList(ctx *fasthttp.RequestCtx) {
	sendOK(ctx, thermostats)
}

func handleThermostatGet(ctx *fasthttp.RequestCtx) {
	id, err := strconv.ParseUint(ctx.UserValue("id").(string), 10, 32)
	if err != nil {
		sendError(ctx, fasthttp.StatusBadRequest,
			"ID field must be integer, got ID='%s'\n",
			ctx.UserValue("id").(string))
		return
	}

	thermostat := getThermostatByID(uint32(id))
	if thermostat == nil {
		sendError(ctx, fasthttp.StatusNotFound,
			"No thermostat with ID=%d\n", id)
		return
	}

	sendOK(ctx, thermostat)
}

func handleThermostatSet(ctx *fasthttp.RequestCtx) {
	id, err := strconv.ParseUint(ctx.UserValue("id").(string), 10, 32)
	if err != nil {
		sendError(ctx, fasthttp.StatusBadRequest,
			"ID field must be integer, got ID='%s'\n",
			ctx.UserValue("id").(string))
		return
	}

	device := getThermostatByID(uint32(id))
	if device == nil {
		sendError(ctx, fasthttp.StatusNotFound,
			"No device with ID=%d\n", id)
		return
	}

	var body struct {
		Name          *string                   `json:"name"`
		OperatingMode *thermostat.OperatingMode `json:"operatingMode"`
		HeatPoint     *temp.Temp                `json:"heatPoint"`
		CoolPoint     *temp.Temp                `json:"coolPoint"`
		FanMode       *thermostat.FanMode       `json:"fanMode"`
	}

	err = jsonUnmarshalStrict(ctx.PostBody(), &body)
	if err != nil {
		sendError(ctx, fasthttp.StatusBadRequest,
			"Could not unmarshal body: %v", err)
		return
	}

	if body.Name != nil {
		device.SetName(*body.Name)
	}

	if body.OperatingMode != nil {
		device.SetOperatingMode(*body.OperatingMode)
	}

	if body.HeatPoint != nil {
		device.SetHeatPoint(*body.HeatPoint)
	}

	if body.CoolPoint != nil {
		device.SetCoolPoint(*body.CoolPoint)
	}

	if body.FanMode != nil {
		device.SetFanMode(*body.FanMode)
	}

	sendOK(ctx, device)
}

// jsonUnmarshalStrict is a helper function to unmarshal with added strictness.
// unlike json.Unmarshal, this function fails if the input data has extra
// fields Note that this functionality is going to be incorperated into the
// standard library in Go 1.10.
func jsonUnmarshalStrict(data []byte, v interface{}) error {
	tmp := map[string]interface{}{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	t := reflect.ValueOf(v).Elem().Type()
	var unwantedFields []string
nextField:
	for k := range tmp {
		for i := 0; i < t.NumField(); i++ {
			if strings.EqualFold(k, t.Field(i).Name) {
				continue nextField
			}
		}
		unwantedFields = append(unwantedFields, k)
	}
	if unwantedFields != nil {
		return fmt.Errorf("unsupported fields '%s'",
			strings.Join(unwantedFields, "', '"))
	}

	return nil
}

func sendOK(ctx *fasthttp.RequestCtx, v interface{}) {
	resp := struct {
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	}{
		Status: "ok",
		Data:   v,
	}

	body, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		sendError(ctx, fasthttp.StatusBadRequest, "Could not encode JSON: %v\n", err)
		return
	}

	ctx.Success("application/json", body)
}

func sendError(ctx *fasthttp.RequestCtx, statusCode int, format string, a ...interface{}) {
	resp := struct {
		Status       string `json:"status"`
		ErrorMessage string `json:"errorMessage"`
	}{
		Status:       "error",
		ErrorMessage: fmt.Sprintf(format, a...),
	}

	j, err := json.MarshalIndent(&resp, "", "  ")
	if err != nil {
		panic(err)
	}

	ctx.Error(string(j), statusCode)
}
