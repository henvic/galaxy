package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/henvic/galaxy"
	log "github.com/sirupsen/logrus"
)

// @Summary Show a DNS location
// @Description Get a DNS location of a sector of the galaxy
// @Tags galaxy
// @ID dns-sector-galaxy
// @Accept json
// @Produce json
// @Param sector_id path int true "Sector ID"
// @Success 200 {object} server.DNSResponse
// @Failure 400 {object} server.ErrorResponse
// @Failure 405 {object} server.ErrorResponse
// @Failure 406 {object} server.ErrorResponse
// @Router /sectors/{sector_id}/dns [post]
func dnsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sid, ok := getSectorID(vars["sector_id"])

	if !ok {
		ErrorHandler(w, r, http.StatusBadRequest, "Sector value must be numeric")
		return
	}

	if r.Method != http.MethodPost {
		ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	if t := r.Header.Get("Content-Type"); !strings.Contains(t, "application/json") {
		ErrorHandler(w, r, http.StatusNotAcceptable)
		return
	}

	d := dnsRequest{}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		ErrorHandler(w, r, http.StatusBadRequest, "cannot decode request body as JSON")
		log.Debugf("bad request: %v", err)
		return
	}

	dns, err := d.geolocation()

	if err != nil {
		ErrorHandler(w, r, http.StatusBadRequest, err.Error())
		log.Debugf("bad request: %v", err)
		return
	}

	resp := DNSResponse{
		Loc: dns.Loc(sid),
	}

	e := json.NewEncoder(w)

	e.SetIndent("", "    ")

	w.Header().Set("Content-Type", "application/json; charset=utf8")

	if err = e.Encode(resp); err != nil {
		ErrorHandler(w, r, -1, fmt.Sprintf("cannot encode response for request: %v", err))
		log.Errorf("cannot encode response for request: %v", err)
		return
	}
}

type dnsRequest struct {
	X   json.Number `json:"x"`
	Y   json.Number `json:"y"`
	Z   json.Number `json:"z"`
	Vel json.Number `json:"vel"`
}

func (d *dnsRequest) geolocation() (dns galaxy.DNS, err error) {
	var failed = []string{}

	if dns.X, err = d.X.Float64(); err != nil {
		failed = append(failed, "x")
	}

	if dns.Y, err = d.Y.Float64(); err != nil {
		failed = append(failed, "y")
	}

	if dns.Z, err = d.Z.Float64(); err != nil {
		failed = append(failed, "z")
	}

	if dns.Vel, err = d.Vel.Float64(); err != nil {
		failed = append(failed, "vel")
	}

	if len(failed) != 0 {
		return dns, fmt.Errorf("invalid coordinates: %v", strings.Join(failed, ", "))
	}

	return dns, nil
}

// DNSResponse with a galaxy sector location.
type DNSResponse struct {
	Loc float64 `json:"loc" example:"27372.229"`
}

func getSectorID(s string) (int, bool) {
	sid, err := strconv.Atoi(s)
	return sid, sid > 0 && err == nil
}
