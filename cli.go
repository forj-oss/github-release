package main

import (
	"context"
	"github.com/alecthomas/kingpin"
	"github.com/google/go-github/github"
)

type ConnectStruct struct {
	Token        *string
	api_uri      *string

}

type RepoStruct struct {
	Org          *string
	Repo         *string
	tag          *string

}

type ManageReleaseCmd struct {
	Cmd          *kingpin.CmdClause
	ConnectStruct
	RepoStruct
	name         *string
	IsDraft      *bool
	IsPreRelease *bool
	desc         *string

}

type DeleteReleaseCmd struct {
	Cmd     *kingpin.CmdClause
	ConnectStruct
	RepoStruct
}


type GithubReleaseApp struct {
	App *kingpin.Application

	Manage ManageReleaseCmd
	Delete DeleteReleaseCmd

	Client *github.Client
	ctxt   context.Context

	release *github.RepositoryRelease
}

func (a *GithubReleaseApp) init() {
	a.App = kingpin.New("github-release", "create/update github release")

	a.setVersion()

	// ----------------- Manage
	a.Manage.Cmd = a.App.Command("release", "create or update a github release")

	a.Manage.Token = a.Manage.Cmd.Flag("github-token", "github token with release access").Required().Envar("GITHUB_TOKEN").String()
	a.Manage.api_uri = a.Manage.Cmd.Flag("github-api-uri", "Github API end point. For github Entreprise use "+
		"https://<server>/api/v3/").Envar("GITHUB_API").Default("https://api.github.com/").String()
	a.Manage.Org = a.Manage.Cmd.Flag("github-user", "User or Organization name").Short('u').Required().Envar("GITHUB_USER").String()
	a.Manage.Repo = a.Manage.Cmd.Flag("github-repo", "Github Repository name.").Short('r').Required().Envar("GITHUB_REPO").String()
	a.Manage.tag = a.Manage.Cmd.Arg("tag", "Tag name to use for release.").Required().String()

	a.Manage.name = a.Manage.Cmd.Flag("name", "Release name. By default, it uses the tag as name.").String()
	a.Manage.IsDraft = a.Manage.Cmd.Flag("draft", "To set release as Draft.").Bool()
	a.Manage.IsPreRelease = a.Manage.Cmd.Flag("pre-release", "to set release as pre-release").Bool()
	a.Manage.desc = a.Manage.Cmd.Flag("description", "Release description").String()

	// ----------------- Delete
	a.Delete.Cmd = a.App.Command("delete", "delete a github release")

	a.Delete.Token = a.Delete.Cmd.Flag("github-token", "github token with release access").Required().Envar("GITHUB_TOKEN").String()
	a.Delete.api_uri = a.Delete.Cmd.Flag("github-api-uri", "Github API end point. For github Entreprise use "+
		"https://<server>/api/v3/").Envar("GITHUB_API").Default("https://api.github.com/").String()
	a.Delete.Org = a.Delete.Cmd.Flag("github-user", "User or Organization name").Short('u').Required().Envar("GITHUB_USER").String()
	a.Delete.Repo = a.Delete.Cmd.Flag("github-repo", "Github Repository name.").Short('r').Required().Envar("GITHUB_REPO").String()
	a.Delete.tag = a.Delete.Cmd.Arg("tag", "Tag name of the release to delete.").Required().String()

}
