package boiler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test1(t *testing.T) {
	bc := New()
	bc.Initialize()
	bc.Start()

	s, err := bc.GetState()
	assert.NoError(t, err)
	assert.Equal(t, s, "on")
}
