package wireguard

import (
	"fmt"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
)

type Client struct {
	PublicKey     string
	PrivateKey    string
	Port          string
	Address       string
	PeerPublicKey string
	PeerAddress   string
	AllowedIPs    []string
}

func NewClient(server *Server, address string, allowedIPS []string) *Client {
	cpub, cpriv := GenerateKeypair()
	c := &Client{
		Address:       address,
		PublicKey:     cpub,
		PrivateKey:    cpriv,
		Port:          server.Port,
		PeerAddress:   server.PublicAddress,
		PeerPublicKey: server.PublicKey,
		AllowedIPs:    allowedIPS,
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
AllowedIPs=%s`
	return fmt.Sprintf(clientConfigTemplate, c.PrivateKey, c.Address, c.Port, c.PeerPublicKey, c.PeerAddress, strings.Join(c.AllowedIPs, ","))
}

func (c *Client) QR(size int) ([]byte, error) {
	return qrcode.Encode(c.Config(), qrcode.Medium, size)
}
