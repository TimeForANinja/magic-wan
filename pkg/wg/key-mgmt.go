package wg

import (
	"fmt"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func MustParsePublicKey(key string) wgtypes.Key {
	publicKey, err := wgtypes.ParseKey(key)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse public key: %v", err))
	}
	return publicKey
}

func MustParsePrivateKey(key string) wgtypes.Key {
	privateKey, err := wgtypes.ParseKey(key)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse private key: %v", err))
	}
	return privateKey
}

func MustGeneratePrivateKey() wgtypes.Key {
	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		panic(fmt.Sprintf("Failed to generate private key: %v", err))
	}
	return privateKey
}

func GenerateKeyPair() (wgtypes.Key, wgtypes.Key, error) {
	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return wgtypes.Key{}, wgtypes.Key{}, err
	}

	publicKey := privateKey.PublicKey()
	return privateKey, publicKey, nil
}

func MustGenerateKeyPair() (wgtypes.Key, wgtypes.Key) {
	privateKey, publicKey, err := GenerateKeyPair()
	if err != nil {
		panic(fmt.Sprintf("Failed to generate key pair: %v", err))
	}
	return privateKey, publicKey
}
