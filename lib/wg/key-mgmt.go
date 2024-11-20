package wg

import (
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"log"
)

func mustParsePublicKey(key string) wgtypes.Key {
	publicKey, err := wgtypes.ParseKey(key)
	if err != nil {
		log.Fatalf("Failed to parse public key: %v", err)
	}
	return publicKey
}

func mustParsePrivateKey(key string) wgtypes.Key {
	privateKey, err := wgtypes.ParseKey(key)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}
	return privateKey
}

func mustGeneratePrivateKey() wgtypes.Key {
	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}
	return privateKey
}

func generateKeyPair() (wgtypes.Key, wgtypes.Key, error) {
	privateKey := mustGeneratePrivateKey()
	publicKey := privateKey.PublicKey()
	return privateKey, publicKey, nil
}

func mustGenerateKeyPair() (wgtypes.Key, wgtypes.Key) {
	privateKey, publicKey, err := generateKeyPair()
	if err != nil {
		log.Fatalf("Failed to generate key pair: %v", err)
	}
	return privateKey, publicKey
}
