package cmd

import (
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func (c *command) initConfiguratorOptionsCmd() (err error) {

	cmd := &cobra.Command{
		Use:   "print-config",
		Short: "Print default or provided configuration in yaml format",
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			if len(args) > 0 {
				return cmd.Help()
			}

			d := c.config.AllSettings()
			ym, err := yaml.Marshal(d)
			if err != nil {
				return err
			}
			cmd.Println(string(ym))
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
