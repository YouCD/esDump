package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)
var (
	Version   string
	commitID  string
	buildTime string
	goVersion string
	buildUser string
	name ="esDump"
)
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: fmt.Sprintf("Print the version number of %s",name),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version:   %s\n", Version)
		fmt.Printf("CommitID:  %s\n", commitID)
		fmt.Printf("BuildTime: %s\n", buildTime)
		fmt.Printf("GoVersion: %s\n", goVersion)
		fmt.Printf("BuildUser: %s\n", buildUser)
	},
}
