package cmd

import (
	"errors"
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

	size int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "esDump",
	Short: "elasticsearch索引导出器",
	Long:  `esDump 是golang编写的一个elasticsearch索引导出器`,
	Example: `	账户认证：
                从ES中导出到文件
			esDump -i http://root:root@127.0.0.1:9200/index -o Output.txt -s 100 
		从文件导入到ES中
			esDump -i Output.txt -o http://root:root@127.0.0.1:9200/index -s 100

	没有账户认证：
		esDump -i Output.txt -o http://127.0.0.1:9200/index -s 100
		esDump -o Output.txt -i http://127.0.0.1:9200/index -s 100`,


	Run: func(cmd *cobra.Command, args []string) {

		if len(input) == 0 {
			_ = cmd.Help()
			return
		}

		if strings.Contains(input, "http://") {
			dumpInfo, err := parFlag(input)
			if err != nil {
				log.Println(err)
				return
			}
			esChannel := make(chan string, size)
			err = esHandler.Exporter(dumpInfo, esChannel)
			if err != nil {
				log.Println(err)
				return
			}
			fileHandler.WriteToFile(output, esChannel)
		} else if strings.Contains(output, "http://") {
			dumpInfo, err := parFlag(output)
			if err != nil {
				log.Println(err)
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
		log.Println(err)
		os.Exit(1)
	}
}

func init() {
	//cobra.OnInitialize(initConfig)
	rootCmd.Flags().StringVarP(&input, "input", "i", "", "input 输入源")
	rootCmd.Flags().IntVarP(&size, "size", "s", 100, "size*10，默认100即可")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "output 输出源")
	rootCmd.Flags().StringVarP(&query, "query", "q", "", "query 查询")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(updateCmd)
}

//解析参数
func parFlag(urlStr string) (dumpInfo *esHandler.DumpInfo, err error) {
	dumpInfoTemp := new(esHandler.DumpInfo)
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	if query != "" {
		dumpInfoTemp.Query = query
	}
	dumpInfoTemp.Size = size * 10
	dumpInfoTemp.Host = "http://" + u.Host
	passwordStr, hasPwd := u.User.Password()
	if u.Scheme == "http" && hasPwd && passwordStr != "" {
		dumpInfoTemp.User = u.User.Username()
		dumpInfoTemp.Password = passwordStr

		path := strings.Split(u.Path, "/")
		if path[1] == "" {
			return dumpInfoTemp, errors.New("please enter the correct path")
		}
		dumpInfoTemp.Index = path[1]
		return dumpInfoTemp, nil
	} else if u.User.Username() == "" && !hasPwd {
		path := strings.Split(u.Path, "/")
		if path[1] == "" {
			return dumpInfoTemp, errors.New("please enter the correct path")
		}
		dumpInfoTemp.Index = path[1]
		return dumpInfoTemp, nil

	} else if u.User.Username() != "" && passwordStr == "" {
		return nil, errors.New("password can not be empty")
	} else if strings.ToUpper(u.Scheme) != "HTTP" {
		return nil, errors.New("only supports http protocol")
	}

	return nil, nil
}
