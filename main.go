package main

import "github.com/MikeTangoEcho/marathon/cmd"

// https://goreleaser.com/cookbooks/using-main.version/
var (
	version string = "dev"
	commit  string = "none"
	date    string = "unknown"
)

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
