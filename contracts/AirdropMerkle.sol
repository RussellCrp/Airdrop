// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {MerkleProof} from "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";

/// @title Merkle Airdrop contract
/// @notice Verify eligibility with Merkle proof and distribute an existing ERC20 token.
contract AirdropMerkle {
    using SafeERC20 for IERC20;

    IERC20 public immutable token;
    address public owner;

    // Current merkle root
    bytes32 public merkleRoot;

    // index => claimed
    mapping(uint256 => bool) public claimed;

    event Claimed(uint256 indexed index, address indexed account, uint256 amount);
    event MerkleRootUpdated(bytes32 oldRoot, bytes32 newRoot);
    event OwnerUpdated(address indexed oldOwner, address indexed newOwner);
    event EmergencyWithdraw(address indexed to, uint256 amount);

    modifier onlyOwner() {
        require(msg.sender == owner, "not owner");
        _;
    }

    constructor(address token_, bytes32 merkleRoot_, address owner_) {
        require(token_ != address(0), "token is zero");
        require(owner_ != address(0), "owner is zero");
        token = IERC20(token_);
        merkleRoot = merkleRoot_;
        owner = owner_;
    }

    /// @notice claim airdrop
    /// @param index index in merkle tree
    /// @param account receiver address (must equal msg.sender)
    /// @param amount token amount
    /// @param merkleProof proof from backend
    function claim(
        uint256 index,
        address account,
        uint256 amount,
        bytes32[] calldata merkleProof
    ) external {
        require(msg.sender == account, "forbidden");
        require(!claimed[index], "already claimed");

        bytes32 node = keccak256(abi.encodePacked(index, account, amount));
        require(MerkleProof.verify(merkleProof, merkleRoot, node), "invalid proof");

        claimed[index] = true;
        token.safeTransfer(account, amount);

        emit Claimed(index, account, amount);
    }

    /// @notice Owner can update root for a new round of airdrop.
    function setMerkleRoot(bytes32 newRoot) external onlyOwner {
        bytes32 old = merkleRoot;
        merkleRoot = newRoot;
        emit MerkleRootUpdated(old, newRoot);
    }

    /// @notice Transfer ownership.
    function transferOwnership(address newOwner) external onlyOwner {
        require(newOwner != address(0), "zero owner");
        address old = owner;
        owner = newOwner;
        emit OwnerUpdated(old, newOwner);
    }

    /// @notice Emergency withdraw remaining tokens.
    function emergencyWithdraw(address to, uint256 amount) external onlyOwner {
        token.safeTransfer(to, amount);
        emit EmergencyWithdraw(to, amount);
    }
}


