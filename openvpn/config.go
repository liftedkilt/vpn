package openvpn

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/liftedkilt/vpn/op"
)

type OpenVPNConfigs map[string]OpenVPNConfig

type OpenVPNConfig struct {
	ACL                  OpenVPNConfigACL `json:"acl"`
	DCO                  bool             `json:"dco"`
	Imported             string           `json:"imported"`        // As a string or time.Time based on usage
	ImportedTimestamp    int64            `json:"imported_tstamp"` // UNIX timestamp
	LastUsed             string           `json:"lastused"`        // As a string or time.Time based on usage
	LastUsedTimestamp    int64            `json:"lastused_tstamp"` // UNIX timestamp
	Name                 string           `json:"name"`
	TransferOwnerSession bool             `json:"transfer_owner_session"`
	UseCount             int              `json:"use_count"`
}

type OpenVPNConfigACL struct {
	LockedDown   bool   `json:"locked_down"`
	Owner        string `json:"owner"`
	PublicAccess bool   `json:"public_access"`
}

func FindConfig(region string) bool {
	cmd := exec.Command("openvpn3", "configs-list", "--json")

	output, _ := cmd.StdoutPipe()

	cmd.Start()

	outputBytes, _ := io.ReadAll(output)

	cmd.Wait()

	configs := OpenVPNConfigs{}

	json.Unmarshal(outputBytes, &configs)

	for _, config := range configs {
		if config.Name == region {
			return true
		}
	}

	return false
}

func ImportConfig(prefix, path, region string) error {
	configFile := path + region + ".ovpn"

	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Config file not found, attempting to download it")

		config, err := op.GetDocument(prefix + " " + region)

		if err != nil {
			return fmt.Errorf("unable to obtain config '%s %s' from 1Password: %s", prefix, region, err)
		}

		err = os.WriteFile(configFile, config, 0644)

		if err != nil {
			return err
		}
	}

	cmd := exec.Command("openvpn3", "config-import", "--config", configFile, "--name", region, "--persistent")

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func DeleteConfigByName(region string) error {
	cmd := exec.Command("openvpn3", "config-remove", "--name", region)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func DeleteConfigByPath(path string) error {
	cmd := exec.Command("openvpn3", "config-remove", "--config", path)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
