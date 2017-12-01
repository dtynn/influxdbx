package main

import (
	"os"

	"github.com/influxdata/influxdb/logger"
)

var zlog = logger.New(os.Stderr)

func main() {
	if err := rootcmd.Execute(); err != nil {
		zlog.Fatal(err.Error())
	}
}
