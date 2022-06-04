package es

import (
	"bytes"
	"encoding/json"
	"github.com/tidwall/gjson"
	"log"
	"os"
)

func MetaData(str string) (mate, data []byte) {
	esDoc := make(map[string]interface{})
	err := json.Unmarshal([]byte(str), &esDoc)
	if err != nil {
		log.Println("数据有误")
		os.Exit(1)
	}
	m := make(map[string]interface{})
	var d interface{}
	for k, v := range esDoc {
		if k == "_source" {
			d = v
		} else if k == "sort" || k == "_index" {
			continue
		} else {
			if v != nil {
				m[k] = v
			}
		}
	}
	var a struct {
		Index interface{} `json:"index"`
	}
	var aa interface{}
	marshal, _ := json.Marshal(m)
	_ = json.Unmarshal(marshal, &aa)
	a.Index = aa
	mb, _ := json.Marshal(a)
	da, _ := json.Marshal(d)

	mate = []byte(string(mb) + "\n")
	data = []byte(string(da) + "\n")
	return
}

func (e *esDump) PipeImporter(ch chan string) {
	var buf bytes.Buffer
	i := 0
	batch := 0
	for {
		select {
		case str, ok := <-ch:
			if ok {
				m, d := MetaData(str)
				if len(m) > 0 && len(d) > 0 {
					buf.Grow(len(m) + len(d))
					buf.Write(m)
					buf.Write(d)
				}
				i++
				if i == e.size {
					batch++
					res, err := e.Client.Bulk(bytes.NewReader(buf.Bytes()), e.Client.Bulk.WithIndex(e.index))
					if err != nil {
						log.Panic(err)
					}
					jsonStr := Read(res.Body)
					defer res.Body.Close()

					str := gjson.Get(jsonStr, "errors").String()
					if str != "false" {
						log.Printf("导入失败: err: %s", jsonStr)
						os.Exit(1)
					}
					buf.Reset()
					i = 0
					log.Printf("导入第%d批数据 ", batch)
				}
			}

			if !ok {
				batch++
				res, err := e.Client.Bulk(bytes.NewReader(buf.Bytes()), e.Client.Bulk.WithIndex(e.index))
				if err != nil {
					log.Panic(err)
				}
				jsonStr := Read(res.Body)
				defer res.Body.Close()
				str := gjson.Get(jsonStr, "errors").String()
				if str != "false" {
					log.Printf("导入失败: err: %s", jsonStr)
					os.Exit(1)
				}
				buf.Reset()
				i = 0
				log.Printf("导入第%d批数据 ", batch)
				return
			}

		}

	}
}
