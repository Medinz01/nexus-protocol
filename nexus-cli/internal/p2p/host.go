package p2p

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

// CreateHost initializes a new libp2p host.
func CreateHost() {
	// This is the simplest way to create a host.
	// It will listen on a random TCP port on all available interfaces.
	host, err := libp2p.New()
	if err != nil {
		panic(err) // For a simple CLI, panicking on critical error is fine.
	}
	defer host.Close()

	// Get the host's Peer ID and addresses.
	// The Multiaddress format includes the Peer ID at the end.
	fmt.Println("Node is listening... Press CTRL+C to stop.")
	for _, addr := range host.Addrs() {
		fmt.Printf("  %s/p2p/%s\n", addr, host.ID())
	}

	// This part is important. It keeps the program running.
	select {} // block forever
}
func PingPeer(targetAddress string) {
	ctx := context.Background()

	// Create a new libp2p host for the pinger.
	host, err := libp2p.New()
	if err != nil {
		panic(err)
	}
	defer host.Close()

	// Parse the target address string into a Multiaddress.
	maddr, err := multiaddr.NewMultiaddr(targetAddress)
	if err != nil {
		fmt.Println("Invalid address:", err)
		return
	}

	// Extract the Peer ID from the Multiaddress.
	peerInfo, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		fmt.Println("Could not get peer info:", err)
		return
	}

	// Connect to the target peer.
	fmt.Println("Pinging...", peerInfo.ID)
	err = host.Connect(ctx, *peerInfo)
	if err != nil {
		fmt.Println("Failed to connect:", err)
		return
	}

	fmt.Println("✅ Connection established! Peer is online.")
}

