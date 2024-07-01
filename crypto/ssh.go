package crypto

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
)

type Dialer struct {
	Host string
	Port int
	User string
}

func NewSession(privateKeyFile string, target *Dialer, jumpList ...*Dialer) (*ssh.Session, error) {
	key, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}

	var conn net.Conn
	var client *ssh.Client

	for i, jump := range jumpList {
		config := &ssh.ClientConfig{
			User: jump.User,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

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

	config := &ssh.ClientConfig{
		User: target.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

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
