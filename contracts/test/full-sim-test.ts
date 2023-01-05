import { expect } from "chai";
import hre from "hardhat";
import "hardhat-change-network";
import { Address, Receipt } from "hardhat-deploy/types";
import { ObscuroBridge, ObscuroERC20, EthereumBridge } from "../typechain-types";
import { l1 } from "../typechain-types/src/bridge";
import { ObsERC20 } from "../typechain-types/src/common/ObsERC20.sol";


describe("Simulation test", function () {
    let l2Bridge: EthereumBridge;
    let l1Owner: Address 
    let l2Owner: Address

    type CrossChainAccount = {
        layer1: Address,
        layer2: Address
    };

    type Tokens = {
        [name: string] : CrossChainAccount;
    };

    let aliceAccounts : CrossChainAccount;
    let bobAccounts : CrossChainAccount;
    let tokens: Tokens = {};

    const eventSignature = "LogMessagePublished(address,uint64,uint32,uint32,bytes,uint8)";
    const topic = hre.ethers.utils.id(eventSignature)
    let eventIface = new hre.ethers.utils.Interface([ `event ${eventSignature}`]);
    function getXChainMessages(result: Receipt) {
        const events = result.logs?.filter((x)=> { 
            return x.topics.find((t: string)=> t == topic) != undefined;
        });

        const messages = events!.map((event)=> {
            const decodedEvent = eventIface.parseLog({
                topics: event!.topics!,
                data: event!.data
            });
        
            const xchainMessage = {
                sender: decodedEvent.args[0],
                sequence: decodedEvent.args[1],
                nonce: decodedEvent.args[2],
                topic: decodedEvent.args[3],
                payload: decodedEvent.args[4],
                consistencyLevel: decodedEvent.args[5]
            };

            return xchainMessage;
        })

        return messages;
    }

    this.beforeAll(async function() {
        const l1Deployment = await hre.companionNetworks.layer1.deployments.get("ObscuroBridge");
        const l2Deployment = await hre.deployments.get("EthereumBridge");

        let L1Accounts = (await hre.companionNetworks.layer1.getNamedAccounts());
        let L2Accounts = (await hre.getNamedAccounts());

        l1Owner = L1Accounts.deployer;
        l2Owner = L2Accounts.deployer;

        aliceAccounts = { 
            layer1 : L1Accounts.alice,
            layer2 : L2Accounts.alice
        };

        bobAccounts = { 
            layer1 : L1Accounts.bob,
            layer2 : L2Accounts.bob
        };

        l2Bridge = await hre.ethers.getContractAt("EthereumBridge", l2Deployment.address);
        

        const hocERC20 = await hre.companionNetworks.layer1.deployments.get("HOCERC20");

        const hocl1Address: Address = hocERC20.address;
        tokens.hoc = {
            layer1: hocl1Address,
            layer2: await l2Bridge.remoteToLocalToken(hocl1Address)
        };

        const pocERC20 = await hre.companionNetworks.layer1.deployments.get("POCERC20");
        const pocl1Address : Address = pocERC20.address;
        tokens.poc = { 
            layer1 : pocl1Address,
            layer2 : await l2Bridge.remoteToLocalToken(pocl1Address)
        };


        await hre.run('add-key', { address: aliceAccounts.layer2 });
        await hre.run('add-key', { address: bobAccounts.layer2 });


        console.log(`Token accounts retrieved = ${JSON.stringify(tokens, null, "\t")}`);
        console.log(`Alice accouunts = ${JSON.stringify(aliceAccounts, null, "\t")}`);
        console.log(`Bob accouunts = ${JSON.stringify(bobAccounts, null, "\t")}`);
    })

    it("Network has been initialized.", async function() {
        const hocRead = await hre.deployments.read("EthereumBridge", {
            from: l2Owner,
        }, "remoteToLocalToken", tokens.hoc.layer1);

        expect(hocRead).is.properAddress;
        expect(hocRead).not.hexEqual("0x0");

        const pocRead = await hre.deployments.read("EthereumBridge", {
            from: l2Owner,
        }, "remoteToLocalToken", tokens.poc.layer1);

        expect(pocRead).is.properAddress;
        expect(pocRead).not.hexEqual("0x0");
    });

    it("Alice is able to deposit.", async function() {
        const l1Deployments = hre.companionNetworks.layer1.deployments;
        const l2Deployments = hre.deployments;

        const result = await l1Deployments.execute("HOCERC20", {
            from: l1Owner, 
            log: false,
        }, "issueFor", aliceAccounts.layer1, 100_000);

        await expect(result).not.reverted;
        
        expect(await l1Deployments.read("HOCERC20", {
            from: aliceAccounts.layer1
        }, "balanceOf", aliceAccounts.layer1)).is.equal(100_000);

        const allowanceResult = await l1Deployments.execute("HOCERC20", {
            from: aliceAccounts.layer1, 
            log: false,
        }, "increaseAllowance", (await l1Deployments.get("ObscuroBridge")).address, 100_000);

        await expect(allowanceResult).not.reverted;

        const bridgeResult = l1Deployments.execute("ObscuroBridge", {
            from: aliceAccounts.layer1,
            log: false,
        }, "sendERC20", tokens.hoc.layer1, 100_000, aliceAccounts.layer2);

        await expect(bridgeResult).not.reverted;

        const messages = getXChainMessages(await bridgeResult);
        await Promise.all(messages.map(async msg=>{
            const relayResult = l2Deployments.execute("CrossChainMessenger", {
                from: l2Owner, 
                log: false,
            }, "relayMessage", msg)
            return expect(relayResult).not.reverted;
        }));

        const L2HOC : ObsERC20 = await hre.ethers.getContractAt("ObsERC20", tokens.hoc.layer2);
        
        const l2Balance = await L2HOC.connect(aliceAccounts.layer2).balanceOf(aliceAccounts.layer2, {
            from: aliceAccounts.layer2
        });

        expect(l2Balance).equals(100_000);
    });
});