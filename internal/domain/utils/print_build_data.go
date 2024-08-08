package utils

import "fmt"

func PrintBuildData(buildVersion string, buildDate string, buildCommit string) {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
