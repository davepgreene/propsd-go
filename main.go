package main

import (
	"fmt"
	"os"

	"github.com/davepgreene/propsd/cmd"
)

func main() {
	if err := cmd.PropsdCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
