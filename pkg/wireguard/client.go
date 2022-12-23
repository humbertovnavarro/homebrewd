package wireguard

import (
	"fmt"
	"os/exec"

	qrcode "github.com/skip2/go-qrcode"
)

type Client struct {
	PublicKey  string
	PrivateKey string
	Port       string
	Address    string

	PeerPublicKey string
	PeerAddress   string
}

func NewClient(server *Server, address string, allowedIPS string) *Client {
	cpub, cpriv := GenerateKeypair()
	exec.Command("wg", "set", "wg0", "peer", cpub, "allowed-ips", allowedIPS).Output()
	c := &Client{
		Address:       address,
		PublicKey:     cpub,
		PrivateKey:    cpriv,
		Port:          server.Port,
		PeerAddress:   server.PublicAddress,
		PeerPublicKey: server.PublicKey,
	}
	return c
}

func (c *Client) Config() string {
	var clientConfigTemplate = `[Interface]
PrivateKey=%s
Address=%s
ListenPort=%s

[Peer]
PublicKey=%s
EndPoint=%s
AllowedIPs=%s
	`
	return fmt.Sprintf(clientConfigTemplate, c.PrivateKey, c.Address, c.Port, c.PeerPublicKey, c.PeerAddress, "0.0.0.0/0")
}

func (c *Client) QR(size int) ([]byte, error) {
	return qrcode.Encode(c.Config(), qrcode.Medium, size)
}
