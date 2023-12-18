import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    return;
    const {deployer} = await hre.getNamedAccounts();
    const ethers = hre.ethers;
    const proxyAdmin =await ethers.getContractAt('IDefaultProxyAdmin', '0xFFa0feA0cb522f3d0CE644B9a111215ACbEb061B', await ethers.getSigner(deployer))
    const tx = await proxyAdmin.changeProxyAdmin('0xEde0fC70C4B67916B8d2037dE24cD18BF26e5069', '0xeA052c9635F1647A8a199c2315B9A66ce7d1e2a7')
    const receipt = await tx.wait();
    if (receipt.status != 1) {
        console.error('Unable to change proxy admin');
        throw Error('Failed: unable to change admin');
    }
};

export default func;
// No dependencies