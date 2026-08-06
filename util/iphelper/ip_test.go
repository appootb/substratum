package iphelper

import (
	"net"
	"testing"
)

func TestIsUsable(t *testing.T) {
	cases := []struct {
		name     string
		ip       string
		wantIPv4 bool
		ok       bool
	}{
		{"v4 private", "192.168.1.10", true, true},
		{"v4 public", "8.8.8.8", true, true},
		{"v4 loopback", "127.0.0.1", true, false},
		{"v4 link-local", "169.254.10.20", true, false},
		{"v4 unspecified", "0.0.0.0", true, false},
		{"v4 when want v6", "192.168.1.10", false, false},
		{"v6 global", "2001:db8::1", false, true},
		{"v6 link-local", "fe80::1", false, false},
		{"v6 loopback", "::1", false, false},
		{"v6 when want v4", "2001:db8::1", true, false},
		{"nil", "", true, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var ip net.IP
			if tc.ip != "" {
				ip = net.ParseIP(tc.ip)
				if ip == nil {
					t.Fatalf("ParseIP %q", tc.ip)
				}
			}
			if got := isUsable(ip, tc.wantIPv4); got != tc.ok {
				t.Fatalf("isUsable(%q, v4=%v)=%v want %v", tc.ip, tc.wantIPv4, got, tc.ok)
			}
		})
	}
}

func TestLocalIPValid(t *testing.T) {
	s := LocalIP()
	if s == "" || s == "<nil>" {
		t.Fatalf("LocalIP=%q", s)
	}
	ip := net.ParseIP(s)
	if ip == nil {
		t.Fatalf("not a valid IP: %q", s)
	}
	// Must never advertise loopback unless that is truly all we have —
	// 127.0.0.1 is the documented last resort and is valid syntax.
	if s != "127.0.0.1" {
		if ip.IsLoopback() || ip.IsLinkLocalUnicast() || !ip.IsGlobalUnicast() {
			t.Fatalf("unsuitable LocalIP=%q", s)
		}
	}
}

func TestLocalIPNeverNilString(t *testing.T) {
	// When no interfaces match, localIP returns errNoAddr (not nil,nil).
	// LocalIP must not return the literal "<nil>".
	s := LocalIP()
	if s == "<nil>" {
		t.Fatal(`LocalIP returned literal "<nil>"`)
	}
}

func TestLocalIPFamilyPreference(t *testing.T) {
	// Enumeration helpers respect family flags.
	if ip, err := localIP(true); err == nil {
		if ip.To4() == nil {
			t.Fatalf("localIP(true) returned non-IPv4 %v", ip)
		}
		if !isUsable(ip, true) {
			t.Fatalf("localIP(true) returned unusable %v", ip)
		}
	}
	if ip, err := localIP(false); err == nil {
		if ip.To4() != nil {
			t.Fatalf("localIP(false) returned IPv4 %v", ip)
		}
		if !isUsable(ip, false) {
			t.Fatalf("localIP(false) returned unusable %v", ip)
		}
	}
}
