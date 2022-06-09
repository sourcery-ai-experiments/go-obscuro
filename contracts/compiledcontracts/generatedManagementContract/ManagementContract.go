// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package generatedManagementContract

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// GeneratedManagementContractMetaData contains all meta data concerning the GeneratedManagementContract contract.
var GeneratedManagementContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"rollupData\",\"type\":\"string\"}],\"name\":\"AddRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"aggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"initSecret\",\"type\":\"bytes\"}],\"name\":\"InitializeNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"requestReport\",\"type\":\"string\"}],\"name\":\"RequestNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"attesterID\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"requesterID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"attesterSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"responseSecret\",\"type\":\"bytes\"}],\"name\":\"RespondNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"attestationRequests\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"attested\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"rollups\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50611fd1806100206000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c8063d4c806641161005b578063d4c80664146100ea578063e0643dfc1461011a578063e0fd84bd1461014d578063e34fbfc8146101695761007d565b80638ef74f8914610082578063981214ba146100b2578063c719bf50146100ce575b600080fd5b61009c60048036038101906100979190610f16565b610185565b6040516100a99190610fdc565b60405180910390f35b6100cc60048036038101906100c79190611133565b610225565b005b6100e860048036038101906100e39190611232565b6103c2565b005b61010460048036038101906100ff9190610f16565b610454565b60405161011191906112ad565b60405180910390f35b610134600480360381019061012f91906112fe565b610474565b6040516101449493929190611375565b60405180910390f35b6101676004803603810190610162919061143c565b6104e1565b005b610183600480360381019061017e91906114d6565b61061e565b005b600160205280600052604060002060009150905080546101a490611552565b80601f01602080910402602001604051908101604052809291908181526020018280546101d090611552565b801561021d5780601f106101f25761010080835404028352916020019161021d565b820191906000526020600020905b81548152906001019060200180831161020057829003601f168201915b505050505081565b6000600260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1690508061028057600080fd5b60006102ae86868560405160200161029a93929190611613565b604051602081830303815290604052610671565b905060006102bc82866106ac565b90508673ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16146102f6886106d3565b6102ff836106d3565b6040516020016103109291906117de565b60405160208183030381529060405290610360576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103579190610fdc565b60405180910390fd5b506001600260008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff02191690831515021790555050505050505050565b600360009054906101000a900460ff16156103dc57600080fd5b6001600360006101000a81548160ff0219169083151502179055506001600260008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908315150217905550505050565b60026020528060005260406000206000915054906101000a900460ff1681565b6000602052816000526040600020818154811061049057600080fd5b9060005260206000209060040201600091509150508060000154908060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060020154908060030154905084565b600260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1661053757600080fd5b600060405180608001604052808881526020018773ffffffffffffffffffffffffffffffffffffffff1681526020018681526020018581525090506000804381526020019081526020016000208190806001815401808255809150506001900390600052602060002090600402016000909190919091506000820151816000015560208201518160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506040820151816002015560608201518160030155505050505050505050565b8181600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020919061066c929190610e01565b505050565b600061067d8251610896565b8260405160200161068f92919061186f565b604051602081830303815290604052805190602001209050919050565b60008060006106bb85856109f7565b915091506106c881610a7a565b819250505092915050565b60606000602867ffffffffffffffff8111156106f2576106f1611008565b5b6040519080825280601f01601f1916602001820160405280156107245781602001600182028036833780820191505090505b50905060005b601481101561088c57600081601361074291906118cd565b600861074e9190611901565b600261075a9190611a8e565b8573ffffffffffffffffffffffffffffffffffffffff1661077b9190611b08565b60f81b9050600060108260f81c6107929190611b46565b60f81b905060008160f81c60106107a99190611b77565b8360f81c6107b79190611bb2565b60f81b90506107c582610c4f565b858560026107d39190611901565b815181106107e4576107e3611be6565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535061081c81610c4f565b85600186600261082c9190611901565b6108369190611c15565b8151811061084757610846611be6565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350505050808061088490611c6b565b91505061072a565b5080915050919050565b606060008214156108de576040518060400160405280600181526020017f300000000000000000000000000000000000000000000000000000000000000081525090506109f2565b600082905060005b600082146109105780806108f990611c6b565b915050600a826109099190611b08565b91506108e6565b60008167ffffffffffffffff81111561092c5761092b611008565b5b6040519080825280601f01601f19166020018201604052801561095e5781602001600182028036833780820191505090505b5090505b600085146109eb5760018261097791906118cd565b9150600a856109869190611cb4565b60306109929190611c15565b60f81b8183815181106109a8576109a7611be6565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600a856109e49190611b08565b9450610962565b8093505050505b919050565b600080604183511415610a395760008060006020860151925060408601519150606086015160001a9050610a2d87828585610c95565b94509450505050610a73565b604083511415610a6a576000806020850151915060408501519050610a5f868383610da2565b935093505050610a73565b60006002915091505b9250929050565b60006004811115610a8e57610a8d611ce5565b5b816004811115610aa157610aa0611ce5565b5b1415610aac57610c4c565b60016004811115610ac057610abf611ce5565b5b816004811115610ad357610ad2611ce5565b5b1415610b14576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b0b90611d60565b60405180910390fd5b60026004811115610b2857610b27611ce5565b5b816004811115610b3b57610b3a611ce5565b5b1415610b7c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b7390611dcc565b60405180910390fd5b60036004811115610b9057610b8f611ce5565b5b816004811115610ba357610ba2611ce5565b5b1415610be4576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610bdb90611e5e565b60405180910390fd5b600480811115610bf757610bf6611ce5565b5b816004811115610c0a57610c09611ce5565b5b1415610c4b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c4290611ef0565b60405180910390fd5b5b50565b6000600a8260f81c60ff161015610c7a5760308260f81c610c709190611f10565b60f81b9050610c90565b60578260f81c610c8a9190611f10565b60f81b90505b919050565b6000807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08360001c1115610cd0576000600391509150610d99565b601b8560ff1614158015610ce85750601c8560ff1614155b15610cfa576000600491509150610d99565b600060018787878760405160008152602001604052604051610d1f9493929190611f56565b6020604051602081039080840390855afa158015610d41573d6000803e3d6000fd5b505050602060405103519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415610d9057600060019250925050610d99565b80600092509250505b94509492505050565b60008060007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60001b841690506000601b60ff8660001c901c610de59190611c15565b9050610df387828885610c95565b935093505050935093915050565b828054610e0d90611552565b90600052602060002090601f016020900481019282610e2f5760008555610e76565b82601f10610e4857803560ff1916838001178555610e76565b82800160010185558215610e76579182015b82811115610e75578235825591602001919060010190610e5a565b5b509050610e839190610e87565b5090565b5b80821115610ea0576000816000905550600101610e88565b5090565b6000604051905090565b600080fd5b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610ee382610eb8565b9050919050565b610ef381610ed8565b8114610efe57600080fd5b50565b600081359050610f1081610eea565b92915050565b600060208284031215610f2c57610f2b610eae565b5b6000610f3a84828501610f01565b91505092915050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610f7d578082015181840152602081019050610f62565b83811115610f8c576000848401525b50505050565b6000601f19601f8301169050919050565b6000610fae82610f43565b610fb88185610f4e565b9350610fc8818560208601610f5f565b610fd181610f92565b840191505092915050565b60006020820190508181036000830152610ff68184610fa3565b905092915050565b600080fd5b600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b61104082610f92565b810181811067ffffffffffffffff8211171561105f5761105e611008565b5b80604052505050565b6000611072610ea4565b905061107e8282611037565b919050565b600067ffffffffffffffff82111561109e5761109d611008565b5b6110a782610f92565b9050602081019050919050565b82818337600083830152505050565b60006110d66110d184611083565b611068565b9050828152602081018484840111156110f2576110f1611003565b5b6110fd8482856110b4565b509392505050565b600082601f83011261111a57611119610ffe565b5b813561112a8482602086016110c3565b91505092915050565b6000806000806080858703121561114d5761114c610eae565b5b600061115b87828801610f01565b945050602061116c87828801610f01565b935050604085013567ffffffffffffffff81111561118d5761118c610eb3565b5b61119987828801611105565b925050606085013567ffffffffffffffff8111156111ba576111b9610eb3565b5b6111c687828801611105565b91505092959194509250565b600080fd5b600080fd5b60008083601f8401126111f2576111f1610ffe565b5b8235905067ffffffffffffffff81111561120f5761120e6111d2565b5b60208301915083600182028301111561122b5761122a6111d7565b5b9250929050565b60008060006040848603121561124b5761124a610eae565b5b600061125986828701610f01565b935050602084013567ffffffffffffffff81111561127a57611279610eb3565b5b611286868287016111dc565b92509250509250925092565b60008115159050919050565b6112a781611292565b82525050565b60006020820190506112c2600083018461129e565b92915050565b6000819050919050565b6112db816112c8565b81146112e657600080fd5b50565b6000813590506112f8816112d2565b92915050565b6000806040838503121561131557611314610eae565b5b6000611323858286016112e9565b9250506020611334858286016112e9565b9150509250929050565b6000819050919050565b6113518161133e565b82525050565b61136081610ed8565b82525050565b61136f816112c8565b82525050565b600060808201905061138a6000830187611348565b6113976020830186611357565b6113a46040830185611348565b6113b16060830184611366565b95945050505050565b6113c38161133e565b81146113ce57600080fd5b50565b6000813590506113e0816113ba565b92915050565b60008083601f8401126113fc576113fb610ffe565b5b8235905067ffffffffffffffff811115611419576114186111d2565b5b602083019150836001820283011115611435576114346111d7565b5b9250929050565b60008060008060008060a0878903121561145957611458610eae565b5b600061146789828a016113d1565b965050602061147889828a01610f01565b955050604061148989828a016113d1565b945050606061149a89828a016112e9565b935050608087013567ffffffffffffffff8111156114bb576114ba610eb3565b5b6114c789828a016113e6565b92509250509295509295509295565b600080602083850312156114ed576114ec610eae565b5b600083013567ffffffffffffffff81111561150b5761150a610eb3565b5b611517858286016113e6565b92509250509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061156a57607f821691505b6020821081141561157e5761157d611523565b5b50919050565b60008160601b9050919050565b600061159c82611584565b9050919050565b60006115ae82611591565b9050919050565b6115c66115c182610ed8565b6115a3565b82525050565b600081519050919050565b600081905092915050565b60006115ed826115cc565b6115f781856115d7565b9350611607818560208601610f5f565b80840191505092915050565b600061161f82866115b5565b60148201915061162f82856115b5565b60148201915061163f82846115e2565b9150819050949350505050565b600081905092915050565b7f7265636f7665726564206164647265737320616e64206174746573746572494460008201527f20646f6e74206d61746368200000000000000000000000000000000000000000602082015250565b60006116b3602c8361164c565b91506116be82611657565b602c82019050919050565b7f0a2045787065637465643a20202020202020202020202020202020202020202060008201527f2020202000000000000000000000000000000000000000000000000000000000602082015250565b600061172560248361164c565b9150611730826116c9565b602482019050919050565b600061174682610f43565b611750818561164c565b9350611760818560208601610f5f565b80840191505092915050565b7f0a202f207265636f7665726564416464725369676e656443616c63756c61746560008201527f643a202000000000000000000000000000000000000000000000000000000000602082015250565b60006117c860248361164c565b91506117d38261176c565b602482019050919050565b60006117e9826116a6565b91506117f482611718565b9150611800828561173b565b915061180b826117bb565b9150611817828461173b565b91508190509392505050565b7f19457468657265756d205369676e6564204d6573736167653a0a000000000000600082015250565b6000611859601a8361164c565b915061186482611823565b601a82019050919050565b600061187a8261184c565b9150611886828561173b565b915061189282846115e2565b91508190509392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006118d8826112c8565b91506118e3836112c8565b9250828210156118f6576118f561189e565b5b828203905092915050565b600061190c826112c8565b9150611917836112c8565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156119505761194f61189e565b5b828202905092915050565b60008160011c9050919050565b6000808291508390505b60018511156119b25780860481111561198e5761198d61189e565b5b600185161561199d5780820291505b80810290506119ab8561195b565b9450611972565b94509492505050565b6000826119cb5760019050611a87565b816119d95760009050611a87565b81600181146119ef57600281146119f957611a28565b6001915050611a87565b60ff841115611a0b57611a0a61189e565b5b8360020a915084821115611a2257611a2161189e565b5b50611a87565b5060208310610133831016604e8410600b8410161715611a5d5782820a905083811115611a5857611a5761189e565b5b611a87565b611a6a8484846001611968565b92509050818404811115611a8157611a8061189e565b5b81810290505b9392505050565b6000611a99826112c8565b9150611aa4836112c8565b9250611ad17fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff84846119bb565b905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000611b13826112c8565b9150611b1e836112c8565b925082611b2e57611b2d611ad9565b5b828204905092915050565b600060ff82169050919050565b6000611b5182611b39565b9150611b5c83611b39565b925082611b6c57611b6b611ad9565b5b828204905092915050565b6000611b8282611b39565b9150611b8d83611b39565b92508160ff0483118215151615611ba757611ba661189e565b5b828202905092915050565b6000611bbd82611b39565b9150611bc883611b39565b925082821015611bdb57611bda61189e565b5b828203905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6000611c20826112c8565b9150611c2b836112c8565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03821115611c6057611c5f61189e565b5b828201905092915050565b6000611c76826112c8565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415611ca957611ca861189e565b5b600182019050919050565b6000611cbf826112c8565b9150611cca836112c8565b925082611cda57611cd9611ad9565b5b828206905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f45434453413a20696e76616c6964207369676e61747572650000000000000000600082015250565b6000611d4a601883610f4e565b9150611d5582611d14565b602082019050919050565b60006020820190508181036000830152611d7981611d3d565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265206c656e67746800600082015250565b6000611db6601f83610f4e565b9150611dc182611d80565b602082019050919050565b60006020820190508181036000830152611de581611da9565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265202773272076616c60008201527f7565000000000000000000000000000000000000000000000000000000000000602082015250565b6000611e48602283610f4e565b9150611e5382611dec565b604082019050919050565b60006020820190508181036000830152611e7781611e3b565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265202776272076616c60008201527f7565000000000000000000000000000000000000000000000000000000000000602082015250565b6000611eda602283610f4e565b9150611ee582611e7e565b604082019050919050565b60006020820190508181036000830152611f0981611ecd565b9050919050565b6000611f1b82611b39565b9150611f2683611b39565b92508260ff03821115611f3c57611f3b61189e565b5b828201905092915050565b611f5081611b39565b82525050565b6000608082019050611f6b6000830187611348565b611f786020830186611f47565b611f856040830185611348565b611f926060830184611348565b9594505050505056fea2646970667358221220a14078ed9b28bb7eac6e1690fdea98610496ded8979d922f067349e3776aba3164736f6c634300080c0033",
}

// GeneratedManagementContractABI is the input ABI used to generate the binding from.
// Deprecated: Use GeneratedManagementContractMetaData.ABI instead.
var GeneratedManagementContractABI = GeneratedManagementContractMetaData.ABI

// GeneratedManagementContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use GeneratedManagementContractMetaData.Bin instead.
var GeneratedManagementContractBin = GeneratedManagementContractMetaData.Bin

// DeployGeneratedManagementContract deploys a new Ethereum contract, binding an instance of GeneratedManagementContract to it.
func DeployGeneratedManagementContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *GeneratedManagementContract, error) {
	parsed, err := GeneratedManagementContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(GeneratedManagementContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &GeneratedManagementContract{GeneratedManagementContractCaller: GeneratedManagementContractCaller{contract: contract}, GeneratedManagementContractTransactor: GeneratedManagementContractTransactor{contract: contract}, GeneratedManagementContractFilterer: GeneratedManagementContractFilterer{contract: contract}}, nil
}

// GeneratedManagementContract is an auto generated Go binding around an Ethereum contract.
type GeneratedManagementContract struct {
	GeneratedManagementContractCaller     // Read-only binding to the contract
	GeneratedManagementContractTransactor // Write-only binding to the contract
	GeneratedManagementContractFilterer   // Log filterer for contract events
}

// GeneratedManagementContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type GeneratedManagementContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GeneratedManagementContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GeneratedManagementContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GeneratedManagementContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GeneratedManagementContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GeneratedManagementContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GeneratedManagementContractSession struct {
	Contract     *GeneratedManagementContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                // Call options to use throughout this session
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// GeneratedManagementContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GeneratedManagementContractCallerSession struct {
	Contract *GeneratedManagementContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                      // Call options to use throughout this session
}

// GeneratedManagementContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GeneratedManagementContractTransactorSession struct {
	Contract     *GeneratedManagementContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                      // Transaction auth options to use throughout this session
}

// GeneratedManagementContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type GeneratedManagementContractRaw struct {
	Contract *GeneratedManagementContract // Generic contract binding to access the raw methods on
}

// GeneratedManagementContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GeneratedManagementContractCallerRaw struct {
	Contract *GeneratedManagementContractCaller // Generic read-only contract binding to access the raw methods on
}

// GeneratedManagementContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GeneratedManagementContractTransactorRaw struct {
	Contract *GeneratedManagementContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGeneratedManagementContract creates a new instance of GeneratedManagementContract, bound to a specific deployed contract.
func NewGeneratedManagementContract(address common.Address, backend bind.ContractBackend) (*GeneratedManagementContract, error) {
	contract, err := bindGeneratedManagementContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GeneratedManagementContract{GeneratedManagementContractCaller: GeneratedManagementContractCaller{contract: contract}, GeneratedManagementContractTransactor: GeneratedManagementContractTransactor{contract: contract}, GeneratedManagementContractFilterer: GeneratedManagementContractFilterer{contract: contract}}, nil
}

// NewGeneratedManagementContractCaller creates a new read-only instance of GeneratedManagementContract, bound to a specific deployed contract.
func NewGeneratedManagementContractCaller(address common.Address, caller bind.ContractCaller) (*GeneratedManagementContractCaller, error) {
	contract, err := bindGeneratedManagementContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GeneratedManagementContractCaller{contract: contract}, nil
}

// NewGeneratedManagementContractTransactor creates a new write-only instance of GeneratedManagementContract, bound to a specific deployed contract.
func NewGeneratedManagementContractTransactor(address common.Address, transactor bind.ContractTransactor) (*GeneratedManagementContractTransactor, error) {
	contract, err := bindGeneratedManagementContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GeneratedManagementContractTransactor{contract: contract}, nil
}

// NewGeneratedManagementContractFilterer creates a new log filterer instance of GeneratedManagementContract, bound to a specific deployed contract.
func NewGeneratedManagementContractFilterer(address common.Address, filterer bind.ContractFilterer) (*GeneratedManagementContractFilterer, error) {
	contract, err := bindGeneratedManagementContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GeneratedManagementContractFilterer{contract: contract}, nil
}

// bindGeneratedManagementContract binds a generic wrapper to an already deployed contract.
func bindGeneratedManagementContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(GeneratedManagementContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GeneratedManagementContract *GeneratedManagementContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GeneratedManagementContract.Contract.GeneratedManagementContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GeneratedManagementContract *GeneratedManagementContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.GeneratedManagementContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GeneratedManagementContract *GeneratedManagementContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.GeneratedManagementContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GeneratedManagementContract *GeneratedManagementContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GeneratedManagementContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GeneratedManagementContract *GeneratedManagementContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GeneratedManagementContract *GeneratedManagementContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.contract.Transact(opts, method, params...)
}

// AttestationRequests is a free data retrieval call binding the contract method 0x8ef74f89.
//
// Solidity: function attestationRequests(address ) view returns(string)
func (_GeneratedManagementContract *GeneratedManagementContractCaller) AttestationRequests(opts *bind.CallOpts, arg0 common.Address) (string, error) {
	var out []interface{}
	err := _GeneratedManagementContract.contract.Call(opts, &out, "attestationRequests", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// AttestationRequests is a free data retrieval call binding the contract method 0x8ef74f89.
//
// Solidity: function attestationRequests(address ) view returns(string)
func (_GeneratedManagementContract *GeneratedManagementContractSession) AttestationRequests(arg0 common.Address) (string, error) {
	return _GeneratedManagementContract.Contract.AttestationRequests(&_GeneratedManagementContract.CallOpts, arg0)
}

// AttestationRequests is a free data retrieval call binding the contract method 0x8ef74f89.
//
// Solidity: function attestationRequests(address ) view returns(string)
func (_GeneratedManagementContract *GeneratedManagementContractCallerSession) AttestationRequests(arg0 common.Address) (string, error) {
	return _GeneratedManagementContract.Contract.AttestationRequests(&_GeneratedManagementContract.CallOpts, arg0)
}

// Attested is a free data retrieval call binding the contract method 0xd4c80664.
//
// Solidity: function attested(address ) view returns(bool)
func (_GeneratedManagementContract *GeneratedManagementContractCaller) Attested(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _GeneratedManagementContract.contract.Call(opts, &out, "attested", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Attested is a free data retrieval call binding the contract method 0xd4c80664.
//
// Solidity: function attested(address ) view returns(bool)
func (_GeneratedManagementContract *GeneratedManagementContractSession) Attested(arg0 common.Address) (bool, error) {
	return _GeneratedManagementContract.Contract.Attested(&_GeneratedManagementContract.CallOpts, arg0)
}

// Attested is a free data retrieval call binding the contract method 0xd4c80664.
//
// Solidity: function attested(address ) view returns(bool)
func (_GeneratedManagementContract *GeneratedManagementContractCallerSession) Attested(arg0 common.Address) (bool, error) {
	return _GeneratedManagementContract.Contract.Attested(&_GeneratedManagementContract.CallOpts, arg0)
}

// Rollups is a free data retrieval call binding the contract method 0xe0643dfc.
//
// Solidity: function rollups(uint256 , uint256 ) view returns(bytes32 ParentHash, address AggregatorID, bytes32 L1Block, uint256 Number)
func (_GeneratedManagementContract *GeneratedManagementContractCaller) Rollups(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (struct {
	ParentHash   [32]byte
	AggregatorID common.Address
	L1Block      [32]byte
	Number       *big.Int
}, error) {
	var out []interface{}
	err := _GeneratedManagementContract.contract.Call(opts, &out, "rollups", arg0, arg1)

	outstruct := new(struct {
		ParentHash   [32]byte
		AggregatorID common.Address
		L1Block      [32]byte
		Number       *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ParentHash = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.AggregatorID = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.L1Block = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)
	outstruct.Number = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Rollups is a free data retrieval call binding the contract method 0xe0643dfc.
//
// Solidity: function rollups(uint256 , uint256 ) view returns(bytes32 ParentHash, address AggregatorID, bytes32 L1Block, uint256 Number)
func (_GeneratedManagementContract *GeneratedManagementContractSession) Rollups(arg0 *big.Int, arg1 *big.Int) (struct {
	ParentHash   [32]byte
	AggregatorID common.Address
	L1Block      [32]byte
	Number       *big.Int
}, error) {
	return _GeneratedManagementContract.Contract.Rollups(&_GeneratedManagementContract.CallOpts, arg0, arg1)
}

// Rollups is a free data retrieval call binding the contract method 0xe0643dfc.
//
// Solidity: function rollups(uint256 , uint256 ) view returns(bytes32 ParentHash, address AggregatorID, bytes32 L1Block, uint256 Number)
func (_GeneratedManagementContract *GeneratedManagementContractCallerSession) Rollups(arg0 *big.Int, arg1 *big.Int) (struct {
	ParentHash   [32]byte
	AggregatorID common.Address
	L1Block      [32]byte
	Number       *big.Int
}, error) {
	return _GeneratedManagementContract.Contract.Rollups(&_GeneratedManagementContract.CallOpts, arg0, arg1)
}

// AddRollup is a paid mutator transaction binding the contract method 0xe0fd84bd.
//
// Solidity: function AddRollup(bytes32 ParentHash, address AggregatorID, bytes32 L1Block, uint256 Number, string rollupData) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactor) AddRollup(opts *bind.TransactOpts, ParentHash [32]byte, AggregatorID common.Address, L1Block [32]byte, Number *big.Int, rollupData string) (*types.Transaction, error) {
	return _GeneratedManagementContract.contract.Transact(opts, "AddRollup", ParentHash, AggregatorID, L1Block, Number, rollupData)
}

// AddRollup is a paid mutator transaction binding the contract method 0xe0fd84bd.
//
// Solidity: function AddRollup(bytes32 ParentHash, address AggregatorID, bytes32 L1Block, uint256 Number, string rollupData) returns()
func (_GeneratedManagementContract *GeneratedManagementContractSession) AddRollup(ParentHash [32]byte, AggregatorID common.Address, L1Block [32]byte, Number *big.Int, rollupData string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.AddRollup(&_GeneratedManagementContract.TransactOpts, ParentHash, AggregatorID, L1Block, Number, rollupData)
}

// AddRollup is a paid mutator transaction binding the contract method 0xe0fd84bd.
//
// Solidity: function AddRollup(bytes32 ParentHash, address AggregatorID, bytes32 L1Block, uint256 Number, string rollupData) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactorSession) AddRollup(ParentHash [32]byte, AggregatorID common.Address, L1Block [32]byte, Number *big.Int, rollupData string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.AddRollup(&_GeneratedManagementContract.TransactOpts, ParentHash, AggregatorID, L1Block, Number, rollupData)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0xc719bf50.
//
// Solidity: function InitializeNetworkSecret(address aggregatorID, bytes initSecret) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactor) InitializeNetworkSecret(opts *bind.TransactOpts, aggregatorID common.Address, initSecret []byte) (*types.Transaction, error) {
	return _GeneratedManagementContract.contract.Transact(opts, "InitializeNetworkSecret", aggregatorID, initSecret)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0xc719bf50.
//
// Solidity: function InitializeNetworkSecret(address aggregatorID, bytes initSecret) returns()
func (_GeneratedManagementContract *GeneratedManagementContractSession) InitializeNetworkSecret(aggregatorID common.Address, initSecret []byte) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.InitializeNetworkSecret(&_GeneratedManagementContract.TransactOpts, aggregatorID, initSecret)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0xc719bf50.
//
// Solidity: function InitializeNetworkSecret(address aggregatorID, bytes initSecret) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactorSession) InitializeNetworkSecret(aggregatorID common.Address, initSecret []byte) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.InitializeNetworkSecret(&_GeneratedManagementContract.TransactOpts, aggregatorID, initSecret)
}

// RequestNetworkSecret is a paid mutator transaction binding the contract method 0xe34fbfc8.
//
// Solidity: function RequestNetworkSecret(string requestReport) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactor) RequestNetworkSecret(opts *bind.TransactOpts, requestReport string) (*types.Transaction, error) {
	return _GeneratedManagementContract.contract.Transact(opts, "RequestNetworkSecret", requestReport)
}

// RequestNetworkSecret is a paid mutator transaction binding the contract method 0xe34fbfc8.
//
// Solidity: function RequestNetworkSecret(string requestReport) returns()
func (_GeneratedManagementContract *GeneratedManagementContractSession) RequestNetworkSecret(requestReport string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.RequestNetworkSecret(&_GeneratedManagementContract.TransactOpts, requestReport)
}

// RequestNetworkSecret is a paid mutator transaction binding the contract method 0xe34fbfc8.
//
// Solidity: function RequestNetworkSecret(string requestReport) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactorSession) RequestNetworkSecret(requestReport string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.RequestNetworkSecret(&_GeneratedManagementContract.TransactOpts, requestReport)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0x981214ba.
//
// Solidity: function RespondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactor) RespondNetworkSecret(opts *bind.TransactOpts, attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte) (*types.Transaction, error) {
	return _GeneratedManagementContract.contract.Transact(opts, "RespondNetworkSecret", attesterID, requesterID, attesterSig, responseSecret)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0x981214ba.
//
// Solidity: function RespondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret) returns()
func (_GeneratedManagementContract *GeneratedManagementContractSession) RespondNetworkSecret(attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.RespondNetworkSecret(&_GeneratedManagementContract.TransactOpts, attesterID, requesterID, attesterSig, responseSecret)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0x981214ba.
//
// Solidity: function RespondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactorSession) RespondNetworkSecret(attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.RespondNetworkSecret(&_GeneratedManagementContract.TransactOpts, attesterID, requesterID, attesterSig, responseSecret)
}
