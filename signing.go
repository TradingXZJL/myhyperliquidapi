package myhyperliquidapi

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/sha3"
)

// Constants
const (
	MainnetChainId = 1337
	TestnetChainId = 1337
)

func FloatToWire(x float64) string {
	roundedStr := fmt.Sprintf("%.8f", x)

	if roundedStr == "-0.00000000" {
		return "0"
	}

	s := strings.TrimRight(roundedStr, "0")
	s = strings.TrimRight(s, ".")
	if s == "" {
		return "0"
	}
	if s == "-0" {
		return "0"
	}
	return s
}

// ActionHash calculates the hash of an action for signing (L1 actions)
func ActionHash[T any](action T, vaultAddress *string, nonce uint64, expiresAfter uint64) []byte {
	// MsgPack Serialize Action
	actionBytes, err := msgpack.Marshal(action)
	if err != nil {
		log.Error(err)
		return nil
	}

	// Prepare Buffer
	var data []byte
	data = append(data, actionBytes...)

	// Append Nonce (Big-endian uint64)
	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, nonce)
	data = append(data, nonceBytes...)

	// Handle Vault Address
	if vaultAddress == nil {
		data = append(data, 0x00)
	} else {
		data = append(data, 0x01)
		// Remove 0x prefix and decode hex
		addrHex := strings.TrimPrefix(*vaultAddress, "0x")
		addrBytes, _ := hex.DecodeString(addrHex)
		data = append(data, addrBytes...)
	}

	// Append Expires After
	if expiresAfter != 0 {
		data = append(data, 0x00)
		expiresAfterBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(expiresAfterBytes, expiresAfter)
		data = append(data, expiresAfterBytes...)
	}

	// Keccak256
	hash := sha3.NewLegacyKeccak256()
	hash.Write(data)
	return hash.Sum(nil)
}

func ConstructPhantomAgent(hash []byte, isMainnet bool) map[string]any {
	source := "a"
	if !isMainnet {
		source = "b"
	}
	return map[string]any{
		"source":       source,
		"connectionId": hash,
	}
}

type Signature struct {
	R *string `json:"r"`
	S *string `json:"s"`
	V *uint8  `json:"v"`
}

// SignInner signs the EIP-712 typed data
func SignInner(wallet *Wallet, data apitypes.TypedData) (Signature, error) {
	// HashStruct logic using apitypes

	domainSeparator, err := data.HashStruct("EIP712Domain", data.Domain.Map())
	if err != nil {
		return Signature{}, fmt.Errorf("failed to hash domain: %w", err)
	}
	typedDataHash, err := data.HashStruct(data.PrimaryType, data.Message)
	if err != nil {
		return Signature{}, fmt.Errorf("failed to hash message: %w", err)
	}

	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	hash := crypto.Keccak256(rawData)

	signature, err := crypto.Sign(hash, wallet.PrivateKey)
	if err != nil {
		return Signature{}, err
	}

	// Signature is [R || S || V] - 65 bytes
	if len(signature) != 65 {
		return Signature{}, fmt.Errorf("invalid signature length: %d", len(signature))
	}

	r := hexutil.Encode(signature[:32])
	s := hexutil.Encode(signature[32:64])
	v := signature[64]

	if v < 27 {
		v += 27
	}

	return Signature{
		R: &r,
		S: &s,
		V: &v,
	}, nil
}

func L1Payload(phantomAgent map[string]any) apitypes.TypedData {
	return apitypes.TypedData{
		Domain: apitypes.TypedDataDomain{
			ChainId:           math.NewHexOrDecimal256(MainnetChainId),
			Name:              "Exchange",
			VerifyingContract: "0x0000000000000000000000000000000000000000",
			Version:           "1",
		},
		Types: apitypes.Types{
			"Agent": []apitypes.Type{
				{Name: "source", Type: "string"},
				{Name: "connectionId", Type: "bytes32"},
			},
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
		},
		PrimaryType: "Agent",
		Message:     phantomAgent,
	}
}

func SignL1Action[T any](wallet *Wallet, action T, activePool *string, nonce uint64, expiresAfter uint64, isMainnet bool) (Signature, error) {
	// 1. Calculate Action Hash
	hash := ActionHash(action, activePool, nonce, expiresAfter)
	// 2. Construct Phantom Agent
	phantomAgent := ConstructPhantomAgent(hash, isMainnet)
	// 3. Construct EIP-712 Typed Data
	data := L1Payload(phantomAgent)

	return SignInner(wallet, data)
}

func UserSignedPayload(action map[string]interface{}, payloadTypes []apitypes.Type, primaryType string) apitypes.TypedData {
	chainId, err := strconv.ParseInt(action["signatureChainId"].(string), 16, 64)
	if err != nil {
		log.Error(err)
		return apitypes.TypedData{}
	}
	return apitypes.TypedData{
		Domain: apitypes.TypedDataDomain{
			Name:              "HyperliquidSignTransaction",
			Version:           "1",
			ChainId:           math.NewHexOrDecimal256(chainId),
			VerifyingContract: "0x0000000000000000000000000000000000000000",
		},
		Types: apitypes.Types{
			primaryType: payloadTypes,
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
		},
		PrimaryType: primaryType,
		Message:     action,
	}
}

// SignUserSignedAction signs "User Signed" actions like Withdraw, Transfer
func SignUserSignedAction(wallet *Wallet, action map[string]interface{}, payloadTypes []apitypes.Type, primaryType string) (Signature, error) {
	// 1. Prepare Action Data
	action["signatureChainId"] = "0x66eee" // 421614

	isMainNet := NowNetType == MAIN_NET
	if isMainNet {
		action["hyperliquidChain"] = "Mainnet"
	} else {
		action["hyperliquidChain"] = "Testnet"
	}

	// 2. Construct EIP-712 Typed Data
	data := UserSignedPayload(action, payloadTypes, primaryType)

	return SignInner(wallet, data)
}

var USD_SEND_SIGN_TYPES = []apitypes.Type{
	{Name: "hyperliquidChain", Type: "string"},
	{Name: "destination", Type: "string"},
	{Name: "amount", Type: "string"},
	{Name: "time", Type: "uint64"},
}

var SPOT_TRANSFER_SIGN_TYPES = []apitypes.Type{
	{Name: "hyperliquidChain", Type: "string"},
	{Name: "destination", Type: "string"},
	{Name: "token", Type: "string"},
	{Name: "amount", Type: "string"},
	{Name: "time", Type: "uint64"},
}

var WITHDRAW_SIGN_TYPES = []apitypes.Type{
	{Name: "hyperliquidChain", Type: "string"},
	{Name: "destination", Type: "string"},
	{Name: "amount", Type: "string"},
	{Name: "time", Type: "uint64"},
}

func SignUsdTransferAction(wallet *Wallet, action map[string]interface{}) (Signature, error) {
	return SignUserSignedAction(wallet, action, USD_SEND_SIGN_TYPES, "HyperliquidTransaction:UsdSend")
}

func SignSpotTransferAction(wallet *Wallet, action map[string]interface{}) (Signature, error) {
	return SignUserSignedAction(wallet, action, SPOT_TRANSFER_SIGN_TYPES, "HyperliquidTransaction:SpotSend")
}

func SignWithdrawFromBridgeAction(wallet *Wallet, action map[string]interface{}) (Signature, error) {
	return SignUserSignedAction(wallet, action, WITHDRAW_SIGN_TYPES, "HyperliquidTransaction:Withdraw")
}
