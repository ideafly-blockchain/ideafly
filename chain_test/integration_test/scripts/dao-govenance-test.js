// We require the Hardhat Runtime Environment explicitly here. This is optional
// but useful for running the script in a standalone fashion through `node <script>`.
//
// When running the script with `npx hardhat run <script>` you'll find the Hardhat
// Runtime Environment's members available in the global scope.
const hre = require("hardhat");
const { expect } = require("chai");
const { BigNumber } = require("ethers");
const axios = require("axios");
const ethers = hre.ethers;
const govConAddr = "0x000000000000000000000000000000000000D003";
const blackConAddr = "0x000000000000000000000000000000000000D004";
const adminAddr = "0x3a696FeAe901DAe50967F28D7A2225577052F394";
const url = hre.config.networks[hre.config.defaultNetwork].url;
const { AddressList_ABI, Governance_ABI } = require("./abi");

async function main() {
    console.log("link to node: ", url);
    // Hardhat always runs the compile task when running scripts with its command
    // line interface.
    //
    // If this script is run directly using `node` you may want to call compile
    // manually to make sure everything is compiled
    // await hre.run('compile');
    const admin = await ethers.getSigner(adminAddr);

    const blackList = await ethers.getContractAt(AddressList_ABI, blackConAddr, admin);
    expect(await blackList.admin()).to.equal(adminAddr);

    const signers = await ethers.getSigners();
    let funder = signers[1];
    let blackFrom = signers[2];
    let blackBoth = signers[3];

    // let fbal = await funder.getBalance();
    // let thousand = ethers.parseEther("1000")
    // if (fbal.lte(thousand)) {
    //     await transferNDT(signers[0],funder.address,thousand)
    // }

    //initial fund
    const tenToken = ethers.utils.parseEther("10.0");
    let initTasks = [];
    let funderNonce = await funder.getTransactionCount("pending");
    for (var i = 2; i < signers.length; i++) {
        let balance = await signers[i].getBalance();
        console.log(signers[i].address, " balance: ", balance.toString());
        if (balance.lte(tenToken)) {
            initTasks.push(transferNDT(funder, signers[i].address, tenToken, funderNonce));
            funderNonce += 1;
        }
    }
    Promise.all(initTasks).then(resuts => resuts.forEach((ok, num) => expect(ok).to.be.true));
    // clear old context
    let resut = await blackList.isBlackAddress(blackFrom.address);
    console.log("Old: is in blacklist", blackFrom.address, resut);
    if (resut[0] === true) {
        expect(await removeBlackAddr(blackList, blackFrom.address, resut[1])).to.be.true;
    }

    const totalSupply = BigNumber.from("0x21e19e0c9bab2400000");
    const halfSupply = totalSupply.div(2);

    // deploy a token
    let depTasks = [];
    const dtFatory = await ethers.getContractFactory("DToken", funder);
    dt = await dtFatory.deploy("dtoken", "dt", totalSupply, funder.address);
    depTasks.push(dt.deployed());
    const mgrFatory = await ethers.getContractFactory("TokenManager", funder);
    const mgr = await mgrFatory.deploy({ gasPrice: 0x12a05f200 });
    depTasks.push(mgr.deployed());
    await Promise.all(depTasks);
    console.log("dt address ", dt.address);
    console.log("mgr address", mgr.address);
    tx = await dt.transfer(blackFrom.address, halfSupply);
    receipt = await tx.wait();
    expect(receipt.status).equal(1);

    // approve before test
    dt = dt.connect(blackFrom);
    tx = await dt.approve(mgr.address, halfSupply);
    await tx.wait();
    //show allowance
    allows = await dt.allowance(blackFrom.address, mgr.address);
    console.log("allowance", allows.toString());

    //===`from` blacklist===
    // 1. add an EOA address into `from` blacklist
    console.log("=*=*=*=*= TEST DirectionFrom =*=*=*=*=");

    //add black DirectionFrom
    console.log("try to add black(From)", blackFrom.address);
    let ok = await addBlackAddr(blackList, blackFrom.address, 0);
    expect(ok).to.be.true;
    console.log("✅ add black address ok");
    // should not be repeatedly added
    ok = await addBlackAddr(blackList, blackFrom.address, 0);
    expect(ok).false;
    console.log("✅ repeatedly add black address failed");

    // test1 : should not send any transaction
    expect(await transferNDT(blackFrom, funder.address, BigNumber.from(100))).to.be.false;
    console.log("✅ address in blacklist From can not send transaction");
    // also should not transfer assets by previously approved spender
    tx = await mgr.takeERC20(dt.address, blackFrom.address, 100, { gasLimit: 800000, gasPrice: 0x12a05f200 });
    try {
        await tx.wait();
        console.error("❌ Freeze ERC20 assets FAILED");
    } catch (e) {
        expect(e.code).to.be.equal("CALL_EXCEPTION");
        expect(e.receipt.status).to.be.equal(0);
        console.log("✅ Can't transfer erc20 even by spender,txHash", e.transactionHash);

        //debug trace API test
        // e.receipt.blockHash
        await mustTraceBlock(e.receipt.blockHash);
    }

    dt = dt.connect(funder);
    bal = await dt.balanceOf(mgr.address);
    console.log("mgr's dt balance", bal);
    expect(bal.eq(BigNumber.from(0))).to.be.true;

    ok = await transferNDT(funder, blackFrom.address, tenToken);
    expect(ok).be.true;
    console.log("✅ can transfer to a DirectionFrom black address");

    // 2. add a contract address into the `from` blacklist,
    // then this contract can not call others contract or transfer native token.
    ok = await addBlackAddr(blackList, mgr.address, 0);
    expect(ok).true;

    try {
        console.log("blackFrom contract call another contract, should failed");
        tx = await mgr.takeERC20(dt.address, blackFrom.address, 100, { gasLimit: 800000 });
        await tx.wait();
    } catch (e) {
        console.log("✅ Got expected error:", e.name, e.message);
    }
    try {
        console.log("blackFrom contract staticcall another contract, should failed");
        tx = await mgr.staticCallOtherContract(dt.address, { gasPrice: 0x12a05f200, gasLimit: 800000 });
        await tx.wait();
    } catch (e) {
        console.log("✅ Got expected error:", e.name, e.message);
    }
    // recover mgr
    ok = await removeBlackAddr(blackList, mgr.address, 0);
    expect(ok).true;
    console.log("=*=*=*=*= End of TEST DirectionFrom =*=*=*=*=\n");

    console.log("=*=*=*=*= TEST DirectionTo =*=*=*=*=");
    // add a contract address into the `to` blacklist
    ok = await addBlackAddr(blackList, dt.address, 1);
    expect(ok).true;
    // should not call this contract
    try {
        dt = dt.connect(funder);
        tx = await dt.transfer(signers[0].address, BigNumber.from(1), { gasPrice: 0x12a05f200, gasLimit: 800000 });
        receipt = await tx.wait();
        expect(receipt.status).equal(0);
    } catch (e) {
        console.log("✅  expected: transact with a DirectionTo black address FAILED.\n", e.name, e.message, e.code);
    }
    // also should not be called from another contract
    try {
        console.log("normal contract call DirectionTo black contract, should failed");
        tx = await mgr.takeERC20(dt.address, blackFrom.address, 100, { gasLimit: 800000 });
        await tx.wait();
    } catch (e) {
        console.log("✅ Got expected error:", e.name, e.message);
    }
    try {
        console.log("normal contract staticcall DirectionTo black contract, should failed");
        tx = await mgr.staticCallOtherContract(dt.address, { gasPrice: 0x12a05f200, gasLimit: 800000 });
        await tx.wait();
    } catch (e) {
        console.log("✅ Got expected error:", e.name, e.message);
    }
    console.log("=*=*=*=*=End of TEST DirectionTo =*=*=*=*=");

    console.log("=*=*=*=*= TEST DirectionBoth =*=*=*=*=");
    expect(await addBlackAddr(blackList, blackBoth.address, 0)).true;
    expect(await addBlackAddr(blackList, blackBoth.address, 1)).true;
    resut = await blackList.isBlackAddress(blackBoth.address);
    expect(resut[0]).true;
    expect(resut[1]).equal(2);
    expect(await transferNDT(blackBoth, funder.address, 1)).false;
    expect(await transferNDT(funder, blackBoth.address, 1)).false;

    expect(await removeBlackAddr(blackList, blackBoth.address, 2)).true;
    expect(await addBlackAddr(blackList, blackBoth.address, 2)).true;
    expect(await transferNDT(blackBoth, funder.address, 1)).false;
    expect(await transferNDT(funder, blackBoth.address, 1)).false;
    console.log("✅  test DirectionBoth ok");
    console.log("=*=*=*=*= End of TEST DirectionBoth =*=*=*=*=");

    console.log();
    console.log("=*=*=*=*= TEST managing rules =*=*=*=*=");
    let erc20transRule = {
        sig: "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
        idx: BigNumber.from(1),
        ct: 1,
    };
    await expect(blackList.removeRule(erc20transRule.sig, erc20transRule.idx))
        .to.emit(blackList, "RuleRemoved")
        .withArgs(erc20transRule.sig, erc20transRule.idx, erc20transRule.ct);
    console.log("✅ removes rule OK");
    expect(await removeBlackAddr(blackList, dt.address, 1)).true;
    tx = await mgr.takeERC20(dt.address, blackFrom.address, 100, { gasLimit: 800000, gasPrice: 0x12a05f200 });
    await tx.wait();
    dt = dt.connect(funder);
    bal = await dt.balanceOf(mgr.address);
    console.log("mgr's dt balance", bal.toString());
    expect(bal.eq(BigNumber.from(100))).to.be.true;
    console.log("✅ transfer erc20 by spender OK after remove rule");

    await expect(
        blackList.addOrUpdateRule(erc20transRule.sig, erc20transRule.idx, erc20transRule.ct, { gasPrice: 0x12a05f200 })
    )
        .to.emit(blackList, "RuleAdded")
        .withArgs(erc20transRule.sig, erc20transRule.idx, erc20transRule.ct);
    console.log("✅ add rule OK");

    tx = await mgr.takeERC20(dt.address, blackFrom.address, 100, { gasLimit: 800000, gasPrice: 0x12a05f200 });
    try {
        await tx.wait();
        console.error("❌ Freeze ERC20 assets FAILED");
    } catch (e) {
        expect(e.code).to.be.equal("CALL_EXCEPTION");
        expect(e.receipt.status).to.be.equal(0);
        console.log("✅ Freeze erc20 again after add rule");
    }

    console.log("=*=*=*=*= End of TEST managing rules =*=*=*=*=");

    console.log("=*=*=*=*= TEST remove black addresses to recover =*=*=*=*=");
    expect(await removeBlackAddr(blackList, blackFrom.address, 0)).true;
    expect(await removeBlackAddr(blackList, blackBoth.address, 2)).true;

    expect(await transferNDT(blackFrom, funder.address, 1)).true;
    tx = await dt.transfer(blackFrom.address, 1);
    receipt = await tx.wait();
    expect(receipt.status).equal(1);
    expect(await transferNDT(blackBoth, funder.address, 1)).true;
    expect(await transferNDT(funder, blackBoth.address, 1)).true;
    console.log("✅ remove black list addresses and then they become normal");
    console.log("=*=*=*=*=End of TEST remove black addresses to recover =*=*=*=*=");

    console.log();
    console.log("=*=*=*=*= TEST governance and debug-trace =*=*=*=*=");
    const govContract = await ethers.getContractAt(Governance_ABI, govConAddr, admin);
    expect(await govContract.admin()).to.equal(adminAddr);

    let puppet = signers[4];
    expect(await transferNDT(admin, puppet.address, 123)).true;
    let puppetBalance = await puppet.getBalance();

    // dt mint: 0x1249c58b  =>  mint()
    tx = await govContract.commitProposal(0, puppet.address, dt.address, 23, "0x1249c58b");
    receipt = await tx.wait();
    expect(receipt.status).equal(1);
    console.log("✅  commitProposal with ht value and contract input OK");
    console.log("proposal txHash", tx.hash, "blockHash", receipt.blockHash, "blockNumber", receipt.blockNumber);
    await mustTraceBlock(receipt.blockHash);
    let reft = await puppet.getBalance();
    expect(reft.eq(puppetBalance.sub(23))).true;
    console.log("✅  the proposal execution result is OK");

    expect(await addBlackAddr(blackList, puppet.address, 2)).true;
    tx = await govContract.commitProposal(0, puppet.address, dt.address, 100, "0x1249c58b");
    receipt = await tx.wait();
    expect(receipt.status).equal(1);
    console.log("✅  commitProposal involves black address with ht value and contract input OK");
    console.log("proposal txHash", tx.hash, "blockHash", receipt.blockHash, "blockNumber", receipt.blockNumber);
    await mustTraceBlock(receipt.blockHash);
    reft = await puppet.getBalance();
    expect(reft.eq(puppetBalance.sub(123))).true;
    console.log("✅  the proposal execution result is OK");

    tx = await govContract.commitProposal(0, dt.address, puppet.address, 123, "0x");
    receipt = await tx.wait();
    reft = await puppet.getBalance();
    expect(reft.eq(puppetBalance)).true;
    console.log("✅  take NDT from a contract by gov-proposal OK");
    // await showReceipts(ethers.provider, receipt.blockHash);

    expect(await removeBlackAddr(blackList, puppet.address, 2)).true;

    tx = await govContract.commitProposal(1, "0x0000000000000000000000000000000000000000", dt.address, 0, "0x");
    receipt = await tx.wait();
    console.log("commitProposal receipt logs ", receipt.logs);
    // await showReceipts(ethers.provider, receipt.blockHash);

    code = await ethers.provider.send("eth_getCode", [dt.address, "latest"]);
    expect(code).to.be.eq("0x");
    console.log("=*=*=*=*=End of TEST governance and debug-trace =*=*=*=*=");
}

// We recommend this pattern to be able o use async/await everywhere
// and properly handle errors.
main()
    .then(() => process.exit(0))
    .catch(error => {
        console.error(error);
        process.exit(1);
    });

async function transferNDT(signer, to, val, nonce) {
    try {
        txa = await signer.populateTransaction({ to: to, value: val, nonce: nonce });
        tx = await signer.sendTransaction(txa);
        // console.log(tx)
        let receipt = await tx.wait(1);
        let bal = await ethers.provider.getBalance(to);
        if (receipt.status === 1) {
            console.log(
                "transfer NDT ok ",
                signer.address,
                " to ",
                to,
                val.toString(),
                "Wei. ",
                "FinalBalance: ",
                bal.toString()
            );
            return true;
        } else {
            console.log(
                "transfer NDT failed: ",
                signer.address,
                " to ",
                to,
                val.toString(),
                "Wei. ",
                "FinalBalance: ",
                bal.toString()
            );
            return false;
        }
    } catch (e) {
        console.log(
            "transfer NDT error: ",
            signer.address,
            " to ",
            to,
            val.toString(),
            "Wei. ",
            " ERROR: ",
            e.name,
            e.message,
            e.transaction
        );
    }
    return false;
}

async function addBlackAddr(contract, a, d) {
    try {
        let tx = await contract.addBlacklist(a, d, { gasPrice: 0x12a05f200, gasLimit: 1000000 });
        let receipt = await tx.wait(1);
        if (receipt.status === 1) {
            console.log("addBlacklist success: ", a, d);
            return true;
        } else {
            console.log("addBlacklist failed: ", a, d);
            return false;
        }
    } catch (e) {
        console.log("addBlacklist error: ", a, d, " ERROR: ", e.name, e.message);
    }
    return false;
}

async function removeBlackAddr(contract, a, d) {
    try {
        let tx = await contract.removeBlacklist(a, d, { gasPrice: 0x12a05f200, gasLimit: 1000000 });
        let receipt = await tx.wait(1);
        if (receipt.status === 1) {
            console.log("removeBlacklist success: ", a, d);
            return true;
        } else {
            console.log("removeBlacklist failed: ", a, d);
            return false;
        }
    } catch (e) {
        console.log("removeBlacklist error: ", a, d, " ERROR: ", e.name, e.message);
    }
    return false;
}

function rpc(method, params) {
    return axios.post(url, {
        id: 1,
        jsonrpc: "2.0",
        method,
        params,
    });
}

async function mustTraceBlock(hash) {
    let res = await rpc("debug_traceBlockByHash", [hash, null]);
    if (res.data.result) {
        console.log("✅ debug_traceBlockByHash trace block OK. blockHash:", hash);
        console.log(res.data.result);
        // console.log(res.data.result[0].result.structLogs[0])
    } else {
        console.error("❌ debug_traceBlockByHash trace block FAILED. blockHash:", hash);
        console.log(res.data.error);
    }
}

async function showReceipts(provider, blockHash) {
    block = await provider.send("eth_getBlockByHash", [blockHash, false]);
    trxLen = block.transactions.length;
    for (let i = 0; i < trxLen; i++) {
        sreceipt = await provider.send("eth_getTransactionReceipt", [block.transactions[i]]);
        console.log("receipt:  ", sreceipt);
    }
}
