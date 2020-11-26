package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/jrockway/jsso2/pkg/client"
	"github.com/jrockway/jsso2/pkg/jssopb"
	"github.com/jrockway/jsso2/pkg/jtesting"
	"github.com/jrockway/jsso2/pkg/types"
	"github.com/jrockway/jsso2/pkg/util/zapwriter"
	"go.uber.org/zap"
)

func checkServer(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:4000", nil)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do: %w", err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("non-200 response: %s", res.Status)
	}
	return nil
}

func TestLiveServer(t *testing.T) {
	if e := os.Getenv("TEST_E2E"); e != "true" && e != "1" {
		t.Skip("TEST_E2E is not set to 'true' or '1'; skipping e2e test")
	}

	jtesting.Run(t, "e2e", jtesting.R{Database: true, Logger: true}, func(t *testing.T, e *jtesting.E) {
		os.Unsetenv("GOMAXPROCS")
		os.Unsetenv("HTTP_ADDRESS")
		os.Unsetenv("DEBUG_ADDRESS")
		os.Unsetenv("GRPC_ADDRESS")
		os.Unsetenv("DATABASE_URL")
		if err := checkServer(e.Context); err == nil {
			t.Fatal("server is already running on localhost:4000")
		}
		servers := exec.CommandContext(e.Context, "npm", "run", "all")
		servers.Env = os.Environ()
		servers.Env = append(servers.Env, "DATABASE_URL="+e.DSN)
		servers.Stdout = zapwriter.New(e.Logger.Named("all.stdout"))
		servers.Stderr = zapwriter.New(e.Logger.Named("all.stderr"))
		servers.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
		if err := servers.Start(); err != nil {
			t.Fatal(err)
		}
		serversDoneCh := make(chan error)
		go func() {
			serversDoneCh <- servers.Wait()
			close(serversDoneCh)
		}()
		killServers := func() {
			e.Logger.Info("killing servers")
			if err := syscall.Kill(-servers.Process.Pid, syscall.SIGINT); err != nil {
				e.Logger.Info("failed to kill servers", zap.Error(err))
			}
		}

		for i := 0; i < 10; i++ {
			if err := checkServer(e.Context); err != nil {
				e.Logger.Info("waiting for servers to start")
				time.Sleep((time.Duration(i) + 1) * 100 * time.Millisecond)
				continue
			}
			break
		}

		cs, err := client.Dial(e.Context, "localhost:4000", &client.Credentials{Root: "root"})
		if err != nil {
			t.Fatalf("dial jsso: %v", err)
		}
		ur, err := cs.UserClient.Edit(e.Context, &jssopb.EditUserRequest{
			User: &types.User{
				Id:       0,
				Username: "the-tests",
			},
		})
		if err != nil {
			t.Fatalf("create user: %v", err)
		}
		er, err := cs.UserClient.GenerateEnrollmentLink(e.Context, &jssopb.GenerateEnrollmentLinkRequest{
			Target: ur.GetUser(),
		})
		if err != nil {
			t.Fatalf("create enrollment link: %v", err)
		}
		link := er.GetUrl()
		e.Logger.Info("enrollment link generated", zap.String("link", link))

		cypress := exec.Command("npm", "run", "cypress:run")
		cypress.Env = os.Environ()
		cypress.Env = append(cypress.Env, "ENROLLMENT_LINK="+link)
		cypress.Stdout = zapwriter.New(e.Logger.Named("cypress.stdout"))
		cypress.Stderr = zapwriter.New(e.Logger.Named("cypress.stderr"))
		servers.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
		if err := cypress.Start(); err != nil {
			t.Fatalf("run cypress: %v", err)
		}
		cypressDoneCh := make(chan error)
		go func() {
			cypressDoneCh <- cypress.Wait()
			close(cypressDoneCh)
		}()
		killCypress := func() {
			e.Logger.Info("killing cypress")
			if err := syscall.Kill(-cypress.Process.Pid, syscall.SIGKILL); err != nil {
				e.Logger.Info("failed to kill cypress", zap.Error(err))
			}
		}

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		select {
		case <-sigCh:
			e.Logger.Info("interrupt")
			killServers()
			select {
			case <-time.After(time.Second):
			case <-serversDoneCh:
			}
			t.Fatal("interrupt")
		case <-e.Context.Done():
			t.Error(e.Context.Err())
			killServers()
			select {
			case <-time.After(time.Second):
			case <-serversDoneCh:
			}
			killCypress()
			select {
			case <-time.After(time.Second):
			case <-cypressDoneCh:
			}
			t.Fatal("timeout")
		case err := <-serversDoneCh:
			if err != nil {
				t.Errorf("servers: %v", err)
			}
			killCypress()
			select {
			case <-time.After(time.Second):
			case <-cypressDoneCh:
			}
			t.Fatal("the servers died")
		case err := <-cypressDoneCh:
			if err != nil {
				t.Errorf("cypress: %v", err)
			}
			killServers()
			select {
			case <-time.After(time.Second):
			case <-serversDoneCh:
			}
		}
	})
}
