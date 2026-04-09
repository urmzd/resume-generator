package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

func initUpdateCmd() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Self-update incipit to the latest release",
	RunE: func(cmd *cobra.Command, args []string) error {
		return selfUpdate(Version)
	},
}

func selfUpdate(current string) error {
	resp, err := http.Get("https://api.github.com/repos/urmzd/incipit/releases/latest")
	if err != nil {
		return fmt.Errorf("fetching release info: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return fmt.Errorf("parsing release info: %w", err)
	}

	tag := release.TagName
	if tag == "v"+current || tag == current {
		fmt.Fprintln(os.Stderr, "already up to date")
		return nil
	}

	goos := runtime.GOOS
	goarch := runtime.GOARCH
	assetName := fmt.Sprintf("incipit-%s-%s", goos, goarch)

	var downloadURL string
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}
	if downloadURL == "" {
		return fmt.Errorf("no asset named %q found in release %s", assetName, tag)
	}

	dlResp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("downloading update: %w", err)
	}
	defer func() { _ = dlResp.Body.Close() }()

	if dlResp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned status %d", dlResp.StatusCode)
	}

	tmp, err := os.CreateTemp("", "incipit-update-*")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	tmpName := tmp.Name()
	defer func() { _ = os.Remove(tmpName) }()

	if _, err := io.Copy(tmp, dlResp.Body); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("writing update: %w", err)
	}
	_ = tmp.Close()

	if err := os.Chmod(tmpName, 0755); err != nil {
		return fmt.Errorf("chmod: %w", err)
	}

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("finding current executable: %w", err)
	}

	if err := os.Rename(tmpName, exe); err != nil {
		return fmt.Errorf("replacing executable: %w", err)
	}

	fmt.Fprintf(os.Stderr, "updated: v%s -> %s\n", current, tag)
	return nil
}
