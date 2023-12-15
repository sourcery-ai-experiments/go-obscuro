// SPDX-License-Identifier: Apache 2
pragma solidity ^0.8.0;


interface IDefaultProxyAdmin {
    function changeProxyAdmin(address proxy, address newAdmin) external;
}