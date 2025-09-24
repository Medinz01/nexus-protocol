/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"nexus-cli/internal/p2p"

	"github.com/spf13/cobra"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Sends a ping to another node on the network.",
	Long: `Pings another node to check if it's online and reachable.
You must provide the full multiaddress of the target node,
for example: /ip4/127.0.0.1/tcp/60802/p2p/12D3KooWJYLLBzL9...`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetAddress := args[0]
		fmt.Printf("Attempting to ping node at: %s\n", targetAddress)
		p2p.PingPeer(targetAddress)
	},
}

func init() {
	networkCmd.AddCommand(pingCmd)
}
