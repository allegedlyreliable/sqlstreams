package cli

import (
	"github.com/spf13/cobra"
)

func newConsumerBindingCmd(g *globalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "binding",
		Short: "Read a consumer's binding declaration",
		Long: `A consumer's binding set comes from ConsumerConfig.Bindings, declared through
client.Stream(name).Consumer(name).Register. To change it, update the
application's declaration and register it again. These commands only read.

A consumer with no binding patterns receives every message on its stream.`,
	}

	cmd.AddCommand(newConsumerBindingGetCmd(g))

	return cmd
}
