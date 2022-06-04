package cmd

import (
	"bufio"
	"fmt"
	es2 "github.com/YouCD/esDump/pkg/es"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"net/url"
	"strings"

	"github.com/spf13/cobra"

	"os"
)

var (
	query       string
	size        int
	endpoint    string
	complexFlag bool
)

func init() {
	rootCmd.Flags().IntVarP(&size, "size", "s", 100, "size*10，默认100即可")
	rootCmd.Flags().StringVarP(&query, "query", "q", "", "query 查询")
	rootCmd.Flags().StringVarP(&endpoint, "endpoint", "e", "", "elasticsearch Url")
	rootCmd.Flags().BoolVar(&complexFlag, "complex", false, "开启复杂查询")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(updateCmd)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "esDump",
	Short: "elasticsearch索引导出器",
	Long:  `esDump 是golang编写的一个elasticsearch索引导出器`,
	Example: `   导出到文件
      esDump -e http://root:root@127.0.0.1:9200/index > Output.txt
   从文件导入到ES
      esDump -e http://root:root@127.0.0.1:9200/index < Output.txt
   没有账户认证  
      esDump -e http://root:root@127.0.0.1:9200/index > Output.txt
   添加队列大小  
      esDump -e http://root:root@127.0.0.1:9200/index -s 100 > Output.txt
   简单查询     
      esDump -e http://root:root@127.0.0.1:9200/index -q 'SomeField:SomeValue' > Output.txt
   复杂查询
      esDump -e http://root:root@127.0.0.1:9200/index --complex -q '{"query":{ "match_all": {} }}' > Output.txt
      esDump -e http://root:root@127.0.0.1:9200/index --complex -q '{"query":{ "range": {"age":{"gt":25}} }}' > Output.txt`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(endpoint) == 0 {
			_ = cmd.Help()
			return
		}
		dumpInfo, err := parFlag(endpoint)
		if err != nil {
			log.Println(err)
			os.Exit(1)
			return
		}
		if dumpInfo == nil {
			log.Println("参数解析失败，不支持的endpoint")
			os.Exit(1)
			return
		}

		esChannel := make(chan string, size)

		switch {
		case !terminal.IsTerminal(0):
			go func() {
				log.Println("开始导入")
				if !terminal.IsTerminal(0) {
					scanner := bufio.NewScanner(os.Stdin)
					for scanner.Scan() {
						esChannel <- scanner.Text()
					}
					close(esChannel)
				}

			}()

			es2.PipeImporter(dumpInfo, esChannel)
			log.Println("导入完成")
		default:
			esChannel := make(chan string, size)
			err = es2.Exporter(dumpInfo, esChannel)
			if err != nil {
				log.Println(err)
				os.Exit(1)
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
		log.Println(err)
		os.Exit(1)
	}
}

//parFlag 解析参数
func parFlag(urlStr string) (dumpInfo *es2.DumpInfo, err error) {
	dumpInfo = new(es2.DumpInfo)
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	if query != "" {
		dumpInfo.Query = query
	}

	dumpInfo.Size = size * 10

	passwordStr, hasPwd := u.User.Password()
	if hasPwd {
		if passwordStr == "" {
			log.Println("password can not be empty.")
			os.Exit(1)
		}
		dumpInfo.Password = passwordStr
		dumpInfo.User = u.User.Username()
	}

	path := strings.Split(u.Path, "/")
	if path[1] == "" {
		log.Println("index can not be empty.")
		os.Exit(1)
	}
	dumpInfo.Index = path[1]
	dumpInfo.Complex = complexFlag

	switch {
	case strings.ToUpper(u.Scheme) == "HTTPS":
		dumpInfo.Host = u.Scheme + "://" + u.Host
	case strings.ToUpper(u.Scheme) == "HTTP":
		dumpInfo.Host = u.Scheme + "://" + u.Host
	default:
		log.Println("only support http or https protocol")
		os.Exit(1)
	}
	return
}
