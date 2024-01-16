package viewingkey

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/ten-protocol/go-ten/go/wallet"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	EIP712Domain          = "EIP712Domain"
	EIP712Type            = "Authentication"
	EIP712DomainName      = "name"
	EIP712DomainVersion   = "version"
	EIP712DomainChainID   = "chainId"
	EIP712EncryptionToken = "Encryption Token"
	// EIP712EncryptionTokenV2 is used to support older versions of third party libraries
	// that don't have the support for spaces in type names
	EIP712EncryptionTokenV2  = "EncryptionToken"
	EIP712DomainNameValue    = "Ten"
	EIP712DomainVersionValue = "1.0"
	UserIDHexLength          = 40
	TenChainID               = 443
)

// EIP712EncryptionTokens is a list of all possible options for Encryption token name
var EIP712EncryptionTokens = [...]string{
	EIP712EncryptionToken,
	EIP712EncryptionTokenV2,
}

// ViewingKey encapsulates the signed viewing key for an account for use in encrypted communication with an enclave
type ViewingKey struct {
	Account    *gethcommon.Address // Account address that this Viewing Key is bound to - Users Pubkey address
	PrivateKey *ecies.PrivateKey   // ViewingKey private key to encrypt data to the enclave
	PublicKey  []byte              // ViewingKey public key in decrypt data from the enclave
	Signature  []byte              // ViewingKey public key signed by the Accounts Private key - Allows to retrieve the Account address
}

// GenerateViewingKeyForWallet takes an account wallet, generates a viewing key and signs the key with the acc's private key
// uses the same method of signature handling as Metamask/geth
func GenerateViewingKeyForWallet(wal wallet.Wallet) (*ViewingKey, error) {
	// generate an ECDSA key pair to encrypt sensitive communications with the obscuro enclave
	vk, err := crypto.GenerateKey()
	viewingPrivateKeyECIES := ecies.ImportECDSA(vk)
	if err != nil {
		return nil, err
	}

	// create encryptionToken and store it in the database with the private key
	ecdsaPublicKey := viewingPrivateKeyECIES.PublicKey.ExportECDSA()
	compressedPubKey := crypto.CompressPubkey(ecdsaPublicKey)
	encToken := CalculateUserID(compressedPubKey)
	// sign public key bytes with the wallet's private key
	signature, err := mmSignViewingKey(hex.EncodeToString(encToken), wal.PrivateKey())
	if err != nil {
		return nil, err
	}

	// encode public key as bytes
	viewingPubKeyBytes := crypto.CompressPubkey(&vk.PublicKey)

	accAddress := wal.Address()
	return &ViewingKey{
		Account:    &accAddress,
		PrivateKey: viewingPrivateKeyECIES,
		PublicKey:  viewingPubKeyBytes,
		Signature:  signature,
	}, nil
}

// mmSignViewingKey takes an encryptionToken and the private key for a wallet, it simulates the back-and-forth to
// MetaMask and returns the signature bytes to register with the enclave
func mmSignViewingKey(encryptionToken string, signerKey *ecdsa.PrivateKey) ([]byte, error) {
	signature, err := Sign(signerKey, encryptionToken)
	if err != nil {
		return nil, fmt.Errorf("failed to sign viewing key: %w", err)
	}

	// We have to transform the V from 0/1 to 27/28, and add the leading "0".
	signature[64] += 27
	signatureWithLeadBytes := append([]byte("0"), signature...)

	// this string encoded signature is what the wallet extension would receive after it is signed by metamask
	sigStr := hex.EncodeToString(signatureWithLeadBytes)
	// and then we extract the signature bytes in the same way as the wallet extension
	outputSig, err := hex.DecodeString(sigStr[2:])
	if err != nil {
		return nil, fmt.Errorf("failed to decode signature string: %w", err)
	}
	// This same change is made in geth internals, for legacy reasons to be able to recover the address:
	//	https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L452-L459
	outputSig[64] -= 27

	return outputSig, nil
}

// Sign takes a users Private key and signs the encryption token
func Sign(userPrivKey *ecdsa.PrivateKey, encryptionToken string) ([]byte, error) {
	messages, err := GenerateAuthenticationEIP712RawDataOptions(encryptionToken, TenChainID)
	if err != nil || len(messages) == 0 {
		return nil, err
	}
	signature, err := crypto.Sign(accounts.TextHash(messages[0]), userPrivKey)
	if err != nil {
		return nil, fmt.Errorf("unable to sign messages - %w", err)
	}
	return signature, nil
}

// getBytesFromTypedData creates EIP-712 compliant hash from typedData.
// It involves hashing the message with its structure, hashing domain separator,
// and then encoding both hashes with specific EIP-712 bytes to construct the final message format.
func getBytesFromTypedData(typedData apitypes.TypedData) ([]byte, error) {
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, err
	}
	// Create the domain separator hash for EIP-712 message context
	domainSeparator, err := typedData.HashStruct(EIP712Domain, typedData.Domain.Map())
	if err != nil {
		return nil, err
	}
	// Prefix domain and message hashes with EIP-712 version and encoding bytes
	rawData := append([]byte("\x19\x01"), append(domainSeparator, typedDataHash...)...)
	return rawData, nil
}

// GenerateAuthenticationEIP712RawDataOptions generates all the options or raw data messages (bytes)
// for an EIP-712 message used to authenticate an address with user
func GenerateAuthenticationEIP712RawDataOptions(userID string, chainID int64) ([][]byte, error) {
	if len(userID) != UserIDHexLength {
		return nil, fmt.Errorf("userID hex length must be %d, received %d", UserIDHexLength, len(userID))
	}
	encryptionToken := "0x" + userID

	domain := apitypes.TypedDataDomain{
		Name:    EIP712DomainNameValue,
		Version: EIP712DomainVersionValue,
		ChainId: (*math.HexOrDecimal256)(big.NewInt(chainID)),
	}

	typedDataList := make([]apitypes.TypedData, 0, len(EIP712EncryptionTokens))
	for _, encTokenName := range EIP712EncryptionTokens {
		message := map[string]interface{}{
			encTokenName: encryptionToken,
		}

		types := apitypes.Types{
			EIP712Domain: {
				{Name: EIP712DomainName, Type: "string"},
				{Name: EIP712DomainVersion, Type: "string"},
				{Name: EIP712DomainChainID, Type: "uint256"},
			},
			EIP712Type: {
				{Name: encTokenName, Type: "address"},
			},
		}

		newTypeElement := apitypes.TypedData{
			Types:       types,
			PrimaryType: EIP712Type,
			Domain:      domain,
			Message:     message,
		}
		typedDataList = append(typedDataList, newTypeElement)
	}

	rawDataOptions := make([][]byte, 0, len(typedDataList))
	for _, typedDataItem := range typedDataList {
		rawData, err := getBytesFromTypedData(typedDataItem)
		if err != nil {
			return nil, err
		}
		rawDataOptions = append(rawDataOptions, rawData)
	}
	return rawDataOptions, nil
}

// CalculateUserIDHex CalculateUserID calculates userID from a public key
// (we truncate it, because we want it to have length 20) and encode to hex strings
func CalculateUserIDHex(publicKeyBytes []byte) string {
	return hex.EncodeToString(CalculateUserID(publicKeyBytes))
}

// CalculateUserID calculates userID from a public key (we truncate it, because we want it to have length 20)
func CalculateUserID(publicKeyBytes []byte) []byte {
	return crypto.Keccak256Hash(publicKeyBytes).Bytes()[:20]
}

// CheckSignatureAndAddress checks if the signature is valid for hash of the message and checks if
// signer is an address provided to the function.
// It returns true if both conditions are true and false otherwise
func CheckSignatureAndAddress(hashBytes []byte, signature []byte, address *gethcommon.Address) bool {
	hash := gethcommon.BytesToHash(hashBytes)
	pubKeyBytes, err := crypto.Ecrecover(hash[:], signature)
	if err != nil {
		return false
	}

	pubKey, err := crypto.UnmarshalPubkey(pubKeyBytes)
	if err != nil {
		return false
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	if !bytes.Equal(recoveredAddr.Bytes(), address.Bytes()) {
		return false
	}

	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:64])

	// Verify the signature and return the result (all the checks above passed)
	return ecdsa.Verify(pubKey, hashBytes, r, s)
}

func VerifySignatureEIP712(userID string, address *gethcommon.Address, signature []byte, chainID int64) (bool, error) {
	var rawDataOptions [][]byte

	rawDataOptions, err := GenerateAuthenticationEIP712RawDataOptions(userID, chainID)
	if err != nil {
		return false, err
	}

	if len(signature) != 65 {
		return false, fmt.Errorf("invalid signaure length: %d", len(signature))
	}

	// We transform the V from 27/28 to 0/1. This same change is made in Geth internals, for legacy reasons to be able
	// to recover the address: https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L452-L459
	if signature[64] == 27 || signature[64] == 28 {
		signature[64] -= 27
	}

	for _, rawData := range rawDataOptions {
		// create a hash of structured message (needed for signature verification)
		hashBytes := crypto.Keccak256(rawData)

		// current signature is valid - return true
		if CheckSignatureAndAddress(hashBytes, signature, address) {
			return true, nil
		}
	}
	return false, errors.New("signature verification failed")
}
