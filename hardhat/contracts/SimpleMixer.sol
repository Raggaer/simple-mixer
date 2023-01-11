// SPDX-License-Identifier: MIT 
pragma solidity ^0.8.16;

/// @title A simple mixer using ECDSA signatures
/// @author Álvaro Carvajal Castro
contract SimpleMixer {
  address private authPublicKey;
  address private feeRecipient;
  uint256 public feeRate = 1000;
  mapping(bytes32 => bool) private usedSalt;

  // Struct used for signature computation
  struct WithdrawAction {
    uint256 amount;
    bytes32 salt;
    address to;
  }

  // Calculate the WithdrawAction struct hash
  bytes32 private immutable WITHDRAW_ACTION_IDENTIFIER = keccak256(
    abi.encodePacked("WithdrawAction(uint256 amount,bytes32 salt,address to)")
  );

  // Calculate the EIP712 domain struct hash
  bytes32 private immutable DOMAIN_IDENTIFIER = keccak256(
    abi.encodePacked(
      keccak256(bytes("EIP712Domain(string name,address verifyingContract)")),
      keccak256(bytes("SimpleMixer")),
      uint256(uint160(address(this)))
    )
  );

  /// SimpleMixer constructor
  /// @param _r Address where the fees will go to
  /// @param _rate Fee rate to apply to all withdraw operations
  constructor(address _r, uint256 _rate) {
    require(_rate > 0 && _rate < 10000, "Invalid fee rate");
    require(_r != address(0), "Invalid fee address");

    authPublicKey = msg.sender;
    feeRecipient = _r;
    feeRate = _rate;
  }

  /// Allows users to deposit Ether into the contract
  function deposit() public payable {
    require(msg.value >= 1 ether);
  }

  /// Retrieves current contract balance
  function getBalance() public view returns (uint256) {
    return address(this).balance;
  }

  /// Withdraws the given amount of Ether
  /// @param _action The WithdrawAction struct to use for checking the signature
  /// @param _signature The ECDSA server generated signature
  function withdraw(WithdrawAction memory _action, bytes memory _signature) public {
    require(!usedSalt[_action.salt], "Salt already used");
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

  /// Builds a typed data hash of the given WithdrawAction struct and EIP712 domain
  /// @param _action The WithdrawAction struct to use for building the final hash
  /// @return The final Identifier + EIP712 domain struct + WithdrawAction struct hash
  function getWithdrawTypedDataHash(WithdrawAction memory _action) private view returns (bytes32) {
    // Hash struct
    bytes32 structHash = keccak256(
      abi.encodePacked(
        WITHDRAW_ACTION_IDENTIFIER,
        _action.amount,
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

  /// Recovers a public key from the given ECDSA signature
  /// @param _hash Value hash used on the signature
  /// @param _signature ECDSA signature
  /// @return The signature public key
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
