// +build swagon

package server

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

var swagon = true

func maybeSwagger(router *mux.Router, params Params) {
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}
