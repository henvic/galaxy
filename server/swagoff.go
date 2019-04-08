// +build !swagon

package server

import "github.com/gorilla/mux"

var swagon = false

func maybeSwagger(router *mux.Router, params Params) {}
