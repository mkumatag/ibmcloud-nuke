package cmd

import (
	"fmt"
	"github.com/mkumatag/ibmcloud-nuke/internal/app"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	command := &cobra.Command{
		Use:          "ibmcloud-nuke",
		Short:        "ibmcloud-nuke removes the mentioned resource from IBM Cloud",
		Long:         `Tool which removes the specified resource from the IBM Cloud account based on the filters`,
		SilenceUsage: true,
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("only one argument is expected which is a config file name but provided: %s", args)
		}
		return app.Run(args[0])
	}
	return command
}
