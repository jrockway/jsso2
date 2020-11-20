package cmd

import (
	"fmt"

	"github.com/jrockway/jsso2/pkg/client"
	"github.com/spf13/cobra"
)

var (
	clientset *client.Set
	noClose   bool
)

func AddClientset(c *cobra.Command) {
	if c.PreRun != nil || c.PreRunE != nil {
		panic("command already has a pre-run function")
	}
	c.PreRunE = func(cmd *cobra.Command, args []string) error {
		if clientset != nil {
			// Skip dialing for the tests.
			return nil
		}
		credentials := &client.Credentials{
			Bearer: bearer,
			Root:   root,
			Token:  session,
		}
		set, err := client.Dial(cmd.Context(), address, credentials)
		if err != nil {
			return fmt.Errorf("connect: %v", err)
		}
		clientset = set
		return nil
	}
	if c.PostRun != nil || c.PostRunE != nil {
		panic("command already has a post-run function")
	}
	c.PostRun = func(cmd *cobra.Command, args []string) {
		if !noClose {
			if err := clientset.Close(); err != nil {
				cmd.PrintErrf("while closing client: %v\n", err)
			}
		}
	}
}
