package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/YouCD/esDump/esHandler"
	"github.com/YouCD/esDump/fileHandler"
	"github.com/spf13/cobra"

	"os"
)

var (
	query      string
	exists     string
	size       int
	endpoint   string
	PathErr    = errors.New("please enter the correct path")
	exportFlag bool
	importFlag bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "esDump",
	Short: "elasticsearch索引导出器",
	Long:  `esDump 是golang编写的一个elasticsearch索引导出器`,
	Example: `   导出到文件
       esDump -o -E http://root:root@127.0.0.1:9200/index >Output.txt

   导入到ES
       esDump -i -E http://root:root@127.0.0.1:9200/index <Output.txt

   没有账户认证：
      esDump -o -E http://root:root@127.0.0.1:9200/index >Output.txt

   添加队列大小：
      esDump -o -E http://root:root@127.0.0.1:9200/index -s 100 >Output.txt

   查询条件：
      esDump -o -E http://root:root@127.0.0.1:9200/index -q 'SomeField:SomeValue' >Output.txt

   存在条件：
      esDump -o -E http://root:root@127.0.0.1:9200/index  -e 'SomeField' >Output.txt
`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(endpoint) == 0 {
			_ = cmd.Help()
			return
		}
		dumpInfo, err := parFlag(endpoint)
		switch {
		case importFlag:
			if strings.Contains(endpoint, "http://") {
				if err != nil {
					fmt.Println(err)
					os.Exit(-1)
					return
				}
				esChannel := make(chan string, size)
				go fileHandler.ReadFromPipe(esChannel)
				esHandler.PipeImporter(dumpInfo, esChannel)
				log.Println("导入完成")
			} else {
				_ = cmd.Help()
				return
			}
		case exportFlag:
			esChannel := make(chan string, size)
			err = esHandler.Exporter(dumpInfo, esChannel)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
				return
			}
			for i := range esChannel {
				fmt.Println(i)
			}
		}

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	//cobra.OnInitialize(initConfig)
	rootCmd.Flags().IntVarP(&size, "size", "s", 100, "size*10，默认100即可")
	rootCmd.Flags().StringVarP(&query, "query", "q", "", "query 查询")
	rootCmd.Flags().StringVarP(&exists, "exists", "e", "", "exists 必须存在字段")
	rootCmd.Flags().StringVarP(&endpoint, "endpoint", "E", "", "elasticsearch Url")
	rootCmd.Flags().BoolVarP(&exportFlag, "export", "o", false, "elasticsearch Url")
	rootCmd.Flags().BoolVarP(&importFlag, "import", "i", false, "elasticsearch Url")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(updateCmd)
}

//parFlag 解析参数
func parFlag(urlStr string) (dumpInfo *esHandler.DumpInfo, err error) {
	dumpInfoTemp := new(esHandler.DumpInfo)
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	if query != "" {
		dumpInfoTemp.Query = query
	}
	if exists != "" {
		dumpInfoTemp.ExistsFilter = exists
	}
	dumpInfoTemp.Size = size * 10
	dumpInfoTemp.Host = "http://" + u.Host
	passwordStr, hasPwd := u.User.Password()

	switch {
	case u.Scheme == "http" && hasPwd && passwordStr != "":
		dumpInfoTemp.User = u.User.Username()
		dumpInfoTemp.Password = passwordStr

		path := strings.Split(u.Path, "/")
		if path[1] == "" {
			return dumpInfoTemp, PathErr
		}
		dumpInfoTemp.Index = path[1]
		return dumpInfoTemp, nil
	case u.User.Username() == "" && !hasPwd:
		path := strings.Split(u.Path, "/")
		if path[1] == "" {
			return dumpInfoTemp, PathErr
		}
		dumpInfoTemp.Index = path[1]
		return dumpInfoTemp, nil
	case u.User.Username() != "" && passwordStr == "":
		return nil, errors.New("password can not be empty")
	case strings.ToUpper(u.Scheme) != "HTTP":
		return nil, errors.New("only support http protocol")
	}
	return nil, nil
}
