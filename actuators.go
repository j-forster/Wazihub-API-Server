package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/j-forster/Wazihub-API"
	"github.com/j-forster/Wazihub-API/tools"
	routing "github.com/julienschmidt/httprouter"
)

func getActuator(deviceId, actuatorId string) *wazihub.Actuator {
	device := getDevice(deviceId)
	if device == nil {
		return nil
	}
	for _, actuator := range device.Actuators {
		if actuator.Id == actuatorId {
			return actuator
		}
	}
	return nil
}

func GetActuator(resp http.ResponseWriter, req *http.Request, params routing.Params) {
	sensor := getActuator(params.ByName("device_id"), params.ByName("actuator_id"))
	if sensor == nil {
		http.Error(resp, "Not Found: Actuator or device not found.", http.StatusNotFound)
		return
	}
	tools.SendJSON(resp, sensor)
}

func GetActuators(resp http.ResponseWriter, req *http.Request, params routing.Params) {
	device := getDevice(params.ByName("device_id"))
	if device == nil {
		http.Error(resp, "Not Found: Device not found.", http.StatusNotFound)
		return
	}
	tools.SendJSON(resp, device.Actuators)
}

func CreateActuator(resp http.ResponseWriter, req *http.Request, params routing.Params) {
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

	actuator := &wazihub.Actuator{}
	err = json.Unmarshal(data, actuator)
	if err != nil {
		http.Error(resp, "Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}
	if actuator.Id == "" {
		// NOT-CONFORM: Create a unique id if no id was given.
		actuator.Id = uuid.New().String()
	}
	device.Actuators = append(device.Actuators, actuator)

	// NOT-CONFORM: Return id on success.
	resp.Header().Set("Content-Type", "text/plain")
	resp.Write([]byte(actuator.Id))
}
