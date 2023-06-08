package enclaveguardian

import (
	"fmt"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/retry"
	"github.com/obscuronet/go-obscuro/go/host/db"
	"github.com/obscuronet/go-obscuro/go/host2/services/l1"
	"github.com/obscuronet/go-obscuro/go/host2/services/l2"
	"github.com/pkg/errors"
	"sync"
	"sync/atomic"
	"time"
)

const (
	_reconnectInterval  = 1 * time.Second
	_monitoringInterval = 1 * time.Second
)

// Service is a service that monitors an enclave container, feeds it data and provides access to it from other services
type Service struct {
	l1 *l1.Service
	l2 *l2.Service
	db *db.DB

	isGenesisEnclave bool
	isSequencer      bool
	enclaveClient    common.Enclave
	logger           gethlog.Logger

	// state
	running atomic.Bool
	state   EnclaveState

	enclaveID gethcommon.Address

	// synchronization
	// todo: optimisation required on these locks, just put them in to avoid risk of races for now
	l1Lock sync.Mutex
	l2Lock sync.Mutex
}

func NewEnclaveGuardian(l1 *l1.Service, l2 *l2.Service, enclaveClient common.Enclave, db *db.DB) *Service {
	return &Service{l1: l1, l2: l2, enclaveClient: enclaveClient, db: db}
}

func (g *Service) Start() error {
	g.running.Store(true)

	// subscribe for data updates for L1 and L2
	g.l1.Subscribe(g)
	g.l2.Subscribe(g)

	go g.mainLoop()

	return nil
}

func (g *Service) Stop() error {
	g.running.Store(false)
	return nil
}

// mainLoop runs until the enclave guardian is stopped. It checks the state of the enclave and takes action as
// required to improve the state (e.g. provide a secret, catch up with L1, etc.)
func (g *Service) mainLoop() {
	for g.running.Load() {
		switch g.state.GetState() {
		case Disconnected, Unavailable:
			time.Sleep(_reconnectInterval)
			g.checkEnclaveStatus()
		case AwaitingSecret:
			err := g.provideSecret()
			if err != nil {
				g.logger.Warn("could not provide secret to enclave", "err", err)
			}
		case L1Catchup:
			err := g.catchupWithL1()
			if err != nil {
				g.logger.Warn("could not catch up with L1", "err", err)
			}
		case L2Catchup:
			err := g.catchupWithL2()
			if err != nil {
				g.logger.Warn("could not catch up with L2", "err", err)
			}
		case Live:
			// todo: should we allow interrupt here so we try to recover from a change in state immediately? Would allow for longer monitoring interval if we also interrupt for Stop().
			time.Sleep(_monitoringInterval)
		}
	}
}

// OnNewL1Block is called by the L1 service when a new L1 block is received. It must only be called by one goroutine at a time.
func (g *Service) OnNewL1Block(block *types.Block) {
	// record the new L1 head
	g.state.OnReceivedBlock(block.Hash())
	if !g.state.InSyncWithL1() {
		// we are not ready to submit the latest L1 blocks to the enclave
		return
	}

	err := g.submitL1Block(block, true)
	if err != nil {
		g.logger.Warn("could not submit L1 block to enclave", "err", err)
	}
}

// OnNewBatch is called by the p2p service when a new L2 block is received. It must only be called by one goroutine at a time.
func (g *Service) OnNewBatch(batch *common.ExtBatch) {
	// record the new L2 head
	g.state.OnReceivedBatch(batch.Hash())
	if g.state.GetState() != Live {
		// we are not ready to submit the latest L2 blocks to the enclave
		return
	}
	err := g.submitL2Batch(batch)
	if err != nil {
		g.logger.Warn("could not submit L2 batch to enclave", "err", err)
	}
}

func (g *Service) OnTxGossip(tx common.EncryptedTx) {
	_, err := g.enclaveClient.SubmitTx(tx)
	if err != nil {
		g.logger.Warn("failed to submit gossiped tx to enclave", log.ErrKey, err)
	}
}

func (g *Service) fetchEncodedAttestation() (common.EncodedAttestationReport, error) {
	attestation, err := g.enclaveClient.Attestation()
	if err != nil {
		// something went wrong, check enclave status and let main loop try again as appropriate
		g.checkEnclaveStatus()
		return nil, errors.Wrap(err, "could not get attestation from enclave")
	}
	if attestation.Owner != g.enclaveID {
		return nil, fmt.Errorf("genesis node has ID %s, but its enclave produced an attestation using ID %s", g.enclaveID.Hex(), attestation.Owner.Hex())
	}

	encodedAttestation, err := common.EncodeAttestation(attestation)
	if err != nil {
		return nil, errors.Wrap(err, "could not encode attestation")
	}
	return encodedAttestation, nil
}

func (g *Service) provideSecret() error {
	// whether we are generating or requesting the secret, we require an attestation from the enclave
	encodedAttestation, err := g.fetchEncodedAttestation()
	if err != nil {
		return errors.Wrap(err, "could not fetch encoded attestation")
	}
	if g.isGenesisEnclave {
		err = g.generateAndPublishSecret(encodedAttestation)
	} else {
		err = g.requestAndAwaitSecret(encodedAttestation)
	}
	if err != nil {
		return err
	}

	// completed successfully, update state
	g.state.OnSecretProvided()
	return nil
}

func (g *Service) generateAndPublishSecret(encodedAttestation common.EncodedAttestationReport) error {
	g.logger.Info("Genesis node: generating secret")

	// Create the network shared secret
	secret, err := g.enclaveClient.GenerateSecret()
	if err != nil {
		// something went wrong, check the enclave status and let the main loop try again when appropriate
		g.checkEnclaveStatus()
		return errors.Wrap(err, "could not generate secret")
	}

	// Publish the secret to the L1 management contract
	if err := g.l1.InitializeSecret(g.enclaveID, encodedAttestation, secret); err != nil {
		return errors.Wrap(err, "could not initialize secret")
	}

	g.logger.Info("Genesis node: secret generated and published")
	return nil
}

func (g *Service) requestAndAwaitSecret(encodedAttestation common.EncodedAttestationReport) error {
	g.logger.Info("Requesting secret")

	// Request the shared secret from the L1 management contract

	// request secret returns the L1 block containing the request, we check L1 blocks from there onwards for response
	awaitFromBlock, err := g.l1.RequestSecret(g.enclaveID, encodedAttestation)
	if err != nil {
		return errors.Wrap(err, "unable to request secret from L1")
	}

	err = retry.Do(func() error {
		nextBlock, _, err := g.l1.FetchNextBlock(awaitFromBlock)
		if err != nil {
			return fmt.Errorf("next block after block=%s not found - %w", awaitFromBlock, err)
		}
		secretRespTxs, err := g.l1.ExtractSecretResponses(nextBlock)
		if err != nil {
			return fmt.Errorf("could not extract secret responses from block=%s - %w", nextBlock.Hash(), err)
		}
		for _, s := range secretRespTxs {
			if s.RequesterID == g.enclaveID {
				err = g.enclaveClient.InitEnclave(s.Secret)
				if err != nil {
					g.logger.Warn("could not initialize enclave with received secret", "err", err)
					continue // try the next secret in the block if there are more
				}
				return nil // successfully initialized enclave with secret
			}
		}
		awaitFromBlock = nextBlock.Hash()
		return errors.New("no valid secret received in block")
		// todo @matt §§§ make these times constants or configurable
	}, retry.NewTimeoutStrategy(60*time.Second, 500*time.Millisecond))
	if err != nil {
		// something went wrong, check the enclave status in case it is an enclave problem and let the main loop try again when appropriate
		g.checkEnclaveStatus()
		return errors.Wrap(err, "no valid secret received for enclave")
	}

	g.logger.Info("Secret received")
	return nil
}

func (g *Service) catchupWithL1() error {
	// while we are behind the L1 head, fetch and submit L1 blocks
	for g.state.GetState() == L1Catchup {
		l1Block, isLatest, err := g.l1.FetchNextBlock(g.state.GetEnclaveL1Head())
		if err != nil {
			if err == l1.ErrNoNextBlock {
				return nil // we are up-to-date
			}
			return errors.Wrap(err, "could not fetch next L1 block")
		}
		err = g.submitL1Block(l1Block, isLatest)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Service) submitL1Block(block *common.L1Block, isLatest bool) error {
	receipts := g.l1.FetchReceipts(block)
	g.l1Lock.Lock()
	resp, err := g.enclaveClient.SubmitL1Block(*block, receipts, isLatest)
	g.l1Lock.Unlock()
	if err != nil {
		// something went wrong, check enclave status and let main loop try again as appropriate
		g.checkEnclaveStatus()
		return errors.Wrap(err, "could not submit L1 block to enclave")
	}
	// successfully processed block, update the state
	g.state.OnProcessedBlock(block.Hash())
	err = g.db.AddBlockHeader(block.Header())
	if err != nil {
		return fmt.Errorf("submitted block to enclave but could not store the block processing result. Cause: %w", err)
	}

	// todo: make sure this doesn't respond to old requests
	err = g.l1.PublishSharedSecretResponses(resp.ProducedSecretResponses)
	if err != nil {
		g.logger.Error("failed to publish response to secret request", log.ErrKey, err)
	}
	return nil
}

func (g *Service) catchupWithL2() error {
	// while we are behind the L2 head, fetch and submit L2 batches
	for g.state.GetState() == L2Catchup {
		l2Batch, err := g.l2.FetchNextBatch(g.state.GetEnclaveL2Head())
		if err != nil {
			return errors.Wrap(err, "could not fetch next L2 batch")
		}
		err = g.submitL2Batch(l2Batch)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Service) submitL2Batch(batch *common.ExtBatch) error {
	g.l2Lock.Lock()
	err := g.enclaveClient.SubmitBatch(batch)
	g.l2Lock.Unlock()
	if err != nil {
		// something went wrong, check enclave status and let main loop try again as appropriate
		g.checkEnclaveStatus()
		return errors.Wrap(err, "could not submit L2 batch to enclave")
	}
	// successfully processed batch, update the state
	g.state.OnProcessedBatch(batch.Hash())
	err = g.db.AddBatch(batch)
	if err != nil {
		return fmt.Errorf("submitted batch to enclave but could not store the batch processing result. Cause: %w", err)
	}
	return nil
}

// requests the status from the enclave and updates the state with the response. Called whenever enclave errors occur suggesting our optimistic view of the enclave state is out-of-date.
func (g *Service) checkEnclaveStatus() {
	status, err := g.enclaveClient.Status()
	if err != nil {
		g.state.OnDisconnected()
		g.logger.Warn("enclave status request failed", log.ErrKey, err)
		return
	}
	g.state.OnEnclaveStatus(status)
}

func (g *Service) GetClientIfHealthy() (common.Enclave, error) {
	if g.state.GetState() != Live {
		return nil, fmt.Errorf("enclave is not ready for requests, status=%s", g.state.GetState())
	}
	return g.enclaveClient, nil
}
