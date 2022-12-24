package wireguard

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gorm.io/gorm"
)

type Server struct {
	gorm.Model
	Interface      string
	PublicKey      string
	PrivateKey     string
	Port           string
	PublicAddress  string
	PrivateAddress string
	Peers          map[string]*Client
}

func NewServer(listenInterface string, listenAddress string) *Server {
	pub, priv := GenerateKeypair()
	return &Server{
		PublicKey:      pub,
		Interface:      listenInterface,
		PrivateKey:     priv,
		Port:           "51280",
		PrivateAddress: "10.0.0.1/8",
		PublicAddress:  listenAddress,
		Peers:          make(map[string]*Client),
	}
}

func (s *Server) Serialize() (path string, err error) {
	os.MkdirAll("/etc/wireguard", 0600)
	var f *os.File
	f, err = os.Create("/etc/wireguard/wg0.conf")
	if err != nil {
		if err := os.Remove("/etc/wireguard/wg0.conf"); err != nil {
			return "", err
		}
		f, err = os.Create("/etc/wireguard/wg0.conf")
		if err != nil {
			return "", err
		}
	}
	_, err = f.WriteString(s.Config())
	return f.Name(), err
}

func (s *Server) Config() string {
	var serverConfigHead = `[Interface]
PrivateKey=%s
ListenPort=%s
Address=%s
SaveConfig=true
PostUp=iptables -A FORWARD -i wg0 -j  ACCEPT; iptables -t nat -A POSTROUTING -o %s -j MASQUERADE;
PostDown=iptables -D FORWARD -i wg0 -j  ACCEPT; iptables -t nat -D POSTROUTING -o %s -j MASQUERADE;
	`
	return fmt.Sprintf(serverConfigHead, s.PrivateKey, s.Port, s.PrivateAddress, s.Interface, s.Interface)
}

func (s *Server) Open() error {
	_, err := s.Serialize()
	if err != nil {
		return err
	}
	s.Close()
	c := exec.Command("wg-quick", "up", "wg0")
	err = c.Run()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Close() error {
	_, err := exec.Command("wg-quick", "down", "wg0").Output()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) AddClient(c *Client) error {
	s.Peers[c.PublicKey] = c
	_, err := exec.Command("wg", "set", "wg0", "peer", c.PublicKey, "allowed-ips", strings.Join(c.AllowedIPs, ",")).Output()
	return err
}

func (s *Server) RemoveClientt(c *Client) error {
	delete(s.Peers, c.PublicKey)
	_, err := exec.Command("wg", "set", "wg0", "peer", c.PublicKey, "remove").Output()
	return err
}
