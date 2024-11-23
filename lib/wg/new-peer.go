package wg

import "golang.zx2c4.com/wireguard/wgctrl/wgtypes"

func initNewPeer() wgtypes.Key {
	masterKey := MustGeneratePrivateKey()
	return masterKey
}
