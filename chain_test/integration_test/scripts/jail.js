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

const nodeAddrPos2 = "0xcf9db3f11b55c22926ae0a103eff7bae81bd5241";

const url = hre.config.networks.local.url;

const { Validators_ABI, VotePool_ABI } = require("./abi");

let admin, manager, receipt, tx;
let  vp_pos2;
const EPOCH = 5;
let bnadd=0;

async function step1(validator) {
    try {
        tx = await validator.addValidator(nodeAddrPos2, managerAddr, 100, 0, {
            gasLimit: 5000000,
            gasPrice: 0x12a05f200,
        });
        receipt = await tx.wait();
        expect(receipt.status).equal(1);
        // topics = receipt.events[0].topics;
        vp_pos2 = BigNumber.from(receipt.events[0].data).toHexString();
        console.log({ vp_pos2 });
        console.log("✅ SUCCESS step1: addValidator", nodeAddrPos2, managerAddr, 100, 0);
    } catch (e) {
        console.error("❌ FAILED step1: addValidator", nodeAddrPos2, managerAddr, 100, 0);
        throw e;
    }

    const vote1 = await ethers.getContractAt(VotePool_ABI, vp_pos2, manager);
    try {
        tx = await vote1.addMargin({
            gasLimit: 5000000,
            gasPrice: 0x12a05f200,
            value: ethers.utils.parseEther("5000.0"),
        });
        receipt = await tx.wait();
        expect(receipt.status).equal(1);

        bnadd = Number(receipt.blockNumber);
        console.log("✅ SUCCESS step1: addMargin 5000.0 eth");
    } catch (e) {
        console.error("❌ FAILED step1: addMargin 5000.0 eth");
        throw e;
    }

    try {
        const st1 = await vote1.state();
            console.log("✅ SUCCESS step1: after staking vote state",st1);
    } catch (e) {
        console.error("❌ FAILED step1: after staking vote state",st1);
        throw e;
    }

    let maxPunishBlk = (parseInt(bnadd / EPOCH) + 2) * EPOCH + 4 * 5 // 4*5 represents `removeThreshold * activeValidatorsLen`
    let currBlock = await ethers.provider.getBlockNumber()
    console.log({ maxPunishBlk, currBlock })
    await sleep((maxPunishBlk - currBlock) * 2000);
    let punishBlk = 0;
    try {
        punishBlk = await vote1.punishBlk();

        if (punishBlk > 0) {
            console.log("✅ SUCCESS step2: pos2 was punished at ", punishBlk.toString());
        } else {
            console.error("❌ FAILED step2: vote2.punishBlk");
            throw new Error("FAILED step2: vote2.punishBlk");
        }
    } catch (e) {
        console.error("❌ FAILED step2: vote2.punishBlk");
        throw e;
    }
    
    try {
        const st2 = await vote1.state();
            console.log("✅ SUCCESS step2: punished vote state",st2);
    } catch (e) {
        console.error("❌ FAILED step2: punished vote state",st2);
        throw e;
    }

    await sleep(20 * 2000);
}


async function step2() {

    const vote2 = await ethers.getContractAt(VotePool_ABI, vp_pos2, manager);

    try {
        tx = await vote2.addMargin({
            gasLimit: 5000000,
            gasPrice: 0x12a05f200,
            value: ethers.utils.parseEther("100.0"),
        });
        receipt = await tx.wait();
        expect(receipt.status).equal(1);
        console.log("✅ SUCCESS step2: addMargin 100.0 eth");
    } catch (e) {
        console.error("❌ FAILED step2: addMargin 100.0 eth");
        throw e;
    }

    try {
        const st3 = await vote2.state();
            console.log("✅ SUCCESS step2: punished reactive vote state",st3);
    } catch (e) {
        console.error("❌ FAILED step2: punished reactive vote state",st3);
        throw e;
    }
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
    await step2();
}

function sleep(ms) {
    console.log({ sleep: `${ms / 1000} s` });
    return new Promise(resolve => setTimeout(resolve, ms));
}

// We recommend this pattern to be able o use async/await everywhere
// and properly handle errors.
main()
    .then(() => process.exit(0))
    .catch(error => {
        console.error(error);
        process.exit(1);
    });

