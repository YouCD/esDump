# esDump
[![Build Status](https://app.travis-ci.com/YouCD/esDump.svg?branch=master)](https://app.travis-ci.com/YouCD/esDump)

esDump 是golang编写的一个elasticsearch索引导出器
速度比`elasticdump`强太多
# 功能

* 支持索引导入导出
* 导出索引支持简单与复杂查询
* 根据tag使用相关的版本: 

# 使用
```shell
esDump 是golang编写的一个elasticsearch索引导出器

Usage:
  esDump [flags]
  esDump [command]

Examples:
   导出到文件
      esDump -o -E https://root:root@127.0.0.1:9200/index > Output.txt
   从文件导入到ES
      esDump -i -E https://root:root@127.0.0.1:9200/index < Output.txt
   没有账户认证  
      esDump -o -E https://root:root@127.0.0.1:9200/index > Output.txt
   添加队列大小  
      esDump -o -E https://root:root@127.0.0.1:9200/index -s 100 > Output.txt
   简单查询     
      esDump -o -E https://root:root@127.0.0.1:9200/index -q 'SomeField:SomeValue' > Output.txt
   复杂查询
      esDump -o -E https://root:root@127.0.0.1:9200/index --complex -q '{"query":{ "match_all": {} }}' > Output.txt
      esDump -o -E https://root:root@127.0.0.1:9200/index --complex -q '{"query":{ "range": {"age":{"gt":25}} }}' > Output.txt

Available Commands:
  help        Help about any command
  update      update the esDump server
  version     Print the version number of esDump

Flags:
      --complex           开启复杂查询
  -e, --endpoint string   elasticsearch Url
  -o, --export            elasticsearch Url
  -h, --help              help for esDump
  -i, --import            elasticsearch Url
  -q, --query string      query 查询
  -s, --size int          size*10，默认100即可 (default 100)


```
