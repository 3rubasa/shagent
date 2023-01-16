package webserver

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test1(t *testing.T) {
	w := New()
	err := w.Initialize()
	assert.NoError(t, err)
	err = w.Start()
	assert.NoError(t, err)

	time.Sleep(time.Second)

	_, err = http.Get("http://127.0.0.1:8888/controllers/boiler/get_state")
	assert.NoError(t, err)

	// body, err := io.ReadAll(r.Body)
	// assert.NoError(t, err)

	// strBody := string(body)
	// assert.Equal(t, strBody, "{}")
}
