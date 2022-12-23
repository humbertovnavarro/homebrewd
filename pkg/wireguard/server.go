package wireguard

import (
	"fmt"
	"os"
	"os/exec"

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
	}
}

func (s *Server) Config() string {
	var serverConfigTemplate = `[Interface]
PrivateKey=%s
ListenPort=%s
Address=%s
SaveConfig=true
PostUp=iptables -A FORWARD -i wg0 -j  ACCEPT; iptables -t nat -A POSTROUTING -o %s -j MASQUERADE;
PostDown=iptables -D FORWARD -i wg0 -j  ACCEPT; iptables -t nat -D POSTROUTING -o %s -j MASQUERADE;
	`
	return fmt.Sprintf(serverConfigTemplate, s.PrivateKey, s.Port, s.PrivateAddress, s.Interface, s.Interface)
}

func (s *Server) serialize() error {
	os.Mkdir("/etc/wireguard", 0600)
	f, err := os.Create("/etc/wireguard/wg0.conf")
	if err != nil {
		return err
	}
	f.WriteString(s.Config())
	return f.Close()
}

func (s *Server) Open() error {
	err := s.serialize()
	if err != nil {
		return err
	}
	_, err = exec.Command("wg-quick", "up", "wg0").Output()
	return err
}

func (s *Server) Close() error {
	_, err := exec.Command("wg-quick", "down", "wg0").Output()
	return err
}
