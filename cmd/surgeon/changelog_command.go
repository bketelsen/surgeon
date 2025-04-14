/*
Copyright © 2025 Brian Ketelsen <bketelsen@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/bketelsen/toolbox/cobra"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/spf13/viper"
)

// NewChangelogCommand creates a new command to generate a changelog for the project
func NewChangelogCommand(config *viper.Viper) *cobra.Command {
	// gendocsCmd represents the gendocs command
	changelogCmd := &cobra.Command{
		Use:    "changelog",
		Hidden: true,
		Short:  "Generates changelog for the project",
		Long: `Generates changelog for the command using the github api.
The changelog is written to the output directory in markdown format.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			lipgloss.DefaultRenderer().SetColorProfile(termenv.Ascii)

			o := config.GetString("changelog.output")
			cmd.Root().DisableAutoGenTag = true
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			target := filepath.Join(wd, o)
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
			r, err := http.Get("https://api.github.com/repos/bketelsen/surgeon/releases")
			if err != nil {
				return err
			}
			defer r.Body.Close()
			if r.StatusCode != http.StatusOK {
				return err
			}
			var changelogs []Changelog
			if err := json.NewDecoder(r.Body).Decode(&changelogs); err != nil {
				return err
			}

			// Create a new file in the target directory
			file, err := os.Create(filepath.Join(target, "changelog.md"))
			if err != nil {
				return err
			}
			defer file.Close()

			// Write the header to the file
			header := `# Changelog

`
			// Write the header to the file
			_, err = file.WriteString(header)
			if err != nil {
				return err
			}
			n := config.GetInt("changelog.last-n")
			fmt.Println("Last N:", n)
			fmt.Println("Changelogs:", len(changelogs))
			if n > len(changelogs) {
				n = len(changelogs)
			}
			for _, changelog := range changelogs[0:n] {
				_, _ = fmt.Fprintf(file, "## %s\n", changelog.TagName)
				_, _ = fmt.Fprintf(file, "### Released %s\n", changelog.PublishedAt.Format("2006-01-02"))
				_, _ = fmt.Fprintf(file, "### %s\n", changelog.Body)
			}
			// Print the file path
			cmd.Println("Changelog written to:", filepath.Join(target, "changelog.md"))
			return nil
		},
	}

	//	gendocsCmd.Flags().StringP("basepath", "b", "inventory", "Base path for the documentation (default is /inventory)")

	changelogCmd.PreRunE = func(cmd *cobra.Command, _ []string) error {
		_ = config.BindPFlag("changelog.output", cmd.Flags().Lookup("output"))
		_ = config.BindPFlag("changelog.last-n", cmd.Flags().Lookup("last-n"))

		return nil
	}
	// Define cobra flags, the default value has the lowest (least significant) precedence
	changelogCmd.Flags().StringP("output", "o", "docs", "Output directory for the documentation (default is docs)")
	changelogCmd.Flags().IntP("last-n", "n", 4, "Number of recent changelogs to include (default is 4)")

	return changelogCmd
}

type Changelog struct {
	URL       string `json:"url"`
	AssetsURL string `json:"assets_url"`
	UploadURL string `json:"upload_url"`
	HTMLURL   string `json:"html_url"`
	ID        int    `json:"id"`
	Author    struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		UserViewType      string `json:"user_view_type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	NodeID          string    `json:"node_id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Draft           bool      `json:"draft"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
	Assets          []struct {
		URL      string `json:"url"`
		ID       int    `json:"id"`
		NodeID   string `json:"node_id"`
		Name     string `json:"name"`
		Label    string `json:"label"`
		Uploader struct {
			Login             string `json:"login"`
			ID                int    `json:"id"`
			NodeID            string `json:"node_id"`
			AvatarURL         string `json:"avatar_url"`
			GravatarID        string `json:"gravatar_id"`
			URL               string `json:"url"`
			HTMLURL           string `json:"html_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			OrganizationsURL  string `json:"organizations_url"`
			ReposURL          string `json:"repos_url"`
			EventsURL         string `json:"events_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			Type              string `json:"type"`
			UserViewType      string `json:"user_view_type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"uploader"`
		ContentType        string    `json:"content_type"`
		State              string    `json:"state"`
		Size               int       `json:"size"`
		DownloadCount      int       `json:"download_count"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		BrowserDownloadURL string    `json:"browser_download_url"`
	} `json:"assets"`
	TarballURL string `json:"tarball_url"`
	ZipballURL string `json:"zipball_url"`
	Body       string `json:"body"`
}
