#!/bin/sh
npx hardhat compile
npx hardhat clear-abi
npx hardhat export-abi

# Copy ABI file to the web part
cp abi/contracts/SimpleMixer.sol/SimpleMixer.json ../web/abi/
