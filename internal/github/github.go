package github

import (
	"context"

	ghrl "github.com/gofri/go-github-ratelimit/github_ratelimit"
	gh "github.com/google/go-github/v54/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

// OrgConfig contains fields used in releasegen's config.yaml file to configure
// its behaviour when generating reports about Github repositories.
type OrgConfig struct {
	Org          string   `mapstructure:"org"`
	Teams        []string `mapstructure:"teams"`
	IgnoredRepos []string `mapstructure:"ignores"`

	ghClient *gh.Client
}

// GithubClient returns either a new instance of the Github client, or a previously
// initialised client.
func (oc *OrgConfig) GithubClient() *gh.Client {
	if oc.ghClient == nil {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: viper.GetString("token")})
		tc := oauth2.NewClient(context.Background(), ts)

		rateLimiter, err := ghrl.NewRateLimitWaiterClient(tc.Transport)
		if err != nil {
			panic(err)
		}

		oc.ghClient = gh.NewClient(rateLimiter)
	}

	return oc.ghClient
}
