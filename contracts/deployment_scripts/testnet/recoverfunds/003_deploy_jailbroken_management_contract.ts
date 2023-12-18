import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const {deployer} = await hre.getNamedAccounts();
    await hre.deployments.deploy('ManagementContract', {
        from: deployer,
        log: true
    })
};

export default func;
// No dependencies