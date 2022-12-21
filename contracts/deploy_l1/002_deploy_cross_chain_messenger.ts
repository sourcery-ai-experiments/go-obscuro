import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';
import { ManagementContract } from '../typechain-types/src/management';

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const { 
        deployments, 
        getNamedAccounts
    } = hre;

    const {deployer} = await getNamedAccounts();
    console.log(`002 - deployer ${deployer}`);

    const messageBusAddress : string = await deployments.read("ManagementContract", {}, "messageBus");

    await deployments.deploy('CrossChainMessenger', {
        from: deployer,
        args: [ messageBusAddress ],
        log: true,
    });
};

export default func;
func.tags = ['CrossChainMessenger', 'CrossChainMessenger_deploy'];
func.dependencies = ['ManagementContract'];
