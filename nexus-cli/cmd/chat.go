/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"nexus-cli/internal/p2p"

	"github.com/spf13/cobra"
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat [multiaddress] [message]",
	Short: "Sends a single chat message to another node.",
	Args:  cobra.ExactArgs(2), // Requires exactly two arguments
	Run: func(cmd *cobra.Command, args []string) {
		targetAddress := args[0]
		message := args[1]
		p2p.SendChatMessage(targetAddress, message)
	},
}

func init() {
	networkCmd.AddCommand(chatCmd)

}
