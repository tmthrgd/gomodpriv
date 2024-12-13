package main // import "go.tmthrgd.dev/gomodpriv"

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

const help = `
%[1]s is a tool to automatically set the GOPRIVATE environment variable.

The GOPRIVATE environment variable is set to a list of all private repositories
owned by the current GitHub user. It reads go.mod files in the root directory,
to handle any custom import paths.

The environment variable is set for go commands with the go env command.

When run for the first time you will be prompted for your GitHub username and
password. Your password is not stored, but instead used to request an Oauth2
token.

If you want the list of private repositories to be refreshed periodically,
ensure %[1]s is run at start-up or upon login. Otherwise it can be run
manually when creating new private repositories.

Usage of %[1]s:
`

func main() {
	if err := main1(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main1() error {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			strings.TrimLeft(help, "\n"), os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	repos, err := githubPrivateRepos(ctx)
	if err != nil {
		return err
	}
	sort.Strings(repos)

	cmd := exec.Command("go", "env", "-w", "GOPRIVATE="+strings.Join(repos, ","))
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd.Run()
}
