/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"nexus-cli/internal/wallet" // <-- FIX #1: Import your internal wallet package

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new network wallet and identity.",
	Long: `Creates a new network wallet by generating a secure cryptographic keypair (Ed25519).

Your private key is stored locally, and you will be given a 12-word recovery phrase.
This phrase is the only way to recover your account. Store it securely and offline.`,
	Run: func(cmd *cobra.Command, args []string) {
		// This part of your code was already correct!
		pubKey, _, mnemonic, err := wallet.GenerateNewKeys()
		if err != nil {
			fmt.Println("Error creating wallet:", err)
			return
		}

		// Convert public key to a more readable format if desired (e.g., hex)
		pubKeyHex := fmt.Sprintf("%x", pubKey)

		fmt.Println("✅ Wallet created successfully!")
		fmt.Println("=================================================================")
		fmt.Println("Your Public Key (Peer ID):", pubKeyHex)
		fmt.Println("\n!!! IMPORTANT: Save this 12-word recovery phrase securely and offline. !!!")
		fmt.Println("=================================================================")
		fmt.Println(mnemonic)
	},
}

func init() {
	walletCmd.AddCommand(createCmd) // <-- FIX #2: Removed single quotes around walletCmd
}
