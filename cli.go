package main

import (
	"context"
	"github.com/alecthomas/kingpin"
	"github.com/google/go-github/github"
)

type GithubReleaseApp struct {
	App *kingpin.Application

	ReleaseCmd *kingpin.CmdClause
	Release    *string
	Token      *string
	api_uri    *string
	Org        *string
	Repo       *string

	tag          *string
	name         *string
	IsDraft      *bool
	IsPreRelease *bool
	desc         *string

	Client *github.Client
	ctxt   context.Context
}

func (a *GithubReleaseApp) init() {
	a.App = kingpin.New("github-release", "create/update github release")

	a.ReleaseCmd = a.App.Command("release", "create or update a github release")

	a.Token = a.ReleaseCmd.Flag("github-token", "github token with release access").Required().Envar("GITHUB_TOKEN").String()
	a.api_uri = a.ReleaseCmd.Flag("github-api-uri", "Github API end point. For github Entreprise use "+
		"https://<server>/api/v3/").Envar("GITHUB_API").Default("https://api.github.com/").String()
	a.Org = a.ReleaseCmd.Flag("github-user", "User or Organization name").Short('u').Required().Envar("GITHUB_USER").String()
	a.Repo = a.ReleaseCmd.Flag("github-repo", "Github Repository name.").Short('r').Required().Envar("GITHUB_REPO").String()

	a.tag = a.ReleaseCmd.Arg("tag", "Tag name to use for release.").Required().String()
	a.name = a.ReleaseCmd.Flag("name", "Release name. By default, it uses the tag as name.").String()
	a.IsDraft = a.ReleaseCmd.Flag("draft", "To set release as Draft.").Bool()
	a.IsPreRelease = a.ReleaseCmd.Flag("pre-release", "to set release as pre-release").Bool()
	a.desc = a.ReleaseCmd.Flag("description", "Release description").String()
}
