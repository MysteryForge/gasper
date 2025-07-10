// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/access/Ownable.sol";

contract BatchFunder is Ownable {
    event Funded(address indexed recipient, uint256 amount);
    event BatchFunded(uint256 totalRecipients, uint256 totalAmount);

    constructor() Ownable(msg.sender) {}

    /**
     * @notice Batch sends ETH to a list of recipients. Only callable by owner.
     * @param recipients Array of recipient addresses.
     * @param amount Amount of ETH to send to each address (in wei).
     */
    function batchSend(address[] calldata recipients, uint256 amount) external payable {
        for (uint256 i = 0; i < recipients.length; i++) {
            require(recipients[i] != address(0), "Invalid recipient");
            (bool success, ) = payable(recipients[i]).call{value: amount}("");
            require(success, "ETH transfer failed");
            emit Funded(recipients[i], amount);
        }

        emit BatchFunded(recipients.length, msg.value);
    }

    /**
     * @notice Withdraws ETH accidentally sent to the contract. Only callable by owner.
     * @param to Address to withdraw to.
     */
    function rescueEth(address payable to) external onlyOwner {
        require(to != address(0), "Invalid address");
        to.transfer(address(this).balance);
    }

    receive() external payable {}
}
