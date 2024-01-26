# NPoS Automated Testing

## Test Environment Requirements

Run the blockchain on `../docker_chain` to facilitate repeatable testing (resetting the chain is necessary for repeated tests).

### Dependencies:

- Docker
- Node.js

### Blockchain Running Configuration

Initial validators (3):

```
0x879b32dbDac37bb91d91B9eC2f5c1b876ae9890f
0xa162eD4644334Fc54a52B0c0dE4022b9bd8706FD
0x7E89120FbFacbbf17a91E533D3b93A836082a2FF
```

Two additional validators added via test scripts:

```
0x1D342EEFD8A93b513a17e8CD9Fa9619C35aA21fA
0x01A9015151C05b5406d57e49fdBD3440b6aB1050
```

For candidate validator testing:

```
0xcF9DB3f11B55C22926AE0a103eff7BAe81Bd5241
```

All validators use the following address as a manager (for ease of test only):

```
0x9E737Ee8bDc132c349dE7801Efbc9e12f4FE99e9
```

Chain constant setting:
- period: 2
- epoch: 5

System contract constants:

`Params.sol`:
```js
    // max active validators
    uint16 public constant MaxValidators = 5;  
    // margin threshold for a PoS type of validator 
    uint public constant PosMinMargin = 5000 ether;
    // margin threshold for a PoA type of validator
    uint public constant PoaMinMargin = 1 ether;

    // The punish amount from validators margin when the validator is jailed
    uint public constant PunishAmount = 100 ether;

    // JailPeriod: a block count, how many blocks a validator should be jailed when it got punished
    uint public constant JailPeriod = 12;
    // when a validator claim to exit, how many blocks its margin should be lock  
    // (after that locking period, the valicator can withdraw its margin)
    uint public constant MarginLockPeriod = 12;
    // when a voter claim to withdraw its stake(vote), how many blocks its token should be lock
    uint public constant WithdrawLockPeriod = 12;
    // when a validator change its commison rate, it should take `PercentChangeLockPeriod` blocks to take effect
    uint public constant PercentChangeLockPeriod = 12;
```

`Punish.sol`:
```js
    // When the missedBlocksCounter reaches `punishThreshold`, the currently unclaimed rewards of the validator will be forfeited.
    uint256 public constant punishThreshold = 2;
    // When the missedBlocksCounter reaches `removeThreshold`, the validator will be jailed 
    uint256 public constant removeThreshold = 4;
    // How many blocks were allowed to missing for a validator in one epoch  
    uint256 public constant decreaseCountPerEpoch = 0;
```

contract address:
validators: `0x000000000000000000000000000000000000d001`


## Test-cases and testing flow

### before all test cases

Update params, function : `function updateParams(uint8 _posCount,uint8 _posBackup,uint8 _poaCount,uint8 _poaBackup) external onlyAdmin`

`updateParams(1,1,4,1)`

### T1: Adding a PoA-Type Validator and Participating in Consensus

1. Add `poa4` (`0x1D342EEFD8A93b513a17e8CD9Fa9619C35aA21fA`) as a PoA-type validator; function: `function addValidator(address _validator, address _manager, uint256 _percent, uint8 _type) returns (address)`, to obtain the votePool contract address `vp-poa4` corresponding to `poa4`.
2. Use `vp-poa4` to add margin for poa4 in the amount of `PoaMinMargin`, using the function: `function addMargin() payable`; record the block number `bnadd` of the transaction. Since there are enough PoA seats, poa4 can enter the set of authoritative validators without needing votes at this point.
3. The epoch number at which poa4 starts participating in consensus (epoch index starting from 0) is: `bnadd/EPOCH + 2`. On block `(bnadd/EPOCH + 2)*EPOCH`, query the validators contract's `function getActiveValidators() view returns (address[])`. The return result should include the address of poa4. **If it does not, then the test case fails**.
4. In the consecutive 4 blocks starting from block `(bnadd/EPOCH + 2)*EPOCH + 1`, one block will be produced by poa4, meaning one block's `header.coinbase` will be equal to poa4.

### T2: Adding a PoS-Type Validator and Enabling Its Participation in Consensus

1. Add pos1 (`0x01A9015151C05b5406d57e49fdBD3440b6aB1050`) as a PoS-type validator; and obtain the votePool contract address `vp-pos1` corresponding to pos1.
2. Use `vp-pos1` to add margin for pos1 in the amount of `PosMinMargin`, using the function: `function addMargin() payable`; record the block number `bnadd` of the transaction. Since there are enough PoS seats, pos1 can enter the set of authoritative validators without needing votes at this point.
3. The epoch number at which pos1 starts participating in consensus (starting from 0) is: `bnadd/EPOCH + 2`. On block `(bnadd/EPOCH + 2)*EPOCH`, query the validators contract's `function getActiveValidators() view returns (address[])`. The return result should include the address of pos1. **If it does not, then the test case fails**.
4. In the consecutive 5 blocks starting from block `(bnadd/EPOCH + 2)*EPOCH + 1`, one block will be produced by pos1, meaning one block's `header.coinbase` will be equal to pos1.

### T3: Adding a Candidate Validator (Using PoS Type as an Example)

1. First, vote for pos1 to clearly confirm the candidate validator. Vote 10 ether for pos1 using `function deposit() payable` (any address can vote for pos1, in the test case, you can directly use the manager's address to vote).
   1. You can query a user's voting information for that validator through the votepool contract's function `function voters(address) view returns (uint256 amount, uint256 rewardDebt, uint256 withdrawPendingAmount, uint256 withdrawExitBlock)`.
2. Add pos2 (`0xcF9DB3f11B55C22926AE0a103eff7BAe81Bd5241`) as a PoS type validator; obtain the votePool contract address vp-pos2 corresponding to pos2.
3. Use vp-pos2 to add collateral for pos2 in the amount of `PosMinMargin`, and record the block number `bnadd` of the transaction.
4. Check the status of pos2, which should be `1`, meaning `Ready` state; `function state() view returns (uint8)`.
5. On block `(bnadd/EPOCH + 2)*EPOCH`, query the candidate validator set, which should only contain the address of pos2, using `function getBackupValidators() view returns (address[])`. If not, the test case fails.


### T4: Candidate Validator Becoming an Authoritative Validator and Penalty for Authoritative Validator's Inactivity

1. Provide more votes to pos2 than pos1, by voting 20 ether for pos2, and record the block number `bnadd` of the transaction.
2. On block `(bnadd/EPOCH + 2)*EPOCH`, query `getActiveValidators`. The return result should include 5 addresses, containing the address of pos2 but not pos1.
3. Query `getBackupValidators`. The return result should only include the address of pos1.
4. Since pos2 does not run the corresponding node to participate in consensus, after `removeThreshold * activeValidatorsLen` i.e., `4*5` blocks, i.e., by no later than block `(bnadd/EPOCH + 2)*EPOCH + 4*5`, pos2 will be penalized. At this time, querying the status of pos2 should return `3`, meaning `Jail` status. Querying the margin of pos2, it should be `4900` ether at this time, using the function `function margin() view returns (uint256)`.
5. The block number when pos2 is penalized can be queried through the `function punishBlk() view returns (uint256)` function of pos2's corresponding votepool contract, denoted as `punishBlk`. After another 2 epochs, i.e., after block `(punishBlk/EPOCH + 2)*EPOCH`, pos1 will transition from a candidate validator back to an authoritative validator, meaning `getActiveValidators` will include the address of pos1.

