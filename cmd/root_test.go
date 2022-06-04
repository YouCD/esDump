package cmd

import (
	"fmt"
	"os"
	"testing"
)

func Test_parFlag(t *testing.T) {

	url := "https://elastic:123.com@127.0.0.1:9200/all_singer"
	dumpInfo, err := parFlag(url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
	if dumpInfo == nil {
		fmt.Println("参数解析失败，不支持的endpoint")
		os.Exit(1)
		return
	}
}
