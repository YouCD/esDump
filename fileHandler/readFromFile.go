package fileHandler

import (
	"bufio"
	"log"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

//ReadFromPipe 读取文件
func ReadFromPipe(ch chan string) {
	log.Println("开始导入")
	if !terminal.IsTerminal(0) {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			ch <- scanner.Text()
		}
		close(ch)
	}
}
