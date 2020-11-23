package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jrockway/jsso2/pkg/client"
	"github.com/jrockway/jsso2/pkg/jssopb"
	"github.com/olekukonko/tablewriter"
)

func Get(ctx context.Context, url string, code int) error {
	fmt.Printf("get %s...", url)
	t := time.Now()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Println("FAIL", err)
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("FAIL", err)
		return err
	}
	res.Body.Close()
	if res.StatusCode != code {
		fmt.Println("FAIL")
		return fmt.Errorf("unsuccessful response %s (want %d)", res.Status, code)
	}
	fmt.Printf("ok (%s)\n", time.Since(t).Round(time.Microsecond))
	return nil
}

func main() {
	ctx, c := context.WithTimeout(context.Background(), 30*time.Second)
	checks := []struct {
		url  string
		code int
	}{
		{"http://localhost:9901/ready", http.StatusOK},
		{"http://localhost:4000/envoy/ready", http.StatusOK},
		{"http://localhost:4000/", http.StatusOK},
		{"http://localhost:4000/build/bundle.js", http.StatusOK},
		{"http://localhost:4000/logout", http.StatusOK},
		{"http://localhost:4000/grpcui", http.StatusOK},
		{"http://localhost:8081/metrics", http.StatusOK},
		{"http://localhost:4000/backend-debug/metrics", http.StatusOK},
		{"http://localhost:8181/metrics", http.StatusOK},
		{"http://localhost:4000/authz-debug/metrics", http.StatusOK},
		{"http://localhost:8280/metrics", http.StatusOK},
		{"http://localhost:8281/metrics", http.StatusOK},
		{"http://localhost:8280/", http.StatusOK},
		{"http://localhost:4000/protected", http.StatusForbidden},
	}

	errors := make(map[string]error)
	for _, check := range checks {
		if err := Get(ctx, check.url, check.code); err != nil {
			errors[check.url] = err
		}
	}

	cs, err := client.Dial(ctx, "dns:///localhost:4000", &client.Credentials{Root: "root"})
	if err != nil {
		errors["dial grpc"] = err
	} else {
		fmt.Printf("grpc whoami...")
		_, err := cs.UserClient.WhoAmI(ctx, &jssopb.WhoAmIRequest{})
		if err != nil {
			fmt.Println("FAIL")
			errors["whoami"] = err
		} else {
			fmt.Println("ok")
		}
	}
	c()

	if len(errors) > 0 {
		fmt.Println("FAIL")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetColWidth(80)
		table.SetRowSeparator("")
		table.SetHeader([]string{"URL", "ERROR"})
		table.SetRowLine(false)
		for url, err := range errors {
			table.Append([]string{url, err.Error()})
		}
		table.Render()
		os.Exit(1)
	}
}
