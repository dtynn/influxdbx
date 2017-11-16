// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"log"

	"github.com/dtynn/influxdbx/raftdemo/service"
	"github.com/spf13/cobra"
)

var srvCfg service.Config

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here

		srv := service.New(srvCfg)
		if err := srv.Run(context.Background()); err != nil {
			log.Println(err)
		}

		log.Println("finished")
	},
}

func init() {
	serviceCmd.PersistentFlags().StringVarP(&srvCfg.Listen, "listen", "", ":10021", "listen addr")
	serviceCmd.PersistentFlags().StringVarP(&srvCfg.DataDir, "dir", "", "./node", "data dir")
	serviceCmd.PersistentFlags().StringSliceVarP(&srvCfg.Join, "join", "", nil, "join address")
	serviceCmd.PersistentFlags().IntVarP(&srvCfg.RaftLogCacheSize, "cache-log-size", "", 256, "raft cache log size")

	RootCmd.AddCommand(serviceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serviceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
