package cli

import "github.com/spf13/cobra"

func newSystemBindingCmd(g *globalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "binding",
		Short: "Inspect binding declarations across the deployment",
	}

	cmd.AddCommand(newSystemBindingListCmd(g))

	return cmd
}
