package cmd

import (
	"fmt"
	"testing"
)

func Test_tagToFloat(t *testing.T) {
	version := tagToInt("v1.9.2")
	fmt.Println(version)
}

func Test_filename(t *testing.T) {
	filename("https://github.com/YouCD/esDump/releases/download/v6.8.23/esDump-linux-amd64.txz")
}
