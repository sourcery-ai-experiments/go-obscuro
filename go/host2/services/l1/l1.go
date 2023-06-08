package l1

import (
	"fmt"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/retry"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
	"github.com/obscuronet/go-obscuro/go/host/db"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/pkg/errors"
	"time"
)

// service that provides all L1 interactions required by the host (block streaming, rollup/secret publishing, etc)

const (
	// Attempts to send secret initialisation, request or response transactions to the L1. Worst-case, equates to 63 seconds, plus time per request.
	l1TxTriesSecret           = 7
	maxWaitForL1Receipt       = 100 * time.Second
	retryIntervalForL1Receipt = 10 * time.Second
)

var (
	ErrNoNextBlock = errors.New("no next block")
)

type Service struct {
	subscribers     []BlockReceiver
	mgmtContractLib mgmtcontractlib.MgmtContractLib
	db              *db.DB
	hostAddress     string // identifying address (p2p) for the host
	ethWallet       wallet.Wallet
	ethClient       ethadapter.EthClient

	logger    gethlog.Logger
	hostID    gethcommon.Address
	isGenesis bool
}

func NewL1Service(config *config.HostConfig, client ethadapter.EthClient, wallet wallet.Wallet, mgmtContractLib mgmtcontractlib.MgmtContractLib, database *db.DB, logger gethlog.Logger) *Service {
	return &Service{
		ethClient:       client,
		ethWallet:       wallet,
		mgmtContractLib: mgmtContractLib,
		db:              database,
		logger:          logger,

		hostID:    config.ID,
		isGenesis: config.IsGenesis,
	}
}

func (s *Service) Start() error {
	// start streaming blocks from the client

	// when a new block is received, call OnNewHead for each subscriber
	return nil
}

func (s *Service) Stop() error {
	return nil
}

func (s *Service) InitializeSecret(id gethcommon.Address, attestation common.EncodedAttestationReport, secret common.EncryptedSharedEnclaveSecret) error {
	l1Tx := &ethadapter.L1InitializeSecretTx{
		AggregatorID:  &id,
		InitialSecret: secret,
		HostAddress:   s.hostAddress,
		Attestation:   attestation,
	}
	var err error
	initialiseSecretTx := s.mgmtContractLib.CreateInitializeSecret(l1Tx, s.ethWallet.GetNonceAndIncrement())
	initialiseSecretTx, err = s.ethClient.EstimateGasAndGasPrice(initialiseSecretTx, s.ethWallet.Address())
	if err != nil {
		s.ethWallet.SetNonce(s.ethWallet.GetNonce() - 1)
		return errors.Wrap(err, "could not estimate gas for initialize secret tx")
	}
	// we block here until we confirm a successful receipt. It is important this is published before the initial rollup.
	err = s.signAndBroadcastL1Tx(initialiseSecretTx, l1TxTriesSecret, true)
	if err != nil {
		return fmt.Errorf("failed to initialise enclave secret. Cause: %w", err)
	}
	s.logger.Info("Node is genesis node. Secret was broadcast.")
	return nil
}

func (s *Service) RequestSecret(enclaveID gethcommon.Address, attestation common.EncodedAttestationReport) (gethcommon.Hash, error) {
	return gethcommon.Hash{}, nil
}

// FetchNextBlock returns block, isLatest bool, error. If prevBlock is current L1 head then ErrNoNextBlock is returned.
func (s *Service) FetchNextBlock(prevBlock gethcommon.Hash) (*common.L1Block, bool, error) {
	return nil, false, nil
}

// BlockReceiver is an interface that can be implemented by any service that wants to be notified of new blocks as they are created.
// Note: new head will not necessarily be the following block from the previous head as their may have been a reorg.
type BlockReceiver interface {
	OnNewL1Block(block *types.Block)
}

func (s *Service) Subscribe(br BlockReceiver) {
	// todo: return a close subscription function based on a sub ID
	s.subscribers = append(s.subscribers, br)
}

// `tries` is the number of times to attempt broadcasting the transaction.
// if awaitReceipt is true then this method will block and synchronously wait to check the receipt, otherwise it is fire
// and forget and the receipt tracking will happen in a separate go-routine
func (s *Service) signAndBroadcastL1Tx(tx types.TxData, tries uint64, awaitReceipt bool) error {
	var err error
	tx, err = s.ethClient.EstimateGasAndGasPrice(tx, s.ethWallet.Address())
	if err != nil {
		return fmt.Errorf("unable to estimate gas limit and gas price - %w", err)
	}

	signedTx, err := s.ethWallet.SignTransaction(tx)
	if err != nil {
		return err
	}

	s.logger.Info("Host issuing l1 tx", log.TxKey, signedTx.Hash(), "size", signedTx.Size()/1024)

	err = retry.Do(func() error {
		return s.ethClient.SendTransaction(signedTx)
	}, retry.NewDoublingBackoffStrategy(time.Second, tries)) // doubling retry wait (3 tries = 7sec, 7 tries = 63sec)
	if err != nil {
		return fmt.Errorf("broadcasting L1 transaction failed after %d tries. Cause: %w", tries, err)
	}
	s.logger.Info("Successfully issued Rollup on L1", "txHash", signedTx.Hash())

	if awaitReceipt {
		// block until receipt is found and then return
		return s.waitForReceipt(signedTx.Hash())
	}

	// else just watch for receipt asynchronously and log if it fails
	go func() {
		// todo (#1624) - consider how to handle the various ways that L1 transactions could fail to improve node operator QoL
		err = s.waitForReceipt(signedTx.Hash())
		if err != nil {
			s.logger.Error("L1 transaction failed", log.ErrKey, err)
		}
	}()

	return nil
}
func (s *Service) waitForReceipt(txHash common.TxHash) error {
	var receipt *types.Receipt
	var err error
	err = retry.Do(
		func() error {
			receipt, err = s.ethClient.TransactionReceipt(txHash)
			if err != nil {
				// adds more info on the error
				return fmt.Errorf("unable to get receipt for tx: %s - %w", txHash.Hex(), err)
			}
			return err
		},
		retry.NewTimeoutStrategy(maxWaitForL1Receipt, retryIntervalForL1Receipt),
	)
	if err != nil {
		return fmt.Errorf("receipt for L1 transaction never found despite 'successful' broadcast - %w", err)
	}

	if err == nil && receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("unsuccessful receipt found for published L1 transaction, status=%d", receipt.Status)
	}
	s.logger.Trace("Successful L1 transaction receipt found.", log.BlockHeightKey, receipt.BlockNumber, log.BlockHashKey, receipt.BlockHash)
	return nil
}

func (s *Service) ExtractSecretResponses(fromBlock *types.Block) ([]ethadapter.L1RespondSecretTx, error) {
	var secretResponses []ethadapter.L1RespondSecretTx
	for _, tx := range fromBlock.Transactions() {
		t := s.mgmtContractLib.DecodeTx(tx)
		if t == nil {
			continue
		}
		if scrtTx, ok := t.(*ethadapter.L1RespondSecretTx); ok {
			secretResponses = append(secretResponses, *scrtTx)
		}
	}
	return secretResponses, nil
}

func (s *Service) FetchReceipts(block *common.L1Block) types.Receipts {
	receipts := make(types.Receipts, 0)

	for _, transaction := range block.Transactions() {
		receipt, err := s.ethClient.TransactionReceipt(transaction.Hash())

		if err != nil || receipt == nil {
			s.logger.Error("Problem with retrieving the receipt on the host!", log.ErrKey, err, log.CmpKey, log.CrossChainCmp)
			continue
		}

		s.logger.Trace("Adding receipt", "status", receipt.Status, log.TxKey, transaction.Hash(),
			log.BlockHashKey, block.Hash(), log.CmpKey, log.CrossChainCmp)

		receipts = append(receipts, receipt)
	}

	return receipts
}

func (s *Service) PublishSharedSecretResponses(scrtResponses []*common.ProducedSecretResponse) error {
	var err error

	for _, scrtResponse := range scrtResponses {
		// todo (#1624) - implement proper protocol so only one host responds to this secret requests initially
		// 	for now we just have the genesis host respond until protocol implemented
		if !s.isGenesis {
			s.logger.Trace("Not genesis node, not publishing response to secret request.",
				"requester", scrtResponse.RequesterID)
			return nil
		}

		l1tx := &ethadapter.L1RespondSecretTx{
			Secret:      scrtResponse.Secret,
			RequesterID: scrtResponse.RequesterID,
			AttesterID:  s.hostID,
			HostAddress: scrtResponse.HostAddress,
		}
		// todo (#1624) - l1tx.Sign(a.attestationPubKey) doesn't matter as the waitSecret will process a tx that was reverted
		respondSecretTx := s.mgmtContractLib.CreateRespondSecret(l1tx, s.ethWallet.GetNonceAndIncrement(), false)
		respondSecretTx, err = s.ethClient.EstimateGasAndGasPrice(respondSecretTx, s.ethWallet.Address())
		if err != nil {
			s.ethWallet.SetNonce(s.ethWallet.GetNonce() - 1)
			return err
		}
		s.logger.Trace("Broadcasting secret response L1 tx.", "requester", scrtResponse.RequesterID)
		// fire-and-forget (track the receipt asynchronously)
		err = s.signAndBroadcastL1Tx(respondSecretTx, l1TxTriesSecret, false)
		if err != nil {
			return fmt.Errorf("could not broadcast secret response. Cause %w", err)
		}
	}
	return nil
}
