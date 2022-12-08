# SimpleMixer - Hardhat

This directory contains all the Solidity related files to test and deploy the SimpleMixer contract.

```shell
npx hardhat accounts
npx hardhat compile
npx hardhat clean
npx hardhat test
npx hardhat node
npx hardhat run scripts/deploy-mixer.js
npx hardhat help
```

## Deployment

The contract constructor expects two arguments:

- Address where fees will be sent to.
- Fee precentage to charge from every withdraw.

On the `scripts/deploy-mixer.js` there is a deeploy example.
