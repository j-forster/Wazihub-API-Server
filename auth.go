package main

import (
	"fmt"
	"net/http"

	routing "github.com/julienschmidt/httprouter"
)

func GetToken(resp http.ResponseWriter, req *http.Request, params routing.Params) {
	// TODO: implement
	fmt.Fprint(resp, "GetToken()")
}

func GetPermissions(resp http.ResponseWriter, req *http.Request, params routing.Params) {
	// TODO: implement
	fmt.Fprint(resp, "GetPermissions()")
}
