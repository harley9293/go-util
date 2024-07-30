package crypto

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
)

type Dialer struct {
	Host       string
	Port       int
	User       string
	PrivateKey []byte
	Password   string
}

func (d *Dialer) GetClientConfig() (*ssh.ClientConfig, error) {
	config := &ssh.ClientConfig{
		User:            d.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if d.Password != "" {
		config.Auth = []ssh.AuthMethod{
			ssh.Password(d.Password),
		}
	} else {
		signer, err := ssh.ParsePrivateKey(d.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("unable to parse private key: %v", err)
		}
		config.Auth = []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		}
	}

	return config, nil
}

func NewSession(target *Dialer, jumpList ...*Dialer) (*ssh.Session, error) {
	var conn net.Conn
	var client *ssh.Client
	var err error

	for i, jump := range jumpList {
		config, err := jump.GetClientConfig()

		addr := fmt.Sprintf("%s:%d", jump.Host, jump.Port)
		if i == 0 {
			client, err = ssh.Dial("tcp", addr, config)
			if err != nil {
				return nil, fmt.Errorf("unable to connect to jump host: %v", err)
			}
		} else {
			conn, err = client.Dial("tcp", addr)
			if err != nil {
				client.Close()
				return nil, fmt.Errorf("unable to dial next jump host: %v", err)
			}

			clientConn, chans, reqs, err := ssh.NewClientConn(conn, addr, config)
			if err != nil {
				conn.Close()
				client.Close()
				return nil, fmt.Errorf("unable to establish SSH connection to next jump host: %v", err)
			}
			client = ssh.NewClient(clientConn, chans, reqs)
		}
	}

	config, err := target.GetClientConfig()
	addr := fmt.Sprintf("%s:%d", target.Host, target.Port)
	conn, err = client.Dial("tcp", addr)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("unable to connect to target host: %v", err)
	}

	clientConn, chans, reqs, err := ssh.NewClientConn(conn, addr, config)
	if err != nil {
		conn.Close()
		client.Close()
		return nil, fmt.Errorf("unable to establish SSH connection to target host: %v", err)
	}
	client = ssh.NewClient(clientConn, chans, reqs)

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("unable to create SSH session: %v", err)
	}

	return session, nil
}
