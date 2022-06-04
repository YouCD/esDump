package es

import (
	"bytes"
	"github.com/tidwall/gjson"
	"log"
	"os"
)

func PipeImporter(dumpInfo *DumpInfo, ch chan string) {
	EsInit(dumpInfo)
	var buf bytes.Buffer
	i := 0
	batch := 0
	for {
		select {
		case str, ok := <-ch:
			if ok {

				m, d := MetaData(str)
				if len(m) > 0 {
					buf.Grow(len(m) + len(d))
					buf.Write(m)
					buf.Write(d)
				}
				i++
				if i == dumpInfo.Size {
					batch++
					res, err := Es.Bulk(bytes.NewReader(buf.Bytes()), Es.Bulk.WithIndex(dumpInfo.Index))
					if err != nil {
						log.Printf("err")
						os.Exit(1)
					}
					jsonStr := read(res.Body)
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
				res, err := Es.Bulk(bytes.NewReader(buf.Bytes()), Es.Bulk.WithIndex(dumpInfo.Index))
				if err != nil {
					log.Panic(err)
				}
				jsonStr := read(res.Body)
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
