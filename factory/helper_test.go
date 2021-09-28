package factory

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	name := Name()
	assert.NotZero(t, name)
}

func TestEmail(t *testing.T) {
	email := Email()
	assert.NotZero(t, email)
}

func TestTime(t *testing.T) {
	tm := Time()
	assert.NotZero(t, tm)
}

func TestTimeBefore(t *testing.T) {
	now := time.Now()
	tm := TimeBefore(now)
	assert.NotZero(t, tm)
	assert.True(t, tm.Before(now))
}

func TestTimeAfter(t *testing.T) {
	now := time.Now()
	tm := TimeAfter(now)
	assert.True(t, tm.After(now))
}

func TestURL(t *testing.T) {
	u := URL()
	assert.NotZero(t, u)
}
