package lib

import (
	"testing"
)

type testRec struct {
	domain string
	includeIps bool
	nResults int
}

var testData = []testRec {
	{ domain: "hello.dylt.dev", includeIps: false, nResults: 5 },
	{ domain: "google.com", includeIps: false, nResults: 0},
}

func TestGetSrvs (t *testing.T) {
	// init data
	// Call target
	for _, d :=  range testData {
		srvs := GetSrvs(d.domain, d.includeIps)
	// Tests
		if (len(srvs) != d.nResults) { t.Fatalf("Expected %d, got %d", d.nResults, len(srvs))}
	}
}
