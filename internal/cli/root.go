package cli

import (
	"time"

	"github.com/davidperjans/opsctl/internal/execx"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "opsctl",
		Short: "Developer ops helper for Go services",
	}

	runner := execx.OSRunner{Timeout: 5 * time.Second}

	root.AddCommand(
		NewDoctorCmd(runner),
		NewEnvCmd(),
		NewCiCmd(runner),
		NewInitCmd(),
		NewVersionCmd(),
	)

	return root
}
