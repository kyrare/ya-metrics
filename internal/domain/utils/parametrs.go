package utils

import (
	"flag"
	"os"
)

// GetParameter функция хелпер для получения параметра, в первую очередь берет значение из параметров из командной строки
// во вторую очередь берет значение из окружения
func GetParameter(flagName string, envName string, configValue, defaultValue string, usage string) *string {
	if configValue != "" {
		defaultValue = configValue
	}

	v := flag.String(flagName, defaultValue, usage)

	if envName != "" {
		if envV := os.Getenv(envName); envV != "" {
			v = &envV
		}
	}

	return v
}
