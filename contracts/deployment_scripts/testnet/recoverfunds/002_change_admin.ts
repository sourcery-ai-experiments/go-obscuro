import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const {deployer} = await hre.getNamedAccounts();
    const ethers = hre.ethers;

    const proxyAdminAddr = '0x9c46fb7c0fe2effc0795046d142c89c4c09fc1c6';
    const proxyAddr = '0x9a3031a39a34887516b66b455c490c6eb8d048ee';
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