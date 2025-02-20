package faucet

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ten-protocol/go-ten/tools/tenscan/backend/config"
	"github.com/ten-protocol/go-ten/tools/tenscan/backend/container"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"github.com/ten-protocol/go-ten/go/enclave/genesis"
	"github.com/ten-protocol/go-ten/go/obsclient"
	"github.com/ten-protocol/go-ten/go/rpc"
	"github.com/ten-protocol/go-ten/go/wallet"
	"github.com/ten-protocol/go-ten/integration/datagenerator"
	"github.com/valyala/fasthttp"

	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/integration"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"github.com/ten-protocol/go-ten/integration/ethereummock"
	"github.com/ten-protocol/go-ten/integration/simulation/network"
	"github.com/ten-protocol/go-ten/integration/simulation/params"
)

func init() { //nolint:gochecknoinits
	testlog.Setup(&testlog.Cfg{
		LogDir:      testLogs,
		TestType:    "tenscan",
		TestSubtype: "test",
		LogLevel:    log.LvlInfo,
	})
}

const (
	testLogs = "../.build/tenscan/"
)

func TestTenscan(t *testing.T) {
	startPort := integration.StartPortTenscanUnitTest
	createObscuroNetwork(t, startPort)

	tenScanConfig := &config.Config{
		NodeHostAddress: fmt.Sprintf("http://127.0.0.1:%d", startPort+integration.DefaultHostRPCHTTPOffset),
		ServerAddress:   fmt.Sprintf("127.0.0.1:%d", startPort+integration.DefaultTenscanHTTPPortOffset),
		LogPath:         "sys_out",
	}
	serverAddress := fmt.Sprintf("http://%s", tenScanConfig.ServerAddress)

	tenScanContainer, err := container.NewTenScanContainer(tenScanConfig)
	require.NoError(t, err)

	err = tenScanContainer.Start()
	require.NoError(t, err)

	// wait for the msg bus contract to be deployed
	time.Sleep(30 * time.Second)

	// make sure the server is ready to receive requests
	err = waitServerIsReady(serverAddress)
	require.NoError(t, err)

	issueTransactions(
		t,
		fmt.Sprintf("ws://127.0.0.1:%d", startPort+integration.DefaultHostRPCWSOffset),
		wallet.NewInMemoryWalletFromConfig(genesis.TestnetPrefundedPK, integration.TenChainID, testlog.Logger()),
		5,
	)

	// Issue tests
	statusCode, body, err := fasthttp.Get(nil, fmt.Sprintf("%s/count/contracts/", serverAddress))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, "{\"count\":2}", string(body))

	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/count/transactions/", serverAddress))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, "{\"count\":6}", string(body))

	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/batch/latest/", serverAddress))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	type itemRes struct {
		Item common.BatchHeader `json:"item"`
	}

	itemObj := itemRes{}
	err = json.Unmarshal(body, &itemObj)
	assert.NoError(t, err)
	batchHead := itemObj.Item

	statusCode, _, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/rollup/latest/", serverAddress))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	statusCode, _, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/batch/%s", serverAddress, batchHead.Hash().String()))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/transactions/?offset=0&size=99", serverAddress))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	type publicTxsRes struct {
		Result common.TransactionListingResponse `json:"result"`
	}

	publicTxsObj := publicTxsRes{}
	err = json.Unmarshal(body, &publicTxsObj)
	assert.NoError(t, err)
	assert.Equal(t, 6, len(publicTxsObj.Result.TransactionsData))
	assert.Equal(t, uint64(6), publicTxsObj.Result.Total)

	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/batches/?offset=0&size=10", serverAddress))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	type batchlisting struct {
		Result common.BatchListingResponse `json:"result"`
	}

	batchlistingObj := batchlisting{}
	err = json.Unmarshal(body, &batchlistingObj)
	assert.NoError(t, err)
	assert.LessOrEqual(t, 9, len(batchlistingObj.Result.BatchesData))
	assert.LessOrEqual(t, uint64(9), batchlistingObj.Result.Total)
	// check results are descending order (latest first)
	assert.LessOrEqual(t, batchlistingObj.Result.BatchesData[1].Number.Cmp(batchlistingObj.Result.BatchesData[0].Number), 0)
	// check "hash" field is included in json response
	assert.Contains(t, string(body), "\"hash\"")

	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/blocks/?offset=0&size=10", serverAddress))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	type blockListing struct {
		Result common.BlockListingResponse `json:"result"`
	}

	blocklistingObj := blockListing{}
	err = json.Unmarshal(body, &blocklistingObj)
	assert.NoError(t, err)
	// assert.LessOrEqual(t, 9, len(blocklistingObj.Result.BlocksData))
	// assert.LessOrEqual(t, uint64(9), blocklistingObj.Result.Total)

	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/batch/%s", serverAddress, batchlistingObj.Result.BatchesData[0].Hash()))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	type batchFetch struct {
		Item *common.ExtBatch `json:"item"`
	}

	batchObj := batchFetch{}
	err = json.Unmarshal(body, &batchObj)
	assert.NoError(t, err)
	assert.Equal(t, batchlistingObj.Result.BatchesData[0].Hash(), batchObj.Item.Hash())

	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/info/obscuro/", serverAddress))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	type configFetch struct {
		Item common.ObscuroNetworkInfo `json:"item"`
	}

	configFetchObj := configFetch{}
	err = json.Unmarshal(body, &configFetchObj)
	assert.NoError(t, err)
	assert.NotEqual(t, configFetchObj.Item.SequencerID, gethcommon.Address{})

	// Gracefully shutdown
	err = tenScanContainer.Stop()
	assert.NoError(t, err)
}

func waitServerIsReady(serverAddr string) error {
	for now := time.Now(); time.Since(now) < 30*time.Second; time.Sleep(500 * time.Millisecond) {
		statusCode, _, err := fasthttp.Get(nil, fmt.Sprintf("%s/health/", serverAddr))
		if err != nil {
			// give it time to boot up
			if strings.Contains(err.Error(), "connection") {
				continue
			}
			return err
		}

		if statusCode == http.StatusOK {
			return nil
		}
	}
	return fmt.Errorf("timed out before server was ready")
}

// Creates a single-node Obscuro network for testing.
func createObscuroNetwork(t *testing.T, startPort int) {
	// Create the Obscuro network.
	wallets := params.NewSimWallets(1, 1, integration.EthereumChainID, integration.TenChainID)
	simParams := params.SimParams{
		NumberOfNodes:    1,
		AvgBlockDuration: 1 * time.Second,
		MgmtContractLib:  ethereummock.NewMgmtContractLibMock(),
		ERC20ContractLib: ethereummock.NewERC20ContractLibMock(),
		Wallets:          wallets,
		StartPort:        startPort,
		WithPrefunding:   true,
	}

	obscuroNetwork := network.NewNetworkOfSocketNodes(wallets)
	t.Cleanup(obscuroNetwork.TearDown)
	_, err := obscuroNetwork.Create(&simParams, nil)
	if err != nil {
		panic(fmt.Sprintf("failed to create test Obscuro network. Cause: %s", err))
	}
}

func issueTransactions(t *testing.T, hostWSAddr string, issuerWallet wallet.Wallet, numbTxs int) {
	ctx := context.Background()

	vk, err := viewingkey.GenerateViewingKeyForWallet(issuerWallet)
	assert.Nil(t, err)
	client, err := rpc.NewEncNetworkClient(hostWSAddr, vk, testlog.Logger())
	assert.Nil(t, err)
	authClient := obsclient.NewAuthObsClient(client)

	balance, err := authClient.BalanceAt(ctx, nil)
	assert.Nil(t, err)

	if balance.Cmp(big.NewInt(0)) <= 0 {
		t.Errorf("not enough balance: has %s has %s ten", issuerWallet.Address().Hex(), balance.String())
	}

	nonce, err := authClient.NonceAt(ctx, nil)
	assert.Nil(t, err)
	issuerWallet.SetNonce(nonce)

	var receipts []gethcommon.Hash
	for i := 0; i < numbTxs; i++ {
		toAddr := datagenerator.RandomAddress()

		estimatedTx := authClient.EstimateGasAndGasPrice(&types.LegacyTx{
			Nonce:    issuerWallet.GetNonceAndIncrement(),
			To:       &toAddr,
			Value:    big.NewInt(100),
			Gas:      uint64(1_000_000),
			GasPrice: gethcommon.Big1,
		})
		assert.Nil(t, err)

		signedTx, err := issuerWallet.SignTransaction(estimatedTx)
		assert.Nil(t, err)

		err = authClient.SendTransaction(ctx, signedTx)
		assert.Nil(t, err)

		fmt.Printf("Issued Tx: %s \n", signedTx.Hash().Hex())
		receipts = append(receipts, signedTx.Hash())
		time.Sleep(1500 * time.Millisecond)
	}

	for _, txHash := range receipts {
		fmt.Printf("Checking for tx receipt for %s \n", txHash)
		var receipt *types.Receipt
		for start := time.Now(); time.Since(start) < time.Minute; time.Sleep(time.Second) {
			receipt, err = authClient.TransactionReceipt(ctx, txHash)
			if err == nil {
				break
			}
			//
			// Currently when a receipt is not available the obscuro node is returning nil instead of err ethereum.NotFound
			// once that's fixed this commented block should be removed
			//if !errors.Is(err, ethereum.NotFound) {
			//	t.Fatal(err)
			//}
			if receipt != nil && receipt.Status == 1 {
				break
			}
			fmt.Printf("no tx receipt after %s - %s\n", time.Since(start), err)
		}

		if receipt == nil {
			t.Fatalf("Did not mine the transaction after %s seconds  - receipt: %+v", 30*time.Second, receipt)
		}
		if receipt.Status == 0 {
			t.Fatalf("Tx Failed")
		}
	}
}
