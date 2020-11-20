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

	generateEnrollmentLinkCmd = &cobra.Command{
		Use:     "generate-enrollment-link",
		Aliases: []string{"enroll"},
		Short:   "Generate a link for a user to enroll their security token.",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			user := &types.User{}
			if id, err := cmd.Flags().GetInt64("id"); err != nil {
				return fmt.Errorf("get id: %w", err)
			} else if id != 0 {
				user.Id = id
			}
			if username, err := cmd.Flags().GetString("username"); err != nil {
				return fmt.Errorf("get username: %w", err)
			} else if username != "" {
				user.Username = username
			}

			req := &jssopb.GenerateEnrollmentLinkRequest{
				Target: user,
			}
			reply, err := clientset.UserClient.GenerateEnrollmentLink(cmd.Context(), req)
			if err != nil {
				return fmt.Errorf("generate enrollment link: %w", err)
			}
			if jsonOutput {
				fmt.Fprintln(cmd.OutOrStdout(), protojson.Format(reply))
				fmt.Fprintln(cmd.ErrOrStderr(), "OK")
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "Click this link to enroll an authenticator: %s\n", reply.GetUrl())
			}
			return nil
		},
	}

	whoAmICmd = &cobra.Command{
		Use:   "whoami",
		Short: "Print some information about your current session.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			reply, err := clientset.UserClient.WhoAmI(cmd.Context(), &jssopb.WhoAmIRequest{})
			if err != nil {
				return fmt.Errorf("whoami: %w", err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), protojson.Format(reply))
			fmt.Fprintln(cmd.ErrOrStderr(), "OK")
			return nil
		},
	}
)

func init() {
	generateEnrollmentLinkCmd.Flags().String("username", "", "the name of the user to enroll")
	generateEnrollmentLinkCmd.Flags().Int64("id", 0, "the id of the user to enroll")
	usersCmd.AddCommand(addUserCmd, generateEnrollmentLinkCmd, whoAmICmd)
	AddClientset(addUserCmd)
	AddClientset(generateEnrollmentLinkCmd)
	AddClientset(whoAmICmd)
}
