import { expect } from "chai";
import hre from "hardhat";
import "hardhat-change-network";
import { Address } from "hardhat-deploy/types";
import { ObscuroBridge, EthereumBridge } from "../typechain-types";


describe("Simulation test", function () {

    let L1Bridge: ObscuroBridge
    let L2Bridge: EthereumBridge
    let L1Owner: Address 
    let L2Owner: Address

    this.beforeAll(async function() {
        const l1Deployment = await hre.companionNetworks.layer1.deployments.get("ObscuroBridge");
        const l2Deployment = await hre.deployments.get("EthereumBridge");

        L1Owner = (await hre.companionNetworks.layer1.getNamedAccounts()).deployer;
        L2Owner = (await hre.getNamedAccounts()).deployer;        

        L1Bridge = await hre.ethers.getContractAt("ObscuroBridge", l1Deployment.address);
        L2Bridge = await hre.ethers.getContractAt("EthereumBridge", l2Deployment.address);
    })

    it("Must not fail.", async function() {
        
    });
});