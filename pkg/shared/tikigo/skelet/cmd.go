package skelet

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

func Daemon[T any](
	project, name, version string,
	flavors FlavorProvider,
	skelet T,
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
			Short: "Runs daemon, terminate with SIGTERM",
			RunE: func(cmd *cobra.Command, args []string) error {

				runner, err := AssembleRunner(cmd, name, flavors, skelet, bones...)
				if err != nil {
					return err
				}

				ctx, cancel := signal.NotifyContext(context.Background(), TerminationSignals()...)
				defer cancel()

				return runner.run(ctx)
			},
		},
	)

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}

// TerminationSignals returns slice with default termination signals
// as defined here https://www.gnu.org/software/libc/manual/html_node/Termination-Signals.html
func TerminationSignals() []os.Signal {
	return []os.Signal{
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGKILL,
		syscall.SIGHUP,
	}
}
