package cli

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/allegedlyreliable/sqlstreams/client"
	"github.com/allegedlyreliable/sqlstreams/pkg/common"
	"github.com/allegedlyreliable/sqlstreams/pkg/schedule"
	"github.com/spf13/cobra"
)

func newScheduleRunCmd(g *globalFlags) *cobra.Command {
	var concurrency string

	cmd := &cobra.Command{
		Use:   "run <name>",
		Short: "Produce a schedule's message immediately",
		Long: `Produce a schedule's stored message immediately, even when suspended. Its
expression and next scheduled time are unchanged.

The message uses parallel concurrency by default, allowing it to run while
a previous message is still being processed. Use --concurrency exclusive to
prevent overlap. The new message supersedes older unclaimed messages under
the schedule's key. Schedules compact their messages, so ordered concurrency
is not supported.`,
		Args: requireScheduleName("run"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			name := args[0]
			f := cmd.Flags()

			// Build sparse options from only the flags that were passed.
			options := &sqlstreams.ScheduleRunOptions{}
			if f.Changed("concurrency") {
				options.Concurrency = common.ConcurrencyPolicy(concurrency)
			}

			// Validate up front for a clean usage error (a bad flag value,
			// exit 2) instead of the raw wrapped error the run returns.
			probe := *options
			probe.WithDefaults()
			if err := probe.Validate(); err != nil {
				return failUsage("invalid options: %s", err)
			}

			connection, err := newConnection(ctx, g.databaseURL, g.schema, slog.LevelError)
			if err != nil {
				return err
			}
			defer connection.Close()
			client := connection.client

			produced, err := client.Scheduler(name).Run(ctx, options)
			if err != nil {
				if errors.Is(err, schedule.ErrScheduleNotFound) {
					return errScheduleNotFound(name)
				}
				return translateAdminError(err)
			}

			if g.jsonOutput() {
				writeJSON(cmd.OutOrStdout(), scheduleRunDocument{Schedule: name, MessageId: produced.Id})
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s produced %q (message id=%d)\n",
				glyphOK(), name, produced.Id)
			return nil
		},
	}

	cmd.Flags().StringVar(&concurrency, "concurrency", "", "concurrency policy for this message: parallel or exclusive (default parallel)")

	return cmd
}

// scheduleRunDocument is schedule run's json result: the handle the run produced.
type scheduleRunDocument struct {
	Schedule  string `json:"schedule"`
	MessageId int64  `json:"message_id"`
}
