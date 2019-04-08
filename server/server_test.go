package server_test

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/henvic/galaxy/server"
)

var params = server.Params{
	Address: "127.0.0.1:9375", // TODO(henvic): use random port
}

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	go func() {
		if err := server.Start(ctx, params); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// TODO(henvic): add channel to verify that the server is ready and remove this sleep call.
	time.Sleep(time.Second)

	flag.Parse()
	os.Exit(m.Run())
}

func TestDNSBadBody(t *testing.T) {
	path := fmt.Sprintf("http://%s/v1/sectors/1/dns", params.Address)
	body := ioutil.NopCloser(bytes.NewReader([]byte("x")))

	resp, err := http.Post(path, "application/json", body)

	if err != nil {
		t.Errorf("Expected no error, got %v instead", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status to be %v, got %v instead", http.StatusBadRequest, resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")

	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected response to be application/json, got %v instead", contentType)
	}
}

func TestDNSMethodNotAllowed(t *testing.T) {
	path := fmt.Sprintf("http://%s/v1/sectors/1/dns", params.Address)

	resp, err := http.Get(path)

	if err != nil {
		t.Errorf("Expected no error, got %v instead", err)
	}

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status to be %v, got %v instead", http.StatusMethodNotAllowed, resp.StatusCode)
	}
}

func TestDNSBadContentType(t *testing.T) {
	path := fmt.Sprintf("http://%s/v1/sectors/1/dns", params.Address)
	body := ioutil.NopCloser(bytes.NewReader([]byte("x")))

	resp, err := http.Post(path, "application/unacceptable", body)

	if err != nil {
		t.Errorf("Expected no error, got %v instead", err)
	}

	if resp.StatusCode != http.StatusNotAcceptable {
		t.Errorf("Expected status to be %v, got %v instead", http.StatusNotAcceptable, resp.StatusCode)
	}
}

func TestDNSNonNumericSectionID(t *testing.T) {
	path := fmt.Sprintf("http://%s/v1/sectors/3x/dns", params.Address)
	body := ioutil.NopCloser(bytes.NewReader([]byte(
		`{"x": "33", "y": "42", "z": "13", "vel": "4.229"}`,
	)))

	resp, err := http.Post(path, "application/json", body)

	if err != nil {
		t.Errorf("Expected no error, got %v instead", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status to be %v, got %v instead", http.StatusBadRequest, resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Errorf("Expected no error, got %v instead", err)
	}

	if !strings.Contains(string(b), "Sector value must be numeric") {
		t.Errorf("Expected sector error missing, got %s instead", b)
	}
}

var invalidCases = []struct {
	name     string
	sectorID string
	request  string
	want     string
}{
	{
		name:    "all invalid",
		request: `{"x": "foo"}`,
		want:    "invalid coordinates: x, y, z, vel",
	},
	{
		name:    "x is invalid",
		request: `{"x": "foo", "y": "10", "z": "4.1", "vel": "2"}`,
		want:    "invalid coordinates: x",
	},
	{
		name:    "y is invalid",
		request: `{"x": "4", "y": "foo", "z": "4.1", "vel": "2"}`,
		want:    "invalid coordinates: y",
	},
	{
		name:    "z is invalid",
		request: `{"x": "2.4", "y": "13.0", "z": "foo", "vel": "2.7"}`,
		want:    "invalid coordinates: z",
	},
	{
		name:    "vel is invalid",
		request: `{"x": "2.4", "y": "13.0", "z": "5.4", "vel": "foo"}`,
		want:    "invalid coordinates: vel",
	},
}

func TestDNSInvalidCoordinates(t *testing.T) {
	for _, tt := range invalidCases {
		t.Run(tt.name, func(t *testing.T) {
			path := fmt.Sprintf("http://%s/v1/sectors/17/dns", params.Address)
			body := ioutil.NopCloser(bytes.NewReader([]byte(tt.request)))

			resp, err := http.Post(path, "application/json", body)

			if err != nil {
				t.Errorf("Expected no error, got %v instead", err)
			}

			if resp.StatusCode != http.StatusBadRequest {
				t.Errorf("Expected status to be %v, got %v instead", http.StatusBadRequest, resp.StatusCode)
			}

			contentType := resp.Header.Get("Content-Type")

			if !strings.Contains(contentType, "application/json") {
				t.Errorf("Expected response to be application/json, got %v instead", contentType)
			}

			b, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				t.Errorf("Expected no error, got %v instead", err)
			}

			var got = string(b)

			if strings.Contains(tt.want, got) {
				t.Errorf("Expected error to contain %v, got %v instead", tt.want, got)
			}
		})
	}
}

func TestDNS(t *testing.T) {
	path := fmt.Sprintf("http://%s/v1/sectors/50/dns", params.Address)
	body := ioutil.NopCloser(bytes.NewReader([]byte(
		`{"x": "33", "y": "42", "z": "13", "vel": "4.229"}`,
	)))

	resp, err := http.Post(path, "application/json", body)

	if err != nil {
		t.Errorf("Expected no error, got %v instead", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status to be %v, got %v instead", http.StatusOK, resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")

	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected response to be application/json, got %v instead", contentType)
	}

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Errorf("Expected no error, got %v instead", err)
	}

	var want = `{
    "loc": 4404.229
}
`

	var got = string(b)

	if want != got {
		t.Errorf("Expected %v, got %v instead", want, got)
	}
}
