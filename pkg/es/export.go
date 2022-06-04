package es

import (
	"bytes"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

var (
	scrollID string
)

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
func (e *esDump) Exporter(ch chan string) (err error) {
	var res *esapi.Response
	switch {
	case e.query != "" && e.complex == false:
		res, err = e.Client.Search(
			e.Client.Search.WithIndex(e.index),
			e.Client.Search.WithSort("_doc"),
			e.Client.Search.WithSize(e.size),
			e.Client.Search.WithScroll(time.Minute),
			e.Client.Search.WithQuery(fmt.Sprintf(`%s`, e.query)),
		)
		if err != nil {
			log.Println(err)
			return err
		}
	case e.query != "" && e.complex:
		res, err = e.Client.Search(
			e.Client.Search.WithIndex(e.index),
			e.Client.Search.WithSort("_doc"),
			e.Client.Search.WithSize(e.size),
			e.Client.Search.WithScroll(time.Minute),
			e.Client.Search.WithBody(strings.NewReader(fmt.Sprintf(`%s`, e.query))),
		)
		if err != nil {
			return err
		}
	case e.query == "":
		res, err = e.Client.Search(
			e.Client.Search.WithIndex(e.index),
			e.Client.Search.WithSort("_doc"),
			e.Client.Search.WithSize(e.size),
			e.Client.Search.WithScroll(time.Minute),
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
			res, err := e.Client.Scroll(e.Client.Scroll.WithScrollID(scrollID), e.Client.Scroll.WithScroll(time.Minute))
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
