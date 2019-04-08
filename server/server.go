package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/hashicorp/errwrap"
	log "github.com/sirupsen/logrus"
)

// Params of the service.
type Params struct {
	Address string

	ExposeDebug bool
}

// Start server.
func Start(ctx context.Context, params Params) error {
	var s = &Server{}

	return s.Serve(ctx, params)
}

// Server for handling requests.
type Server struct {
	ctx context.Context

	params Params

	http *http.Server
	ec   chan error
}

// Serve handlers.
func (s *Server) Serve(ctx context.Context, params Params) error {
	s.ctx = ctx
	s.params = params

	var router = mux.NewRouter()
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/v1/sectors/{sector_id}/dns", dnsHandler)
	maybeSwagger(router, params)

	if swagon {
		log.Infof("warning: swagger available at %s/swagger/index.html", getAddr(params.Address))
	}

	s.http = &http.Server{
		Handler: router,
	}

	return s.serve()
}

func getAddr(a string) string {
	l := strings.LastIndex(a, ":")

	if l == -1 && len(a) <= l {
		return a
	}

	return "http://localhost:" + a[l+1:]
}

// Serve HTTP requests.
func (s *Server) serve() error {
	s.ec = make(chan error, 1)
	go s.listen()
	go s.waitShutdown()

	err := <-s.ec

	if err == http.ErrServerClosed {
		fmt.Println()
		log.Info("Server shutting down gracefully.")
	}

	return err
}

func (s *Server) waitShutdown() {
	<-s.ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.http.Shutdown(ctx); err != nil && err != context.Canceled {
		s.ec <- errwrap.Wrapf("can't shutdown server properly: {{err}}", err)
	}
}

func (s *Server) listen() {
	l, err := net.Listen("tcp", s.params.Address)

	if err != nil {
		s.ec <- err
		return
	}

	log.Infof("Starting server on %v", getAddr(l.Addr().String()))

	err = s.http.Serve(l)

	s.ec <- err
}

// @Summary Home
// @Description Healthcheck endpoint
// @Tags galaxy
// @ID home-galaxy
// @Produce plain
// @Success 200 {string} string	"Documentation message"
// @Router / [get]
func homeHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintln(w, "galaxy is running. Docs in https://github.com/henvic/galaxy")
}
