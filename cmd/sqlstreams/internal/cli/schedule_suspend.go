package cli

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/allegedlyreliable/sqlstreams/pkg/schedule"
	"github.com/spf13/cobra"
)

func newScheduleSuspendCmd(g *globalFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "suspend <name>",
		Short: "Suspend automatic message production for a schedule",
		Long: `Stop automatic message production until the schedule is unsuspended. Already
produced messages remain available. Manual production is still allowed with
sqlstreams scheduler run.`,
		Args: requireScheduleName("suspend"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			name := args[0]

			connection, err := newConnection(ctx, g.databaseURL, g.schema, slog.LevelError)
			if err != nil {
				return err
			}
			defer connection.Close()
			client := connection.client

			if err := client.Scheduler(name).Suspend(ctx); err != nil {
				if errors.Is(err, schedule.ErrScheduleNotFound) {
					return errScheduleNotFound(name)
				}
				return translateAdminError(err)
			}

			if g.jsonOutput() {
				writeJSON(cmd.OutOrStdout(), scheduleSuspendedDocument{Schedule: name, Suspended: true})
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s schedule %q suspended\n", glyphOK(), name)
			return nil
		},
	}
}

func newScheduleUnsuspendCmd(g *globalFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "unsuspend <name>",
		Short: "Resume automatic message production for a schedule",
		Long: `Resume automatic message production at the next scheduled time after now.
Scheduled times missed during suspension are skipped.`,
		Args: requireScheduleName("unsuspend"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			name := args[0]

			connection, err := newConnection(ctx, g.databaseURL, g.schema, slog.LevelError)
			if err != nil {
				return err
			}
			defer connection.Close()
			client := connection.client

			if err := client.Scheduler(name).Unsuspend(ctx); err != nil {
				if errors.Is(err, schedule.ErrScheduleNotFound) {
					return errScheduleNotFound(name)
				}
				return translateAdminError(err)
			}

			row, err := client.Scheduler(name).Get(ctx)
			if err != nil || row == nil {
				// the unsuspend itself succeeded -- report that even if the
				// follow-up read for the next-scheduled-time detail didn't cooperate
				if g.jsonOutput() {
					writeJSON(cmd.OutOrStdout(), scheduleSuspendedDocument{Schedule: name, Suspended: false})
					return nil
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s schedule %q unsuspended\n", glyphOK(), name)
				return nil
			}

			if g.jsonOutput() {
				writeJSON(cmd.OutOrStdout(), scheduleSuspendedDocument{
					Schedule:        name,
					Suspended:       false,
					NextScheduledAt: &row.NextScheduledAt,
				})
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s schedule %q unsuspended, next scheduled time %s\n",
				glyphOK(), name, timeCell(row.NextScheduledAt))
			return nil
		},
	}
}

// scheduleSuspendedDocument is suspend/unsuspend's json result;
// next_scheduled_at is null while suspended, and after an unsuspend whose
// follow-up read did not cooperate.
type scheduleSuspendedDocument struct {
	Schedule        string     `json:"schedule"`
	Suspended       bool       `json:"suspended"`
	NextScheduledAt *time.Time `json:"next_scheduled_at"`
}
