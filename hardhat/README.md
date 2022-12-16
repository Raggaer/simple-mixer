# SimpleMixer - Hardhat

This directory contains all the Solidity related files to test and deploy the SimpleMixer contract.
Any user is able to call the contract `deposit` method, retrieve a signature from the expected signer and then call `withdraw` method.

- `deposit` Just a simple payable function.
- `withdraw` Generates a typed data hash and alongside the provided ECDSA signature checks if the public key matches.

## Deployment

The contract constructor expects two arguments:

- Address where fees will be sent to.
- Fee precentage to charge from every withdraw.

On the `scripts/deploy-mixer.js` there is a deployment example.

## Testing

The following command will execute all the repository tests:

```shell
npx hardhat test
```
