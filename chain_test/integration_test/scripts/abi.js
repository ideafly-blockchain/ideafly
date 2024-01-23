const fs = require("fs");

const Validators_ABI = loadABI("Validators");
const AddressList_ABI = loadABI("AddressList");
const Governance_ABI = loadABI("Governance");
const VotePool_ABI = loadABI("VotePool");

function loadABI(name) {
    let data = fs.readFileSync(`./abi/${name}.json`, "utf-8");
    return JSON.parse(data).abi;
}

module.exports = {
    Validators_ABI,
    AddressList_ABI,
    Governance_ABI,
    VotePool_ABI,
};
