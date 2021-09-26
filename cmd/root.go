package cmd

import (
	"errors"
	"fmt"
	"github.com/YouCD/esDump/esHandler"
	"github.com/YouCD/esDump/fileHandler"
	"github.com/spf13/cobra"
	"log"
	"net/url"
	"strings"
	//"github.com/spf13/viper"
	"os"
)

//var cfgFile string
var (
	output string
	input  string
	query  string
	exists string
	size   int

	PathErr = errors.New("please enter the correct path")
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "esDump",
	Short: "elasticsearch索引导出器",
	Long:  `esDump 是golang编写的一个elasticsearch索引导出器`,
	Example: `   导出到文件
       esDump -i http://root:root@127.0.0.1:9200/index -o Output.txt

   导入到ES
       esDump -i Output.txt -o http://root:root@127.0.0.1:9200/index

   没有账户认证：
      esDump -i Output.txt -o http://127.0.0.1:9200/index

   添加队列大小：
      esDump -i Output.txt -o http://127.0.0.1:9200/index -s 100

   查询条件：
      esDump -o Output.txt -i http://127.0.0.1:9200/index  -q 'SomeField:SomeValue'

   存在条件：
      esDump -o Output.txt -i http://127.0.0.1:9200/index  -e 'SomeField'
`,

	Run: func(cmd *cobra.Command, args []string) {

		if len(input) == 0 {
			_ = cmd.Help()
			return
		}

		if strings.Contains(input, "http://") {
			dumpInfo, err := parFlag(input)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
				return
			}
			esChannel := make(chan string, size)
			err = esHandler.Exporter(dumpInfo, esChannel)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
				return
			}
			fileHandler.WriteToFile(output, esChannel)
		} else if strings.Contains(output, "http://") {
			dumpInfo, err := parFlag(output)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
				return
			}

			esChannel := make(chan string, size)
			lineCount := fileHandler.ReadFromFile(input, esChannel)
			esHandler.Importer(dumpInfo, lineCount, esChannel)
			log.Println("导入完成")
		} else {
			_ = cmd.Help()
			return
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
	rootCmd.Flags().StringVarP(&input, "input", "i", "", "input 输入源")
	rootCmd.Flags().IntVarP(&size, "size", "s", 100, "size*10，默认100即可")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "output 输出源")
	rootCmd.Flags().StringVarP(&query, "query", "q", "", "query 查询")
	rootCmd.Flags().StringVarP(&exists, "exists", "e", "", "exists 必须存在字段")

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

	switch  {
	case u.Scheme == "http" && hasPwd && passwordStr != "":
		dumpInfoTemp.User = u.User.Username()
		dumpInfoTemp.Password = passwordStr

		path := strings.Split(u.Path, "/")
		if path[1] == "" {
			return dumpInfoTemp, PathErr
		}
		dumpInfoTemp.Index = path[1]
		return dumpInfoTemp, nil
	case u.User.Username() == "" && !hasPwd :
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
