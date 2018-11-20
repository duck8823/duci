package main

import (
	"github.com/duck8823/duci/application/cmd"
	"os"
)

func main() {
	cmd.Execute(os.Args[1:])
}
