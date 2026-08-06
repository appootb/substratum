// Package iphelper resolves a preferred non-loopback local IP address.
package iphelper

import (
	"errors"
	"net"
)

var errNoAddr = errors.New("iphelper: no suitable address")

// LocalIP returns a usable host IP for service advertisement:
//  1. address of the default outbound route (UDP dial, no real traffic required),
//  2. else first suitable IPv4 on an up interface,
//  3. else first suitable global IPv6,
//  4. else "127.0.0.1".
//
// Link-local (169.254/16, fe80::/10), loopback, and unspecified addresses are skipped.
func LocalIP() string {
	if ip, err := fromOutbound(); err == nil {
		return ip.String()
	}
	if ipv4, err := localIP(true); err == nil {
		return ipv4.String()
	}
	if ipv6, err := localIP(false); err == nil {
		return ipv6.String()
	}
	return "127.0.0.1"
}

// fromOutbound discovers the local address used for default routing by "dialing"
// well-known destinations over UDP (no packets need to be exchanged successfully).
func fromOutbound() (net.IP, error) {
	for _, network := range []struct {
		proto string
		dst   string
		v4    bool
	}{
		{"udp4", "8.8.8.8:80", true},
		{"udp6", "[2001:4860:4860::8888]:80", false},
	} {
		conn, err := net.Dial(network.proto, network.dst)
		if err != nil {
			continue
		}
		la, ok := conn.LocalAddr().(*net.UDPAddr)
		_ = conn.Close()
		if !ok || !isUsable(la.IP, network.v4) {
			continue
		}
		return la.IP, nil
	}
	return nil, errNoAddr
}

// isUsable reports whether ip is suitable for external advertisement.
// wantIPv4 selects address family; when false only pure IPv6 is accepted.
func isUsable(ip net.IP, wantIPv4 bool) bool {
	if ip == nil || len(ip) == 0 {
		return false
	}
	if wantIPv4 {
		if ip.To4() == nil {
			return false
		}
	} else if ip.To4() != nil {
		// IPv4 or IPv4-mapped IPv6 when caller asked for IPv6.
		return false
	}
	// IsGlobalUnicast excludes loopback, link-local (incl. 169.254/16, fe80::/10),
	// unspecified, and multicast.
	return ip.IsGlobalUnicast()
}

func localIP(v4 bool) (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 ||
			iface.Flags&net.FlagLoopback != 0 ||
			iface.Flags&net.FlagPointToPoint != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			// Skip this interface; try others.
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if !isUsable(ip, v4) {
				continue
			}
			return ip, nil
		}
	}
	return nil, errNoAddr
}
