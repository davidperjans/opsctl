package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/davidperjans/opsctl/internal/execx"
	"github.com/spf13/cobra"
)

type ciStep struct {
	Name string
	Run  func(ctx context.Context) (string, error)
}

func NewCiCmd(r execx.Runner) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ci",
		Short: "CI helpers (run standard checks locally)",
	}

	cmd.AddCommand(NewCiRunCmd(r))
	return cmd
}

func NewCiRunCmd(r execx.Runner) *cobra.Command {
	var skipFmt bool
	var skipBuild bool
	var race bool
	var verbose bool
	var timeout time.Duration

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run fmt + test + build (similar to CI)",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			if timeout > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, timeout)
				defer cancel()
			}

			fmt.Println("opsctl ci run\n")

			steps := make([]ciStep, 0, 3)

			if !skipFmt {
				steps = append(steps, ciStep{
					Name: "go fmt ./...",
					Run: func(ctx context.Context) (string, error) {
						return r.Run(ctx, "go", "fmt", "./...")
					},
				})
			}

			testArgs := []string{"test"}
			if race {
				testArgs = append(testArgs, "-race")
			}
			testArgs = append(testArgs, "./...")

			steps = append(steps, ciStep{
				Name: "go " + joinArgs(testArgs),
				Run: func(ctx context.Context) (string, error) {
					return r.Run(ctx, "go", testArgs...)
				},
			})

			if !skipBuild {
				steps = append(steps, ciStep{
					Name: "go build ./...",
					Run: func(ctx context.Context) (string, error) {
						return r.Run(ctx, "go", "build", "./...")
					},
				})
			}

			for i, step := range steps {
				fmt.Printf("[%d/%d] %s\n", i+1, len(steps), step.Name)

				out, err := step.Run(ctx)
				if err != nil {
					fmt.Printf("❌ Failed: %s\n", step.Name)
					if out != "" {
						fmt.Printf("   %s\n", oneLine(out))
					}
					os.Exit(1)
				}

				fmt.Printf("✅ OK\n")
				if verbose && out != "" {
					fmt.Printf("   %s\n", oneLine(out))
				}
				fmt.Println()
			}

			fmt.Println("CI checks passed ✅")
			return nil
		},
	}

	cmd.Flags().BoolVar(&skipFmt, "skip-fmt", false, "Skip go fmt")
	cmd.Flags().BoolVar(&skipBuild, "skip-build", false, "Skip go build")
	cmd.Flags().BoolVar(&race, "race", false, "Run tests with -race (may be slower)")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "Print command output on success")
	cmd.Flags().DurationVar(&timeout, "timeout", 10*time.Second, "Command timeout (best-effort)")

	return cmd
}
