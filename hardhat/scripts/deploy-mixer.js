const hre = require("hardhat");
const mixerFee = 1000;

async function main() {
  const [owner, addr1, addr2, addr3] = await ethers.getSigners();

  const mixerFactory = await hre.ethers.getContractFactory("SimpleMixer");
  const mixer = await mixerFactory.deploy(owner.address, 1000);

  await mixer.deployed();

  console.log("SimpleMixer mixing fee:", (mixerFee / 100) + "%");
  console.log("SimpleMixer fee destination address:", owner.address);
  console.log("SimpleMixer deployed to:", mixer.address);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
