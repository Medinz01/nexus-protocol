/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
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
		ctx := context.Background()
		host, err := p2p.CreateHost(ctx)
		if err != nil {
			panic(err)
		}
		defer host.Close()

		// Set the stream handler for our chat protocol.
		host.SetStreamHandler(p2p.ChatProtocolID, p2p.ChatStreamHandler)

		fmt.Println("Node is listening... Press CTRL+C to stop.")
		for _, addr := range host.Addrs() {
			fmt.Printf("  %s/p2p/%s\n", addr, host.ID())
		}

		select {}
	},
}

func init() {
	networkCmd.AddCommand(listenCmd)
}
