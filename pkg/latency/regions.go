package latency

import "time"

type RegionLatency struct {
	Min time.Duration
	Max time.Duration
}

// Regions map defines average latency added for each region relative to "US-Central" (simulated origin).
var Regions = map[string]RegionLatency{
	"us-east-1":      {Min: 10 * time.Millisecond, Max: 50 * time.Millisecond},
	"us-west-1":      {Min: 40 * time.Millisecond, Max: 80 * time.Millisecond},
	"eu-west-1":      {Min: 80 * time.Millisecond, Max: 120 * time.Millisecond},
	"ap-northeast-1": {Min: 180 * time.Millisecond, Max: 250 * time.Millisecond}, // Tokyo
	"ap-southeast-2": {Min: 200 * time.Millisecond, Max: 300 * time.Millisecond}, // Sydney
	"sa-east-1":      {Min: 150 * time.Millisecond, Max: 220 * time.Millisecond}, // Sao Paulo
}

func GetLatency(region string) (RegionLatency, bool) {
	val, ok := Regions[region]
	return val, ok
}
