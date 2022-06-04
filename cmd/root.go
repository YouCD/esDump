package cmd

import (
	"bufio"
	"fmt"
	"github.com/YouCD/esDump/pkg/es"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"log"

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
      esDump -e https://root:root@127.0.0.1:9200/index > Output.txt
   从文件导入到ES
      esDump -e https://root:root@127.0.0.1:9200/index < Output.txt
   没有账户认证  
      esDump -e https://root:root@127.0.0.1:9200/index > Output.txt
   添加队列大小  
      esDump -e https://root:root@127.0.0.1:9200/index -s 100 > Output.txt
   简单查询     
      esDump -e https://root:root@127.0.0.1:9200/index -q 'SomeField:SomeValue' > Output.txt
   复杂查询
      esDump -e https://root:root@127.0.0.1:9200/index --complex -q '{"query":{ "match_all": {} }}' > Output.txt
      esDump -e https://root:root@127.0.0.1:9200/index --complex -q '{"query":{ "range": {"age":{"gt":25}} }}' > Output.txt`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(endpoint) == 0 {
			_ = cmd.Help()
			return
		}

		esDump := es.NewEsDump(endpoint, query, size, complexFlag)

		if esDump == nil {
			log.Println("初始化失败")
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

			esDump.PipeImporter(esChannel)
			log.Println("导入完成")
		default:
			esChannel := make(chan string, size)
			err := esDump.Exporter(esChannel)
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
