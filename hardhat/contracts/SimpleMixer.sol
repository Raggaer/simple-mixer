//SPDX-License-Identifier: Unlicense
pragma solidity ^0.8.0;

import "hardhat/console.sol";

contract SimpleMixer {
  address private authPublicKey;
  address private feeRecipient;
  uint256 public feeRate = 1000;
  mapping(uint256 => bool) private usedSalt;

  struct WithdrawAction {
    uint256 amount;
    uint256 deadline;
    uint256 salt;
    address to;
  }

  // Action identifier to build the EIP712 hash
  bytes32 private WITHDRAW_ACTION_IDENTIFIER = keccak256(
    abi.encodePacked("WithdrawAction(uint256 amount,uint256 deadline,uint256 salt,address to)")
  );

  // Domain identifier for the EIP712 hash
  bytes32 private DOMAIN_IDENTIFIER = keccak256(
    abi.encodePacked(
      keccak256(abi.encodePacked("EIP712Domain(string name,address verifyingContract)")),
      keccak256(bytes("SimpleMixer")),
      uint256(uint160(address(this)))
    )
  );

  constructor(address _r, uint256 _rate) {
    require(_rate > 0 && _rate < 10000, "Invalid fee rate");

    authPublicKey = msg.sender;
    feeRecipient = _r;
    feeRate = _rate;
  }

  function deposit() public payable {
    require(msg.value >= 1 ether);
  }

  function getBalance() public view returns (uint256) {
    return address(this).balance;
  }

  function withdraw(WithdrawAction memory _action, bytes memory _signature) public {
    require(!usedSalt[_action.salt], "Salt already used");
    require(_action.deadline > block.timestamp, "Withdraw deadline expired");
    require(address(this).balance >= _action.amount, "Contract does not have enough balance");

    usedSalt[_action.salt] = true;

    bytes32 h = getWithdrawTypedDataHash(_action);
    address pub = recover(h, _signature);

    // The signature should be signed by the expected public key
    require(pub == authPublicKey, "Invalid provided signature");

    // From the funds retrieve mixer fee
    uint256 feeAmount = (_action.amount * feeRate) / 10000;

    bool success;
    (success, ) = _action.to.call{value: _action.amount - feeAmount}("");
    require(success, "Unable to send value");
    (success, ) = feeRecipient.call{value: feeAmount}("");
    require(success, "Unable to send fee value");
  }

  function getWithdrawTypedDataHash(WithdrawAction memory _action) private view returns (bytes32) {
    // Hash struct
    bytes32 structHash = keccak256(
      abi.encodePacked(
        WITHDRAW_ACTION_IDENTIFIER,
        _action.amount,
        _action.deadline,
        _action.salt,
        uint256(uint160(_action.to))
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

  function recover(bytes32 _hash, bytes memory _signature) internal pure returns (address) {
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
