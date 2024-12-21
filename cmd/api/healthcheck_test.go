package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/zukofett/greenlight/internal/assert"
)

func TestHealthcheck(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	status, _, body := ts.get(t, "/v1/healthcheck")

	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	want, err := json.MarshalIndent(env, "", "\t")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, status, http.StatusOK)
	assert.Equal(t, body, string(want))
}
