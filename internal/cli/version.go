package cli

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var (
	version   = "dev"
	commitSHA = ""
	buildDate = ""
)

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			v, c, d := resolveVersionInfo()

			// Alltid visa version
			fmt.Printf("opsctl %s\n", v)

			// Visa bara commit/built om de finns
			if c != "" {
				fmt.Printf("commit: %s\n", c)
			}
			if d != "" {
				fmt.Printf("built:  %s\n", d)
			}
		},
	}
}

func resolveVersionInfo() (string, string, string) {
	v := version
	c := commitSHA
	d := buildDate

	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return v, c, d
	}

	// Version fallback (när ldflags inte satt den)
	if (v == "dev" || v == "" || v == "(devel)") && bi.Main.Version != "" && bi.Main.Version != "(devel)" {
		v = bi.Main.Version
	}

	// Commit/time fallback (finns ibland, men inte alltid vid go install)
	for _, s := range bi.Settings {
		switch s.Key {
		case "vcs.revision":
			if c == "" {
				c = s.Value
			}
		case "vcs.time":
			if d == "" {
				d = s.Value
			}
		}
	}

	return v, c, d
}
