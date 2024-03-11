package helpful

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ten-protocol/go-ten/integration/networktest/userwallet"
	"golang.org/x/net/context"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/wallet"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/integration/networktest/env"
)

const (
	_sepoliaChainID = 11155111

	SepoliaRPCAddress1 = "wss://sepolia.infura.io/ws/v3/<api-key>" // seq
	SepoliaRPCAddress2 = "wss://sepolia.infura.io/ws/v3/<api-key>" // val
	SepoliaRPCAddress3 = "wss://sepolia.infura.io/ws/v3/<api-key>" // tester

	_sepoliaSequencerPK  = "<pk>" // account 0x<acc>
	_sepoliaValidator1PK = "<pk>" // account 0x<acc>

	_tenChainID = 443
	_gatewayURL = "http://localhost:11180"
	_tenPK      = "<pk>"
)

func TestGatewayAddUser(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.EnsureTestLogsSetUp("local-gateway-user")
	tenWallet := wallet.NewInMemoryWalletFromConfig(_tenPK, _tenChainID, testlog.Logger())
	err := addWalletsToGateway(append([]wallet.Wallet{}, tenWallet))
	if err != nil {
		t.Fatal(err)
	}
}

func addWalletsToGateway(wallets []wallet.Wallet) error {
	for _, w := range wallets {
		_, err := userwallet.NewGatewayUser(w, _gatewayURL, testlog.Logger())
		if err != nil {
			return err
		}
		fmt.Println("Ten Sim Wallet:\n", "pk:", privateKeyToHex(w.PrivateKey()), "account:", w.Address().String())
	}
	return nil
}

func TestRunLocalNetwork(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.EnsureTestLogsSetUp("local-geth-network")
	networkConnector, cleanUp, err := env.LocalDevNetwork(env.WithTenGateway(), env.WithSimWallets(3)).Prepare()
	if err != nil {
		t.Fatal(err)
	}
	wallets, err := networkConnector.GetNetworkWallets()
	ethWallets := wallets.AllEthWallets()
	obsWallets := wallets.AllObsWallets()
	simWallets := wallets.SimObsWallets

	if len(_tenPK) == 64 { // a pk const is entered so create wallet and append to sims
		tenWallet := wallet.NewInMemoryWalletFromConfig(_tenPK, _tenChainID, testlog.Logger())
		simWallets = append(simWallets, tenWallet)
		// fund this wallet
		err = networkConnector.AllocateFaucetFunds(context.Background(), tenWallet.Address())
		if err != nil {
			t.Log(err)
		}
	}

	err = addWalletsToGateway(simWallets)
	if err != nil {
		t.Fatal(err)
	}

	// Print the Ethereum and Obscuro wallet addresses
	for _, w := range ethWallets {
		fmt.Println("Ethereum Wallet:", w.Address().Hex())
	}
	for _, w := range obsWallets {
		fmt.Println("Obscuro Wallet:", w.Address().Hex())
	}
	defer cleanUp()

	keepRunning(networkConnector)
}

func TestRunLocalGatewayAgainstRemoteTestnet(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.EnsureTestLogsSetUp("local-faucet-remote-testnet")

	// set the testnet the gateway will connect to here
	netw := env.SepoliaTestnet(env.WithLocalTenGateway())
	networkConnector, cleanUp, err := netw.Prepare()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanUp()

	keepRunning(networkConnector)
}

func TestRunLocalNetworkAgainstSepolia(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.EnsureTestLogsSetUp("local-sepolia-network")

	l1DeployerWallet := wallet.NewInMemoryWalletFromConfig(_sepoliaSequencerPK, _sepoliaChainID, testlog.Logger())
	checkBalance("sequencer", l1DeployerWallet, SepoliaRPCAddress1)

	val1Wallet := wallet.NewInMemoryWalletFromConfig(_sepoliaValidator1PK, _sepoliaChainID, testlog.Logger())
	checkBalance("validator1", val1Wallet, SepoliaRPCAddress2)

	validatorWallets := []wallet.Wallet{val1Wallet}
	networktest.EnsureTestLogsSetUp("local-network-live-l1")
	networkConnector, cleanUp, err := env.LocalNetworkLiveL1(l1DeployerWallet, validatorWallets, []string{SepoliaRPCAddress1, SepoliaRPCAddress2}).Prepare()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanUp()

	keepRunning(networkConnector)
}

func checkBalance(walDesc string, wal wallet.Wallet, rpcAddress string) {
	client, err := ethadapter.NewEthClientFromURL(rpcAddress, 20*time.Second, common.HexToAddress("0x0"), testlog.Logger())
	if err != nil {
		panic("unable to create live L1 eth client, err=" + err.Error())
	}

	bal, err := client.BalanceAt(wal.Address(), nil)
	if err != nil {
		panic(fmt.Errorf("failed to get balance for %s (%s): %w", walDesc, wal.Address(), err))
	}
	fmt.Println(walDesc, "wallet balance", wal.Address(), bal)

	if bal.Cmp(big.NewInt(0)) <= 0 {
		panic(fmt.Errorf("%s wallet has no funds: %s", walDesc, wal.Address()))
	}
}

func keepRunning(networkConnector networktest.NetworkConnector) {
	gatewayURL, err := networkConnector.GetGatewayURL()
	hasGateway := err == nil

	fmt.Println("----")
	fmt.Println("Sequencer RPC", networkConnector.SequencerRPCAddress())
	for i := 0; i < networkConnector.NumValidators(); i++ {
		fmt.Println("Validator  ", i, networkConnector.ValidatorRPCAddress(i))
	}
	if hasGateway {
		fmt.Println("Gateway      ", gatewayURL)
	}
	fmt.Println("----")

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Network running until test is stopped...")
	<-done // Will block here until user hits ctrl+c
}

func privateKeyToHex(priv *ecdsa.PrivateKey) string {
	return hex.EncodeToString(priv.D.Bytes())
}
