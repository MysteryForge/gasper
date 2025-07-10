// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract DummyStorage {
    uint256 public slot1;
    uint256 public slot2;
    uint256 public slot3;

    /// @param _initial Initial value to set in slot 0
    constructor(uint256 _initial) {
        slot1 = _initial;
        slot2 = _initial;
        slot3 = _initial;
    }

    /// @notice Adds `_value` to `value` (slot 0)
    function addValue(uint256 _value) external {
        slot1 += _value;
        slot2 += _value;
        slot3 += _value;
    }
}