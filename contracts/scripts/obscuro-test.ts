import { task } from "hardhat/config";
import 'hardhat/types/config';

import process, {exit} from 'process';
import { HardhatNetworkUserConfig } from "hardhat/types/config";

import { spawn } from 'node:child_process';
import * as path from "path";
import { HardhatRuntimeEnvironment } from "hardhat/types";

async function startGethNodes(hre: HardhatRuntimeEnvironment) {
    const integrationDir = path.resolve(hre.config.paths.root, "../integration");
    const gethNetworkMain = path.resolve(integrationDir, "gethnetwork/main");
    const gethnetwork = spawn('go', [
        'run', 
        gethNetworkMain, 
        "--prefundedAddrs", "0x323AefbFC16159655514846a9e5433C457de9389,0x0654D8B60033144D567f25bF41baC1FB0D60F23B,0x13E23Ca74DE0206C56ebaE8D51b5622EFF1E9944",
        "--startPort=8000",
        "--websocketStartPort=9000",
        "--numNodes=3"
    ]);

    await new Promise(resolve=>{
        gethnetwork.stdout.on('data', (data: string) => {
           if (data.includes("Geth network started")) {
            resolve(true)
           }
        });
    })
    
}

async function startEnclave(hre: HardhatRuntimeEnvironment) {
    const enclaveDir = path.resolve(hre.config.paths.root, "../go/enclave/main");
    const enclaveProc = spawn('go', [
        'run', 
        enclaveDir,
        "--hostID=0x0654D8B60033144D567f25bF41baC1FB0D60F23B",
        "--address=127.0.0.1:11000",
        "--nodeType=sequencer",
        "--useInMemoryDB=true",
        "--managementContractAddress=0xeDa66Cc53bd2f26896f6Ba6b736B1Ca325DE04eF",
        "--erc20ContractAddresses=0x6d2994ACb911CFceaeE6C36D881cbDFE2F9553B0,0x26c62148Cf06C9742b8506A2BCEcd7d72E51A206",
        "--profilerEnabled=false",
        "--hostAddress=127.0.0.1:10000",
        "--logPath=sys_out",
        "--logLevel=4",
        "--sequencerID=0x0654D8B60033144D567f25bF41baC1FB0D60F23B",
        "--messageBusAddress=0xFD03804faCA2538F4633B3EBdfEfc38adafa259B"
    ], { detached: false });

    process.on('beforeExit', ()=> { 
        console.log(`Stopping`)
        enclaveProc.kill() 
    });

    return new Promise((resolve)=>{
        const timer = setTimeout(resolve, 60_000);
        enclaveProc.stdout.on('data', (data: string) => {
           // console.log(data.toString());
            if (data.includes("Obscuro enclave service started.")) {
                clearTimeout(timer);
                resolve(true);
            }
        });
    })   
}

async function startHost(hre: HardhatRuntimeEnvironment) {
    const hostDir = path.resolve(hre.config.paths.root, "../go/host/main");
    const hostProc = spawn('go', [
        "run",
        hostDir,
        "--l1NodeHost=127.0.0.1",
        "--l1NodePort=9000",
        "--enclaveRPCAddress=127.0.0.1:11000",
        "--rollupContractAddress=0xeDa66Cc53bd2f26896f6Ba6b736B1Ca325DE04eF",
        "--privateKey=8ead642ca80dadb0f346a66cd6aa13e08a8ac7b5c6f7578d4bac96f5db01ac99",
        "--clientRPCHost=127.0.0.1",
        "--isGenesis=true",
        "--nodeType=sequencer",
        "--logPath=sys_out",
        "--logLevel=4",
        "--profilerEnabled=false",
        "--p2pPublicAddress=127.0.0.1:10000"
    ]);

    process.on('beforeExit', ()=> { 
        console.log(`Stopping host ...`)
        hostProc.kill() 
    });

    return new Promise((resolve)=>{
        const timer = setTimeout(resolve, 60_000);
        hostProc.stdout.on('data', (data: string) => {
           // console.log(data.toString())
            if (data.includes("Started P2P networking")) {
                clearTimeout(timer);
                resolve(true);
            }
        });
    }) 
}

task("sim-test", "Runs the tests. Takes into account if the network has obscuroEncRpcUrl")
.addFlag("runObscuroNode", "Will start a docker container with the node and test against it.")
.addFlag("runGethNodes", "Will start geth nodes to use as layer 1.")
.setAction(async function(args, hre, runSuper) {
    const rpcURL = (hre.network.config as HardhatNetworkUserConfig).obscuroEncRpcUrl;
    if (!rpcURL) {
        await runSuper();
        return;
    } 

    let promises: Promise<any>[] = []

    if (args.runGethNodes) {
        console.log(`Booting up GETH...`);
        await startGethNodes(hre);

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
        console.log(`Booting up Obscuro...`)
        const enclavePromise = startEnclave(hre);
        const hostPromise = startHost(hre).then(()=>hre.run("run-wallet-extension", {rpcUrl : rpcURL}));
        promises = promises.concat([enclavePromise, hostPromise])
    }

    process.on('exit', ()=>{
        hre.run("stop-wallet-extension")
    });
    

    process.on('SIGINT', ()=> { 
        exit(1) 
    });

    try {

        await Promise.all(promises);

        const { deployer } = await hre.getNamedAccounts();
        await hre.run('add-key', { address: deployer })
                .then(()=>hre.run('deploy', { noWallet: true, reset: true }));

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