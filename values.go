package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/j-forster/Wazihub-API/tools"
	routing "github.com/julienschmidt/httprouter"
)

type Value struct {
	DeviceId     string
	SensorId     string
	TimeRecieved time.Time
	Value        interface{}
}

func PostValue(resp http.ResponseWriter, req *http.Request, params routing.Params) {

	deviceId := params.ByName("device_id")
	sensorId := params.ByName("sensor_id")

	value := Value{
		DeviceId:     deviceId,
		SensorId:     sensorId,
		TimeRecieved: time.Now(),
	}

	data, err := tools.ReadAll(req.Body)
	if err != nil {
		http.Error(resp, "Request Error: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(data, &value.Value)
	if err != nil {
		http.Error(resp, "Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if collection == nil {
		http.Error(resp, "Database not available.", http.StatusServiceUnavailable)
		return
	}

	//_, err = collection.InsertOne(req.Context(), &value)
	err = collection.Insert(&value)
	if err != nil {
		http.Error(resp, "Internal Database Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func PostValues(resp http.ResponseWriter, req *http.Request, params routing.Params) {

	deviceId := params.ByName("device_id")
	sensorId := params.ByName("sensor_id")

	var values []interface{}

	data, err := tools.ReadAll(req.Body)
	if err != nil {
		http.Error(resp, "Request Error: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(data, &values)
	if err != nil {
		http.Error(resp, "Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now()
	points := make([]interface{}, len(values))
	for i, _ := range points {
		points[i] = &Value{
			DeviceId:     deviceId,
			SensorId:     sensorId,
			TimeRecieved: now,
			Value:        values[i],
		}
	}

	if collection == nil {
		http.Error(resp, "Database not available.", http.StatusServiceUnavailable)
		return
	}

	//_, err = collection.InsertOne(req.Context(), &value)
	err = collection.Insert(points...)
	if err != nil {
		http.Error(resp, "Internal Database Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
