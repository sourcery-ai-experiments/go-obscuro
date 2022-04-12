package simulation

import (
	"math/rand"
	"testing"
	"time"

	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	stats2 "github.com/obscuronet/obscuro-playground/integration/simulation/stats"

	"github.com/google/uuid"
)

// testSimulation encapsulates the shared logic for simulating and testing various types of nodes.
func testSimulation(t *testing.T, netw network.Network, params params.SimParams) {
	rand.Seed(time.Now().UnixNano())
	uuid.EnableRandPool()

	stats := stats2.NewStats(params.NumberOfNodes) // todo - temporary object used to collect metrics. Needs to be replaced with something better

	ethClients, obscuroNodes, p2pAddrs := netw.Create(params, stats)

	txInjector := NewTransactionInjector(params.NumberOfWallets, uint64(params.AvgBlockDurationUSecs), stats, ethClients, obscuroNodes)

	simulation := Simulation{
		EthClients:       ethClients,
		ObscuroNodes:     obscuroNodes,
		ObscuroP2PAddrs:  p2pAddrs,
		AvgBlockDuration: uint64(params.AvgBlockDurationUSecs),
		TxInjector:       txInjector,
		SimulationTime:   params.SimulationTime,
		Stats:            stats,
		Params:           &params,
	}

	// execute the simulation
	simulation.Start()

	// run tests
	checkNetworkValidity(t, &simulation)

	simulation.Stop()

	// generate and print the final stats
	t.Logf("Simulation results:%+v", NewOutputStats(&simulation))
	netw.TearDown()
}
