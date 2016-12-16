package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "netlifyctl v0.custom"

var versionCmd = &cobra.Command{
	Run: showVersion,
	Use: "version",
}

func showVersion(cmd *cobra.Command, args []string) {
	fmt.Println(Version)
}
