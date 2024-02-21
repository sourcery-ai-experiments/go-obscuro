package vkhandler

import (
	"crypto/rand"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/accounts"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"gitlab.com/NebulousLabs/fastrand"

	"github.com/ethereum/go-ethereum/crypto/ecies"
)

// Used when the result to an eth_call is equal to nil. Attempting to encrypt then decrypt nil using ECIES throws an exception.
var placeholderResult = []byte("0x")

// AuthenticatedViewingKey - the enclave side of the viewing key. Used for authenticating requests and for encryption
type AuthenticatedViewingKey struct {
	rpcVK          *viewingkey.RPCSignedViewingKey
	AccountAddress *gethcommon.Address
	ecdsaKey       *ecies.PublicKey
	UserID         string
}

func VerifyViewingKey(rpcVK *viewingkey.RPCSignedViewingKey, chainID int64) (*AuthenticatedViewingKey, error) {
	vkPubKey, err := crypto.DecompressPubkey(rpcVK.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("could not decompress viewing key bytes - %w", err)
	}

	rvk := &AuthenticatedViewingKey{
		AccountAddress: rpcVK.Account,
		rpcVK:          rpcVK,
		ecdsaKey:       ecies.ImportECDSAPublic(vkPubKey),
	}

	// 2. Authenticate
	recoveredAccountAddress, err := checkViewingKeyAndRecoverAddress(rvk, chainID)
	if err != nil {
		return nil, err
	}

	rvk.AccountAddress = recoveredAccountAddress
	return rvk, nil
}

// this method is unnecessarily complex due to a legacy signing format
func checkViewingKeyAndRecoverAddress(vk *AuthenticatedViewingKey, chainID int64) (*gethcommon.Address, error) {
	// get userID from viewingKey public key
	userID := viewingkey.CalculateUserIDHex(vk.rpcVK.PublicKey)
	vk.UserID = userID

	// check signature and check if it matches the account address (assuming the message was signed with EIP712) and return address if it does
	isValidAndMatchesAddress := viewingkey.CheckIfSignatureIsValidAndMatchesAddress(userID, vk.rpcVK.SignatureWithAccountKey, chainID, vk.AccountAddress)
	if isValidAndMatchesAddress {
		return vk.AccountAddress, nil
	}

	// TODO @Ziga - this must be removed once the legacy format is no longer supported
	// signature is either invalid or it might have been using the legacy format

	// check signature and recover the address assuming the message was signed with EIP712
	legacyMessageHash := accounts.TextHash([]byte(viewingkey.GenerateSignMessage(vk.rpcVK.PublicKey)))
	legacyMessageSignerAddress, err := viewingkey.CheckSignatureAndReturnAccountAddress(legacyMessageHash, vk.rpcVK.SignatureWithAccountKey)
	if err != nil {
		return nil, fmt.Errorf("invalid vk signature")
	}

	if legacyMessageSignerAddress.Hex() != vk.AccountAddress.Hex() {
		return nil, fmt.Errorf("invalid VK")
	}

	return vk.AccountAddress, err
}

// crypto.rand is quite slow. When this variable is true, we will use a fast CSPRNG algorithm
const useFastRand = true

func rndSource() io.Reader {
	rndSource := rand.Reader
	if useFastRand {
		rndSource = fastrand.Reader
	}
	return rndSource
}

// Encrypt returns the payload encrypted with the viewingKey
func (vk *AuthenticatedViewingKey) Encrypt(bytes []byte) ([]byte, error) {
	if len(bytes) == 0 {
		bytes = placeholderResult
	}
	encryptedBytes, err := ecies.Encrypt(rndSource(), vk.ecdsaKey, bytes, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to encrypt with given public VK - %w", err)
	}

	return encryptedBytes, nil
}
