package cmd

import (
	"bytes"
	"testing"

	"github.com/jrockway/jsso2/pkg/client"
	"github.com/jrockway/jsso2/pkg/jssopb"
	"github.com/jrockway/jsso2/pkg/jtesting"
	"github.com/jrockway/jsso2/pkg/testserver"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestUsersAdd(t *testing.T) {
	jtesting.Run(t, "add user", jtesting.R{Logger: true, Database: true, GRPC: testserver.Setup}, func(t *testing.T, e *jtesting.E) {
		clientset = client.FromCC(e.ClientConn)
		rootCmd.SetArgs([]string{"users", "add", "test"})

		out := new(bytes.Buffer)
		rootCmd.SetOut(out)
		err := new(bytes.Buffer)
		rootCmd.SetErr(err)

		if err := addUserCmd.ExecuteContext(e.Context); err != nil {
			t.Fatalf("execute: %v", err)
		}
		user := &jssopb.EditUserReply{}
		if err := protojson.Unmarshal(out.Bytes(), user); err != nil {
			t.Errorf("unmarshal output: %v", err)
		}
		if got, want := user.GetUser().GetUsername(), "test"; got != want {
			t.Errorf("username:\n  got: %v\n want: %v", got, want)
		}
		if got, want := err.String(), "OK\n"; got != want {
			t.Errorf("status output:\n  got: %v\n want: %v", got, want)
		}

		out.Reset()
		err.Reset()
		if err := addUserCmd.ExecuteContext(e.Context); err == nil {
			t.Errorf("expected second call to fail")
		}
	})
}
