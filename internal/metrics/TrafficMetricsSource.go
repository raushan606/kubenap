package metrics

import "time"

// TrafficMetricsSource is a pluggable interface that provides the last observed
// request time for a given service. This is used to determine application idleness.

type TrafficMetricsSource interface {

	// LastRequestTime returns the timestamp of the last observed request for the given service name.
	// It returns an error if the value cannot be determined or retrieved.
	LastRequestTime(serviceName string) (time.Time, error)
}
