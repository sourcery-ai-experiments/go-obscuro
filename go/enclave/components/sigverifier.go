package components

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ten-protocol/go-ten/go/enclave/storage"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	// signature contains R, S and V (1 byte recovery ID)
	_ECDSASignatureLength = 65
)

type SignatureValidator struct {
	SequencerID gethcommon.Address
	attestedKey *ecdsa.PublicKey
	storage     storage.Storage
}

func NewSignatureValidator(seqID gethcommon.Address, storage storage.Storage) (*SignatureValidator, error) {
	// todo (#718) - sequencer identities should be retrieved from the L1 management contract
	return &SignatureValidator{
		SequencerID: seqID,
		storage:     storage,
		attestedKey: nil,
	}, nil
}

// CheckSequencerSignature - verifies the signature against the registered sequencer
func (sigChecker *SignatureValidator) CheckSequencerSignature(headerHash gethcommon.Hash, signature []byte) error {
	if signature == nil || len(signature) != _ECDSASignatureLength {
		return fmt.Errorf("missing signature on batch")
	}

	if sigChecker.attestedKey == nil {
		attestedKey, err := sigChecker.storage.FetchAttestedKey(sigChecker.SequencerID)
		if err != nil {
			return fmt.Errorf("could not retrieve attested key for aggregator %s. Cause: %w", sigChecker.SequencerID, err)
		}
		sigChecker.attestedKey = attestedKey
	}

	if !crypto.VerifySignature(crypto.FromECDSAPub(sigChecker.attestedKey), headerHash.Bytes(), signature[:_ECDSASignatureLength-1]) {
		return fmt.Errorf("could not verify ECDSA signature")
	}
	return nil
}
