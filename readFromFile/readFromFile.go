package readFromFile

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
)

//读取文件
func ReadFromFile(filename string, ch chan (string))(lineCout int ){
	log.Println("开始导入")
	go func() {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			ch <- scanner.Text()
		}
		close(ch)
	}()
	lineCout,err:=lineCounter(filename)
	if err != nil {
		log.Fatal(err)
	}
	return
}

//获取文件的总行数
func lineCounter(filename string) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	defer file.Close()
	for {
		c, err := file.Read(buf)
		count += bytes.Count(buf[:c], lineSep)
		switch {
		case err == io.EOF:
			return count, nil
		case err != nil:
			return count, err
		}
	}
}