// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "forge-std/Test.sol";

import {AirdropToken} from "../src/AirdropToken.sol";
import {AirdropDistributor} from "../src/AirdropDistributor.sol";

contract AirdropDistributorTest is Test {
    AirdropToken internal token;
    AirdropDistributor internal distributor;
    address internal admin = address(0xA11CE);
    address internal alice = address(0xBEEF1);
    address internal bob = address(0xBEEF2);
    uint256 internal constant ROUND_ID = 1;

    function setUp() public {
        vm.startPrank(admin);
        token = new AirdropToken("AirdropToken", "ADP", admin);
        distributor = new AirdropDistributor(address(token), admin);
        token.mint(address(distributor), 1_000_000 ether);
        vm.stopPrank();

        bytes32 leafAlice = _leaf(ROUND_ID, alice, 100 ether);
        bytes32 leafBob = _leaf(ROUND_ID, bob, 200 ether);
        bytes32 root = _parent(leafAlice, leafBob);
        vm.prank(admin);
        distributor.startRound(ROUND_ID, root, uint64(block.timestamp + 1 days));
    }

    function testClaimTransfersTokens() public {
        bytes32[] memory proof = new bytes32[](1);
        proof[0] = _leaf(ROUND_ID, bob, 200 ether);
        vm.prank(alice);
        distributor.claim(ROUND_ID, 100 ether, proof);
        assertEq(token.balanceOf(alice), 100 ether);
        assertTrue(distributor.claimed(ROUND_ID, alice));
    }

    function testCannotClaimTwice() public {
        bytes32[] memory proof = new bytes32[](1);
        proof[0] = _leaf(ROUND_ID, bob, 200 ether);
        vm.prank(alice);
        distributor.claim(ROUND_ID, 100 ether, proof);

        vm.prank(alice);
        vm.expectRevert("claimed");
        distributor.claim(ROUND_ID, 100 ether, proof);
    }

    function testInvalidProofReverts() public {
        bytes32[] memory proof = new bytes32[](1);
        proof[0] = bytes32(uint256(123));
        vm.prank(bob);
        vm.expectRevert("invalid proof");
        distributor.claim(ROUND_ID, 200 ether, proof);
    }

    function testExpiredRoundReverts() public {
        bytes32 leafAlice = _leaf(2, alice, 1 ether);
        vm.prank(admin);
        distributor.startRound(2, leafAlice, uint64(block.timestamp + 1));
        vm.warp(block.timestamp + 2);

        bytes32[] memory proof = new bytes32[](0);
        vm.prank(alice);
        vm.expectRevert("expired");
        distributor.claim(2, 1 ether, proof);
    }

    function _leaf(uint256 roundId, address account, uint256 amount) internal pure returns (bytes32) {
        return keccak256(abi.encodePacked(roundId, account, amount));
    }

    function _parent(bytes32 left, bytes32 right) internal pure returns (bytes32) {
        return keccak256(abi.encodePacked(left, right));
    }
}
