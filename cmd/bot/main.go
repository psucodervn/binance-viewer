
package bot

import (
  "github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bot",
		Short: "bot description",
		Run:   run,
	}
	return cmd
}

func run(cmd *cobra.Command, args []string) {
	// TODO: implement
}
