/*
Copyright © 2024 William Forsyth william.forsyth@liferay.com
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

 1. Redistributions of source code must retain the above copyright notice,
    this list of conditions and the following disclaimer.

 2. Redistributions in binary form must reproduce the above copyright notice,
    this list of conditions and the following disclaimer in the documentation
    and/or other materials provided with the distribution.

 3. Neither the name of the copyright holder nor the names of its contributors
    may be used to endorse or promote products derived from this software
    without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vpn",
	Short: "manage your vpn connections",
	Long: `vpn is a CLI tool to manage your vpn connections.
It can list, connect, and disconnect from vpn connections.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	data, err := os.ReadFile("version.txt")

	version := string(data)

	if err != nil {
		cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
		output, err := cmd.Output()
		if err != nil {
			version = "dev-unknown"
		} else {
			version = "dev-" + strings.TrimSpace(string(output))
		}
	}

	rootCmd.Version = version

	home, err := os.UserHomeDir()

	if err != nil {
		// Handle error if the home directory cannot be determined
		home = "/tmp"
	}

	confPath := home + "/vpn/"

	rootCmd.PersistentFlags().StringP("configpath", "c", confPath, "Directory containing VPN configs")

	viper.SetConfigName(".vpn")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(confPath)

	err = viper.ReadInConfig()

	if err != nil {
		// Prompt user for default region
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter default region: ")
		region, _ := reader.ReadString('\n')

		viper.Set("default-region", strings.TrimSpace(region))

		err := viper.WriteConfigAs(confPath + ".vpn.yaml")
		if err != nil {
			fmt.Printf("Error writing to %s.vpn.yaml: %s\n", confPath, err)
		}
	}

	// Set region flag from Viper
	rootCmd.PersistentFlags().StringP("region", "r", viper.GetString("default-region"), "VPN Region")

}
