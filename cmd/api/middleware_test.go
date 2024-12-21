package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zukofett/greenlight/internal/assert"
)

func TestRateLimit(t *testing.T) {
    app := newTestApplication(t)

    r := httptest.NewRequest(http.MethodGet, "/", nil)
    next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("OK"))
    })

    retStatuses := make([]int, 40) 
    
    for i := 0; i < len(retStatuses); i++ {
        rr := httptest.NewRecorder()
        app.rateLimit(next).ServeHTTP(rr, r)
        retStatuses[i] = rr.Result().StatusCode
    }

    assert.Contains(t, retStatuses, http.StatusTooManyRequests)
}
