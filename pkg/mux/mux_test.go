package mux_test

import (
	"math/rand"
	"testing"

	"proto.zip/studio/mux/pkg/host"
	"proto.zip/studio/mux/pkg/mux"
)

func TestNewMux(t *testing.T) {
	m := mux.New[any, any]()

	if m == nil {
		t.Error("expected new mux")
	}
}

func BenchmarkDomain(b *testing.B) {
	dn := "this.is.a.domain.for.benchmarking"
	m := mux.New[any, any]()
	m.NewHost(dn)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		d, _ := m.Host(dn)
		if d == nil {
			b.Error("got nil domain")
			return
		}
	}
}

func BenchmarkDomain_Large(b *testing.B) {
	m := mux.New[any, any]()

	tlds := []string{
		"com",
		"net",
		"org",
	}

	charset := "abcdefghijklmnopqrstuvwxyz"
	subDomains := make([]string, 200)
	for i := range subDomains {
		subDomain := make([]byte, 15)
		for j := 0; j < 15; j++ {
			subDomain[j] = charset[rand.Intn(len(charset))]
		}
		subDomains[i] = string(subDomain)
	}

	size := len(subDomains) * len(subDomains) * len(tlds)
	domainHandlers := make([]*host.Host[any, any], size)
	domains := make([]string, size)

	c := 0
	for _, tld := range tlds {
		for _, sub1 := range subDomains {
			for _, sub2 := range subDomains {
				domains[c] = sub1 + "." + sub2 + ".example." + tld
				domainHandlers[c], _ = m.NewHost(domains[c])
				c++
			}
		}
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		idx := n % len(domainHandlers)
		domainToFind := domainHandlers[idx]
		if domainToFind == nil {
			b.Errorf("domain to find is nil at %d", idx)
			return
		}
		d, _ := m.Host(domains[idx])

		if d != domainToFind {
			b.Error("got mismatched domain")
			return
		}
	}
}
