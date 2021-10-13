package esHandler

import (
	"bytes"
	"encoding/json"
	"github.com/tidwall/gjson"
	"log"
)

type EsIndex struct {
	Index string `json:"_index"`
	Type  string `json:"_type"`
	Id    string `json:"_id"`
}

func MetaData(str string, dumpInfo *DumpInfo) (mate, data []byte) {

	id := gjson.GetMany(str, "_index", "_type", "_id", "_score")[2].String()
	typeTemp := gjson.GetMany(str, "_index", "_type", "_id", "_score")[1].String()

	if id != "" && typeTemp != "" {
		var esIndex = EsIndex{
			Index: dumpInfo.Index,
			Type:  gjson.GetMany(str, "_index", "_type", "_id", "_score")[1].String(),
			Id:    gjson.GetMany(str, "_index", "_type", "_id", "_score")[2].String(),
		}

		dataBytes, err := json.Marshal(esIndex)
		if err != nil {
			log.Panic(err)
		}

		mate = []byte(`{"index":` + string(dataBytes) + `}` + "\n")
		data = []byte(gjson.Get(str, "_source").String() + "\n")
		return
	}
	return
}

func PipeImporter(dumpInfo *DumpInfo, ch chan string) {
	EsInit(*dumpInfo)
	var buf bytes.Buffer
	i := 0
	for {
		select {
		case str, ok := <-ch:
			m, d := MetaData(str, dumpInfo)
			if len(m) > 0 {
				buf.Grow(len(m) + len(d))
				buf.Write(m)
				buf.Write(d)
			}
			i++
			if i == dumpInfo.Size {
				res, err := Es.Bulk(bytes.NewReader(buf.Bytes()), Es.Bulk.WithIndex(dumpInfo.Index))
				if err != nil {
					log.Panic(err)
				}
				jsonStr := Read(res.Body)
				defer res.Body.Close()
				str := gjson.Get(jsonStr, "errors").String()
				if str != "false" {
					log.Printf("导入失败...")
				}
				buf.Reset()
				i = 0
			}

			if !ok {
				res, err := Es.Bulk(bytes.NewReader(buf.Bytes()), Es.Bulk.WithIndex(dumpInfo.Index))
				if err != nil {
					log.Panic(err)
				}
				jsonStr := Read(res.Body)
				defer res.Body.Close()
				str := gjson.Get(jsonStr, "errors").String()
				if str != "false" {
					log.Printf("导入失败...")
				}
				buf.Reset()
				i = 0
				return
			}

		}

	}
}
