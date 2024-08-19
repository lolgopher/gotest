package main

import "github.com/coreos/go-semver/semver"

var versions = []string{
	"0.19.0",
	"0.23.0",
	"0.25.0",
	"0.34.0",
	"0.35.0",
	"0.36.0",
	// New Version Here.
}

func main() {
	var newVersion = "0.36.0"

	lastVersion := semver.New(LastSpecVersion())
	println("version check: ", !semver.New(newVersion).LessThan(*lastVersion))
}

func LastSpecVersion() string {
	return versions[len(versions)-1]
}
