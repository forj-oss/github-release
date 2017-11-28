package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"net/url"
)

func (a *GithubReleaseApp) release() {
	if err := a.github_connect(); err != nil {
		fmt.Printf("Connection issue: %s", err)
	}
	fmt.Print("Connected.\n")

	if tags, resp, err := a.Client.Repositories.ListTags(a.ctxt, *a.Org, *a.Repo, nil); err != nil {
		fmt.Printf("Tags not found in %s/%s. %s\n", *a.Org, *a.Repo, err)
		return
	} else if resp.StatusCode != 200 {
		fmt.Printf("Tags not found in %s/%s. HTTP code %d", *a.Org, *a.Repo, resp.StatusCode)
	} else {
		fmt.Printf("Tags found on %s/%s: %d\n", *a.Org, *a.Repo, len(tags))
		for _, tag := range tags {
			fmt.Printf("- tag '%s'\n", *tag.Name)
		}
	}
}

func (a *GithubReleaseApp) github_connect() (err error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *a.Token})
	a.ctxt = context.Background()
	tc := oauth2.NewClient(a.ctxt, ts)

	a.Client = github.NewClient(tc)

	a.Client.BaseURL, err = url.Parse(*a.api_uri)
	if err != nil {
		return
	}

	fmt.Printf("Github API URL used : %s\n", a.Client.BaseURL)

	if user, _, err := a.Client.Users.Get(a.ctxt, ""); err != nil {
		return fmt.Errorf("Unable to get the owner of the token given. %s", err)
	} else {
		fmt.Printf("Connection successful. Token given by user '%s'\n", *user.Login)
	}

	return
}
