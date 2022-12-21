import { task } from "hardhat/config";
import 'hardhat/types/config';

import process, { exit} from 'process';
import { HardhatNetworkUserConfig } from "hardhat/types/config";
import "hardhat-change-network";


task("deploy", "Prepares for deploying.")
.addFlag("noWallet")
.addOptionalParam("environment")
.setAction(async function(args, hre, runSuper) {

    if (args.environment) {
        hre.changeNetwork(args.environment);
    }

    const rpcURL = (hre.network.config as HardhatNetworkUserConfig).obscuroEncRpcUrl;

    if (!rpcURL || args.noWallet) {
        await runSuper();
        return;
    } 

    process.on('exit', ()=>{
        hre.run("stop-wallet-extension")
    });
    
    await hre.run("run-wallet-extension", {rpcUrl : rpcURL});

    process.on('SIGINT', ()=>exit(1));

    try {
        const { deployer } = await hre.getNamedAccounts();
        await hre.run('add-key', { address: deployer });
        await runSuper();
    } 
    finally 
    {
        await hre.run("stop-wallet-extension");
    }

});


declare module 'hardhat/types/config' {
    interface HardhatNetworkUserConfig {
        obscuroEncRpcUrl?: string
      }    
}