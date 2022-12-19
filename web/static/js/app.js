// Setup RPC ReadOnly providers
const rpcProvider = new ethers.providers.JsonRpcProvider(rpcURL);
const contract_RO = new ethers.Contract(mixerContractAddress, mixerContractABI, rpcProvider);

document.getElementById("deposit").addEventListener("click", async function() {
  hideAlert();

  let provider;
  let signer;
  try {
    [provider, signer] = await getProviderAndSigner();
  } catch(e) {
    showAlert("Unable to retrieve Web3 signer. Please try again.");
    return;
  }

  // Deposit funds and get txid
  const txid = await depositFunds(signer);

  // Sign message and send to the server
  const signature = await getServerSignedMessage(txid, signer);
});

// Sends a signed message to the server
// The server returns a valid EIP-712 signature for later use
async function getServerSignedMessage(txid, signer) {
  const msg = await signer.signMessage(txid);

  const addr = await signer.getAddress();

  // Send message to the server
  const data = await fetch("/api/sign", {
    method: "POST",
    headers: {
      "Content-Type": "application/x-www-form-urlencoded"
    },
    body: "msg=" + msg + "&signer=" + addr  + "&tx=" + txid
  });
}

// Calls contract deposit and returns the transaction hash
async function depositFunds(signer) {
  let depositAmount;

  try {
    depositAmount = ethers.utils.parseEther(document.getElementById("amount").value);
  } catch(e) {
    showAlert("Invalid Ether deposit amount. Please try again.");
    return;
  }

  let transactionId;

  // Deposit some Ether
  try {
    const contract_WO = new ethers.Contract(mixerContractAddress, mixerContractABI, signer);
    const receipt = await contract_WO.deposit({
      value: depositAmount
    });
    await receipt.wait();

    transactionId = receipt.hash;
  } catch(e) {
    showAlert("Unable to call contract deposit. Please try again.");
    return;
  }

  return transactionId;
}

// Retrieve connected window provider (ie: MetaMask) and any account as a signer
async function getProviderAndSigner() {
  const provider = new ethers.providers.Web3Provider(window.ethereum, "any");
  await provider.send("eth_requestAccounts", []);

  const signer = provider.getSigner();
  
  return [provider, signer];
}

// Retrieve current contract balance
async function getBalance(r) {
  const balance = await r.getBalance();
  return ethers.utils.formatEther(balance);
}

function hideAlert() {
  const alert = document.getElementById("err");
  alert.style.display = "none";
}

function showAlert(msg) {
  const alert = document.getElementById("err");
  alert.innerHTML = msg;
  alert.style.display = "block";
}
