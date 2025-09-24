package wallet

import (
	"crypto/ed25519"

	"github.com/tyler-smith/go-bip39"
)

// GenerateNewKeys creates a new keypair and a mnemonic phrase.
func GenerateNewKeys() (ed25519.PublicKey, ed25519.PrivateKey, string, error) {
	// Generate 32 bytes of entropy, which is the size of an Ed25519 private key seed.
	// This is secure enough for a 24-word mnemonic, but we'll generate 12 words for simplicity.
	entropy, err := bip39.NewEntropy(128) // 128 bits for a 12-word phrase
	if err != nil {
		return nil, nil, "", err
	}

	// Generate a mnemonic phrase from the entropy.
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, nil, "", err
	}

	// Create a seed from the mnemonic. This is what we'll use for the keypair.
	// We use an empty password for the seed.
	seed := bip39.NewSeed(mnemonic, "")

	// Ed25519 keys can be generated from a 32-byte seed.
	// The seed from BIP39 is 64 bytes, so we take the first 32.
	privateKey := ed25519.NewKeyFromSeed(seed[:32])
	publicKey := privateKey.Public().(ed25519.PublicKey)

	return publicKey, privateKey, mnemonic, nil
}
