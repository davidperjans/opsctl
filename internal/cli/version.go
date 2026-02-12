package cli

import (
	"fmt"
	"runtime/debug"

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
			printVersion()
		},
	}
}

func printVersion() {
	v := version
	c := commitSHA
	d := buildDate

	// Fallback for binaries built via `go install ...@vX.Y.Z`
	if bi, ok := debug.ReadBuildInfo(); ok {
		// Module version (often v0.1.0 when installed with @v0.1.0)
		if (v == "dev" || v == "") && bi.Main.Version != "" && bi.Main.Version != "(devel)" {
			v = bi.Main.Version
		}

		// VCS settings (revision/time)
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.revision":
				if c == "none" || c == "" {
					c = s.Value
				}
			case "vcs.time":
				if d == "unknown" || d == "" {
					d = s.Value
				}
			}
		}
	}

	fmt.Printf("opsctl %s\n", v)
	fmt.Printf("commit: %s\n", c)
	fmt.Printf("built:  %s\n", d)
}
