package wireguard

import (
	"fmt"
	"os/exec"
	"strings"
)

func GenerateKeypair() (string, string) {
	output, _ := exec.Command("wg", "genkey").Output()
	privateKey := strings.TrimSpace(string(output))
	cmd := fmt.Sprintf("echo %s | wg pubkey", privateKey)
	output, _ = exec.Command("bash", "-c", cmd).Output()
	publicKey := strings.TrimSpace(string(output))
	return privateKey, publicKey
}
