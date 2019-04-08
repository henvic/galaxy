package main

import (
	"context"
	_ "expvar"
	"flag"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/henvic/ctxsignal"
	"github.com/henvic/galaxy/server"
	log "github.com/sirupsen/logrus"
)

// @title galaxy API
// @version 1.0
// @description galaxy offers a drone navigation service (DNS).
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:9000
// @BasePath /v1

var params = server.Params{}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.Parse()

	var debug = (os.Getenv("DEBUG") != "")

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	if debug || params.ExposeDebug {
		go profiler()
	}

	ctx, cancel := ctxsignal.WithTermination(context.Background())
	defer cancel()

	if err := server.Start(ctx, params); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func profiler() {
	// let expvar and pprof be exposed here indirectly through http.DefaultServeMux

	// TODO(henvic): handle termination gracefully.
	// https://dave.cheney.net/practical-go/presentations/qcon-china.html#_never_start_a_goroutine_without_when_it_will_stop

	log.Info("Exposing expvar and pprof on localhost:8081")
	log.Fatal(http.ListenAndServe("localhost:8081", nil))
}

func init() {
	flag.StringVar(&params.Address, "addr", "127.0.0.1:9000", "Serving address")
	flag.BoolVar(&params.ExposeDebug, "expose-debug", true, "Expose debugging tools over HTTP (on port 8081)")
}
