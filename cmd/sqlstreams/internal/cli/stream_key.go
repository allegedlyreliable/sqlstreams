package cli

import (
	"github.com/spf13/cobra"
)

func newStreamKeyCmd(g *globalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "key",
		Short: "Inspect retained messages for a stream key",
		Long: `Inspect a key's compaction head and retained messages. Set a message key
when producing a message with ProduceOptions.MessageKey.

Payloads are displayed as JSON.`,
	}

	cmd.AddCommand(newStreamKeyCompactionHeadCmd(g))
	cmd.AddCommand(newStreamKeyMessagesCmd(g))

	return cmd
}
