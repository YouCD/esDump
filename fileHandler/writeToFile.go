package fileHandler

import (
	"log"
	"os"
)

//WriteToFile 写入文件
func WriteToFile(filename string, ch chan string) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println("文件打开失败", err.Error())
		os.Exit(-1)
	}
	defer file.Close()
	for i := range ch {
		file.WriteString(i + "\n")
	}
	log.Println("导出完成")
}
