package cmd

import (
	"github.com/spf13/cobra"
)

func (c *command) initInitCmd() (err error) {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialise a Swarm node",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) > 0 {
				return cmd.Help()
			}

			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return c.config.BindPFlags(cmd.Flags())
		},
	}

	c.setAllFlags(cmd)
	c.root.AddCommand(cmd)
	return nil
}
