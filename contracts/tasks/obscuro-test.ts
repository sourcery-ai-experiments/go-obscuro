import { task } from "hardhat/config";
import 'hardhat/types/config';

import process, {exit} from 'process';
import { HardhatNetworkUserConfig } from "hardhat/types/config";

import { spawn } from 'node:child_process';

task("create-test-env", "Creates environment for testing")
.setAction(async function(args, hre) {
    process.on('SIGINT', ()=> { 
        exit(1) 
    });

    const rpcURL = (hre.network.config as HardhatNetworkUserConfig).obscuroEncRpcUrl;
    await Promise.all([
        hre.run('run-geth-nodes')
            .then((gethNetworkProc)=>{
                const deployProc = spawn("npx", [
                    "hardhat",
                    "deploy",
                    "--network", "simGeth",
                    "--reset"
                ])

                return new Promise(resolve=>{
                    deployProc.on('exit',function() {
                        resolve(true);
                    })
                });
            }),
        hre.run('run-enclave'),
        hre.run('run-host')
            .then(()=>hre.run("run-wallet-extension", {rpcUrl : rpcURL}))
            .then(()=>hre.getNamedAccounts())
            .then((accounts)=>hre.run('add-key', { address: accounts.deployer}))
    ]);
    console.log(`Deploying L2...`);
    await hre.run('deploy', { noWallet: true, reset: true });
});

task("sim-test", "Runs the tests. Takes into account if the network has obscuroEncRpcUrl")
.addFlag("runObscuroNode", "Will start a docker container with the node and test against it.")
.addFlag("runGethNodes", "Will start geth nodes to use as layer 1.")
.setAction(async function(args, hre, runSuper) {
    const rpcURL = (hre.network.config as HardhatNetworkUserConfig).obscuroEncRpcUrl;
    if (!rpcURL) {
        console.log(`Network is not configured as an obscuro network. Please include "obscuroEncRpcUrl" inside the network config.`);
        return;
    } 
    try {

        let promises: Promise<any>[] = []

        if (args.runGethNodes) {
            await hre.run('run-geth-nodes');

            console.log(`Deploying GETH contracts...`);
            const res = spawn("npx", [
                "hardhat",
                "deploy",
                "--network", "simGeth",
                "--reset"
            ])

            promises = promises.concat(new Promise(resolve=>{
                res.on('exit',function() {
                    resolve(true);
                })
            }));
        }

        if (args.runObscuroNode) {
            const enclavePromise = hre.run('run-enclave');
            const hostPromise = hre.run('run-host')
            .then(()=>hre.run("run-wallet-extension", {
                rpcUrl : rpcURL
            }));
            promises = promises.concat([enclavePromise, hostPromise])
        }

        process.on('exit', ()=>{
            hre.run("stop-wallet-extension")
        });
        

        process.on('SIGINT', ()=> { 
            exit(1) 
        });

        await Promise.all(promises);

        const { deployer } = await hre.getNamedAccounts();
        await hre.run('add-key', { 
            address: deployer 
        }).then(()=>hre.run('deploy', { 
            noWallet: true, 
            reset: true 
        }));

        await hre.run("test");
    } 
    finally 
    {
        await hre.run("stop-wallet-extension");
    }
});


declare module 'hardhat/types/config' {
    interface HardhatNetworkUserConfig {
        isSimTest?: boolean
      }    
}