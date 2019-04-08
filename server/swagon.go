// +build swagon

package server

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

var swagon = true

func maybeSwagger(router *mux.Router, params Params) {
	log.Infof("warning: swagger available at %s/swagger/index.html", getAddr(params.Address))
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}
