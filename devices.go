package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/j-forster/Wazihub-API"
	"github.com/j-forster/Wazihub-API/tools"

	routing "github.com/julienschmidt/httprouter"
)

var devices map[string]*wazihub.Device = map[string]*wazihub.Device{
	wazihub.CurrentDeviceId(): &wazihub.Device{
		Id:   wazihub.CurrentDeviceId(),
		Name: "Gateway " + wazihub.CurrentDeviceId(),
	},
}

////////////////////

func GetDevices(resp http.ResponseWriter, req *http.Request, params routing.Params) {

	buf := bytes.Buffer{}
	buf.Write([]byte{'['})
	n := 0
	for _, device := range devices {
		if n != 0 {
			buf.Write([]byte{','})
		}
		n++
		// NOT-CONFORM: Returns also last_value.
		data, err := json.MarshalIndent(device, "", "  ")
		if err != nil {
			http.Error(resp, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		buf.Write(data)
	}
	buf.Write([]byte{']'})
	resp.Write(buf.Bytes())
}

func CreateDevice(resp http.ResponseWriter, req *http.Request, params routing.Params) {

	data, err := tools.ReadAll(req.Body)
	if err != nil {
		http.Error(resp, "Request Error: "+err.Error(), http.StatusBadRequest)
		return
	}

	device := &wazihub.Device{}
	err = json.Unmarshal(data, device)
	if err != nil {
		http.Error(resp, "Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}
	if device.Id == "" {
		// NOT-CONFORM: Create a unique id if no id was given.
		device.Id = uuid.New().String()
	}
	devices[device.Id] = device

	// NOT-CONFORM: Return id on success.
	resp.Header().Set("Content-Type", "text/plain")
	resp.Write([]byte(device.Id))
}

func getDevice(id string) *wazihub.Device {
	return devices[id]
}

func GetDevice(resp http.ResponseWriter, req *http.Request, params routing.Params) {

	device := getDevice(params.ByName("device_id"))
	if device == nil {
		http.Error(resp, "Not Found: Device not found.", http.StatusNotFound)
		return
	}

	tools.SendJSON(resp, device)
}
