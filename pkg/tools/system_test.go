package tools

import (
	"fmt"
	"runtime"
	"testing"
)

func TestCheckOS(t *testing.T) {
	v, err := CheckOS()
	if err != nil {
		t.Fatalf("failed to get data: %v", err)
	}

	if v != runtime.GOOS {
		t.Errorf("expected %s, got %s", runtime.GOOS, v)
	}
}

func TestCheckArch(t *testing.T) {
	v, err := CheckArch()
	if err != nil {
		t.Fatalf("failed to get data: %v", err)
	}

	if v != runtime.GOARCH {
		t.Errorf("expected %s, got %s", runtime.GOARCH, v)
	}
}

func TestHostID(t *testing.T) {
	host := Hostname()
	mac := RealMacAddr()
	expected := JoinStringsSep("|", host, mac)
	hostID := HostID()

	if hostID != expected {
		t.Errorf("expected %s, got %s", expected, hostID)
	}
}

func TestHostIDHex(t *testing.T) {
	expected := fmt.Sprintf("%x", HostID())
	hostIDHash := HostIDHex()

	if hostIDHash != expected {
		t.Errorf("expected %s, got %s", expected, hostIDHash)
	}
}

func TestMemoryUsage(t *testing.T) {
	s, n := MemoryUsage()

	if !(s > 0) || !(n >= 0) {
		t.Errorf("invalid values: %d %d", s, n)
	}
}
