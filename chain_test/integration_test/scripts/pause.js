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

    const vote1 = await ethers.getContractAt(VotePool_ABI, vp_poa4, manager);
    try {
        const st1 = await vote1.state();
            console.log("✅ SUCCESS step1: before pause vote state",st1);
    } catch (e) {
        console.error("❌ FAILED step1: before pause vote state",st1);
        throw e;
    }

    try {
        tx = await validator.updateValidatorState(nodeAddrPoa4, true);
        receipt = await tx.wait();
        expect(receipt.status).equal(1);
        // console.log(JSON.stringify(receipt));

        console.log("✅ SUCCESS step1: updateValidatorState", nodeAddrPoa4, managerAddr);
    } catch (e) {
        console.error("❌ FAILED step1: updateValidatorState", nodeAddrPoa4, managerAddr);
        throw e;
    }

    try {
        const st2 = await vote1.state();
            console.log("✅ SUCCESS step1: after pause vote state",st2);
    } catch (e) {
        console.error("❌ FAILED step1: after pause vote state",st2);
        throw e;
    }
}

async function step2(validator) {

    try {
        tx = await validator.updateValidatorState(nodeAddrPoa4, false);
        receipt = await tx.wait();
        expect(receipt.status).equal(1);
        // console.log(JSON.stringify(receipt));

        console.log("✅ SUCCESS step2: updateValidatorState", nodeAddrPoa4, managerAddr);
    } catch (e) {
        console.error("❌ FAILED step2: updateValidatorState", nodeAddrPoa4, managerAddr);
        throw e;
    }

    const vote2 = await ethers.getContractAt(VotePool_ABI, vp_poa4, manager);
    try {
        const st3 = await vote2.state();
            console.log("✅ SUCCESS step2: after reactive vote state",st3);
    } catch (e) {
        console.error("❌ FAILED step2: after reactive vote state",st3);
        throw e;
    }
}


async function step3() {
    const vote = await ethers.getContractAt(VotePool_ABI, vp_poa4, manager);

    try {
        const tx = await vote.exit();
        receipt = await tx.wait();
        expect(receipt.status).equal(1);
        e=receipt.events[0].event

        console.log("✅ SUCCESS step3: vote.exit event",e);
    } catch (e) {
        console.error("❌ FAILED step3: vote.exit event",e);
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
    await step2(validator);
    await step3();
}

// We recommend this pattern to be able o use async/await everywhere
// and properly handle errors.
main()
    .then(() => process.exit(0))
    .catch(error => {
        console.error(error);
        process.exit(1);
    });

