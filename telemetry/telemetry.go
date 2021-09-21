// Package telemetry can be used for implementing observability using OpenTelemetry API.
// It aims to unify the three pillars of observability in one single package that is easy-to-use and hard-to-misuse.
// It offers a unified developer experience for implementing observability.
package telemetry

// The singleton probe
var singleton Probe

// init initializes the singleton probe with a void probe.
// init function will be only called once in runtime regardless of how many times the package is imported.
func init() {
	singleton = NewVoidProbe()
}

// Get returns the singleton probe.
func Get() Probe {
	return singleton
}

// Set sets the singleton probe.
func Set(o Probe) {
	singleton = o
}
