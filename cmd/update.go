package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	path             string
	Force            bool
	GitHubReleaseUrl = "https://api.github.com/repos/youcd/esDump/releases"
)

type author struct {
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
	SiteAdmin         bool   `json:"site_admin"`
}
type uploader struct {
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
	SiteAdmin         bool   `json:"site_admin"`
}

type assets struct {
	URL                string    `json:"url"`
	ID                 int       `json:"id"`
	NodeID             string    `json:"node_id"`
	Name               string    `json:"name"`
	Label              string    `json:"label"`
	Uploader           uploader  `json:"uploader"`
	ContentType        string    `json:"content_type"`
	State              string    `json:"state"`
	Size               int       `json:"size"`
	DownloadCount      int       `json:"download_count"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	BrowserDownloadURL string    `json:"browser_download_url"`
}
type GithubRelease struct {
	//URL             string      `json:"url"`
	//AssetsURL       string      `json:"assets_url"`
	//UploadURL       string      `json:"upload_url"`
	//HTMLURL         string      `json:"html_url"`
	//ID              int         `json:"id"`
	//Author          author      `json:"author"`
	//NodeID          string      `json:"node_id"`
	TagName string `json:"tag_name"`
	//TargetCommitish string      `json:"target_commitish"`
	//Name            interface{} `json:"name"`
	//Draft           bool        `json:"draft"`
	//Prerelease      bool        `json:"prerelease"`
	//CreatedAt       time.Time   `json:"created_at"`
	//PublishedAt     time.Time   `json:"published_at"`
	Assets []assets `json:"assets"`
	//TarballURL      string      `json:"tarball_url"`
	//ZipballURL      string      `json:"zipball_url"`
	//Body            interface{} `json:"body"`
}

type ReleaseVersion struct {
	TagName     string `json:"tage_name"`
	DownloadUrl string `json:"download_url"`
}

func init() {
	updateCmd.Flags().BoolVarP(&Force, "force", "f", false, "force updating.")
}
func GetRelease(OS string) (v ReleaseVersion) {

	resp, err := http.Get(GitHubReleaseUrl)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	var githubReleases []GithubRelease
	err = json.Unmarshal(bytes, &githubReleases)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	var ReleaseVersions []ReleaseVersion

	for _, i := range githubReleases {
		for _, a := range i.Assets {
			var r ReleaseVersion
			r.TagName = i.TagName
			r.DownloadUrl = a.BrowserDownloadURL
			ReleaseVersions = append(ReleaseVersions, r)
		}
	}
	//count := gjson.Get(string(bytes), "assets.#").Int()

	//for i := 0; i < int(count); i++ {
	//	v.TagName = gjson.Get(string(bytes), "tag_name").Str
	//	v.SwName = gjson.Get(string(bytes), fmt.Sprintf("assets.%d.name", i)).Str
	//	v.DownloadUrl = gjson.Get(string(bytes), fmt.Sprintf("assets.%d.browser_download_url", i)).Str
	//	vList = append(vList, v)
	//}
	switch OS {
	case "linux":
		return getReleaseVersion(ReleaseVersions, "linux")
	case "darwin":
		return getReleaseVersion(ReleaseVersions, "darwin")
	case "windows":
		return getReleaseVersion(ReleaseVersions, "windows")
	}
	return v
}

func getReleaseVersion(ReleaseVersions []ReleaseVersion, os string) (v ReleaseVersion) {
	var list []string
	listMap := make(map[string]string)
	for _, v := range ReleaseVersions {
		if strings.Contains(v.DownloadUrl, os) {
			list = append(list, v.TagName)
			listMap[v.TagName] = v.DownloadUrl
		}
	}
	version := chooseVersion(list)
	v.DownloadUrl = listMap[version]
	v.TagName = version
	return
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: fmt.Sprintf("update the %s server", name),
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		os.Rename(path+".tmp", path)
	},
	Run: func(cmd *cobra.Command, args []string) {
		//系统类型
		OS := runtime.GOOS
		v := GetRelease(OS)
		if Force {
			path, _ = os.Executable()
			file := filename(v.DownloadUrl)
			DownloadFileProgress(v.DownloadUrl, file)
			return
		} else if Version != v.TagName {
			file := filename(v.DownloadUrl)
			DownloadFileProgress(v.DownloadUrl, file)

		} else {
			log.Println(fmt.Sprintf("version: %s. The version is latest version.", Version))
			return
		}

	},
}

type Reader struct {
	io.Reader
	Total   int64
	Current int64
}

func filename(url string) string {
	split := strings.Split(url, "/")

	return split[len(split)-1]

}
func DownloadFileProgress(url, filename string) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer func() { _ = resp.Body.Close() }()
	f, err := os.Create(filename)
	// 更改权限
	err = f.Chmod(0644)
	if err != nil {
		panic(err)
	}
	defer func() { _ = f.Close() }()

	reader := io.LimitReader(resp.Body, resp.ContentLength)

	p := mpb.New(
		//mpb.WithWidth(60),
		mpb.WithRefreshRate(180 * time.Millisecond),
	)

	bar := p.New(resp.ContentLength,
		mpb.BarStyle().Rbound("|"),
		mpb.PrependDecorators(
			decor.CountersKibiByte("% .2f / % .2f"),
		),
		mpb.AppendDecorators(
			decor.EwmaETA(decor.ET_STYLE_GO, 90),
			decor.Name(" ] "),
			decor.EwmaSpeed(decor.UnitKiB, "% .2f", 60),
		),
	)

	// create proxy reader
	proxyReader := bar.ProxyReader(reader)
	defer proxyReader.Close()

	// copy from proxyReader, ignoring errors
	_, _ = io.Copy(f, proxyReader)

	p.Wait()
	log.Printf("文件 %s 依稀在完毕.", filename)
}

func tagToInt(tag string) int {
	s := strings.Split(strings.ToLower(tag), "v")
	if len(s) < 2 {
		return 0
	}
	strFloat := strings.Split(strings.ToLower(tag), "v")[1]
	no := 0
	nos := strings.Split(strFloat, ".")
	if len(nos) != 3 {
		return -1
	}
	for i, n := range nos {
		j, err := strconv.Atoi(n)
		if err != nil || j >= 100 || j < 0 {
			return -1
		}
		no += j * int(pow(100, (2-i)))
	}
	return no

}

func pow(x int, n int) int {
	if x == 0 || n < 0 {
		return 0
	}
	if n == 0 {
		return 1
	}
	result := 1
	for i := 0; i < n; i++ {
		result *= x
	}
	return result
}
func chooseVersion(list []string) (ns string) {
	prompt := promptui.Select{
		Label: "请选择要下载的版本",
		Items: list,
		Size:  10,
	}
	_, ns, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	return ns
}
