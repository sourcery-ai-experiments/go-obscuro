package p2p

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/common"
	hostcommon "github.com/obscuronet/go-obscuro/go/common/host"
)

type Service struct {
	p2p hostcommon.P2P
}

func (s *Service) FetchNextBatches(prevBatch gethcommon.Hash) ([]*common.ExtBatch, error) {
	request := &common.BatchRequest{
		Requester:        s.publicAddress,
		CurrentHeadBatch: &prevBatch,
	}
	err := s.p2p.RequestBatchesFromSequencer(request)
	if err != nil {
		return nil, err
	}

}
