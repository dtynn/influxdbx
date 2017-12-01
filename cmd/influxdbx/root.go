package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dtynn/influxdbx/cmd/influxdbx/server"
	"github.com/spf13/cobra"
)

var cfgFile string

var rootcmd = &cobra.Command{
	Use:   "influxdbx",
	Short: "influxdbx",
	Long:  `influxdbx`,
	Run: func(cmd *cobra.Command, args []string) {
		if cfgFile == "" {
			zlog.Fatal("no configuration provided, using default settings")
		}

		cfg := server.NewConfig()
		err := cfg.Load(cfgFile)

		if err != nil {
			zlog.Fatal(fmt.Sprintf("fail to load config: %s", err))
		}

		server, err := server.NewServer(cfg)
		if err != nil {
			zlog.Fatal(fmt.Sprintf("fail to init server: %s", err))
		}

		if err := server.Open(); err != nil {
			zlog.Fatal(fmt.Sprintf("fail to open server: %s", err))
		}

		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
		zlog.Info("Listening for signals")

		// Block until one of the signals above is received
		<-signalCh
		zlog.Info("Signal received, initializing clean shutdown...")

		server.Close()
	},
}

func init() {
	rootcmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
}
