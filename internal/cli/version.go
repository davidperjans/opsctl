package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version   = "dev"
	commitSHA = "none"
	buildDate = "unknown"
)

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("opsctl %s\n", version)
			fmt.Printf("commit: %s\n", commitSHA)
			fmt.Printf("built:  %s\n", buildDate)
		},
	}
}
