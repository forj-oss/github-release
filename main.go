package main

import (
	"github.com/alecthomas/kingpin"
	"os"
	"forjj-modules/trace"
)

var cliApp GithubReleaseApp

func main() {
	cliApp.init()

	switch kingpin.MustParse(cliApp.App.Parse(os.Args[1:])) {
	case "release":
		cliApp.do_release()
	case "delete":
		cliApp.do_delete()
	default:
		kingpin.Usage()
	}
}

func (a *GithubReleaseApp) do_release() {
	if err := a.github_connect(a.Manage.ConnectStruct); err != nil {
		gotrace.Error("Connection issue: %s\n", err)
		os.Exit(1)
	}
	gotrace.Info("Connected.\n")

	if _, err := a.search_tag(a.Manage.RepoStruct) ; err != nil {
		gotrace.Error("Unable to create/update a release. %s\n", err)
		os.Exit(1)
	}

	if err := a.manage_release() ; err != nil {
		gotrace.Error("Unable to create/update a release. %s\n", err)
		os.Exit(1)
	}

	gotrace.Info("Done.\n")
}

func (a *GithubReleaseApp) do_delete() {
	if err := a.github_connect(a.Delete.ConnectStruct); err != nil {
		gotrace.Error("Connection issue: %s\n", err)
		os.Exit(1)
	}
	gotrace.Info("Connected.\n")


	if found, err := a.search_release(a.Delete.RepoStruct) ; err != nil {
		gotrace.Error("Unable to delete a release. %s\n", err)
		os.Exit(1)
	} else {
		if ! found {
			gotrace.Warning("No release found for tag '%s'", *a.Delete.tag)
		}
	}

	if err := a.delete_release() ; err != nil {
		gotrace.Error("Unable to create/update a release. %s\n", err)
		os.Exit(1)
	}

	gotrace.Info("Done.\n")
}
