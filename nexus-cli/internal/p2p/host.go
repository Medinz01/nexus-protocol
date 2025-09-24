package p2p

import (
	"bufio"
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

const chatProtocolID = "/nexus/chat/1.0.0"

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
	host.SetStreamHandler(chatProtocolID, chatStreamHandler)
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

// chatStreamHandler reads incoming messages and prints them to the console.
func chatStreamHandler(stream network.Stream) {
	// A buffered reader is useful for reading data line-by-line.
	reader := bufio.NewReader(stream)

	// Read the message from the stream.
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from stream:", err)
		return
	}

	// Print the message from the remote peer.
	remotePeer := stream.Conn().RemotePeer()
	fmt.Printf("Received message from %s: %s", remotePeer.ShortString(), message)

	// Close the stream when we're done.
	stream.Close()
}
func SendChatMessage(targetAddress string, message string) {
	ctx := context.Background()
	host, err := libp2p.New()
	if err != nil {
		panic(err)
	}
	defer host.Close()

	maddr, err := multiaddr.NewMultiaddr(targetAddress)
	if err != nil {
		fmt.Println("Invalid address:", err)
		return
	}

	peerInfo, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		fmt.Println("Could not get peer info:", err)
		return
	}

	// Connect to the target peer.
	err = host.Connect(ctx, *peerInfo)
	if err != nil {
		fmt.Println("Failed to connect:", err)
		return
	}

	// Open a new stream to the peer for our chat protocol.
	stream, err := host.NewStream(ctx, peerInfo.ID, chatProtocolID)
	if err != nil {
		fmt.Println("Failed to open stream:", err)
		return
	}
	defer stream.Close()

	// Write the message to the stream.
	_, err = stream.Write([]byte(message + "\n")) // Add a newline to separate messages
	if err != nil {
		fmt.Println("Failed to send message:", err)
		return
	}

	fmt.Println("✅ Message sent successfully!")
}
