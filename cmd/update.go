package cmd

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
)

var (
	path             string
	Force            bool
	GitHubReleaseUrl = "https://api.github.com/repos/youcd/esDump/releases/latest"
)

type ReleaseVersion struct {
	TagName     string `json:"tag_nmae"`
	SwName      string `json:"sw_name"`
	DownloadUrl string `json:"download_url"`
}

func init() {
	updateCmd.Flags().BoolVarP(&Force, "force", "f", false, "force updating.")
}
func GetRelease(OS string) (v ReleaseVersion) {

	resp, err := http.Get(GitHubReleaseUrl)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	vList := make([]ReleaseVersion, 0)
	count := gjson.Get(string(bytes), "assets.#").Int()

	for i := 0; i < int(count); i++ {
		v.TagName = gjson.Get(string(bytes), "tag_name").Str
		v.SwName = gjson.Get(string(bytes), fmt.Sprintf("assets.%d.name", i)).Str
		v.DownloadUrl = gjson.Get(string(bytes), fmt.Sprintf("assets.%d.browser_download_url", i)).Str
		vList = append(vList, v)
	}
	switch OS {
	case "linux":
		for _, v := range vList {
			if strings.Contains(v.SwName, "linux") {
				return v
			}
		}
	case "darwin":
		for _, v := range vList {
			if strings.Contains(v.SwName, "darwin") {
				return v
			}
		}
	case "windows":
		for _, v := range vList {
			if strings.Contains(v.SwName, "windows") {
				return v
			}
		}
	}
	return v
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update the WorkReport server",
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		os.Rename(path+".tmp", path)
	},
	Run: func(cmd *cobra.Command, args []string) {
		//系统类型
		OS := runtime.GOOS
		v := GetRelease(OS)
		if Force {
			path, _ = os.Executable()
			DownloadFileProgress(v.DownloadUrl, path+".tmp")
			return
		} else if Version != v.TagName {
			path, _ = os.Executable()
			DownloadFileProgress(v.DownloadUrl, path+".tmp")

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

func DownloadFileProgress(url, filename string) {
	r, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer func() { _ = r.Body.Close() }()
	f, err := os.Create(filename)
	// 更改权限
	err = f.Chmod(0775)
	if err != nil {
		panic(err)
	}
	defer func() { _ = f.Close() }()
	bar := progressbar.DefaultBytes(
		r.ContentLength,
		"下载中",
	)
	io.Copy(io.MultiWriter(f, bar), r.Body)
}
