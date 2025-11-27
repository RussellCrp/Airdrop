// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "forge-std/Script.sol";
import {AirdropToken} from "../src/AirdropToken.sol";
import {AirdropDistributor} from "../src/AirdropDistributor.sol";

contract Deploy is Script {
    function run() external {
        uint256 deployerKey = vm.envUint("PRIVATE_KEY");
        vm.startBroadcast(deployerKey);
        AirdropToken token = new AirdropToken("AirdropToken", "ADP", vm.addr(deployerKey));
        AirdropDistributor distributor = new AirdropDistributor(address(token), vm.addr(deployerKey));
        token.mint(address(distributor), 1_000_000 ether);
        vm.stopBroadcast();
    }
}
