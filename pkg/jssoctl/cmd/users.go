package cmd

import (
	"fmt"

	"github.com/jrockway/jsso2/pkg/jssopb"
	"github.com/jrockway/jsso2/pkg/types"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
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
			req := &jssopb.EditUserRequest{
				User: &types.User{
					Id:       0,
					Username: args[0],
				},
			}
			reply, err := clientset.UserClient.Edit(cmd.Context(), req)
			if err != nil {
				return fmt.Errorf("add user: %w", err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), protojson.Format(reply))
			fmt.Fprintln(cmd.ErrOrStderr(), "OK")
			return nil
		},
	}
)

func init() {
	usersCmd.AddCommand(addUserCmd)
	AddClientset(addUserCmd)
}
