package simulation

import (
	"testing"
	"time"

	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"

	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
)

// This test creates a network of L2 nodes, then injects transactions, and finally checks the resulting output blockchain
// The L2 nodes communicate with each other via sockets, and with their enclave servers via RPC.
// All nodes and enclaves live in the same process, and the Ethereum nodes are mocked out.
func TestSocketNodesMonteCarloSimulation(t *testing.T) {
	setupTestLog()

	simParams := params.SimParams{
		NumberOfNodes:             10,
		NumberOfWallets:           5,
		AvgBlockDuration:          250 * time.Millisecond,
		SimulationTime:            25 * time.Second,
		L1EfficiencyThreshold:     0.2,
		L2EfficiencyThreshold:     0.3,
		L2ToL1EfficiencyThreshold: 0.4,
		TxHandler:                 ethereum_mock.NewMockTxHandler(),
	}
	simParams.AvgNetworkLatency = simParams.AvgBlockDuration / 15
	simParams.AvgGossipPeriod = simParams.AvgBlockDuration / 4

	testSimulation(t, network.NewBasicNetworkOfSocketNodes(), simParams)
}
