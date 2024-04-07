package utils

import (
	"flag"
	"os"
)

func GetParameter(flagName string, envName string, defaultValue string, usage string) *string {
	v := flag.String(flagName, defaultValue, usage)

	if envName != "" {
		if envV := os.Getenv(envName); envV != "" {
			v = &envV
		}
	}

	return v
}
