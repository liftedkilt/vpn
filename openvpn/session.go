package openvpn

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

type VPNSession struct {
	Path        string
	Created     string
	PID         string
	Owner       string
	Device      string
	ConfigName  string
	SessionName string
	Status      string
}

func CleanSessions() error {
	sessions, err := GetSessions()

	if err != nil {
		return err
	}

	for _, session := range sessions {
		if session.Status != "Connection, Client connected" {
			cmd := exec.Command("openvpn3", "session-manage", "--disconnect", "--path", session.Path)

			cmd.Stderr = os.Stderr

			cmd.Run()
		}
	}

	return nil
}

func StartSession(region string) error {
	return exec.Command("openvpn3", "session-start", "--config", region).Run()
}

func StopSessions(region string) error {
	sessions, err := GetSessions()

	if err != nil {
		return err
	}

	if len(sessions) == 0 {
		return fmt.Errorf("VPN session '%s' not found", region)
	}

	for _, session := range sessions {
		if session.ConfigName == region {

			cmd := exec.Command("openvpn3", "session-manage", "--disconnect", "--path", session.Path)

			cmd.Stderr = os.Stderr

			cmd.Run()
		}
	}

	return nil
}

func ListSessions() error {
	sessions, err := GetSessions()

	if err != nil {
		return err
	}

	for _, session := range sessions {
		if session.Device != "" {
			var ip net.IP

			int, _ := net.InterfaceByName(session.Device)

			addrs, err := int.Addrs()
			if err != nil {
				return err
			}

			for _, addr := range addrs {
				ipNet, ok := addr.(*net.IPNet)
				if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
					ip = ipNet.IP
				}
			}

			fmt.Printf("VPN region '%s' in status: '%s', IP: %s\n", session.ConfigName, session.Status, ip)
		} else {
			fmt.Printf("VPN region '%s' in status: '%s'\n", session.ConfigName, session.Status)
		}
	}

	return nil

}

func GetSessions() ([]VPNSession, error) {
	cmd := exec.Command("openvpn3", "sessions-list")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		return nil, err
	}

	return parseSessions(out.String()), nil
}

func parseSessions(output string) []VPNSession {
	sessions := []VPNSession{}
	scanner := bufio.NewScanner(strings.NewReader(output))
	var session VPNSession
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Path:") {
			if session.Path != "" { // Save the previous session before starting a new one
				sessions = append(sessions, session)
				session = VPNSession{} // Reset the struct
			}
			session.Path = parseField(line)
		} else if strings.Contains(line, "Created:") {
			session.Created = parseFieldWithIndex(line, 1)
			session.PID = parseFieldWithIndex(line, 3)
		} else if strings.Contains(line, "Owner:") {
			session.Owner = parseField(line)
			if fields := strings.Fields(line); len(fields) > 3 {
				session.Device = fields[3] // Get the device from the part that contains "Device: xxx"
			}
		} else if strings.Contains(line, "Config name:") {
			session.ConfigName = parseField(line)
		} else if strings.Contains(line, "Session name:") {
			session.SessionName = parseField(line)
		} else if strings.Contains(line, "Status:") {
			session.Status = parseField(line)
		}
	}
	if session.Path != "" { // Save the last session
		sessions = append(sessions, session)
	}
	return sessions
}

// parseField extracts the part after the colon and trims spaces
func parseField(line string) string {
	parts := strings.Split(line, ":")
	if len(parts) > 1 {
		return strings.TrimSpace(parts[1])
	}
	return ""
}

// parseFieldWithIndex extracts fields separated by spaces and returns the field at index if available
func parseFieldWithIndex(line string, index int) string {
	parts := strings.Fields(line)
	if len(parts) > index {
		return parts[index]
	}
	return ""
}
