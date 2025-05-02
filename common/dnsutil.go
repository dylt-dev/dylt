package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
)


type Srv struct {
	net.SRV
	Ips []string
}


func GetCname (host string) string {
	cname, err := net.LookupCNAME(host)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ""
	}

	return cname
}


func GetSrvs (domain string, service string, proto string, includeIps bool) ([]Srv, error) {
	var srvs []Srv
	// Use "net" package to do a DNS lookup of SRVs 
	_, dnsSrvs, err := net.LookupSRV(service, proto, domain)
	if err != nil {
		var dnsError *net.DNSError
		if errors.As(err, &dnsError) {
			slog.Warn("This is a DNS error", "errmsg", err)
			j, _ := json.Marshal(dnsError)
			slog.Warn(string(j))
		} else {
			slog.Warn("Error during net.LookupSRV()", "errmsg", err)
		}
	}
	// For each DNS SRV record, create an Srv object
	// Optionally populate the srv Ips field if --include-ips was set
	for _, dnsSrv := range dnsSrvs {
		srv := Srv{SRV: *dnsSrv}
		if includeIps {
			ips, err := net.LookupHost(dnsSrv.Target)
			if err != nil {
				slog.Warn("Error during net.LookupHost()", "err", err.Error())
				srv.Ips = []string{}
			} else {
				srv.Ips = ips
			}
		}
		srvs = append(srvs, srv)
	}

	return srvs, err
}

func GetSrvsEtcdServer (domain string, includeIps bool) ([]Srv, error) {
	return GetSrvs(domain, "etcd-server", "tcp", includeIps)
}

func GetTxts (domain string) []string {
	txts, err := net.LookupTXT(domain)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}

	return txts
}
