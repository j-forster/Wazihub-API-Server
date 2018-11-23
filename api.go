package main

import (
	routing "github.com/julienschmidt/httprouter"
)

var router = routing.New()

func init() {

	router.POST("/auth/token", GetToken)
	router.GET("/auth/permissions", GetPermissions)

	router.GET("/devices", GetDevices)
	router.POST("/devices", CreateDevice)
	router.GET("/devices/:device_id", GetDevice)

	router.GET("/devices/:device_id/sensors", GetSensors)
	router.POST("/devices/:device_id/sensors", CreateSensor)
	router.GET("/devices/:device_id/sensors/:sensor_id", GetSensor)

	router.POST("/devices/:device_id/sensors/:sensor_id/values", PostValues)

	router.GET("/devices/:device_id/actuators", GetActuators)
	router.POST("/devices/:device_id/actuators", CreateActuator)
	router.GET("/devices/:device_id/actuators/:actuator_id", GetActuator)
}
