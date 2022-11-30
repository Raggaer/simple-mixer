const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("SimpleMixer test", function () {
  it("Mix some coins properly", async function () {
    const [owner, addr1, addr2, addr3] = await ethers.getSigners();

    const mixerFactory = await ethers.getContractFactory("SimpleMixer");
    const mixer = await mixerFactory.deploy(owner.address, 1000);
    await mixer.deployed();

    // User 1 deposits some coins
    await mixer.connect(addr1).deposit({
      value: ethers.utils.parseEther("5.0")
    });

    expect(await mixer.getBalance()).to.equal(ethers.utils.parseEther("5.0"));

    // User now has a signature from the main server
    // addr2 sends the transaction and addr3 should get the mixer coins
    const startingBalance3 = await ethers.provider.getBalance(addr3.address);
    const startingBalanceFee = await ethers.provider.getBalance(owner.address);

    const [v, s] = await generateMainServerSignature(owner, mixer.address, addr3.address);
    await mixer.connect(addr2).withdraw(v, s);

    // Check addr3 balance
    const finalBalance3 = await ethers.provider.getBalance(addr3.address);
    expect(finalBalance3).to.equal(startingBalance3.add(ethers.utils.parseEther("4.5")));

    // Check 10% fee
    const finalBalanceFee = await ethers.provider.getBalance(owner.address);
    expect(finalBalanceFee).to.equal(startingBalanceFee.add(ethers.utils.parseEther("0.5")));
  });

  it("Try to re-use same signature salt", async function () {
    const [owner, addr1, addr2, addr3] = await ethers.getSigners();

    const mixerFactory = await ethers.getContractFactory("SimpleMixer");
    const mixer = await mixerFactory.deploy(owner.address, 1000);
    await mixer.deployed();

    // User 1 deposits some coins
    await mixer.connect(addr1).deposit({
      value: ethers.utils.parseEther("5.0")
    });

    expect(await mixer.getBalance()).to.equal(ethers.utils.parseEther("5.0"));

    // User now has a signature from the main server
    // addr2 sends the transaction and addr3 should get the mixer coins
    const startingBalance3 = await ethers.provider.getBalance(addr3.address);
    const startingBalanceFee = await ethers.provider.getBalance(owner.address);

    const [v, s] = await generateMainServerSignature(owner, mixer.address, addr3.address);
    await mixer.connect(addr2).withdraw(v, s);

    // Check addr3 balance
    const finalBalance3 = await ethers.provider.getBalance(addr3.address);
    expect(finalBalance3).to.equal(startingBalance3.add(ethers.utils.parseEther("4.5")));

    // Check 10% fee
    const finalBalanceFee = await ethers.provider.getBalance(owner.address);
    expect(finalBalanceFee).to.equal(startingBalanceFee.add(ethers.utils.parseEther("0.5")));

    await expect(mixer.connect(addr2).withdraw(v, s)).to.be.revertedWith("Salt already used");
  });
});

// Generates an EIP-712 valid signature
async function generateMainServerSignature(server, _contract, _to) {
  const domain = {
    name: "SimpleMixer",
    verifyingContract: _contract
  };

  const types = {
    WithdrawAction: [
      { name: "amount", type: "uint256" },
      { name: "deadline", type: "uint256" },
      { name: "salt", type: "uint256" },
      { name: "to", type: "address" },
    ]
  };

  const value = {
    amount: ethers.utils.parseEther("5.0"),
    deadline: 2670810385,
    salt: 420,
    to: _to,
  };

  const s = await server._signTypedData(domain, types, value);
  return [value, s];
}
