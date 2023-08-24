package releases

import (
	"context"
	"log"

	"github.com/gofri/go-github-ratelimit/github_ratelimit"
	"github.com/google/go-github/v54/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

var ghClient *github.Client

func githubClient() *github.Client {
	if ghClient == nil {
		log.Println("creating new Github client")
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: viper.Get("token").(string)},
		)
		tc := oauth2.NewClient(ctx, ts)

		rateLimiter, err := github_ratelimit.NewRateLimitWaiterClient(tc.Transport)
		if err != nil {
			panic(err)
		}

		ghClient = github.NewClient(rateLimiter)
	}
	return ghClient
}

func parseGithubApiError(err error) string {
	if _, ok := err.(*github.RateLimitError); ok {
		return "rate limit exceeded"
	}
	if _, ok := err.(*github.AbuseRateLimitError); ok {
		return "secondary rate limit exceeded"
	}
	return err.Error()
}
