/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"nexus-cli/internal/p2p"

	"github.com/spf13/cobra"
)

// listenCmd represents the listen command
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Starts the Nexus node and listens for incoming connections.",
	Long: `This command initializes your node on the P2P network.
It will print your Peer ID and the network addresses it's listening on.
Keep this command running to stay connected to the network.`,
	Run: func(cmd *cobra.Command, args []string) {
		p2p.CreateHost()
	},
}

func init() {
	networkCmd.AddCommand(listenCmd)
}
