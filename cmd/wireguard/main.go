package main

import (
	"fmt"

	"github.com/humbertovnavarro/homebrewd/pkg/wireguard"
)

func main() {
	s := wireguard.NewServer("enp2s0f0u2", "68.127.190.132")
	c0 := wireguard.NewClient(s, "10.0.0.2", "0.0.0.0/0")
	c1 := wireguard.NewClient(s, "10.0.0.3", "0.0.0.0/0")
	fmt.Println(c0.Config())
	fmt.Println(c1.Config())
}
