// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

contract StateFiller {
    uint[] private items;
    address public owner;

    constructor(uint256 count) {
        owner = msg.sender;
        // Fill storage with 'count' number of items
        for (uint i = 0; i < count; i++) {
            items.push(i);
        }
    }

    function deleteRandomState() external {
        require(items.length > 0, "No more state to delete");

        // Generate pseudo-random index using block data and caller
        uint index = uint(keccak256(abi.encodePacked(
            items.length,
            block.timestamp,
            msg.sender
        ))) % items.length;

        // Move last element to the random position and pop
        items[index] = items[items.length - 1];
        items.pop();
    }

    // View function to get item at index
    function getItems(uint256 index) external view returns (uint256) {
        return items[index];
    }

    function size() external view returns (uint256) {
        return items.length;
    }
}