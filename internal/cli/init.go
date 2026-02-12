package cli

import (
	"fmt"
	"path/filepath"

	"github.com/davidperjans/opsctl/internal/scaffold"
	"github.com/davidperjans/opsctl/internal/templates"
	"github.com/spf13/cobra"
)

func NewInitCmd() *cobra.Command {
	var dir string
	var force bool
	var modulePath string

	cmd := &cobra.Command{
		Use:   "init [service-name]",
		Short: "Scaffold a new Go service",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]
			targetDir := filepath.Join(dir, serviceName)

			if modulePath == "" {
				// default module path convention
				modulePath = fmt.Sprintf("github.com/davidperjans/%s", serviceName)
			}

			fmt.Println("opsctl init")
			fmt.Printf("Service: %s\nDir:     %s\nModule:  %s\n\n", serviceName, targetDir, modulePath)

			err := scaffold.GenerateFromFS(
				templates.GoServiceFS,
				"go-service",
				scaffold.Options{
					ServiceName: serviceName,
					ModulePath:  modulePath,
					TargetDir:   targetDir,
					Force:       force,
				},
			)

			if err != nil {
				return err
			}

			fmt.Println("Created ✅")
			fmt.Println("Next:")
			fmt.Printf("  cd %s\n", targetDir)
			fmt.Println("  make run")
			return nil
		},
	}

	cmd.Flags().StringVar(&dir, "dir", ".", "Directory to create the service in")
	cmd.Flags().BoolVar(&force, "force", false, "Overwrite existing files")
	cmd.Flags().StringVar(&modulePath, "module", "", "Go module path (default: github.com/davidperjans/<service-name>)")

	return cmd
}
