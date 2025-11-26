// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Script.sol";
import {AirdropMerkle} from "../contracts/AirdropMerkle.sol";

contract DeployAirdrop is Script {
    function run() external {
        address token = vm.envAddress("AIRDROP_TOKEN");
        bytes32 root = vm.envBytes32("AIRDROP_MERKLE_ROOT");

        vm.startBroadcast();
        new AirdropMerkle(token, root, msg.sender);
        vm.stopBroadcast();
    }
}


