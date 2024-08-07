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

const initNodeAddr1 = "0x879b32dbdac37bb91d91b9ec2f5c1b876ae9890f";
const initNodeAddr2 = "0xa162ed4644334fc54a52b0c0de4022b9bd8706fd";
const initNodeAddr3 = "0x7e89120fbfacbbf17a91e533d3b93a836082a2ff";
const nodeAddrPos1 = "0x01a9015151c05b5406d57e49fdbd3440b6ab1050";
const nodeAddrPos2 = "0xcf9db3f11b55c22926ae0a103eff7bae81bd5241";
const nodeAddrPoa4 = "0x1d342eefd8a93b513a17e8cd9fa9619c35aa21fa";

const url = hre.config.networks.local.url;

const { Validators_ABI, VotePool_ABI } = require("./abi");

let admin, manager, receipt, tx;
let vp_pos1, vp_pos2, vp_poa4;
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

    let bnadd = 0;
    let next = 0;
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
        const addrs = await validator.getActiveValidators();
        const addrLowers = addrs.map(addr => addr.toLowerCase());
        // console.log({ addrLowers });
        if (addrLowers.includes(nodeAddrPoa4)) {
            console.log("✅ SUCCESS step1: getActiveValidators", addrLowers, nodeAddrPoa4);
        } else {
            console.error("❌ FAILED step1: getActiveValidators", addrLowers, nodeAddrPoa4);
            throw new Error("FAILED step1: getActiveValidators");
        }
    } catch (e) {
        console.error("❌ FAILED step1: getActiveValidators", nodeAddrPoa4);
        throw e;
    }
}
async function step2(validator) {
    try {
        tx = await validator.addValidator(nodeAddrPos1, managerAddr, 100, 0, {
            gasLimit: 5000000,
            gasPrice: 0x12a05f200,
        });
        receipt = await tx.wait();
        expect(receipt.status).equal(1);
        // topics = receipt.events[0].topics;
        vp_pos1 = BigNumber.from(receipt.events[0].data).toHexString();
        console.log({ vp_pos1 });
        console.log("✅ SUCCESS step2: addValidator", nodeAddrPos1, managerAddr, 100, 0);
    } catch (e) {
        console.error("❌ FAILED step2: addValidator", nodeAddrPos1, managerAddr, 100, 0);
        throw e;
    }

    const vote = await ethers.getContractAt(VotePool_ABI, vp_pos1, manager);

    let bnadd = 0;
    let next = 0;
    try {
        tx = await vote.addMargin({
            gasLimit: 5000000,
            gasPrice: 0x12a05f200,
            value: ethers.utils.parseEther("5000.0"),
        });
        receipt = await tx.wait();
        expect(receipt.status).equal(1);
        // console.log(JSON.stringify(receipt));
        bnadd = Number(receipt.blockNumber);
        next = (parseInt(bnadd / EPOCH) + 2) * EPOCH + 1
        console.log({ bnadd, next });
        expect(bnadd).gt(0);
        console.log("✅ SUCCESS step2: addMargin 5000.0 eth");
    } catch (e) {
        console.error("❌ FAILED step2: addMargin 5000.0 eth");
        throw e;
    }
    const ds = next + 1 - bnadd;
    await sleep(ds * 2000);
    let generated = false;
    for (let i = next; i <= next + 5; i += 1) {
        await sleep(2000);
        const block = await ethers.provider.getBlock(i);
        console.log({ blockNumber: i, miner: block.miner });
        if (block.miner.toLowerCase() == nodeAddrPos1.toLowerCase()) {
            generated = true;
            break;
        }
    }
    if (generated) {
        console.log("✅ SUCCESS step2: generate");
    } else {
        console.error("❌ FAILED step2: generate");
    }

    try {
        const addrs = await validator.getActiveValidators();
        const addrLowers = addrs.map(addr => addr.toLowerCase());
        // console.log({ addrLowers });
        if (addrLowers.includes(nodeAddrPos1)) {
            console.log("✅ SUCCESS step2: getActiveValidators", addrLowers, nodeAddrPos1);
        } else {
            console.error("❌ FAILED step2: getActiveValidators", addrLowers, nodeAddrPos1);
            throw new Error("FAILED step2: getActiveValidators");
        }
    } catch (e) {
        console.error("❌ FAILED step2: getActiveValidators", nodeAddrPos1);
        throw e;
    }
}

async function step3(validator) {
    const vote1 = await ethers.getContractAt(VotePool_ABI, vp_pos1, manager);
    try {
        tx = await vote1.deposit({
            gasLimit: 5000000,
            gasPrice: 0x12a05f200,
            value: ethers.utils.parseEther("10.0"),
        });
        receipt = await tx.wait();
        expect(receipt.status).equal(1);
        console.log("✅ SUCCESS step3: vote1.deposit 10.0 eth");
    } catch (e) {
        console.error("❌ FAILED step3: vote1.deposit 10.0 eth");
        throw e;
    }

    try {
        const voteInfos = await vote1.voters(managerAddr);
        console.log("voteInfos.amount", ethers.utils.formatEther(voteInfos.amount));
        if (voteInfos.amount.eq(ethers.utils.parseEther("10.0"))) {
            console.log("✅ SUCCESS step3: vote1.voters 10.0 eth");
        } else {
            console.error("❌ FAILED step3: vote1.voters 10.0 eth");
            throw new Error("FAILED step3: vote1.voters 10.0 eth");
        }
    } catch (e) {
        console.error("❌ FAILED step3: vote1.voters 10.0 eth");
        throw e;
    }

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
        console.log("✅ SUCCESS step3: addValidator", nodeAddrPos2, managerAddr, 100, 0);
    } catch (e) {
        console.error("❌ FAILED step3: addValidator", nodeAddrPos2, managerAddr, 100, 0);
        throw e;
    }

    const vote2 = await ethers.getContractAt(VotePool_ABI, vp_pos2, manager);

    let bnadd = 0;
    try {
        tx = await vote2.addMargin({
            gasLimit: 5000000,
            gasPrice: 0x12a05f200,
            value: ethers.utils.parseEther("5000.0"),
        });
        receipt = await tx.wait();
        expect(receipt.status).equal(1);
        // console.log(JSON.stringify(receipt));
        bnadd = Number(receipt.blockNumber);
        console.log({ bnadd, next: (parseInt(bnadd / EPOCH) + 2) * EPOCH + 1 });
        expect(bnadd).gt(0);
        console.log("✅ SUCCESS step3: vote2.addMargin 5000.0 eth");
    } catch (e) {
        console.error("❌ FAILED step3: vote2.addMargin 5000.0 eth");
        throw e;
    }

    try {
        const st = await vote2.state();
        if (st == 1) {
            console.log("✅ SUCCESS step3: vote2 state");
        } else {
            console.error("❌ FAILED step3: vote2 state");
            throw new Error("FAILED step3: vote2 state");
        }
        expect(st).equal(1);
    } catch (e) {
        console.error("❌ FAILED step3: vote2 state");
        throw e;
    }

    await sleep(30000);

    try {
        const addrs = await validator.getBackupValidators();
        const addrLowers = addrs.map(addr => addr.toLowerCase());
        console.log({ addrLowers });
        if (addrLowers.length == 1 && addrLowers.includes(nodeAddrPos2)) {
            console.log("✅ SUCCESS step3: getBackupValidators", addrLowers, nodeAddrPos2);
        } else {
            console.error("❌ FAILED step3: getBackupValidators", addrLowers, nodeAddrPos2);
            throw new Error("FAILED step3: getBackupValidators");
        }
    } catch (e) {
        console.error("❌ FAILED step3: getBackupValidators", nodeAddrPos2);
        throw e;
    }
}

async function step4(validator) {
    const vote2 = await ethers.getContractAt(VotePool_ABI, vp_pos2, manager);

    // 1
    let bnadd = 0;
    let next = 0;
    try {
        tx = await vote2.deposit({
            gasLimit: 5000000,
            gasPrice: 0x12a05f200,
            value: ethers.utils.parseEther("20.0"),
        });
        receipt = await tx.wait();
        expect(receipt.status).equal(1);
        // console.log(JSON.stringify(receipt));
        bnadd = Number(receipt.blockNumber);
        next = (parseInt(bnadd / EPOCH) + 2) * EPOCH + 1
        console.log("pos2 deposit blockNum: ", { bnadd, next });
        expect(bnadd).gt(0);
        console.log("✅ SUCCESS step4: vote2.deposit 20.0 eth");
    } catch (e) {
        console.error("❌ FAILED step4: vote2.deposit 20.0 eth");
        throw e;
    }
    const ds = next + 1 - bnadd;
    await sleep(2000 * ds);

    // 2
    try {
        const addrs = await validator.getActiveValidators();
        const addrLowers = addrs.map(addr => addr.toLowerCase());
        console.log({ addrLowers });
        if (addrLowers.includes(nodeAddrPos2)) {
            console.log("✅ SUCCESS step4: getActiveValidators nodeAddrPos2", addrLowers, nodeAddrPos2);
        } else {
            console.error("❌ FAILED step4: getActiveValidators nodeAddrPos2", addrLowers, nodeAddrPos2);
            throw new Error('FAILED step4: getActiveValidators nodeAddrPos2');
        }
    } catch (e) {
        console.error("❌ FAILED step4: getActiveValidators nodeAddrPos2", nodeAddrPos2);
        throw e;
    }

    // 3
    try {
        const addrs = await validator.getBackupValidators();
        const addrLowers = addrs.map(addr => addr.toLowerCase());
        console.log({ addrLowers });
        if (addrLowers.length == 1 && addrLowers.includes(nodeAddrPos1)) {
            console.log("✅ SUCCESS step4: getBackupValidators", addrLowers, nodeAddrPos1);
        } else {
            console.error("❌ FAILED step4: getBackupValidators", addrLowers, nodeAddrPos1);
            throw new Error('FAILED step4: getBackupValidators');
        }
    } catch (e) {
        console.error("❌ FAILED step4: getBackupValidators", nodeAddrPos1);
        throw e;
    }

    // 4
    let maxPunishBlk = (parseInt(bnadd / EPOCH) + 2) * EPOCH + 4 * 5 // 4*5 represents `removeThreshold * activeValidatorsLen`
    let currBlock = await ethers.provider.getBlockNumber()
    console.log({ maxPunishBlk, currBlock })
    await sleep((maxPunishBlk - currBlock) * 2000);
    let punishBlk = 0;
    try {
        punishBlk = await vote2.punishBlk();

        if (punishBlk > 0) {
            console.log("✅ SUCCESS step4: pos2 was punished at ", punishBlk.toString());
        } else {
            console.error("❌ FAILED step4: vote2.punishBlk");
            throw new Error("FAILED step4: vote2.punishBlk");
        }
    } catch (e) {
        console.error("❌ FAILED step4: vote2.punishBlk");
        throw e;
    }

    // 5
    try {
        const amount = await vote2.margin();
        console.log(BigNumber.from(amount).toString());

        if (ethers.utils.parseEther("4900.0").eq(BigNumber.from(amount))) {
            console.log("✅ SUCCESS step4: pos2 was punished as expected");
        } else {
            console.error("❌ FAILED step4: vote2.margin");
            throw new Error("FAILED step4: vote2.margin");
        }
    } catch (e) {
        console.error("❌ FAILED step4: vote2.margin");
        throw e;
    }

    currBlock = await ethers.provider.getBlockNumber()
    next = (parseInt(punishBlk / EPOCH) + 2) * EPOCH + 1
    await sleep((next - currBlock) * 2000);

    try {
        const addrs = await validator.getActiveValidators();
        const addrLowers = addrs.map(addr => addr.toLowerCase());
        console.log({ addrLowers });
        if (addrLowers.includes(nodeAddrPos1)) {
            console.log("✅ SUCCESS step4: getActiveValidators recovered, pos1 becomes active validator again", addrLowers, nodeAddrPos1);
        } else {
            console.error("❌ FAILED step4: getActiveValidators recovered", addrLowers, nodeAddrPos1);
            throw new Error("FAILED step4: getActiveValidators recovered");
        }
    } catch (e) {
        console.error("❌ FAILED step4: getActiveValidators recovered", nodeAddrPos1);
        throw e;
    }
}

async function main() {
    console.log("link to node: ", url);
    admin = await ethers.getSigner(adminAddr);
    manager = await ethers.getSigner(managerAddr);
    console.log(`admin addr: ${adminAddr}`);
    console.log(`manager addr: ${managerAddr}`);
    let currBlock = await ethers.provider.getBlockNumber();
    let i = 0;
    while(i < 5 && currBlock < 2) {
        console.log(`Wait untill blockchain running, current block number: ${currBlock}`);
        await sleep(2000);
        i += 1;
        if (i >= 5) {
            console.log("Chain not running!");
            process.exit(1);
        }
        currBlock = await ethers.provider.getBlockNumber();
    }

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
    await step3(validator);
    await step4(validator);
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
