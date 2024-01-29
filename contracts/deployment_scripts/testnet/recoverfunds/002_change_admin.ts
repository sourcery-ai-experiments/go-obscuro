import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const {deployer} = await hre.getNamedAccounts();
    const ethers = hre.ethers;

    const proxyAdminAddr = '0xEb3710985693057D9164723a9B6F320569A48b24';
    const proxyAddr = '0x52a8D76CF9840cb7b9797d296b23042E031F35F4';
    const safeWalletAddr = '0xeA052c9635F1647A8a199c2315B9A66ce7d1e2a7';

    const proxyAdmin =await ethers.getContractAt('IDefaultProxyAdmin', proxyAdminAddr, await ethers.getSigner(deployer))
    const tx = await proxyAdmin.changeProxyAdmin(proxyAddr, safeWalletAddr)
    const receipt = await tx.wait();
    if (receipt.status != 1) {
        console.error('Unable to change proxy admin');
        throw Error('Failed: unable to change admin');
    }
};

export default func;
// No dependencies