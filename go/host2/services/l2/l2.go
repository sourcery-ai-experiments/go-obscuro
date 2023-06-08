package l2

import (
	"fmt"
	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
	hostcommon "github.com/obscuronet/go-obscuro/go/common/host"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/host/db"
	"github.com/pkg/errors"
)

// service that provides host-level L2 data functionality. Including:
// - ensuring batch data is complete
// - streaming new batches to subscribers (e.g. enclave guardians)
// - providing missing batch lookup for enclaves and for peers
// - receiving txs from p2p and passing them on to subscribers (e.g. enclave guardians)

type Service struct {
	l2Subscribers []DataSubscriber
	p2p           hostcommon.P2P
	db            *db.DB
	logger        gethlog.Logger
}

func NewL2Service(config *config.HostConfig, p2p hostcommon.P2P, database *db.DB, logger gethlog.Logger) *Service {
	return &Service{
		p2p:    p2p,
		db:     database,
		logger: logger,
	}
}

func (s *Service) Start() error {
	s.p2p.StartListening(s)
	return nil
}

func (s *Service) Stop() error {
	return nil
}

// DataSubscriber is an interface that can be implemented by any service that wants to be notified of new blocks as they are created.
// Note: new head will not necessarily be the following block from the previous head as their may have been a reorg.
type DataSubscriber interface {
	OnNewBatch(block *common.ExtBatch)
	OnTxGossip(tx common.EncryptedTx)
}

func (s *Service) Subscribe(br DataSubscriber) {
	s.l2Subscribers = append(s.l2Subscribers, br)
}

// FetchNextBatch returns the next batch that an enclave should be fed if it thinks `fromHead` is the current head.
// This might be the child batch of `fromHead` or it might be an earlier ancestor if there has been a reorg.
func (s *Service) FetchNextBatch(fromHead gethcommon.Hash) (*common.ExtBatch, error) {
	batch, err := s.db.GetBatch(fromHead)
	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve batch")
	}

	canonicalBatchHashAtSameHeight, err := s.db.GetBatchHash(batch.Header.Number)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve canonical batch hash, height=%d - %w", batch.Header.Number, err)
	}

	// If the batch's hash does not match the canonical batch's hash at the same height, we need to keep walking back.
	if batch.Hash() != *canonicalBatchHashAtSameHeight {
		return s.FetchNextBatch(batch.Header.ParentHash)
	}
	return batch, nil
}

// The following methods are called by the p2p service when it receives messages from other nodes

// ReceiveTx is called whenever a new transaction has arrived from p2p
func (s *Service) ReceiveTx(tx common.EncryptedTx) {
	for _, sub := range s.l2Subscribers {
		// notify subscribers in a new goroutine, so we don't block. Up to subscribers to handle concurrency
		go sub.OnTxGossip(tx)
	}
}

// ReceiveBatches is called whenever batches have arrived from p2p
func (s *Service) ReceiveBatches(encBatchMsg common.EncodedBatchMsg) {
	var batchMsg *hostcommon.BatchMsg
	err := rlp.DecodeBytes(encBatchMsg, &batchMsg)
	if err != nil {
		s.logger.Warn("unable to decode batch message", "err", err)
		return
	}

	// we don't currently do any checks at this stage, just add them to DB/cache and make them available to enclaves
	for _, b := range batchMsg.Batches {
		err := s.db.AddBatch(b)
		if err != nil {
			s.logger.Warn("unable to add batch to db", "err", err)
			return // probably don't want to continue with other batches from same msg after this fails
		}

		if !batchMsg.IsCatchUp {
			// don't notify subscribers for catch-up
			s.notifySubscribers(b)
		}
	}

}

// ReceiveBatchRequest is called whenever a request for batches has arrived from p2p
func (s *Service) ReceiveBatchRequest(batchRequest common.EncodedBatchRequest) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) notifySubscribers(b *common.ExtBatch) {
	for _, sub := range s.l2Subscribers {
		// notify subscribers in a new goroutine, so we don't block. Up to subscribers to handle concurrency
		go sub.OnNewBatch(b)
	}
}
