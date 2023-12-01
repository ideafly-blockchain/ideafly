package systemcontract

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

// ValidatorsInteractiveABI contains all methods to interactive with validator contracts.
const ValidatorsInteractiveABI = `[
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "validator",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "address",
          "name": "votePool",
          "type": "address"
        }
      ],
      "name": "AddValidator",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "admin",
          "type": "address"
        }
      ],
      "name": "ChangeAdmin",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "address",
          "name": "foundation",
          "type": "address"
        }
      ],
      "name": "UpdateFoundationAddress",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "uint8",
          "name": "posCount",
          "type": "uint8"
        },
        {
          "indexed": false,
          "internalType": "uint8",
          "name": "posBackup",
          "type": "uint8"
        },
        {
          "indexed": false,
          "internalType": "uint8",
          "name": "poaCount",
          "type": "uint8"
        },
        {
          "indexed": false,
          "internalType": "uint8",
          "name": "poaBackup",
          "type": "uint8"
        }
      ],
      "name": "UpdateParams",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "burnRate",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "foundationRate",
          "type": "uint256"
        }
      ],
      "name": "UpdateRates",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "address",
          "name": "receiver",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "WithdrawFoundationReward",
      "type": "event"
    },
    {
      "inputs": [],
      "name": "JailPeriod",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "MarginLockPeriod",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "MaxValidators",
      "outputs": [
        {
          "internalType": "uint16",
          "name": "",
          "type": "uint16"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "PercentChangeLockPeriod",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "PoaMinMargin",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "PosMinMargin",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "PunishAmount",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "WithdrawLockPeriod",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_validator",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "_manager",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "_percent",
          "type": "uint256"
        },
        {
          "internalType": "enum ValidatorType",
          "name": "_type",
          "type": "uint8"
        }
      ],
      "name": "addValidator",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "admin",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "name": "allValidators",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "enum ValidatorType",
          "name": "",
          "type": "uint8"
        }
      ],
      "name": "backupCount",
      "outputs": [
        {
          "internalType": "uint8",
          "name": "",
          "type": "uint8"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "burnRate",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_newAdmin",
          "type": "address"
        }
      ],
      "name": "changeAdmin",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "enum ValidatorType",
          "name": "",
          "type": "uint8"
        }
      ],
      "name": "count",
      "outputs": [
        {
          "internalType": "uint8",
          "name": "",
          "type": "uint8"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "distributeBlockReward",
      "outputs": [],
      "stateMutability": "payable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "foundation",
      "outputs": [
        {
          "internalType": "address payable",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "foundationRate",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "foundationReward",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getActiveValidators",
      "outputs": [
        {
          "internalType": "address[]",
          "name": "",
          "type": "address[]"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getAllValidatorsLength",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getBackupValidators",
      "outputs": [
        {
          "internalType": "address[]",
          "name": "",
          "type": "address[]"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getTopValidators",
      "outputs": [
        {
          "internalType": "address[]",
          "name": "",
          "type": "address[]"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "improveRanking",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address[]",
          "name": "_validators",
          "type": "address[]"
        },
        {
          "internalType": "address[]",
          "name": "_managers",
          "type": "address[]"
        },
        {
          "internalType": "address",
          "name": "_admin",
          "type": "address"
        }
      ],
      "name": "initialize",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "initialized",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "lowerRanking",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "contract IVotePool",
          "name": "",
          "type": "address"
        }
      ],
      "name": "pendingReward",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "punishContract",
      "outputs": [
        {
          "internalType": "contract IPunish",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "removeRanking",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address[]",
          "name": "newSet",
          "type": "address[]"
        },
        {
          "internalType": "uint256",
          "name": "epoch",
          "type": "uint256"
        }
      ],
      "name": "updateActiveValidatorSet",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address payable",
          "name": "_foundation",
          "type": "address"
        }
      ],
      "name": "updateFoundation",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint8",
          "name": "_posCount",
          "type": "uint8"
        },
        {
          "internalType": "uint8",
          "name": "_posBackup",
          "type": "uint8"
        },
        {
          "internalType": "uint8",
          "name": "_poaCount",
          "type": "uint8"
        },
        {
          "internalType": "uint8",
          "name": "_poaBackup",
          "type": "uint8"
        }
      ],
      "name": "updateParams",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "_burnRate",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "_foundationRate",
          "type": "uint256"
        }
      ],
      "name": "updateRates",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_validator",
          "type": "address"
        },
        {
          "internalType": "bool",
          "name": "pause",
          "type": "bool"
        }
      ],
      "name": "updateValidatorState",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "validatorsContract",
      "outputs": [
        {
          "internalType": "contract IValidators",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "name": "votePools",
      "outputs": [
        {
          "internalType": "contract IVotePool",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "withdrawFoundationReward",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "withdrawReward",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    }
  ]`
const PunishInteractiveABI = `
[
	{
		"inputs": [],
		"name": "initialize",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
		  {
			"internalType": "address",
			"name": "val",
			"type": "address"
		  }
		],
		"name": "punish",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
		  {
			"internalType": "uint256",
			"name": "epoch",
			"type": "uint256"
		  }
		],
		"name": "decreaseMissedBlocksCounter",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	  }
]
`

const SysGovInteractiveABI = `
[
    {
		"inputs": [
			{
				"internalType": "uint256",
				"name": "id",
				"type": "uint256"
			}
		],
		"name": "finishProposalById",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint32",
				"name": "index",
				"type": "uint32"
			}
		],
		"name": "getPassedProposalByIndex",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "id",
				"type": "uint256"
			},
			{
        		"internalType": "uint256",
        		"name": "action",
        		"type": "uint256"
        	},
			{
				"internalType": "address",
				"name": "from",
				"type": "address"
			},
			{
				"internalType": "address",
				"name": "to",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "value",
				"type": "uint256"
			},
			{
				"internalType": "bytes",
				"name": "data",
				"type": "bytes"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "getPassedProposalCount",
		"outputs": [
			{
				"internalType": "uint32",
				"name": "",
				"type": "uint32"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "_admin",
				"type": "address"
			}
		],
		"name": "initialize",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	}
]`

const AddrListInteractiveABI = `
[
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "newAdmin",
          "type": "address"
        }
      ],
      "name": "AdminChanged",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "newAdmin",
          "type": "address"
        }
      ],
      "name": "AdminChanging",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "addr",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "enum AddressList.Direction",
          "name": "d",
          "type": "uint8"
        }
      ],
      "name": "BlackAddrAdded",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "addr",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "enum AddressList.Direction",
          "name": "d",
          "type": "uint8"
        }
      ],
      "name": "BlackAddrRemoved",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "addr",
          "type": "address"
        }
      ],
      "name": "DeveloperAdded",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "addr",
          "type": "address"
        }
      ],
      "name": "DeveloperRemoved",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "bool",
          "name": "newState",
          "type": "bool"
        }
      ],
      "name": "EnableStateChanged",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "bytes32",
          "name": "eventSig",
          "type": "bytes32"
        },
        {
          "indexed": false,
          "internalType": "uint128",
          "name": "checkIdx",
          "type": "uint128"
        },
        {
          "indexed": false,
          "internalType": "enum AddressList.CheckType",
          "name": "t",
          "type": "uint8"
        }
      ],
      "name": "RuleAdded",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "bytes32",
          "name": "eventSig",
          "type": "bytes32"
        },
        {
          "indexed": false,
          "internalType": "uint128",
          "name": "checkIdx",
          "type": "uint128"
        },
        {
          "indexed": false,
          "internalType": "enum AddressList.CheckType",
          "name": "t",
          "type": "uint8"
        }
      ],
      "name": "RuleRemoved",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "bytes32",
          "name": "eventSig",
          "type": "bytes32"
        },
        {
          "indexed": false,
          "internalType": "uint128",
          "name": "checkIdx",
          "type": "uint128"
        },
        {
          "indexed": false,
          "internalType": "enum AddressList.CheckType",
          "name": "t",
          "type": "uint8"
        }
      ],
      "name": "RuleUpdated",
      "type": "event"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "a",
          "type": "address"
        },
        {
          "internalType": "enum AddressList.Direction",
          "name": "d",
          "type": "uint8"
        }
      ],
      "name": "addBlacklist",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bytes32",
          "name": "sig",
          "type": "bytes32"
        },
        {
          "internalType": "uint128",
          "name": "checkIdx",
          "type": "uint128"
        },
        {
          "internalType": "enum AddressList.CheckType",
          "name": "tp",
          "type": "uint8"
        }
      ],
      "name": "addOrUpdateRule",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "admin",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "blackLastUpdatedNumber",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "newAdmin",
          "type": "address"
        }
      ],
      "name": "commitChangeAdmin",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "confirmChangeAdmin",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getBlacksFrom",
      "outputs": [
        {
          "internalType": "address[]",
          "name": "",
          "type": "address[]"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getBlacksTo",
      "outputs": [
        {
          "internalType": "address[]",
          "name": "",
          "type": "address[]"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint32",
          "name": "i",
          "type": "uint32"
        }
      ],
      "name": "getRuleByIndex",
      "outputs": [
        {
          "internalType": "bytes32",
          "name": "",
          "type": "bytes32"
        },
        {
          "internalType": "uint128",
          "name": "",
          "type": "uint128"
        },
        {
          "internalType": "enum AddressList.CheckType",
          "name": "",
          "type": "uint8"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bytes32",
          "name": "sig",
          "type": "bytes32"
        },
        {
          "internalType": "uint128",
          "name": "checkIdx",
          "type": "uint128"
        }
      ],
      "name": "getRuleByKey",
      "outputs": [
        {
          "internalType": "bytes32",
          "name": "",
          "type": "bytes32"
        },
        {
          "internalType": "uint128",
          "name": "",
          "type": "uint128"
        },
        {
          "internalType": "enum AddressList.CheckType",
          "name": "",
          "type": "uint8"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_admin",
          "type": "address"
        }
      ],
      "name": "initialize",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "initialized",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "a",
          "type": "address"
        }
      ],
      "name": "isBlackAddress",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        },
        {
          "internalType": "enum AddressList.Direction",
          "name": "",
          "type": "uint8"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "pendingAdmin",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "a",
          "type": "address"
        },
        {
          "internalType": "enum AddressList.Direction",
          "name": "d",
          "type": "uint8"
        }
      ],
      "name": "removeBlacklist",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bytes32",
          "name": "sig",
          "type": "bytes32"
        },
        {
          "internalType": "uint128",
          "name": "checkIdx",
          "type": "uint128"
        }
      ],
      "name": "removeRule",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "rulesLastUpdatedNumber",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "rulesLen",
      "outputs": [
        {
          "internalType": "uint32",
          "name": "",
          "type": "uint32"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    }
  ]`

var (
	BlackLastUpdatedNumberPosition = common.BytesToHash([]byte{0x07})
	RulesLastUpdatedNumberPosition = common.BytesToHash([]byte{0x08})
)

var (
	ValidatorsContractName  = "validators"
	PunishContractName      = "punish"
	SysGovContractName      = "governance"
	AddressListContractName = "address_list"
	ValidatorsContractAddr  = common.HexToAddress("0x000000000000000000000000000000000000d001")
	PunishContractAddr      = common.HexToAddress("0x000000000000000000000000000000000000D002")
	SysGovContractAddr      = common.HexToAddress("0x000000000000000000000000000000000000D003")
	AddressListContractAddr = common.HexToAddress("0x000000000000000000000000000000000000D004")
	// SysGovToAddr is the To address for the system governance transaction, NOT contract address
	SysGovToAddr = common.HexToAddress("0x000000000000000000000000000000000000ffff")
	// engine caller is a dedicated address for the Engine code to interactive with the system contracts.
	EngineCaller = common.HexToAddress("0x0000000000000000004E506F5320456E67696e65")

	abiMap map[string]abi.ABI
)

func init() {
	abiMap = make(map[string]abi.ABI, 0)
	tmpABI, _ := abi.JSON(strings.NewReader(ValidatorsInteractiveABI))
	abiMap[ValidatorsContractName] = tmpABI
	tmpABI, _ = abi.JSON(strings.NewReader(PunishInteractiveABI))
	abiMap[PunishContractName] = tmpABI
	tmpABI, _ = abi.JSON(strings.NewReader(SysGovInteractiveABI))
	abiMap[SysGovContractName] = tmpABI
	tmpABI, _ = abi.JSON(strings.NewReader(AddrListInteractiveABI))
	abiMap[AddressListContractName] = tmpABI
}

func GetInteractiveABI() map[string]abi.ABI {
	return abiMap
}

func GetValidatorAddr(blockNum *big.Int, config *params.ChainConfig) *common.Address {
	return &ValidatorsContractAddr
}

func GetPunishAddr(blockNum *big.Int, config *params.ChainConfig) *common.Address {
	return &PunishContractAddr
}
