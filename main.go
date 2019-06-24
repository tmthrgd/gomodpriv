package main // import "go.tmthrgd.dev/gomodpriv"

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const help = `
%[1]s is a tool to automatically set the GOPRIVATE environment variable.

The GOPRIVATE environment variable is set to a list of all private repositories
owned by the current GitHub user. It reads go.mod files in the root directory,
to handle any custom import paths.

When run for the first time you will be prompted for your GitHub username and
password. Your password is not stored, but instead used to request an Oauth2
token.

By default the list of private repositories are cached for faster terminal
start-ups. The cache can be refreshed by running %[1]s directly, or
adding -cache=false to the .bashrc line below.

Add the following to your .bashrc file to set the GOPRIVATE environment
variable when a terminal is opened:
  export GOPRIVATE=$(%[1]s -env)

If you want the list of private repositories to be refreshed periodically,
ensure %[1]s is run at start-up or upon login. Otherwise it can be run
manually when creating new private repositories.

Usage of %[1]s:
`

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			strings.TrimLeft(help, "\n"), os.Args[0])
		flag.PrintDefaults()
	}

	env := flag.Bool("env", false, "print the list of repositories for the GOPRIVATE environment variable")
	cache := flag.Bool("cache", true, "use a cache of private repos; refresh cache by running without arguments")
	flag.Parse()

	if !*env && !*cache {
		flag.Usage()
		os.Exit(1)
	}

	var cacheFilePath string
	if *cache {
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			log.Fatal(err)
		}
		cacheFilePath = filepath.Join(cacheDir, "gomodpriv.env")

		if *env {
			cached, err := ioutil.ReadFile(cacheFilePath)
			switch {
			case os.IsNotExist(err):
			default:
				log.Fatal(err)
			case err == nil:
				os.Stdout.Write(cached)
				return
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	repos, err := githubPrivateRepos(ctx)
	if err != nil {
		log.Fatal(err)
	}

	sort.Strings(repos)
	reposStr := strings.Join(repos, ",")

	if *cache {
		if err := ioutil.WriteFile(cacheFilePath, []byte(reposStr), 0600); err != nil {
			log.Fatal(err)
		}
	}

	if *env {
		os.Stdout.WriteString(reposStr)
	}
}
