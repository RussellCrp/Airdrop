// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import {AirdropMerkle} from "../contracts/AirdropMerkle.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract MockToken is IERC20 {
    string public name = "Mock";
    string public symbol = "MCK";
    uint8 public decimals = 18;

    mapping(address => uint256) public override balanceOf;
    mapping(address => mapping(address => uint256)) public override allowance;
    uint256 public override totalSupply;

    function transfer(address to, uint256 amount) external override returns (bool) {
        _transfer(msg.sender, to, amount);
        return true;
    }

    function approve(address spender, uint256 amount) external override returns (bool) {
        allowance[msg.sender][spender] = amount;
        emit Approval(msg.sender, spender, amount);
        return true;
    }

    function transferFrom(address from, address to, uint256 amount) external override returns (bool) {
        uint256 allowed = allowance[from][msg.sender];
        require(allowed >= amount, "allowance");
        allowance[from][msg.sender] = allowed - amount;
        _transfer(from, to, amount);
        return true;
    }

    function _transfer(address from, address to, uint256 amount) internal {
        require(balanceOf[from] >= amount, "balance");
        balanceOf[from] -= amount;
        balanceOf[to] += amount;
        emit Transfer(from, to, amount);
    }

    function mint(address to, uint256 amount) external {
        totalSupply += amount;
        balanceOf[to] += amount;
        emit Transfer(address(0), to, amount);
    }
}

contract AirdropMerkleTest is Test {
    MockToken token;
    AirdropMerkle airdrop;

    address user = address(0x1234);

    function setUp() public {
        token = new MockToken();
        // simple one-leaf merkle tree where leaf = keccak256(abi.encodePacked(0, user, 100))
        uint256 index = 0;
        uint256 amount = 100;
        bytes32 leaf = keccak256(abi.encodePacked(index, user, amount));
        bytes32 root = leaf;

        airdrop = new AirdropMerkle(address(token), root, address(this));
        token.mint(address(airdrop), amount);
    }

    function testClaim() public {
        uint256 index = 0;
        uint256 amount = 100;
        bytes32[] memory proof = new bytes32[](0);

        vm.prank(user);
        airdrop.claim(index, user, amount, proof);

        assertEq(token.balanceOf(user), amount);
        assertTrue(airdrop.claimed(index));
    }
}


