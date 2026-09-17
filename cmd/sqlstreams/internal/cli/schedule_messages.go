package cli

import (
	"log/slog"

	"github.com/allegedlyreliable/sqlstreams/pkg/schedule"
	"github.com/spf13/cobra"
)

func newScheduleMessagesCmd(g *globalFlags) *cobra.Command {
	var limit int
	cmd := &cobra.Command{
		Use:   "messages <name>",
		Short: "List outcomes for a schedule's newest retained messages",
		Long: `Show one outcome per message and matching consumer group, newest message first.
The limit counts messages, not output rows: five messages with two matching
groups can produce ten rows. With no matching groups, no rows are returned.`,
		Args: requireScheduleName("messages"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if limit <= 0 {
				return failUsage("--limit must be > 0, got %d", limit)
			}
			ctx := cmd.Context()
			connection, err := newConnection(ctx, g.databaseURL, g.schema, slog.LevelError)
			if err != nil {
				return err
			}
			defer connection.Close()

			rows, err := connection.client.Scheduler(args[0]).Messages(ctx, limit)
			if err != nil {
				return translateAdminError(err)
			}
			if g.jsonOutput() {
				if rows == nil {
					rows = make([]*schedule.ScheduleMessageStatus, 0)
				}
				writeJSON(cmd.OutOrStdout(), rows)
			} else {
				printScheduleMessages(cmd.OutOrStdout(), rows)
			}
			return nil
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 20, "maximum number of retained messages. Each may have a row per consumer group")
	return cmd
}
