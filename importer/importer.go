package importer

import (
	"bytes"
	"encoding/json"
	"github.com/tidwall/gjson"
	. "github.com/YouCD/esDump/exporter"
	"log"
)

type EsIndex struct {
	Index string `json:"_index"`
	Type  string `json:"_type"`
	Id    string `json:"_id"`
}

func metaData(str string, dumpInfo DumptInfo) (mate, data []byte) {
	var esIndex = EsIndex{
		Index: dumpInfo.Index,
		Type:  gjson.GetMany(str, "_index", "_type", "_id", "_score")[1].String(),
		Id:    gjson.GetMany(str, "_index", "_type", "_id", "_score")[2].String(),
	}

	dataBytes, err := json.Marshal(esIndex)
	if err != nil {
		log.Panic(err)
	}

	mate = []byte(`{"index":` + string(dataBytes) + `}` + "\n" )
	data = []byte(gjson.Get(str, "_source").String() + "\n")
	return
}

func Importer(dumpInfo DumptInfo,lineCount int, ch chan string, ) {
	EsInit(dumpInfo)
	var buf bytes.Buffer
	i := 0
	batchCount:=lineCount/dumpInfo.Size
	batch:=0
	for str := range ch {
		m, d := metaData(str, dumpInfo)
		buf.Grow(len(m) + len(d))
		buf.Write(m)
		buf.Write(d)
		i++
		if i == dumpInfo.Size &&batch<batchCount{
			res, err := Es.Bulk(bytes.NewReader(buf.Bytes()), Es.Bulk.WithIndex(dumpInfo.Index))
			if err != nil {
				log.Panic(err)
			}
			json := Read(res.Body)
			res.Body.Close()
			count := len(gjson.Get(json, "items").Array())
			str := gjson.Get(json, "errors").String()
			if str == "false" {
				log.Printf("索引 %s 导入成功,导入数量%d，本次导入大小%d,导入批次%d/%d", dumpInfo.Index, count, buf.Len(),batch,batchCount)
			} else {
				log.Printf("导入失败...")
			}
			buf.Reset()
			i = 0
			batch++
		}else if i==lineCount%dumpInfo.Size &&batch==batchCount {
			res, err := Es.Bulk(bytes.NewReader(buf.Bytes()), Es.Bulk.WithIndex(dumpInfo.Index))
			if err != nil {
				log.Panic(err)
			}
			json := Read(res.Body)
			res.Body.Close()
			count := len(gjson.Get(json, "items").Array())
			str := gjson.Get(json, "errors").String()
			if str == "false" {
				log.Printf("索引 %s 导入成功,导入数量%d，本次导入大小%d,导入批次%d/%d", dumpInfo.Index, count, buf.Len(),batch,batchCount)
			} else {
				log.Printf("导入失败...")
			}
		}
	}
}
//
//
//func main() {
//	ch := make(chan string, 50)
//	var X = DumptInfo{
//		User:     "",
//		Password: "",
//		Host:     "http://127.0.0.1:9200",
//		Index:    "fwfdwefwe",
//		Size:     1500,
//	}
//
//	lineCount:=readFromFile.ReadFromFile("/Users/ycd/xxxx", ch)
//	Importer(X,lineCount, ch)
//
//
//}
