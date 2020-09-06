package cmd

import (
	"github.com/YouCD/esDump/exporter"
	"github.com/YouCD/esDump/readFromFile"
	"github.com/YouCD/esDump/importer"
	"log"
	"github.com/YouCD/esDump/writeToFile"
	"strings"
	"github.com/spf13/cobra"
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
	Version: `v0.1`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(input) == 0 {
			cmd.Help()
			return
		}

		if strings.Contains(input, "http:") {
			parFlag(input)
			dampinfo := exporter.DumptInfo{
				user,
				password,
				host,
				index,
				size,
			}
			eschannel := make(chan string, size)
			exporter.Exporter(dampinfo, eschannel)
			writeToFile.WriteToFile(output, eschannel)
		} else if strings.Contains(output, "http:") {
			parFlag(output)
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

}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

//解析参数
func parFlag(flag string) {

	if strings.Contains(flag, "http://") && strings.Contains(flag, "@") && strings.Contains(flag, "/") {
		//截取baseauth认证信息
		uriAndUserinfo := strings.Split(flag, "@")[0]
		userinfo := strings.Split(uriAndUserinfo, "//")[1]
		//获取user
		user = strings.Split(userinfo, ":")[0]
		//获取密码
		password = strings.Split(userinfo, ":")[1]
		// 拼接url
		uri := strings.Split(uriAndUserinfo, "/")[0]
		host = uri + "//" + strings.Split(strings.Split(flag, "/")[2], "@")[1]
		//获取index
		index = strings.Split(flag, "/")[len(strings.Split(flag, "/"))-1]

	} else if strings.Contains(flag, "@") == false {
		// 拼接url
		uri := strings.Split(flag, "/")[0]
		host = uri + "//" + strings.Split(strings.Split(flag, "/")[2], "@")[0]
		//获取index
		index = strings.Split(flag, "/")[len(strings.Split(flag, "/"))-1]
	} else {
		log.Println("参数错误")
	}

}

func main() {
	flag := "http://172.0.0.1:9200/index"
	parFlag(flag)
}

//func initFlag(flag string)bool  {
//	if strings.ContainsAny(flag, "http")&&strings.ContainsAny(flag, "@")&&strings.ContainsAny(flag, "/"){
//		parFlag(flag)
//		return true
//	}else {
//		parFlag(flag)
//		return false
//	}
//}
//
//

//func initConfig() {
//	if cfgFile != "" {
//		viper.SetConfigFile(cfgFile)
//	} else {
//		home, err := homedir.Dir()
//		if err != nil {
//			fmt.Println(err)
//			os.Exit(1)
//		}
//		viper.AddConfigPath(home)
//		viper.SetConfigName(".demo")
//	}
//
//	viper.AutomaticEnv()
//	if err := viper.ReadInConfig(); err == nil {
//		fmt.Println("Using config file:", viper.ConfigFileUsed())
//	}
//}
