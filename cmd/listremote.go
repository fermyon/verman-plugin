package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var listRemoteCmd = &cobra.Command{
	Use:     "list-remote",
	Aliases: []string{"ls-remote"},
	Short:   "Lists all available versions of Spin",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := listRemote()
		if err != nil {
			fmt.Printf("Error while loading available Spin version\n%v\n", err)
			os.Exit(1)
		}

		return nil
	},
}

const (
	spinReleasesUrl   = "https://api.github.com/repos/fermyon/spin/releases"
	githubTokenEnvVar = "GH_TOKEN"
)

func listRemote() error {
	fmt.Fprintf(os.Stderr, "Fetching available Spin releases ...\n\n")
	releases, err := loadSpinReleases()
	if err != nil {
		return err
	}
	for _, release := range *releases {
		tagName := strings.Replace(release.TagName, "v", "", 1)
		fmt.Printf("%s\n", tagName)
	}
	return nil
}

func loadSpinReleases() (*[]spinRelease, error) {
	req, err := http.NewRequest("GET", spinReleasesUrl, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	token := os.Getenv(githubTokenEnvVar)
	if len(token) > 0 {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to load available Spin releases")
	}
	defer resp.Body.Close()

	// the value stored in env GH_TOKEN is a bad credential
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("Unauthorized: Bad credentials. Please check your GitHub token (%s).", githubTokenEnvVar)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %v", err)
	}
	var releases []spinRelease
	err = json.Unmarshal(body, &releases)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal JSON: %v", err)
	}
	return &releases, nil
}

type spinRelease struct {
	TagName string `json:"tag_name"`
}
