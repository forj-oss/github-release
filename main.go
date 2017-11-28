package main

import (
	"github.com/alecthomas/kingpin"
	"os"
)

var cliApp GithubReleaseApp

func main() {
	cliApp.init()

	switch kingpin.MustParse(cliApp.App.Parse(os.Args[1:])) {
	case "release":
		cliApp.release()
	default:
		kingpin.Usage()
	}
}
