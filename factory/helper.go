package factory

import (
	"fmt"
	"math/rand"
	"net/url"
	"time"
)

const (
	maxUint64 = ^uint64(0)
	maxInt64  = int64(maxUint64 >> 1)
)

var (
	names = []string{
		"Alice", "Aria", "Aurora", "Ava",
		"Bob",
		"Charlotte",
		"Emily",
		"Gabriel", "Grace",
		"Hailey",
		"James", "Jane", "John",
		"Liam",
		"Madison", "Mason", "Milad", "Mona",
		"Olivia",
		"Ryan",
		"Sarah", "Scarlett",
		"Tina",
	}
	domainNames = []string{"placeholder", "example", "alias"}
	domainTLDs  = []string{"com", "net", "org", "info", "io", "dev"}
)

// Name generates a random human name.
func Name() string {
	name := randPick(names...)
	return name
}

// Email generates a random email address.
func Email() string {
	name := randPick(names...)
	domainName := randPick(domainNames...)
	domainTLD := randPick(domainTLDs...)
	email := fmt.Sprintf("%s@%s.%s", name, domainName, domainTLD)
	return email
}

// Time generates a random time.Time value.
func Time() time.Time {
	nsec := rand.Int63()
	return time.Unix(0, nsec)
}

// TimeBefore generates a random time.Time value before a given time.
func TimeBefore(t time.Time) time.Time {
	max := t.UnixNano()
	nsec := rand.Int63n(max)
	return time.Unix(0, nsec)
}

// TimeAfter generates a random time.Time value after a given time.
func TimeAfter(t time.Time) time.Time {
	min := t.UnixNano() + 1
	nsec := min + rand.Int63n(maxInt64-min)
	return time.Unix(0, nsec)
}

// URL generates a random url.URL value.
func URL() url.URL {
	domainName := randPick(domainNames...)
	domainTLD := randPick(domainTLDs...)
	path := randPick("api", "app", "user", "team", "organization", "profile")
	u, _ := url.Parse(fmt.Sprintf("https://%s.%s/%s", domainName, domainTLD, path))
	return *u
}
