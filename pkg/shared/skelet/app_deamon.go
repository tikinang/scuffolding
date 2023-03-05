package skelet

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

func Daemon[T any](
	project, name, version string,
	flavors FlavorProvider,
	bones ...any,
) {

	cmd := &cobra.Command{
		Version: version,
		Use:     fmt.Sprintf("%s-%s", project, name),
		Long:    fmt.Sprintf("Command-line interface for %s-%s.", project, name),
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "run",
			Short: "Runs daemon, terminate with SIGTERM.",
			RunE: func(cmd *cobra.Command, args []string) error {
				runner, _, err := AssembleRunner[T](cmd, name, flavors, bones...)
				if err != nil {
					return err
				}

				ctx, cancel := signal.NotifyContext(
					context.Background(),
					syscall.SIGTERM,
					syscall.SIGINT,
				)
				defer cancel()

				return runner.Run(ctx)
			},
		},
	)

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
