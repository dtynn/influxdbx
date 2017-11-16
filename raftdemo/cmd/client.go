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
	"log"
	"net"
	"strconv"
	"time"

	"github.com/dtynn/influxdbx/raftdemo/rpc"
	"github.com/dtynn/influxdbx/raftdemo/rpc/proto"
	"github.com/dtynn/influxdbx/tcp"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			log.Fatalln("expected 3 args")
		}

		// TODO: Work your own magic here
		dialer := func(address string, timeout time.Duration) (net.Conn, error) {
			return tcp.Dial("tcp", address, timeout, 11)
		}

		cli, err := rpc.NewClusterClient(args[2], dialer)
		if err != nil {
			log.Fatalln(err)
		}

		now := time.Now()

		resp, err := cli.Add(context.Background(), &proto.ClusterAddReq{
			Key:   args[0] + " " + strconv.FormatInt(now.Unix(), 10),
			Value: args[1] + " " + strconv.FormatInt(now.Unix(), 10),
		})

		if err != nil {
			log.Fatalln(err)
		}

		log.Println("get code", resp.GetCode())
	},
}

func init() {
	RootCmd.AddCommand(clientCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
