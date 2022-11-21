//SPDX-License-Identifier: Unlicense
pragma solidity ^0.8.0;

import "hardhat/console.sol";

contract Greeter {
  address private authPublicKey;
  address private feeRecipient;
  uint256 public feeRate = 1000;
  mapping(bytes32 => bool) private usedSignatures; 

  struct WithdrawAction {
    uint256 amount;
    uint256 deadline;
  }

  // Action identifier to build the EIP712 hash
  bytes32 private WITHDRAW_ACTION_IDENTIFIER = keccak256(
    abi.encodePacked("WithdrawAction(uint256 amount,uint256 deadline)")
  );

  // Domain identifier for the EIP712 hash
  bytes32 private DOMAIN_IDENTIFIER = keccak256(
    abi.encodePacked(
      keccak256(abi.encodePacked("EIP712Domain(string name,address verifyingContract)")
      keccak256(bytes("SimpleMixer")),
      uint256(uint160(address(this)))
    )
  );

  constructor(address _r, _rate) {
    require(_rate > 0 && _rate < 10000, "Invalid fee rate");

    authPublicKey = msg.sender;
    feeRecipient = _r;
    feeRate = _rate;
  }

  function deposit() public payable {
    require(msg.value >= 1 ether);
  }

  function withdraw(WithdrawAction _action, bytes memory _signature, address _to) public {
    // Check deadline
    require(_action.deadline > block.timestamp, "Withdraw deadline expired");

    require(address(this).balance > _action.amount, "Contract does not ahve enough balance");

    bytes32 h = getWithdrawTypedDataHash(_action);
    address pub = recover(h, _signature);

    // The signature should be signed by the expected public key
    require(pub == authPublicKey, "Invalid provided signature");

    usedSignatures[_signature] = true;

    // From the funds retrieve mixer fee
    uint256 feeAmount = (_action.value * feeRate) / 10000;
    msg.sender.call{value: _action.value - feeAmount}("");
    feeRecipient.call{value: feeAmount}("");
  }

  function getWithdrawTypedDataHash(WithdrawAction memory _action) private view returns (bytes32) {
    // Hash struct
    bytes32 structHash = keccak256(
      abi.encodePacked(
        WITHDRAW_ACTION_IDENTIFIER,
        _action.amount,
        _action.deadline
      )
    );

    return keccak256(
      abi.encodePacked(
        "\x19\x01",
        DOMAIN_IDENTIFIER,
        structHash
      )
    );
  }

  function recover(bytes32 _hash, bytes32 _signature) internal returns (address) {
    require(_signature.length == 65, "Invalid signature length");

    bytes32 r;
    bytes32 s;
    uint8 v;

    assembly {
      r := mload(add(_signature, 0x20))
      s := mload(add(_signature, 0x40))
      v := byte(0, mload(add(_signature, 0x60)))
    }

    return ecrecover(_hash, v, r, s);
  }
}
