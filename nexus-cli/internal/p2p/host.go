package p2p

import (
	"bufio"
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/multiformats/go-multiaddr"
)

// ChatProtocolID is the unique identifier for our chat protocol.
// It's capitalized to be public.
const ChatProtocolID = "/nexus/chat/1.0.0"

// bootstrapPeers are the public nodes we connect to for peer discovery.
var bootstrapPeers = []string{
	"/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
	"/dnsaddr/bootstrap.libp2p.io/p2p/QmQCU2EcUXa2kE4CRCS3iS0weA8EMvPkbEhwQroutvXaWp",
}

// ChatStreamHandler handles incoming chat messages.
// It's capitalized to be public.
func ChatStreamHandler(stream network.Stream) {
	reader := bufio.NewReader(stream)
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from stream:", err)
		return
	}
	remotePeer := stream.Conn().RemotePeer()
	fmt.Printf("Received message from %s: %s", remotePeer.ShortString(), message)
	stream.Close()
}

// CreateHost initializes a new libp2p host.
func CreateHost(ctx context.Context) (host.Host, error) {
	opts := []libp2p.Option{
		libp2p.EnableHolePunching(),
		libp2p.EnableRelay(),
	}

	host, err := libp2p.New(opts...)
	if err != nil {
		return nil, err
	}

	kademliaDHT, err := dht.New(ctx, host)
	if err != nil {
		return nil, err
	}

	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		return nil, err
	}

	for _, peerAddrStr := range bootstrapPeers {
		// FIX #1: Handle the error from NewMultiaddr before using its result.
		maddr, err := multiaddr.NewMultiaddr(peerAddrStr)
		if err != nil {
			fmt.Println("Error parsing bootstrap peer address:", err)
			continue
		}
		p, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			fmt.Println("Error getting peer info from address:", err)
			continue
		}

		ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
		host.Connect(ctxWithTimeout, *p)
		cancel()
	}

	routingDiscovery := routing.NewRoutingDiscovery(kademliaDHT)
	dutil.Advertise(ctx, routingDiscovery, ChatProtocolID)

	return host, nil
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
	stream, err := host.NewStream(ctx, peerInfo.ID, ChatProtocolID)
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
