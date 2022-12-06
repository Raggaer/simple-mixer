# SimpleMixer

SimpleMixer is a simple centralized Ethereum mixer that uses ECDSA signatures for withdrawing funds.
Users are able to deposit funds into the contract but they need a signature from the central server in order to retrieve the funds back.

The central server will emit a new EIP712 signed message. 
The user can now use that signature from a different public address to interact with the contract **withdraw** method.
This method will check if the signature was signed by the central address and send the funds (minus a configurable fee).

## Repository structure

- hardhat: Directory containing all Solidity related files. Tests, deployment and smart-contracts.
