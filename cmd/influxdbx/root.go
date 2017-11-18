package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/dtynn/influxdbx/cmd/influxdbx/run"
	"github.com/spf13/cobra"
)

var cfgFile string

var rootcmd = &cobra.Command{
	Use:   "influxdbx",
	Short: "influxdbx",
	Long:  `influxdbx`,
	Run: func(cmd *cobra.Command, args []string) {
		var cfg *run.Config
		var err error

		if cfgFile == "" {
			zlog.Info("no configuration provided, using default settings")
			cfg, err = run.NewDemoConfig()
		} else {

			cfg = run.NewConfig()
			err = cfg.Load(cfgFile)
		}

		if err != nil {
			zlog.Fatal(err.Error())
		}

		server, err := run.NewServer(cfg)
		if err != nil {
			zlog.Info("!!!")
			zlog.Fatal(err.Error())
		}

		if err := server.Open(); err != nil {
			zlog.Fatal(err.Error())
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
