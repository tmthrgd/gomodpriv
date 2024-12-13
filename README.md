# gomodpriv

**gomodpriv** is a tool to automatically manage the `GOPRIVATE` environment
variable used by go1.13+.

## Installation

Installation is simple and no different to any Go tool. The only requirement is
a working [Go](https://golang.org/) install.

```
go get go.tmthrgd.dev/gomodpriv
```

## Usage

The GOPRIVATE environment variable is set to a list of all private repositories
owned by the current GitHub user. It reads `go.mod` files in the root
directory, to handle any custom import paths.

The environment variable is set for go commands with the `go env` command.

When run for the first time you will be prompted for your GitHub username and
password. Your password is not stored, but instead used to request an Oauth2
token.

If you want the list of private repositories to be refreshed periodically,
ensure **gomodpriv** is run at start-up or upon login. Otherwise it can be run
manually when creating new private repositories.

## License

[BSD 3-Clause License](LICENSE)

## Note

**gomodpriv** (mis)uses [github/hub](https://github.com/github/hub) to handle
authentication. If hub is already installed, you may not be prompted for your
username and password. To have this work with GitHub enterprise follow the
instructions found here:
[hub.github.com/hub.1.html#github-enterprise](https://hub.github.com/hub.1.html#github-enterprise).
