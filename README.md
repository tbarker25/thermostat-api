# thermostat-api
Toy implementation of a thermostat API in Golang

## usage
	go get -u github.com/tbarker25/thermostat-api
	thermostat-api -address=127.0.0.1:8080

Note that GOPATH should be in your path. Otherwise the binary will be
found in:
	${GOPATH:-${HOME}/go}/bin/thermostat-api

## routes

To list all thermostats:
	curl -X GET --location 'http://localhost:8080/thermostat'
```json
{
  "status": "ok",
  "data": [
    {
      "id": 1,
      "name": "thermostat-1",
      "currentTemp": "57.7°F",
      "operatingMode": "off",
      "heatPoint": "65.0°F",
      "coolPoint": "80.0°F",
      "fanMode": "auto"
    },
    {
      "id": 2,
      "name": "thermostat-2",
      "currentTemp": "68.7°F",
      "operatingMode": "off",
      "heatPoint": "65.0°F",
      "coolPoint": "80.0°F",
      "fanMode": "auto"
    }
  ]
}
```

To retrieve a single thermostat
	curl -X GET --location 'http://localhost:8080/thermostat/1'
```json
{
  "status": "ok",
  "data": {
    "id": 1,
    "name": "thermostat-1",
    "currentTemp": "64.8°F",
    "operatingMode": "off",
    "heatPoint": "65.0°F",
    "coolPoint": "80.0°F",
    "fanMode": "auto"
  }
}
```

To modify the state of a thermostat
	curl -X PATCH --location 'http://localhost:8080/thermostat/2' --data '{"name": "downstairs kitchen"}'
```json
{
  "status": "ok",
  "data": {
    "id": 2,
    "name": "downstairs kitchen",
    "currentTemp": "92.9°F",
    "operatingMode": "off",
    "heatPoint": "65.0°F",
    "coolPoint": "80.0°F",
    "fanMode": "auto"
  }
}
```
