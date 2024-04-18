package op

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GetConfig(configPath, region string) error {
	regionUpper := strings.ToUpper(region)

	doc := "[OpenVPN] " + regionUpper

	filePath := configPath + region + ".ovpn"

	fmt.Printf("Downloading config for %s to %s\n", region, filePath)

	cmd := exec.Command("op", "document", "get", doc, "--out-file", filePath)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
