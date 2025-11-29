// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "./MerkleProof.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/// @title Merkle-based airdrop distributor with replay protection.
contract AirdropDistributor is Ownable {
    using SafeERC20 for IERC20;

    struct RoundConfig {
        bytes32 merkleRoot;
        uint64 claimDeadline;
        bool active;
    }

    IERC20 public immutable token;
    mapping(uint256 => RoundConfig) public rounds;
    mapping(uint256 => mapping(address => bool)) public claimed;

    event RoundStarted(uint256 indexed roundId, bytes32 indexed merkleRoot, uint64 claimDeadline);
    event RoundClosed(uint256 indexed roundId);
    event Claimed(uint256 indexed roundId, address indexed account, uint256 amount, bytes32 leaf);

    constructor(address token_, address admin_) Ownable(admin_) {
        require(token_ != address(0), "token required");
        token = IERC20(token_);
    }

    function startRound(uint256 roundId, bytes32 merkleRoot, uint64 claimDeadline) external onlyOwner {
        require(roundId > 0, "round=0");
        require(merkleRoot != bytes32(0), "root required");
        require(claimDeadline > block.timestamp, "deadline invalid");
        rounds[roundId] = RoundConfig({merkleRoot: merkleRoot, claimDeadline: claimDeadline, active: true});
        emit RoundStarted(roundId, merkleRoot, claimDeadline);
    }

    function closeRound(uint256 roundId) external onlyOwner {
        RoundConfig storage config = rounds[roundId];
        require(config.active, "round not active");
        config.active = false;
        emit RoundClosed(roundId);
    }

    function claim(uint256 roundId, uint256 amount, bytes32[] calldata proof) external {
        RoundConfig memory config = rounds[roundId];
        require(config.active, "inactive");
        require(block.timestamp <= config.claimDeadline, "expired");
        require(!claimed[roundId][msg.sender], "claimed");
        require(amount > 0, "amount=0");

        bytes32 leaf = keccak256(abi.encodePacked(roundId, msg.sender, amount));
        require(MerkleProof.verify(proof, config.merkleRoot, leaf), "invalid proof");

        claimed[roundId][msg.sender] = true;
        token.safeTransfer(msg.sender, amount);
        emit Claimed(roundId, msg.sender, amount, leaf);
    }
}
