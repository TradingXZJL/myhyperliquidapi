package myhyperliquidapi

import (
	"crypto/ecdsa"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

type Wallet struct {
	PrivateKey    *ecdsa.PrivateKey
	WalletAddress string
	VaultAddress  string
}

// privateKeyHex 钱包的私钥
// walletAddress 钱包的地址
// vaultAddress 主账号的钱包地址
func NewWallet(privateKeyHex string, walletAddress string, vaultAddress string) *Wallet {
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyHex, "0x"))
	if err != nil {
		log.Error(err)
		return nil
	}
	return &Wallet{
		PrivateKey:    privateKey,
		WalletAddress: walletAddress,
		VaultAddress:  vaultAddress,
	}
}
