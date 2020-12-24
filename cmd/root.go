package cmd

import (
	"errors"
	"fmt"
	"github.com/YouCD/esDump/exporter"
	"github.com/YouCD/esDump/importer"
	"github.com/YouCD/esDump/readFromFile"
	"github.com/YouCD/esDump/writeToFile"
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

	size     int
	index    string
	user     string
	password string
	host     string

	//version info
	version   string
	commitID  string
	buildTime string
	goVersion string
	buildUser string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "esDump",
	Short: "elasticsearch索引导出器",
	Long:  `esDump 是golang编写的一个elasticsearch索引导出器`,
	Example: `	账户认证：
		esDump -i http://root:root@127.0.0.1:9200/index -o Output.txt -s 100
		esDump -i Output.txt -o http://root:root@127.0.0.1:9200/index -s 100

	没有账户认证：
		esDump -i Output.txt -o http://127.0.0.1:9200/index -s 100
		esDump -o Output.txt -i http://127.0.0.1:9200/index -s 100`,

	Run: func(cmd *cobra.Command, args []string) {

		if len(input) == 0 {
			cmd.Help()
			return
		}

		if strings.Contains(input, "http://") {
			err := parFlag(input)
			if err != nil {
				log.Println(err)
				return
			}
			dampinfo := exporter.DumptInfo{
				user,
				password,
				host,
				index,
				size,
			}
			eschannel := make(chan string, size)
			err = exporter.Exporter(dampinfo, eschannel)
			if err != nil {
				return
			}
			writeToFile.WriteToFile(output, eschannel)
		} else if strings.Contains(output, "http://") {
			err := parFlag(output)
			if err != nil {
				log.Println(err)
				return
			}
			dampinfo := exporter.DumptInfo{
				user,
				password,
				host,
				index,
				size * 10,
			}
			eschannel := make(chan string, size)
			lineCount := readFromFile.ReadFromFile(input, eschannel)
			importer.Importer(dampinfo, lineCount, eschannel)
			log.Println("导入完成")
		} else {
			cmd.Help()
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
	rootCmd.Flags().StringVarP(&input, "input", "i", "", "es->file：http://root:root@127.0.0.1:9200/index;  file->es: 文件名")
	rootCmd.Flags().IntVarP(&size, "size", "s", 100, "es->file：每次获取数据量;                         file->es: size*10")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "es->file：文件名;                                 file->es: http://root:root@127.0.0.1:9200/index")
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of esDump",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version:   %s\n", version)
		fmt.Printf("CommitID:  %s\n", commitID)
		fmt.Printf("BuildTime: %s\n", buildTime)
		fmt.Printf("GoVersion: %s\n", goVersion)
		fmt.Printf("BuildUser: %s\n", buildUser)
	},
}

//解析参数
func parFlag(flag string) (err error) {
	u, err := url.Parse(flag)
	if err != nil {
		return err
	}
	p, hasPwd := u.User.Password()
	if u.Scheme == "http" && hasPwd && p != "" {
		user = u.User.Username()
		password = p
		host = u.Host
	} else if !hasPwd {
		return errors.New("Please enter password or check url format")
	} else if p == "" {
		return errors.New("Password can not be empty")
	} else if u.Scheme != strings.ToUpper("HTTP") {
		return errors.New("Only supports http protocol")
	}
	path := strings.Split(u.Path, "/")
	if path[1] == "" {
		return errors.New("Please enter the correct path")
	}
	index = path[1]
	return nil
}
