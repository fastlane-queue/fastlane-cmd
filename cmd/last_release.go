package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/blang/semver"
)

// ReleaseAsset represents an asset in a release
type ReleaseAsset struct {
	Name string
	Size string
	URL  string
}

// Release represents a release in Fastlane
type Release struct {
	Version    semver.Version
	URL        string
	Name       string
	Body       string
	ZipballURL string
	TarballURL string
	Assets     []*ReleaseAsset
}

func getLastRelease() *Release {
	url := "https://api.github.com/repos/fastlane-queue/fastlane-cmd/releases"
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	resp, _ := netClient.Get(url)
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Println("Failed to close body.")
		}
	}()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(au.Red("Failed to verify if there's a new version of fastlane-cmd."), err)
			return nil
		}
		var result []interface{}
		err = json.Unmarshal(bodyBytes, &result)
		if err != nil {
			fmt.Println(au.Red("Failed to verify if there's a new version of fastlane-cmd."), err)
			return nil
		}
		if len(result) == 0 {
			return nil
		}
		releaseData := result[0].(map[string]interface{})

		release := &Release{}
		v, err := semver.Make(releaseData["tag_name"].(string))
		release.Version = v
		release.URL = releaseData["url"].(string)
		release.Name = releaseData["name"].(string)
		release.Body = releaseData["body"].(string)
		release.TarballURL = releaseData["tarball_url"].(string)
		release.ZipballURL = releaseData["zipball_url"].(string)

		release.Assets = []*ReleaseAsset{}
		for _, assetData := range releaseData["assets"].([]interface{}) {
			asset := assetData.(map[string]interface{})
			release.Assets = append(
				release.Assets,
				&ReleaseAsset{
					Name: asset["name"].(string),
					Size: fmt.Sprintf("%d", asset["size"].(float64)),
					URL:  asset["browser_download_url"].(string),
				},
			)
		}
		return release
	}

	fmt.Println(
		au.Sprintf(
			au.Red("Failed to verify if there's a new version of fastlane-cmd (Status Code: %d)."), resp.StatusCode,
		),
	)
	return nil
}
