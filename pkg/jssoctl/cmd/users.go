package cmd

import (
	"fmt"

	"github.com/jrockway/jsso2/pkg/jssopb"
	"github.com/spf13/cobra"
)

var (
	usersCmd = &cobra.Command{
		Use:     "users",
		Aliases: []string{"user"},
		Short:   "Manage users",
	}

	addUserCmd = &cobra.Command{
		Use:          "add [username]",
		Short:        "Add a new user",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			req := &jssopb.AddUserRequest{
				Username: args[0],
			}
			if _, err := clientset.UserClient.Add(cmd.Context(), req); err != nil {
				return fmt.Errorf("add user: %w", err)
			}
			cmd.Println("OK")
			return nil
		},
	}
)

func init() {
	usersCmd.AddCommand(addUserCmd)
	AddClientset(addUserCmd)
}
