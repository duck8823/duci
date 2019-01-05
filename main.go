package main

import (
	"github.com/duck8823/duci/presentation/cmd"
	"os"
)

func main() {
	cmd.Execute(os.Args[1:])
}
