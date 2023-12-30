import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const {deployer} = await hre.getNamedAccounts();
    const ethers = hre.ethers;

    const proxyAdminAddr = '0xc9c9a36b00dbe5e04e29220f011d4dbeb874f1e2';
    const proxyAddr = '0x80e95A1f064c79aC3CDAECdE46f2877cb8fa6290';
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