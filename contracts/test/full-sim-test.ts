import { expect } from "chai";
import hre from "hardhat";
import "hardhat-change-network";
import { ObscuroBridge, ObscuroL2Bridge } from "../typechain-types";


describe("Simulation test", function () {

    let L1Bridge: ObscuroBridge
    let L2Bridge: ObscuroL2Bridge

    this.beforeAll(async function() {
        const l1Deployment = await hre.companionNetworks.layer1.deployments.get("ObscuroBridge");
        const l2Deployment = await hre.deployments.get("ObsuroL2Bridge");

        L1Bridge = await hre.ethers.getContractAt("ObscuroBridge", l1Deployment.address);
        L2Bridge = await hre.ethers.getContractAt("ObscuroL2Bridge", l2Deployment.address);
    })

    it("Must not fail.", async function() {
        
    });
});