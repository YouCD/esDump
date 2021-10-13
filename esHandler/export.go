package esHandler

import (
	"bytes"
	"fmt"
	"github.com/elastic/go-elasticsearch/v6/esapi"
	"strings"

	"github.com/tidwall/gjson"
	"io"
	"log"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v6"
)

var (
	scrollID string
	Es       *elasticsearch.Client
)

type DumpInfo struct {
	User     string
	Password string
	Host     string
	Index    string
	Size     int
	Query    string
	ExistsFilter string
}

func EsInit(dumpInfo DumpInfo) {
	var err error
	config := elasticsearch.Config{}
	if dumpInfo.User != "" && dumpInfo.Password != "" {
		config.Addresses = []string{dumpInfo.Host}
		config.Username = dumpInfo.User
		config.Password = dumpInfo.Password
		Es, err = elasticsearch.NewClient(config)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	} else if dumpInfo.User == "" && dumpInfo.Password == "" {
		config.Addresses = []string{dumpInfo.Host}
		Es, err = elasticsearch.NewClient(config)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}
}

func Read(r io.Reader) string {
	var b bytes.Buffer
	_, err := b.ReadFrom(r)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return b.String()
}

// Exporter 从elasticsearch中导出索引
func Exporter(dumpInfo *DumpInfo, ch chan string) (err error) {

	EsInit(*dumpInfo)
	var res *esapi.Response
	switch  {
	case dumpInfo.Query != ""&&dumpInfo.ExistsFilter=="":
		res, err = Es.Search(
			Es.Search.WithIndex(dumpInfo.Index),
			Es.Search.WithSort("_doc"),
			Es.Search.WithSize(dumpInfo.Size),
			Es.Search.WithScroll(time.Minute),
			Es.Search.WithQuery(dumpInfo.Query),
		)
		if err != nil {
			return err
		}
	case dumpInfo.Query == ""&&dumpInfo.ExistsFilter=="":
		res, err = Es.Search(
			Es.Search.WithIndex(dumpInfo.Index),
			Es.Search.WithSort("_doc"),
			Es.Search.WithSize(dumpInfo.Size),
			Es.Search.WithScroll(time.Minute),
		)
		if err != nil {
			return err
		}
	case dumpInfo.ExistsFilter!="":
		res, err = Es.Search(
			Es.Search.WithIndex(dumpInfo.Index),
			Es.Search.WithSort("_doc"),
			Es.Search.WithSize(dumpInfo.Size),
			Es.Search.WithScroll(time.Minute),
			Es.Search.WithBody(strings.NewReader(fmt.Sprintf(`{"query":{"exists":{"field":"%s"}}}`,dumpInfo.ExistsFilter))),
		)
		if err != nil {
			return err
		}
	}



	//先做查询初始化并获取scrollID
	json := Read(res.Body)
	defer res.Body.Close()
	//将scrollID保存至全局变量中
	scrollID = gjson.Get(json, "_scroll_id").String()
	res.Body.Close()

	// 从response中提取scrollID
	scrollID = gjson.Get(json, "_scroll_id").String()

	// 提取搜索结果
	hits := gjson.Get(json, "hits.hits")

	go func() {
		for _, v := range hits.Array() {
			ch <- v.String()
		}
		for {
			//执行滚动请求并传递scrollID和滚动持续时间
			res, err := Es.Scroll(Es.Scroll.WithScrollID(scrollID), Es.Scroll.WithScroll(time.Minute))
			if err != nil {
				log.Fatalf("Error: %s", err)
			}
			if res.IsError() {
				log.Fatalf("Error response: %s", res)
			}

			json := Read(res.Body)
			res.Body.Close()

			// 从response中提取scrollID
			scrollID = gjson.Get(json, "_scroll_id").String()

			// 提取搜索结果
			hits := gjson.Get(json, "hits.hits")

			//没有结果时跳出循环
			if len(hits.Array()) < 1 {
				break
			} else {
				//遍历每次查询结果得到一条数据并推送至管道
				for _, v := range hits.Array() {
					ch <- v.String()
				}
			}
		}
		close(ch)
	}()
	return nil
}
