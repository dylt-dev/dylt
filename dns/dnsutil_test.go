package dns

import (
	"testing"

	"github.com/stretchr/testify/require"
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


func TestGetA (t *testing.T) {
	const host = "www.eff.org"
	var addrs []string
	var err error
	addrs, err = GetA(host)
	require.NoError(t, err)
	require.NotEmpty(t, addrs)
	for _, addr := range addrs {
		t.Log(addr)
	}
}

func TestGetSrvs (t *testing.T) {
	// init data
	// Call target
	for _, d :=  range testData {
		t.Logf("Testing GetSrvs() [d.domain=%s, d.includeIps=%t]", d.domain, d.includeIps)
		srvs, err := GetSrvs(d.domain, "etcd-server", "tcp", d.includeIps)
	// Tests
		if err != nil { t.Log("Failure calling GetSrvs()", err, d.domain) }
		if (len(srvs) != d.nResults) { t.Fatalf("Expected %d, got %d", d.nResults, len(srvs))}
	}
}

func TestGetSrvsEtcdServer (t *testing.T) {
	// init data
	// Call target
	for _, d :=  range testData {
		t.Logf("Testing GetSrvs() [d.domain=%s, d.includeIps=%t]", d.domain, d.includeIps)
		srvs, err := GetSrvsEtcdServer(d.domain, d.includeIps)
	// Tests
		if err != nil { t.Log("Failure calling GetSrvs()", err, d.domain) }
		if (len(srvs) != d.nResults) { t.Fatalf("Expected %d, got %d", d.nResults, len(srvs))}
	}
}

