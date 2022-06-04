package es

import (
	"bytes"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v6"
	"io"
	"log"
	"os"
)

var (
	Es *elasticsearch.Client
)

type DumpInfo struct {
	User     string
	Password string
	Host     string
	Index    string
	Size     int
	Query    string
	Complex  bool
}

func EsInit(dumpInfo *DumpInfo) {
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

	return
}

func read(r io.Reader) string {
	var b bytes.Buffer
	_, err := b.ReadFrom(r)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return b.String()
}
