package enclaveguardian

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/common"
	"sync"
)

// This state machine compares the state of the enclave to the state of the world and is used to determine what actions can be taken with the enclave
// It records the last known status code of the enclave. It also records the l1 head and the l2 head that it believes the
// enclave has processed, optimistically updating these after successful actions and verifying the status when errors occur.

// Usage notes:
// - The status is updated by the host when the enclave successfully processed blocks and batches
// - The status is updated when we receive a status from the enclave
// - The status is **not** updated immediately when it receives blocks/batches from the outside world
// - The state should be notified of a live block/batch arrival before notifying if it is successfully processed
// - If unexpected error occurs when interacting with the enclave, then status should be requested and this state updated with the result

// EnclaveStatus is the status of the enclave from the host's perspective (including what it knows of the outside world)
type EnclaveStatus int

const (
	// Live - enclave is up-to-date with known external data. It can process L1 and L2 blocks as they arrive and respond to requests.
	Live EnclaveStatus = iota
	// Disconnected - enclave is unreachable or not returning a valid status (this overrides state calculations)
	Disconnected
	// Unavailable - enclave responding with 'Unavailable' status code
	Unavailable
	// AwaitingSecret - enclave is waiting for host to request and provide secret
	AwaitingSecret
	// L1Catchup - enclave is behind on L1 data, host should submit L1 blocks to catch up
	L1Catchup
	// L2Catchup - enclave is behind on L2 data, host should request and submit L2 batches to catch up
	L2Catchup
)

func (es EnclaveStatus) String() string {
	return [...]string{"Live", "Disconnected", "Unavailable", "AwaitingSecret", "L1Catchup", "L2Catchup"}[es]
}

// EnclaveState is the state machine for the enclave
type EnclaveState struct {
	// status is the cached status of the enclave
	// It is a function of the properties below and recalculated when any of them change
	status EnclaveStatus

	// enclave states (updated when enclave returns Status and optimistically after successful actions)
	enclaveStatusCode common.StatusCode
	enclaveL1Head     gethcommon.Hash
	enclaveL2Head     gethcommon.Hash

	// latest seen heads of L1 and L2 chains from external sources
	hostL1Head gethcommon.Hash
	hostL2Head gethcommon.Hash

	m *sync.RWMutex
}

func (s *EnclaveState) GetState() EnclaveStatus {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.status
}

func (s *EnclaveState) OnProcessedBlock(enclL1Head gethcommon.Hash) {
	s.m.Lock()
	s.enclaveL1Head = enclL1Head
	s.m.Unlock()
	s.recalculateStatus()
}
func (s *EnclaveState) OnReceivedBlock(l1Head gethcommon.Hash) {
	s.m.Lock()
	defer s.m.Unlock()
	s.hostL1Head = l1Head
}

func (s *EnclaveState) OnProcessedBatch(enclL2Head gethcommon.Hash) {
	s.m.Lock()
	s.enclaveL2Head = enclL2Head
	s.m.Unlock()
	s.recalculateStatus()
}
func (s *EnclaveState) OnReceivedBatch(l2Head gethcommon.Hash) {
	s.m.Lock()
	defer s.m.Unlock()
	s.hostL2Head = l2Head
}

func (s *EnclaveState) OnSecretProvided() {
	s.m.Lock()
	if s.enclaveStatusCode == common.AwaitingSecret {
		s.enclaveStatusCode = common.Running
	}
	s.m.Unlock()
	s.recalculateStatus()
}

func (s *EnclaveState) OnEnclaveStatus(es common.Status) {
	s.m.Lock()
	s.enclaveStatusCode = es.StatusCode
	s.enclaveL1Head = es.L1Head
	s.enclaveL2Head = es.L2Head
	s.m.Unlock()

	s.recalculateStatus()
}

// OnDisconnected is called if the enclave is unreachable/not returning a valid Status
func (s *EnclaveState) OnDisconnected() {
	s.m.Lock()
	defer s.m.Unlock()
	s.status = Disconnected
}

// when enclave is operational, this method will update the status based on comparison of current chain heads with enclave heads
func (s *EnclaveState) recalculateStatus() {
	s.m.Lock()
	defer s.m.Unlock()
	switch s.enclaveStatusCode {
	case common.AwaitingSecret:
		s.status = AwaitingSecret
	case common.Unavailable:
		s.status = Unavailable
	case common.Running:
		if s.hostL1Head != s.enclaveL1Head {
			s.status = L1Catchup
			return
		}
		if s.hostL2Head != s.enclaveL2Head {
			s.status = L2Catchup
			return
		}
		s.status = Live
	}
}

// InSyncWithL1 returns true if the enclave is up-to-date with L1 data so can process L1 blocks as they arrive
func (s *EnclaveState) InSyncWithL1() bool {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.status == Live || s.status == L2Catchup
}

func (s *EnclaveState) GetEnclaveL1Head() gethcommon.Hash {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.enclaveL1Head
}

func (s *EnclaveState) GetEnclaveL2Head() gethcommon.Hash {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.enclaveL2Head
}
