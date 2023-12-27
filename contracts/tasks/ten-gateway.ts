import { task } from "hardhat/config";

import * as url from 'node:url';

import { spawn } from 'node:child_process';
import * as path from "path";
// @ts-ignore
import http from 'http';


task("ten:gateway:start:local")
    .addFlag('withStdOut')
    .addParam('rpcurl', "Node rpc endpoint where the Ten gateway should connect to.")
    .addOptionalParam('port', "Port that the Ten gateway will open for incoming requests.", "3001")
    .setAction(async function(args, hre) {
        const nodeUrl = url.parse(args.rpcurl);
        const tenGatewayPath = path.resolve(hre.config.paths.root, "../tools/walletextension/bin/wallet_extension_linux");
        const weProcess = spawn(tenGatewayPath, [
            `-portWS`, `${args.port}`,
            `-nodeHost`, `${nodeUrl.hostname}`,
            `-nodePortWS`, `${nodeUrl.port}`
        ]);

        console.log("Waiting for Ten gateway to start");
        await new Promise((resolve, fail)=>{
            const timeoutSchedule = setTimeout(fail, 60_000);
            weProcess.stdout.on('data', (data: string) => {
                if (args.withStdOut) {
                    console.log(data.toString());
                }

                if (data.includes("Ten Gateway started")) {
                    clearTimeout(timeoutSchedule);
                    resolve(true)
                }
            });

            weProcess.stderr.on('data', (data: string) => {
                console.log(data.toString());
            });
        });

        console.log("Ten gateway started successfully");
        return weProcess;
    });


// This is not to be used for internal development. It is targeted at external devs when the obscuro hh plugin is finished!
task("ten:gateway:start:docker", "Starts up the Ten gateway docker container.")
    .addFlag('wait')
    .addParam('dockerImage',
        'The docker image to use for the Ten gateway',
        'testnetobscuronet.azurecr.io/obscuronet/walletextension') // TODO (@ziga) - change after renaming it
    .addParam('rpcUrl', "Which network to pick the node connection info from?")
    .setAction(async function(args, hre) {
        const docker = new dockerApi.Docker({ socketPath: '/var/run/docker.sock' });

        const parsedUrl = url.parse(args.rpcUrl)

        const container = await docker.container.create({
            Image: args.dockerImage,
            Cmd: [
                "--port=3000",
                "--portWS=3001",
                `--nodeHost=${parsedUrl.hostname}`,
                `--nodePortWS=${parsedUrl.port}`
            ],
            ExposedPorts: { "3000/tcp": {}, "3001/tcp": {}, "3000/udp": {}, "3001/udp": {} },
            PortBindings:  { "3000/tcp": [{ "HostPort": "3000" }], "3001/tcp": [{ "HostPort": "3001" }] }
        })


        process.on('SIGINT', ()=>{
            container.stop();
        })

        await container.start();

        const stream: any = await container.logs({
            follow: true,
            stdout: true,
            stderr: true
        })

        console.log(`\nTen gateway{ ${container.id.slice(0, 5)} } %>\n`);
        const startupPromise = new Promise((resolveInner)=> {
            stream.on('data', (msg: any)=> {
                const message = msg.toString();

                console.log(message);
                if(message.includes("Ten Gateway started")) {
                    console.log(`Wallet - success!`);
                    resolveInner(true);
                }
            });

            setTimeout(resolveInner, 20_000);
        });

        await startupPromise;
        console.log("\n[ . . . ]\n");


        if (args.wait) {
            await container.wait();
        }
    });

// This is not to be used for internal development. It is targeted at external devs when the obscuro hh plugin is finished!
task("ten:gateway:stop:docker", "Stops the docker container with matching image name.")
    .addParam('dockerImage',
        'The docker image to use for the Ten gateway',
        'testnetobscuronet.azurecr.io/obscuronet/walletextension') // TODO @ziga change when renaming wallet extensions
    .setAction(async function(args, hre) {
        const docker = new dockerApi.Docker({ socketPath: '/var/run/docker.sock' });
        const containers = await docker.container.list();

        const container = containers.find((c)=> {
            const data : any = c.data;
            return data.Image == 'testnetobscuronet.azurecr.io/obscuronet/walletextension'
        })

        await container?.stop()
    });

const { URL } = require('url');

task("ten:gateway:join-authenticate", "Joins and authenticates the gateway for a specific address")
    .addParam("address", "The address which to use for authentication")
    .setAction(async function(args, hre) {
        async function joinGateway(url = 'http://127.0.0.1:3000/v1/join') {
            return new Promise((resolve, reject) => {
                http.get(url, (response) => {
                    if (response.statusCode >= 300 && response.statusCode < 400 && response.headers.location) {
                        // Resolve the new location against the original URL
                        const newUrl = new URL(response.headers.location, url).toString();
                        return resolve(joinGateway(newUrl));
                    }

                    if (response.statusCode !== 200) {
                        return reject(new Error(`Server responded with status code: ${response.statusCode}`));
                    }

                    let chunks = [];
                    response.on('data', (chunk) => chunks.push(chunk));
                    response.on('end', () => resolve(Buffer.concat(chunks).toString()));
                }).on('error', reject);
            });
        }


        // authenticateWithGateway function to authenticate using the token
        async function authenticateWithGateway(encryptionToken) {
            const typedData = {
                types: {
                    EIP712Domain: [
                        { name: "name", type: "string" },
                        { name: "version", type: "string" },
                        { name: "chainId", type: "uint256" },
                    ],
                    Authentication: [
                        { name: "Encryption Token", type: "address" },
                    ],
                },
                primaryType: "Authentication",
                domain: {
                    name: "Ten",
                    version: "1.0",
                    chainId: 443, // TODO @ziga - can we get this from some config in typescript?
                },
                message: {
                    "Encryption Token": "0x"+encryptionToken
                },
            };

            const messageData = JSON.stringify(typedData);
            const signer =  await hre.ethers.getSigner(args.address);

            const signature = await signer.provider.send('eth_signTypedData_v4', [
                signer.address,
                messageData
            ]);

            const signedData = { "signature": signature, "address": args.address };

            // Call the authenticate function with the new URL and signed data
            const url = `http://127.0.0.1:3000/v1/authenticate?token=${encryptionToken}`;
            return authenticate(url, signedData);
        }

        // authenticate function to make the POST request
        async function authenticate(url, signedData) {
            return new Promise((resolve, reject) => {
                const makeRequest = (url) => {
                    const options = {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json'
                        }
                    };

                    const req = http.request(url, options, (response) => {
                        if (response.statusCode === 301 || response.statusCode === 302) {
                            // Handle redirect
                            const newUrl = new URL(response.headers.location, url).toString();
                            makeRequest(newUrl); // Resend POST request to new URL
                        } else if (response.statusCode < 200 || response.statusCode >= 300) {
                            reject(new Error(`Server responded with status code: ${response.statusCode}`));
                        } else {
                            let chunks = [];
                            response.on('data', (chunk) => chunks.push(chunk));
                            response.on('end', () => resolve(Buffer.concat(chunks).toString()));
                        }
                    });

                    req.on('error', reject);
                    req.write(JSON.stringify(signedData));
                    req.end();
                };

                makeRequest(url);
            });
            }


        try {
            let encryptionToken = await joinGateway();
            console.log("Encryption token: ", encryptionToken);

            let authenticationResult = await authenticateWithGateway(encryptionToken);
            console.log("Authentication result: ", authenticationResult);
        } catch (error) {
            console.error("Error: ", error);
        }
    });
