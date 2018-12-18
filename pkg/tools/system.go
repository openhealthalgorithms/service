package tools

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

// CheckOS performs check of os on which agent is being running.
//
// Returns an os and an error.
// Error is not nil when os is not darwin, linux or windows.
func CheckOS() (string, error) {
	p := runtime.GOOS
	if p != "darwin" && p != "linux" && p != "windows" {
		return p, errors.Errorf("OS is not supported - %s", p)
	}

	return p, nil
}

// CheckArch performs check of arch on which agent is being running.
//
// Returns an arch and an error.
// Error is not nil when os is not amd64 or 386.
func CheckArch() (string, error) {
	a := runtime.GOARCH
	if a != "amd64" && a != "386" {
		return a, errors.Errorf("Arch is not supported - %s", a)
	}

	return a, nil
}

// HostID returns string `hostname|mac address` as HostID.
func HostID() string {
	return JoinStringsSep("|", Hostname(), RealMacAddr())
}

// HostIDHex is md5(hostname + mac address).
func HostIDHex() string {
	return fmt.Sprintf("%x", JoinStringsSep("|", Hostname(), RealMacAddr()))
}

// MemoryUsage returns memory usage by a program.
func MemoryUsage() (uint64, uint32) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return m.Sys, m.NumGC
}

// GetEnvFromList returns non-empty env vars got from evaluation of list of vars.
func GetEnvFromList(list []string) []string {
	res := make([]string, 0)
	for _, v := range list {
		d := os.Getenv(v)
		if len(d) == 0 {
			continue
		}

		res = append(res, d)
	}

	return res
}

// EnvVarsMap returns a map of vars and values from a given slice.
func EnvVarsMap(raw []string, filterEmpty bool) map[string]string {
	info := make(map[string]string, 0)

	for _, entry := range raw {
		pair := strings.SplitN(entry, "=", 2)
		if len(pair) < 2 {
			continue
		}

		k := pair[0]
		if len(k) == 0 {
			continue
		}

		v := pair[1]
		if filterEmpty && len(v) == 0 {
			continue
		}

		info[k] = v
	}

	return info
}

// GetCurrentDirectory function
func GetCurrentDirectory() string {
	pwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	return pwd
}
