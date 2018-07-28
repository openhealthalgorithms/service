package tools

import (
	"bytes"
	"net"
	"os"
	"testing"
)

func TestHostname(t *testing.T) {
	expected, err := os.Hostname()
	if err != nil {
		t.Fatalf("failed to get ethalon info: %v", err)
	}

	actual := Hostname()
	if actual != expected {
		t.Errorf("expected: %s, got %s", expected, actual)
	}
}

func TestDtoi(t *testing.T) {
	big := 0xFFFFFF

	type diTest struct {
		input  string
		num    int
		consum int
		ok     bool
	}

	diTests := []diTest{
		{"123", 123, 3, true},
		{"", 0, 0, false},
		{"4294967295", big, 7, false},
	}

	for _, dit := range diTests {
		n, i, o := Dtoi(dit.input)
		if n != dit.num || i != dit.consum || o != dit.ok {
			t.Errorf("expected (%s): %d %d %v, got %d %d %v",
				dit.input, dit.num, dit.consum, dit.ok, n, i, o,
			)
		}
	}
}

func TestParseIPPrefix(t *testing.T) {
	type ipPrefixTest struct {
		input  string
		output int
	}

	ipPrefixTests := []ipPrefixTest{
		{"Hello world", 0},
		{"Hello/world", 0},
		{"124.65.33.240/16", 16},
	}

	for _, ipf := range ipPrefixTests {
		actual := ParseIPPrefix(ipf.input)
		if actual != ipf.output {
			t.Errorf("expected %d, got %d", ipf.output, actual)
		}
	}
}

func TestRealMacAddr(t *testing.T) {
	interfaces, err := net.Interfaces()
	if err != nil {
		t.Fatalf("failed to get ethalon info: %v", err)
	}

	macAddr := RealMacAddr()

	hwAddr := make(map[string]bool)
	for _, i := range interfaces {
		hwAddr[i.HardwareAddr.String()] = i.Flags&net.FlagUp != 0 && !bytes.Equal(i.HardwareAddr, nil)
	}

	if !hwAddr[macAddr] {
		t.Errorf("mac address not found: %s", macAddr)
	}
}
