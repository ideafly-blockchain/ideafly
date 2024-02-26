// We require the Hardhat Runtime Environment explicitly here. This is optional
// but useful for running the script in a standalone fashion through `node <script>`.
//
// When running the script with `npx hardhat run <script>` you'll find the Hardhat
// Runtime Environment's members available in the global scope.
const hre = require("hardhat");
const { expect } = require("chai");
const { BigNumber } = require("ethers");
const ethers = hre.ethers;

const validatorConAddr = "0x000000000000000000000000000000000000d001";
const adminAddr = "0x3a696FeAe901DAe50967F28D7A2225577052F394";
const managerAddr = "0x9E737Ee8bDc132c349dE7801Efbc9e12f4FE99e9";

const nodeAddrPoa4 = "0x1d342eefd8a93b513a17e8cd9fa9619c35aa21fa";

const url = hre.config.networks.local.url;

const { Validators_ABI, VotePool_ABI } = require("./abi");

let admin, manager, receipt, tx;
let vp_poa4;
const EPOCH = 5;

async function step1(validator) {
    try {
        tx = await validator.addValidator(nodeAddrPoa4, managerAddr, 100, 1, {
            gasLimit: 5000000,
            gasPrice: 0x12a05f200,
        });
        receipt = await tx.wait();
        expect(receipt.status).equal(1);
        // console.log(JSON.stringify(receipt));

        vp_poa4 = BigNumber.from(receipt.events[0].data).toHexString();
        console.log({ vp_poa4 });
        console.log("✅ SUCCESS step1: addValidator", nodeAddrPoa4, managerAddr, 100, 1);
    } catch (e) {
        console.error("❌ FAILED step1: addValidator", nodeAddrPoa4, managerAddr, 100, 1);
        throw e;
    }

    const vote = await ethers.getContractAt(VotePool_ABI, vp_poa4, manager);
    try {
        tx = await vote.addMargin({
            gasLimit: 5000000,
            gasPrice: 0x12a05f200,
            value: ethers.utils.parseEther("1.0"),
        });
        receipt = await tx.wait();
        expect(receipt.status).equal(1);
        // console.log(JSON.stringify(receipt));
        bnadd = Number(receipt.blockNumber);
        next = (parseInt(bnadd / EPOCH) + 2) * EPOCH + 1
        console.log({ bnadd, next });
        expect(bnadd).gt(0);
        console.log("✅ SUCCESS step1: addMargin 1.0 eth");
    } catch (e) {
        console.error("❌ FAILED step1: addMargin 1.0 eth");
        throw e;
    }

    const ds = next + 1 - bnadd;
    await sleep(ds * 2000);
    let generated = false;
    for (let i = next; i <= next + 4; i += 1) {
        await sleep(2000);
        const block = await ethers.provider.getBlock(i);
        console.log({ blockNumber: i, miner: block.miner });
        if (block.miner.toLowerCase() == nodeAddrPoa4.toLowerCase()) {
            generated = true;
            break;
        }
    }
    if (generated) {
        console.log("✅ SUCCESS step1: generate");
    } else {
        console.error("❌ FAILED step1: generate");
    }
    expect(generated).eq(true);


    try {
        const amount = await vote.getValidatorPendingReward();
        console.log(BigNumber.from(amount).toString());

        console.log("✅ SUCCESS step1: poa4 before rewarded",BigNumber.from(amount).toString());

        // if (ethers.utils.parseEther("5000.0").eq(BigNumber.from(amount))) {
        //     console.log("✅ SUCCESS step1: poa4 was rewarded as expected");
        // } else {
        //     console.error("❌ FAILED step1: vote.margin");
        //     throw new Error("FAILED step1: vote.margin");
        // }
    } catch (e) {
        console.error("❌ FAILED step1: vote.getValidatorPendingReward");
        throw e;
    }
}

function sleep(ms) {
    console.log({ sleep: `${ms / 1000} s` });
    return new Promise(resolve => setTimeout(resolve, ms));
}

async function main() {
    console.log("link to node: ", url);
    admin = await ethers.getSigner(adminAddr);
    manager = await ethers.getSigner(managerAddr);
    console.log(`admin addr: ${adminAddr}`);
    console.log(`manager addr: ${managerAddr}`);

    const validator = await ethers.getContractAt(Validators_ABI, validatorConAddr, admin);
    console.log(".");
    try {
        tx = await validator.updateParams(1, 1, 4, 1, {
            gasLimit: 5000000,
            gasPrice: 0x12a05f200,
        });
        console.log(".");
        receipt = await tx.wait();
        console.log(".");
        expect(receipt.status).equal(1);
        console.log("✅ SUCCESS init: set updateParams 1,1,4,1");
    } catch (e) {
        console.error("❌ FAILED init: set updateParams 1,1,4,1");
        throw e;
    }

    await step1(validator);
}

// We recommend this pattern to be able o use async/await everywhere
// and properly handle errors.
main()
    .then(() => process.exit(0))
    .catch(error => {
        console.error(error);
        process.exit(1);
    });

