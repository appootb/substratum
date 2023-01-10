package ssh

import (
	"net"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type Tunnel struct {
	client *ssh.Client
}

// NewTunnel returns ssh tunnel client
// param: <user>:<password>@<host>[:port]
func NewTunnel(v string) *Tunnel {
	var (
		username, password, hostname string
	)
	//
	parts := strings.Split(v, "@")
	if len(parts) > 1 {
		hostname = parts[1]
	}
	parts = strings.Split(parts[0], ":")
	if len(parts) > 1 {
		password = parts[1]
	}
	username = parts[0]
	//
	cfg := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		BannerCallback:  ssh.BannerDisplayStderr(),
		Timeout:         time.Second * 5,
	}
	client, err := ssh.Dial("tcp", hostname, cfg)
	if err != nil {
		panic("SSH_TUNNEL err:" + err.Error())
	}
	return &Tunnel{
		client: client,
	}
}

func (m *Tunnel) Dial(network, addr string) (net.Conn, error) {
	return m.client.Dial(network, addr)
}
