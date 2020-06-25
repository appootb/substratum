package ip

import (
	"net"
)

func LocalIP(v4 bool) (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || // Interface down
			iface.Flags&net.FlagLoopback != 0 || // Loopback
			iface.Flags&net.FlagPointToPoint != 0 { // Point to point
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			if v4 && ip.To4() == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, nil
}
