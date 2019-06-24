package main

import (
	"context"
	"strings"

	"github.com/github/hub/github"
	"github.com/shurcooL/githubv4"
	"go.tmthrgd.dev/gomodpriv/internal/modfile"
	"golang.org/x/oauth2"
)

func githubPrivateRepos(ctx context.Context) ([]string, error) {
	host, err := github.CurrentConfig().DefaultHost()
	if err != nil {
		return nil, err
	}

	token := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: host.AccessToken,
	})
	c := githubv4.NewClient(oauth2.NewClient(ctx, token))

	var modules []string

	variables := map[string]interface{}{
		"cursor": (*githubv4.String)(nil),
	}
	for {
		var query struct {
			Viewer struct {
				Repositories struct {
					PageInfo struct {
						EndCursor   githubv4.String
						HasNextPage githubv4.Boolean
					}
					Nodes []struct {
						URL   githubv4.URI
						GoMod struct {
							Blob struct {
								Text githubv4.String
							} `graphql:"... on Blob"`
						} `graphql:"gomod: object(expression: \"HEAD:go.mod\")"`
					}
				} `graphql:"repositories(first: 100, after: $cursor, privacy: PRIVATE, affiliations: OWNER, orderBy: {field: NAME, direction: ASC})"`
			}
		}
		if err := c.Query(ctx, &query, variables); err != nil {
			return nil, err
		}

		for _, node := range query.Viewer.Repositories.Nodes {
			repoPath := strings.TrimPrefix(node.URL.String(), "https://")
			modules = append(modules, repoPath)

			module := modfile.ModulePath([]byte(node.GoMod.Blob.Text))
			if module != "" && module != repoPath {
				modules = append(modules, module)
			}
		}

		if !query.Viewer.Repositories.PageInfo.HasNextPage {
			return modules, nil
		}

		variables["cursor"] = githubv4.NewString(query.Viewer.Repositories.PageInfo.EndCursor)
	}
}
