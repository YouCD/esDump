package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)
var (
	//version info
	Version   string
	commitID  string
	buildTime string
	goVersion string
	buildUser string
)
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of esDump",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version:   %s\n", Version)
		fmt.Printf("CommitID:  %s\n", commitID)
		fmt.Printf("BuildTime: %s\n", buildTime)
		fmt.Printf("GoVersion: %s\n", goVersion)
		fmt.Printf("BuildUser: %s\n", buildUser)
	},
}
