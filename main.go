package main

import (
	"fmt"
	"os"

	"github.com/anhdt-vnpay/f5_fulltext_search/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
