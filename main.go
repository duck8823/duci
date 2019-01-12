package main

import (
	"github.com/duck8823/duci/presentation/cmd"
	"github.com/sirupsen/logrus"
	"os"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "02/Jan/2006 15:04:05.000",
		FullTimestamp:   true,
	})
}

func main() {
	cmd.Execute(os.Args[1:])
}
