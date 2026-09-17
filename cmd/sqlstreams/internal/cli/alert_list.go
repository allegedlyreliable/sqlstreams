package cli

import (
	"fmt"
	"io"
	"log/slog"
	"text/tabwriter"

	"github.com/allegedlyreliable/sqlstreams/client"
	"github.com/allegedlyreliable/sqlstreams/pkg/common"
	"github.com/spf13/cobra"
)

func newAlertListCmd(g *globalFlags) *cobra.Command {
	var (
		quiet        bool
		streamName   string
		consumerName string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List the latest retained alert for each name and owner",
		Long: `List the latest retained alert for each name and owner, active or resolved.

Show all owners by default. Use --stream for alerts owned by one stream,
or --stream and --consumer for alerts owned by one consumer group.`,
		Example: `sqlstreams alert list
sqlstreams alert list --stream orders.created
sqlstreams alert list --stream orders.created --consumer billing`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			out := cmd.OutOrStdout()

			if quiet && g.jsonOutput() {
				return failUsage("--quiet and --output json cannot be combined")
			}
			if consumerName != "" && streamName == "" {
				return failUsage("--consumer requires --stream")
			}

			connection, err := newConnection(ctx, g.databaseURL, g.schema, slog.LevelError)
			if err != nil {
				return err
			}
			defer connection.Close()
			client := connection.client

			var alerts []*sqlstreams.Alert
			switch {
			case consumerName != "":
				alerts, err = client.Stream[sqlstreams.RawPayload](streamName).Consumer(consumerName).Alerts().Latest(ctx)
			case streamName != "":
				alerts, err = client.Stream[sqlstreams.RawPayload](streamName).Alerts().Latest(ctx)
			default:
				alerts, err = client.System().Alerts().Latest(ctx)
			}
			if err != nil {
				return translateAdminError(err)
			}

			if g.jsonOutput() {
				writeJSON(out, alerts)
				return nil
			}

			if quiet {
				printAlertKeys(out, alerts)
			} else {
				printAlertsTable(out, alerts)
			}
			return nil
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&quiet, "quiet", "q", false, "alert and owner only, one per line. Incompatible with --output json")
	f.StringVar(&streamName, "stream", "", "only alerts owned by this stream")
	f.StringVar(&consumerName, "consumer", "", "only alerts owned by this consumer group. Requires --stream")
	return cmd
}

// ownerCell - "stream/orders", "system/system".
func ownerCell(owner *common.Owner) string {
	return fmt.Sprintf("%s/%s", owner.Kind(), owner.Name)
}

func printAlertKeys(w io.Writer, alerts []*sqlstreams.Alert) {
	for _, published := range alerts {
		fmt.Fprintf(w, "%s %s\n", published.Name, ownerCell(published.Owner))
	}
}

func printAlertsTable(w io.Writer, alerts []*sqlstreams.Alert) {
	if len(alerts) == 0 {
		fmt.Fprintln(w, "no alerts published")
		return
	}

	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tOWNER\tSTATUS\tSEVERITY\tSINCE\tMESSAGE")
	for _, published := range alerts {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\n",
			published.Name, ownerCell(published.Owner), published.Status,
			published.Severity, timeCell(published.At), published.Message)
	}
	tw.Flush()

	fmt.Fprintf(w, "\n%s\n", pluralize(len(alerts), "alert"))
}
