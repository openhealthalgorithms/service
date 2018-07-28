package tools

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// Bigger than we need, not too big to worry about overflow.
const big = 0xFFFFFF

// Hostname returns hostname.
func Hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return JoinStringsSep("_", "unknown", RandomString(10))
	}

	return hostname
}

// RealMacAddr returns first real mac addr of a system.
func RealMacAddr() string {
	var addr string
	interfaces, err := net.Interfaces()
	if err != nil {
		return JoinStringsSep("_", "unknown", RandomString(10))
	}

	for _, i := range interfaces {
		if i.Flags&net.FlagUp != 0 && !bytes.Equal(i.HardwareAddr, nil) {
			addr = i.HardwareAddr.String()
			break
		}
	}

	return addr
}

// Dtoi converts decimal string to int.
//
// Returns number, characters consumed, success.
func Dtoi(s string) (n int, i int, ok bool) {
	n = 0
	for i = 0; i < len(s) && '0' <= s[i] && s[i] <= '9'; i++ {
		n = n*10 + int(s[i]-'0')
		if n >= big {
			return big, i, false
		}
	}
	if i == 0 {
		return 0, 0, false
	}
	return n, i, true
}

// ParseIPPrefix parses ip_or_network/prefix and returns the prefix as int.
func ParseIPPrefix(s string) int {
	i := strings.IndexByte(s, '/')
	if i < 0 {
		return 0
	}

	prefix, _, ok := Dtoi(s[i+1:])
	if !ok {
		return 0
	}

	return prefix
}

// ParseIPMask parses string ip_or_network/prefix and returns address and network mask.
func ParseIPMask(netAddr string) (string, string, bool) {
	ipAddr, ipNet, err := net.ParseCIDR(netAddr)
	if err != nil {
		ipAddr := net.ParseIP(netAddr)
		if ipAddr == nil {
			return "", "", false
		}

		addr := ipAddr.String()
		mask := "255.255.255.255"

		return addr, mask, true
	}

	addr := ipAddr.String()
	mask := fmt.Sprintf("%d.%d.%d.%d", ipNet.Mask[0], ipNet.Mask[1], ipNet.Mask[2], ipNet.Mask[3])

	return addr, mask, true
}

// CheckConnectionTimeout tries to connect with to host and port with timeout.
func CheckConnectionTimeout(network, address string, timeout time.Duration) error {
	if _, err := net.DialTimeout(network, address, timeout); err != nil {
		return err
	}

	return nil
}
