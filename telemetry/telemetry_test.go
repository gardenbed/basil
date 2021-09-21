package telemetry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingleton(t *testing.T) {
	assert.NotNil(t, singleton)

	probe := new(probe)
	Set(probe)

	assert.Equal(t, probe, Get())
}
