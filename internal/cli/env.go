package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/davidperjans/opsctl/internal/envcheck"
	"github.com/spf13/cobra"
)

func NewEnvCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env",
		Short: "Validate .env files",
	}

	cmd.AddCommand(NewEnvCheckCmd())
	return cmd
}

func NewEnvCheckCmd() *cobra.Command {
	var examplePath string
	var envPath string

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check that .env contains all keys from .env.example",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("opsctl env check\n\n")
			fmt.Printf("Example: %s\nEnv:     %s\n\n", examplePath, envPath)

			res, err := envcheck.Check(examplePath, envPath)
			if err != nil {
				// treat missing file etc as "environment not ready"
				fmt.Printf("❌ %s\n", err.Error())
				os.Exit(2)
			}

			if len(res.Missing) == 0 && len(res.Extra) == 0 {
				fmt.Println("All good ✅")
				return nil
			}

			if len(res.Missing) > 0 {
				fmt.Printf("❌ Missing keys (%d):\n", len(res.Missing))
				for _, k := range res.Missing {
					fmt.Printf(" - %s\n", k)
				}
				fmt.Println()
			} else {
				fmt.Println("✅ No missing keys")
				fmt.Println()
			}

			if len(res.Extra) > 0 {
				fmt.Printf("⚠️  Extra keys (%d) (not in example):\n", len(res.Extra))
				for _, k := range res.Extra {
					fmt.Printf(" - %s\n", k)
				}
				fmt.Println()
			}

			if len(res.Missing) > 0 {
				fmt.Printf("Fix: Add the missing keys to %s (you can copy them from %s).\n", envPath, examplePath)
				os.Exit(1)
			}

			// only extras => still OK
			fmt.Printf("Environment OK ✅ (extra keys present: %s)\n", strings.Join(res.Extra, ", "))
			return nil
		},
	}

	cmd.Flags().StringVar(&examplePath, "example", ".env.example", "Path to example env file")
	cmd.Flags().StringVar(&envPath, "env", ".env", "Path to env file")
	return cmd
}
