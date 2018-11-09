package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/j-forster/Wazihub-API"
	"github.com/j-forster/Wazihub-API/tools"
	routing "github.com/julienschmidt/httprouter"
)

func getSensor(deviceId, sensorId string) *wazihub.Sensor {
	device := getDevice(deviceId)
	if device == nil {
		return nil
	}
	for _, sensor := range device.Sensors {
		if sensor.Id == sensorId {
			return sensor
		}
	}
	return nil
}

func GetSensor(resp http.ResponseWriter, req *http.Request, params routing.Params) {
	sensor := getSensor(params.ByName("device_id"), params.ByName("sensor_id"))
	if sensor == nil {
		http.Error(resp, "Not Found: Sensor or device not found.", http.StatusNotFound)
		return
	}
	tools.SendJSON(resp, sensor)
}

func GetSensors(resp http.ResponseWriter, req *http.Request, params routing.Params) {
	device := getDevice(params.ByName("device_id"))
	if device == nil {
		http.Error(resp, "Not Found: Device not found.", http.StatusNotFound)
		return
	}
	tools.SendJSON(resp, device.Sensors)
}

func CreateSensor(resp http.ResponseWriter, req *http.Request, params routing.Params) {
	device := getDevice(params.ByName("device_id"))
	if device == nil {
		http.Error(resp, "Not Found: Device not found.", http.StatusNotFound)
		return
	}

	data, err := tools.ReadAll(req.Body)
	if err != nil {
		http.Error(resp, "Request Error: "+err.Error(), http.StatusBadRequest)
		return
	}

	sensor := &wazihub.Sensor{}
	err = json.Unmarshal(data, sensor)
	if err != nil {
		http.Error(resp, "Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}
	replaced := false
	if sensor.Id == "" {
		// NOT-CONFORM: Create a unique id if no id was given.
		sensor.Id = uuid.New().String()
	} else {
		for i, existing := range device.Sensors {
			if existing.Id == sensor.Id {
				device.Sensors[i] = sensor
				replaced = true
				break
			}
		}
	}

	if !replaced {
		device.Sensors = append(device.Sensors, sensor)
	}

	// NOT-CONFORM: Return id on success.
	resp.Header().Set("Content-Type", "text/plain")
	resp.Write([]byte(sensor.Id))
}
