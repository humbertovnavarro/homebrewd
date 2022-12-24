package main

import (
	"fmt"
	"os"

	"github.com/humbertovnavarro/homebrewd/pkg/wireguard"
)

func main() {
	s := wireguard.NewServer("enp2s0f0u2", "68.127.190.132")
	err := s.Open()
	if err != nil {
		fmt.Println(err)
	}
	c := wireguard.NewClient(s, "10.0.0.2/8", []string{"0.0.0.0/0"})
	qr, err := c.QR(256)
	if err != nil {
		panic(err)
	}
	f, err := os.Create("./qr.png")
	if err != nil {
		panic(err)
	}
	f.Write(qr)
	f.Close()
	s.AddClient(c)
}
