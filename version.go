package main

import "fmt"

var (
	build_branch string
	build_commit string
	build_date   string
	build_tag    string
)

const (
	PRERELEASE = false
	VERSION    = "1.0.3"
	APP        = "github_release"
)

func (a *GithubReleaseApp) setVersion() {
	var version string
	if PRERELEASE {
		version = APP + " pre-release V" + VERSION
	} else {
		version = APP + " V" + VERSION
	}

	if build_branch != "master" {
		version += fmt.Sprintf(" branch %s", build_branch)
	}
	if build_tag == "false" {
		version += fmt.Sprintf(" patched - %s - %s", build_date, build_commit)
	}

	version += " - Github Repo: https://github.com/forj-oss/github-release"
	cliApp.App.Version(version)
}
