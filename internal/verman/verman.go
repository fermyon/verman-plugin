package verman

import (
	"fmt"
	"os"
	"strings"
)

const (
	spinVersionFileName = ".spin-version"
)

func GetDesiredVersionForSet(args []string) (string, error) {
	// explicitly provided ver has higher priority
	if len(args) > 0 {
		return args[0], nil
	}
	rcVersion := getVersionFromSpinVersionFile()

	// if rc version is empty, return an error
	if len(rcVersion) == 0 {
		return "", fmt.Errorf("you must indicate the version of Spin you wish to set")
	}
	return rcVersion, nil
}

func GetDesiredVersionsForGet(args []string) ([]string, error) {
	// explicitly provided ver has higher priority
	if len(args) > 0 {
		return args, nil
	}
	rcVersion := getVersionFromSpinVersionFile()

	// if rc version is empty, return an error
	if len(rcVersion) == 0 {
		return nil, fmt.Errorf("you must indicate the version of Spin you wish to set")
	}
	return []string{rcVersion}, nil
}

func getVersionFromSpinVersionFile() string {
	_, err := os.Stat(spinVersionFileName)
	if os.IsNotExist(err) {
		return ""
	}

	content, err := os.ReadFile(spinVersionFileName)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(content))
}
