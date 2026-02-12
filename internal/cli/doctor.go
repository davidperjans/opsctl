package cli

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/davidperjans/opsctl/internal/execx"
	"github.com/spf13/cobra"
)

type depCheck struct {
	Name        string
	Command     string
	Args        []string
	Required    bool
	InstallHint string
}

func NewDoctorCmd(r execx.Runner) *cobra.Command {
	var strict bool

	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check dependencies",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			checks := []depCheck{
				{Name: "Go", Command: "go", Args: []string{"version"}, Required: true, InstallHint: "Install Go from https://go.dev/dl/"},
				{Name: "Git", Command: "git", Args: []string{"--version"}, Required: true, InstallHint: "Install Git from https://git-scm.com/downloads"},
				{Name: "Docker", Command: "docker", Args: []string{"--version"}, Required: false, InstallHint: "Install Docker Desktop / Engine"},
				{Name: "golangci-lint", Command: "golangci-lint", Args: []string{"version"}, Required: false, InstallHint: "Install from https://golangci-lint.run/"},
			}

			fmt.Printf("opsctl doctor — %s/%s (%s)\n\n", runtime.GOOS, runtime.GOARCH, runtime.Version())

			var missingRequired []depCheck
			var missingOptional []depCheck

			for _, c := range checks {
				out, err := r.Run(ctx, c.Command, c.Args...)
				if err != nil {
					fmt.Printf("❌ %-14s %s\n", c.Name, requiredLabel(c.Required, strict))
					if msg := oneLine(out); msg != "" {
						fmt.Printf("   %s\n", msg)
					}
					fmt.Printf("   Hint: %s\n\n", c.InstallHint)

					if c.Required || strict {
						missingRequired = append(missingRequired, c)
					} else {
						missingOptional = append(missingOptional, c)
					}
					continue
				}

				fmt.Printf("✅ %-14s %s\n", c.Name, requiredLabel(c.Required, strict))
				fmt.Printf("   %s\n\n", oneLine(out))
			}

			if len(missingRequired) > 0 {
				fmt.Println("Missing required dependencies:")
				for _, m := range missingRequired {
					fmt.Printf(" - %s: %s\n", m.Name, m.InstallHint)
				}
				os.Exit(2)
			}

			if len(missingOptional) > 0 {
				names := make([]string, 0, len(missingOptional))
				for _, m := range missingOptional {
					names = append(names, m.Name)
				}

				fmt.Printf("Environment OK ✅ (optional tools missing: %s)\n", strings.Join(names, ", "))
				return nil
			}

			fmt.Println("All good ✅")
			return nil
		},
	}
	cmd.Flags().BoolVar(&strict, "strict", false, "Treat optional dependencies as required")
	return cmd
}

func requiredLabel(required bool, strict bool) string {
	if required || strict {
		return "(required)"
	}
	return "(optional)"
}
