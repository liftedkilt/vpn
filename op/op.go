package op

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func GetConfig(configPath, region string) error {
	doc := "[OpenVPN] " + region

	filePath := configPath + region + ".ovpn"

	fmt.Printf("Downloading config for %s to %s\n", region, filePath)

	cmd := exec.Command("op", "document", "get", doc, "--out-file", filePath)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func GetDocument(doc string) ([]byte, error) {
	cmd := exec.Command("op", "document", "get", doc)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error executing command: %v, stderr: %s", err, stderr.String())
	}
	return output, nil
}
